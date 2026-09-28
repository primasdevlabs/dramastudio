// Package appwiring is the composition root shared by the api, worker,
// and scheduler binaries. It builds repositories honoring the configured
// store backend, wires services, and registers provider adapters (§6).
package appwiring

import (
	"context"
	"fmt"
	"log/slog"

	"dramastudio/configs"

	agentsApp "dramastudio/internal/agents/application/services"
	"dramastudio/internal/agents/director"

	aiApp "dramastudio/internal/ai/application"
	aidomain "dramastudio/internal/ai/domain"
	intelApp "dramastudio/internal/intelligence/application"

	analyticsApp "dramastudio/internal/analytics/application/services"

	canonApp "dramastudio/internal/canon/application/services"

	charsApp "dramastudio/internal/characters/application/services"

	contApp "dramastudio/internal/continuity/application/services"

	identApp "dramastudio/internal/identity/application/services"

	mediaApp "dramastudio/internal/media/application/services"

	postApp "dramastudio/internal/postproduction/application/services"
	postprodffmpeg "dramastudio/internal/postproduction/infrastructure/ffmpeg"

	prodApp "dramastudio/internal/production/application/services"

	projApp "dramastudio/internal/projects/application/services"

	pubApp "dramastudio/internal/publishing/application/services"
	pubinstagram "dramastudio/internal/publishing/infrastructure/instagram"
	pubtiktok "dramastudio/internal/publishing/infrastructure/tiktok"
	pubyoutube "dramastudio/internal/publishing/infrastructure/youtube"

	storyApp "dramastudio/internal/story/application/services"

	worldApp "dramastudio/internal/world/application/services"

	"dramastudio/internal/platform/ai/provider"
	"dramastudio/internal/platform/ai/routing"
	"dramastudio/internal/platform/audit"
	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/platform/database/sqlite"
	"dramastudio/internal/platform/events"
	platformhttp "dramastudio/internal/platform/http"

	// Adapter packages register their factories via init(). Adding a new
	// provider kind to the catalog requires no changes here.
	_ "dramastudio/internal/platform/ai/provider/anthropic"
	mockprovider "dramastudio/internal/platform/ai/provider/mock"
	_ "dramastudio/internal/platform/ai/provider/musicgen"
	_ "dramastudio/internal/platform/ai/provider/openai_compatible"
	_ "dramastudio/internal/platform/ai/provider/voicegen"
	_ "dramastudio/internal/platform/ai/provider/wan"
	"dramastudio/internal/platform/security"
	"dramastudio/internal/platform/storage"
	"dramastudio/internal/platform/workflow"
	"dramastudio/internal/platform/workflow/activities"
	wflocal "dramastudio/internal/platform/workflow/local"
	"dramastudio/internal/platform/workflow/temporal"
)

// Container holds every wired service plus closable infrastructure.
type Container struct {
	DB     *postgres.DB
	SQLite *sqlite.DB
	Store  storage.ObjectStorage
	JWT    *security.JWTService

	// Events is the shared domain event bus; Broker fans events out to SSE
	// clients (§48, §58). Audit records append-only action logs (§63).
	Events            *events.Bus
	Broker            *platformhttp.SSEBroker
	Audit             audit.Logger
	WebhookDeliveries *events.WebhookDeliveryStore // nil without SQL

	Resolver *routing.Resolver

	Identity      *identApp.IdentityService
	Projects      *projApp.ProjectService
	Story         *storyApp.StoryService
	Canon         *canonApp.CanonService
	Characters    *charsApp.CharacterService
	World         *worldApp.WorldService
	Agents        *agentsApp.AgentService
	Production    *prodApp.ProductionService
	Continuity    *contApp.ContinuityService
	Media         *mediaApp.MediaService
	Postprod      *postApp.PostproductionService
	Publishing    *pubApp.PublishingService
	Analytics     *analyticsApp.AnalyticsService
	AIPolicies    *aiApp.PolicyService
	AIAdmin       *aiApp.ProviderAdminService
	AIExecutor    *aiApp.ExecuteGenerationHandler
	ModelRegistry aidomain.ModelRegistryRepository
	Intel         *intelApp.IntelligenceService
	Activities    *activities.EpisodeActivities
	Engine        workflow.Engine
}

// Close releases pooled infrastructure.
func (c *Container) Close() {
	if c.Engine != nil {
		c.Engine.Close()
	}
	if c.DB != nil {
		c.DB.Close()
	}
	if c.SQLite != nil {
		c.SQLite.Close()
	}
}

// Build constructs the full dependency graph for a process.
func Build(ctx context.Context, cfg *configs.Config) (*Container, error) {
	c := &Container{}

	store, err := buildStorage(cfg)
	if err != nil {
		return nil, err
	}
	c.Store = store

	r, err := buildRepos(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if r.db != nil {
		c.DB = r.db
	}
	if r.sqldb != nil {
		c.SQLite = r.sqldb
	}

	// --- domain events + realtime ---
	// One bus per process: services emit, the log persists (§48), and the
	// SSE broker streams to studio clients (§58).
	var eventLog events.EventLogStore = events.NoopEventLog{}
	var querier postgres.Querier
	if r.db != nil {
		querier = r.db.Pool
	}
	if r.sqldb != nil {
		querier = r.sqldb
	}
	if querier != nil {
		eventLog = events.NewSQLEventLog(querier)
	}
	c.Events = events.NewBus(eventLog)
	c.Broker = platformhttp.NewSSEBroker()
	c.Audit = audit.NewLogger(querier)
	c.WebhookDeliveries = events.NewWebhookDeliveryStore(querier)
	c.Events.SubscribeAll(func(_ context.Context, e events.DomainEvent) error {
		c.Broker.Publish(e.Type, e)
		return nil
	})

	// --- AI provider routing ---
	c.Resolver = routing.NewResolver(aiApp.NewRoutingRegistry(r.modelRegistry))
	c.AIAdmin = aiApp.NewProviderAdminService(r.modelRegistry, c.Resolver)

	// Seed the model catalog into an empty registry (editable data, not
	// code): providers, models, and system-scope policies for every
	// production capability.
	if path := envOr("MODEL_CATALOG", "configs/model-catalog.json"); path != "" {
		if n, err := aiApp.SeedCatalog(ctx, r.modelRegistry, path); err != nil {
			slog.Warn("model catalog seed skipped", "err", err)
		} else if n > 0 {
			slog.Info("model catalog seeded", "models", n)
		}
	}
	c.registerProviderAdapters(ctx, r)

	// --- security ---
	jwtSvc, err := security.NewJWTService(cfg.Security.JWTSecret, cfg.Security.AccessTokenTTL)
	if err != nil {
		if cfg.Security.RequireAuth {
			return nil, fmt.Errorf("jwt: %w", err)
		}
		jwtSvc = nil
	}
	c.JWT = jwtSvc

	c.buildServices(cfg, r)
	if err := c.buildIntelligence(r); err != nil {
		return nil, fmt.Errorf("intelligence: %w", err)
	}
	if err := c.buildEngine(ctx, cfg); err != nil {
		return nil, err
	}
	return c, nil
}

// buildEngine selects the workflow engine adapter (§83). Production code
// only sees workflow.Engine; Temporal is just one implementation.
func (c *Container) buildEngine(ctx context.Context, cfg *configs.Config) error {
	switch cfg.Workflow.Engine {
	case configs.EngineTemporal:
		tc, err := temporal.Dial(ctx, temporal.Config{
			HostPort:  cfg.Workflow.HostPort,
			Namespace: cfg.Workflow.Namespace,
			APIKey:    cfg.Workflow.APIKey,
			UseTLS:    cfg.Workflow.UseTLS,
		})
		if err != nil {
			return fmt.Errorf("workflow temporal: %w", err)
		}
		c.Engine = temporal.NewOrchestrator(tc)
	default:
		c.Engine = wflocal.New(c.Activities)
	}
	c.Production.SetEngine(c.Engine)
	return nil
}

// registerProviderAdapters builds a live adapter for every provider row
// whose type has a registered factory. Rows with pending adapter kinds
// (e.g. fal, replicate before their native adapters land) are skipped
// with a warning — they remain registry data and appear in the UI as
// unavailable until an adapter registers their kind.
func (c *Container) registerProviderAdapters(ctx context.Context, r repoBundle) {
	providers, err := r.modelRegistry.ListProviders(ctx)
	if err != nil {
		slog.Warn("list providers for adapter wiring", "err", err)
		return
	}
	for _, p := range providers {
		if err := c.AIAdmin.RegisterAdapter(p); err != nil {
			slog.Warn("provider adapter unavailable", "provider", p.ID, "type", p.Type, "err", err)
		}
	}
	// Mock is always available for local/dev even without a catalog row.
	if _, ok := c.Resolver.AdapterFor("mock"); !ok {
		c.Resolver.RegisterProvider("mock", provider.Adapter{
			Sync:  mockprovider.NewSync("mock"),
			Async: mockprovider.NewAsync("mock"),
		})
	}
}

// buildStorage selects the object storage backend (§46).
func buildStorage(cfg *configs.Config) (storage.ObjectStorage, error) {
	if cfg.Storage.Backend == configs.StorageS3 {
		s3, err := storage.NewS3(cfg.Storage.Endpoint, cfg.Storage.Region, cfg.Storage.Bucket,
			cfg.Storage.AccessKey, cfg.Storage.SecretKey, cfg.Storage.UseSSL)
		if err != nil {
			return nil, fmt.Errorf("storage s3: %w", err)
		}
		return s3, nil
	}
	l, err := storage.NewLocal(cfg.Storage.LocalDir)
	if err != nil {
		return nil, fmt.Errorf("storage local: %w", err)
	}
	return l, nil
}

// buildServices instantiates application services over the repositories.
func (c *Container) buildServices(cfg *configs.Config, r repoBundle) {
	c.Identity = identApp.NewIdentityService(r.identity)
	c.Projects = projApp.NewProjectService(r.projects)
	c.Story = storyApp.NewStoryService(r.story)
	c.Canon = canonApp.NewCanonService(r.canon)
	c.Characters = charsApp.NewCharacterService(r.characters)
	c.World = worldApp.NewWorldService(r.world)
	c.Production = prodApp.NewProductionService(r.production)
	c.Continuity = contApp.NewContinuityService(r.continuity)
	c.Media = mediaApp.NewMediaService(r.media, c.Resolver)
	c.Postprod = postApp.NewPostproductionService(r.postprod,
		postprodffmpeg.NewFFmpegAdapter("ffmpeg"), c.Store, "")
	c.Publishing = pubApp.NewPublishingService(r.publishing)
	c.Analytics = analyticsApp.NewAnalyticsService(r.analytics)
	c.AIPolicies = aiApp.NewPolicyService(r.modelRegistry)
	c.AIExecutor = aiApp.NewExecuteGenerationHandler(r.modelRegistry, c.Resolver)
	c.ModelRegistry = r.modelRegistry

	// Domain event emitters (§48): every state transition that matters to
	// operators or other modules lands on the shared bus.
	c.Projects.SetEvents(c.Events)
	c.Story.SetEvents(c.Events)
	c.Canon.SetEvents(c.Events)
	c.Characters.SetEvents(c.Events)
	c.World.SetEvents(c.Events)
	c.Production.SetEvents(c.Events)
	c.Continuity.SetEvents(c.Events)
	c.Media.SetEvents(c.Events)
	c.Postprod.SetEvents(c.Events)
	c.Publishing.SetEvents(c.Events)

	// Append-only audit trail (§63): identity actions + human decisions.
	c.Identity.SetAudit(c.Audit)
	c.Production.SetAudit(c.Audit)

	// Publishing adapters — channel config (account/token refs) comes from
	// the channel row; base URLs from env so no secrets live in code.
	c.buildAdapters(cfg, r)
}

// buildAdapters wires external-facing adapters: publishing channels, the
// lead director observer, and the shared workflow activities.
func (c *Container) buildAdapters(cfg *configs.Config, r repoBundle) {
	c.Publishing.RegisterAdapter(pubyoutube.New(envOr("YOUTUBE_BASE_URL", ""), envOr("YOUTUBE_API_KEY", "")))
	c.Publishing.RegisterAdapter(pubtiktok.New(envOr("TIKTOK_BASE_URL", ""), envOr("TIKTOK_API_KEY", "")))
	c.Publishing.RegisterAdapter(pubinstagram.New(envOr("INSTAGRAM_BASE_URL", ""), envOr("INSTAGRAM_API_KEY", "")))

	// Lead director observes real production + continuity state.
	observer := NewProductionObserver(c.Production, c.Continuity)
	c.Agents = agentsApp.NewAgentService(r.agents, director.NewLeadDirector(observer, envOr("AGENT_MODE", "monitored")))
	c.Agents.SetEvents(c.Events)

	// --- workflow activities ---
	c.Activities = activities.NewEpisodeActivities(activities.Deps{
		Story:       c.Story,
		Production:  c.Production,
		Media:       c.Media,
		Continuity:  c.Continuity,
		Postprod:    c.Postprod,
		AI:          c.AIExecutor,
		CallbackURL: cfg.Server.PublicURL + "/v1/webhooks/providers/",
	})
}
