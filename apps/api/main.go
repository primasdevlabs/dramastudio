package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dramastudio/configs"
	agentsHTTP "dramastudio/internal/agents/interfaces/http"
	aiHTTP "dramastudio/internal/ai/interfaces/http"
	analyticsHTTP "dramastudio/internal/analytics/interfaces/http"
	"dramastudio/internal/appwiring"
	canonHTTP "dramastudio/internal/canon/interfaces/http"
	charsHTTP "dramastudio/internal/characters/interfaces/http"
	contHTTP "dramastudio/internal/continuity/interfaces/http"
	"dramastudio/internal/identity/infrastructure/sso"
	identHTTP "dramastudio/internal/identity/interfaces/http"
	intelHTTP "dramastudio/internal/intelligence/interfaces/http"
	mediaHTTP "dramastudio/internal/media/interfaces/http"
	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/platform/observability/tracing"
	"dramastudio/internal/platform/security"
	postHTTP "dramastudio/internal/postproduction/interfaces/http"
	prodHTTP "dramastudio/internal/production/interfaces/http"
	projectsdomain "dramastudio/internal/projects/domain"
	projHTTP "dramastudio/internal/projects/interfaces/http"
	pubHTTP "dramastudio/internal/publishing/interfaces/http"
	storyHTTP "dramastudio/internal/story/interfaces/http"
	worldHTTP "dramastudio/internal/world/interfaces/http"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}
	level := slog.LevelInfo
	_ = level.UnmarshalText([]byte(cfg.Observability.LogLevel))
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// OpenTelemetry provider: stdout exporter when configured, no-op
	// otherwise (§62).
	shutdownTrace, err := tracing.InitProvider(ctx, cfg.Observability.ServiceName, cfg.Observability.TraceExporter, os.Stderr)
	if err != nil {
		slog.Error("tracing", "err", err)
		os.Exit(1)
	}
	defer func() { _ = shutdownTrace(context.Background()) }()

	container, err := appwiring.Build(ctx, cfg)
	if err != nil {
		slog.Error("wiring", "err", err)
		os.Exit(1)
	}
	defer container.Close()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      routes(cfg, container),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("api listening", "addr", server.Addr, "store", cfg.Store, "storage", cfg.Storage.Backend)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("api", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdown); err != nil {
		slog.Error("shutdown", "err", err)
	}
}

// routes assembles the full handler graph: health, SSE, then every
// bounded-context handler under the middleware chain.
func routes(cfg *configs.Config, c *appwiring.Container) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		platformhttp.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok", "service": cfg.Observability.ServiceName,
		})
	})
	mux.Handle("GET /v1/events", c.Broker)

	var googleVerifier *sso.GoogleVerifier
	if cfg.Security.GoogleClientID != "" {
		googleVerifier = sso.NewGoogleVerifier(cfg.Security.GoogleClientID)
	}
	identHTTP.NewIdentityHandler(c.Identity, c.JWT, googleVerifier, cfg.Environment == "development").RegisterRoutes(mux)
	projHTTP.NewProjectsHandler(c.Projects).RegisterRoutes(mux)
	storyHTTP.NewStoryHandler(c.Story).RegisterRoutes(mux)
	canonHTTP.NewCanonHandler(c.Canon).RegisterRoutes(mux)
	charsHTTP.NewCharactersHandler(c.Characters).RegisterRoutes(mux)
	worldHTTP.NewWorldHandler(c.World).RegisterRoutes(mux)
	agentsHTTP.NewAgentsHandler(c.Agents).RegisterRoutes(mux)
	prodHTTP.NewProductionHandler(c.Production).RegisterRoutes(mux)
	contHTTP.NewContinuityHandler(c.Continuity).RegisterRoutes(mux)
	mediaH := mediaHTTP.NewMediaHandler(c.Media, webhookSecrets())
	mediaH.SetDeliveries(c.WebhookDeliveries)
	mediaH.RegisterRoutes(mux)
	postH := postHTTP.NewPostproductionHandler(c.Postprod)
	postH.SetEpisodeResolver(c.Story.ProjectOfEpisode)
	postH.RegisterRoutes(mux)
	pubH := pubHTTP.NewPublishingHandler(c.Publishing)
	pubH.SetEpisodeResolver(c.Story.ProjectOfEpisode)
	pubH.RegisterRoutes(mux)
	aiHTTP.NewAIHandler(c.AIPolicies, c.AIAdmin, c.AIExecutor, c.ModelRegistry).RegisterRoutes(mux)
	projectOrg := func(ctx context.Context, projectID string) (string, error) {
		proj, err := c.Projects.GetProject(ctx, projectsdomain.ProjectID(projectID))
		if err != nil {
			return "", err
		}
		return proj.OrgID, nil
	}
	intelHTTP.NewIntelligenceHandler(c.Intel, platformhttp.ProjectResolver(projectOrg)).RegisterRoutes(mux)
	analyticsHTTP.NewAnalyticsHandler(c.Analytics).RegisterRoutes(mux)

	var authn platformhttp.Authenticator
	if c.JWT != nil {
		authn = c.JWT
	}
	return platformhttp.Chain(mux,
		platformhttp.RequestID,
		platformhttp.Recoverer,
		platformhttp.CORS(cfg.Server.AllowedOrigins),
		platformhttp.Auth(platformhttp.AuthConfig{
			RequireAuth: cfg.Security.RequireAuth,
			DevUserID:   cfg.Security.DevUserID,
			DevOrgID:    cfg.Security.DevOrgID,
			APIKeys: func(ctx context.Context, key string) (security.Principal, error) {
				cred, err := c.Identity.AuthenticateAPIKey(ctx, key)
				if err != nil {
					return security.Principal{}, err
				}
				return security.Principal{
					UserID:      cred.UserID,
					OrgID:       cred.OrgID,
					Roles:       []string{"service"},
					Permissions: cred.Permissions,
					Service:     true,
				}, nil
			},
		}, authn),
		// Org-scope every /v1/projects/{id}/... route across all modules (§8).
		platformhttp.ProjectGuard(platformhttp.ProjectResolver(projectOrg)),
	)
}

// webhookSecrets maps provider name -> HMAC secret env value (§40).
func webhookSecrets() map[string]string {
	out := map[string]string{}
	for _, p := range []string{"wan", "openai", "mock"} {
		if v := os.Getenv("WEBHOOK_SECRET_" + p); v != "" {
			out[p] = v
		}
	}
	return out
}
