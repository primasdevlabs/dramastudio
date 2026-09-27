package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	agentsApp "dramastudio/internal/agents/application/services"
	agentsInfra "dramastudio/internal/agents/infrastructure/persistence"
	agentsHTTP "dramastudio/internal/agents/interfaces/http"

	canonApp "dramastudio/internal/canon/application/services"
	canonInfra "dramastudio/internal/canon/infrastructure/persistence"
	canonHTTP "dramastudio/internal/canon/interfaces/http"

	charsApp "dramastudio/internal/characters/application/services"
	charsInfra "dramastudio/internal/characters/infrastructure/persistence"
	charsHTTP "dramastudio/internal/characters/interfaces/http"

	contApp "dramastudio/internal/continuity/application/services"
	contInfra "dramastudio/internal/continuity/infrastructure/persistence"
	contHTTP "dramastudio/internal/continuity/interfaces/http"

	identApp "dramastudio/internal/identity/application/services"
	identInfra "dramastudio/internal/identity/infrastructure/persistence"
	identHTTP "dramastudio/internal/identity/interfaces/http"

	mediaApp "dramastudio/internal/media/application/services"
	mediaInfra "dramastudio/internal/media/infrastructure/persistence"
	mediaHTTP "dramastudio/internal/media/interfaces/http"

	postApp "dramastudio/internal/postproduction/application/services"
	postInfra "dramastudio/internal/postproduction/infrastructure/persistence"
	postHTTP "dramastudio/internal/postproduction/interfaces/http"

	prodApp "dramastudio/internal/production/application/services"
	prodInfra "dramastudio/internal/production/infrastructure/persistence"
	prodHTTP "dramastudio/internal/production/interfaces/http"

	projApp "dramastudio/internal/projects/application/services"
	projInfra "dramastudio/internal/projects/infrastructure/persistence"
	projHTTP "dramastudio/internal/projects/interfaces/http"

	pubApp "dramastudio/internal/publishing/application/services"
	pubInfra "dramastudio/internal/publishing/infrastructure/persistence"
	pubHTTP "dramastudio/internal/publishing/interfaces/http"

	storyApp "dramastudio/internal/story/application/services"
	storyInfra "dramastudio/internal/story/infrastructure/persistence"
	storyHTTP "dramastudio/internal/story/interfaces/http"

	worldApp "dramastudio/internal/world/application/services"
	worldInfra "dramastudio/internal/world/infrastructure/persistence"
	worldHTTP "dramastudio/internal/world/interfaces/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Initialize Infrastructure Repositories
	identityRepo := identInfra.NewInMemoryUserRepository()
	projectsRepo := projInfra.NewInMemoryProjectRepository()
	storyRepo := storyInfra.NewInMemoryStoryRepository()
	canonRepo := canonInfra.NewInMemoryCanonRepository()
	characterRepo := charsInfra.NewInMemoryCharacterRepository()
	worldRepo := worldInfra.NewInMemoryWorldRepository()
	agentRepo := agentsInfra.NewInMemoryAgentRepository()
	productionRepo := prodInfra.NewInMemoryProductionRepository()
	continuityRepo := contInfra.NewInMemoryContinuityRepository()
	mediaRepo := mediaInfra.NewInMemoryMediaRepository()
	postproductionRepo := postInfra.NewInMemoryPostproductionRepository()
	publishingRepo := pubInfra.NewInMemoryPublishingRepository()

	// 2. Initialize Application Services
	identitySvc := identApp.NewIdentityService(identityRepo)
	projectsSvc := projApp.NewProjectService(projectsRepo)
	storySvc := storyApp.NewStoryService(storyRepo)
	canonSvc := canonApp.NewCanonService(canonRepo)
	characterSvc := charsApp.NewCharacterService(characterRepo)
	worldSvc := worldApp.NewWorldService(worldRepo)
	agentSvc := agentsApp.NewAgentService(agentRepo)
	productionSvc := prodApp.NewProductionService(productionRepo)
	continuitySvc := contApp.NewContinuityService(continuityRepo)
	mediaSvc := mediaApp.NewMediaService(mediaRepo, nil)
	postproductionSvc := postApp.NewPostproductionService(postproductionRepo)
	publishingSvc := pubApp.NewPublishingService(publishingRepo)

	// 3. Initialize HTTP Handlers
	identityHandler := identHTTP.NewIdentityHandler(identitySvc)
	projectsHandler := projHTTP.NewProjectsHandler(projectsSvc)
	storyHandler := storyHTTP.NewStoryHandler(storySvc)
	canonHandler := canonHTTP.NewCanonHandler(canonSvc)
	charactersHandler := charsHTTP.NewCharactersHandler(characterSvc)
	worldHandler := worldHTTP.NewWorldHandler(worldSvc)
	agentsHandler := agentsHTTP.NewAgentsHandler(agentSvc)
	productionHandler := prodHTTP.NewProductionHandler(productionSvc)
	continuityHandler := contHTTP.NewContinuityHandler(continuitySvc)
	mediaHandler := mediaHTTP.NewMediaHandler(mediaSvc)
	postproductionHandler := postHTTP.NewPostproductionHandler(postproductionSvc)
	publishingHandler := pubHTTP.NewPublishingHandler(publishingSvc)

	// 4. Setup Router Mux
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","system":"DramaStudio API","version":"1.0.0"}`))
	})

	mux.Handle("/v1/identity/", identityHandler)
	mux.Handle("/v1/projects/", projectsHandler)
	mux.Handle("/v1/projects", projectsHandler)
	mux.Handle("/v1/story/", storyHandler)
	mux.Handle("/v1/canon/", canonHandler)
	mux.Handle("/v1/characters/", charactersHandler)
	mux.Handle("/v1/characters", charactersHandler)
	mux.Handle("/v1/world/", worldHandler)
	mux.Handle("/v1/agents/", agentsHandler)
	mux.Handle("/v1/production/", productionHandler)
	mux.Handle("/v1/continuity/", continuityHandler)
	mux.Handle("/v1/media/", mediaHandler)
	mux.Handle("/v1/postproduction/", postproductionHandler)
	mux.Handle("/v1/publishing/", publishingHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting DramaStudio API server on port %s...", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down API server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced shutdown: %v", err)
	}
	log.Println("API server stopped.")
}
