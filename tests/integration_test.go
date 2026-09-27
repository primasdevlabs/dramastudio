package tests

import (
	"context"
	"testing"

	agentsApp "dramastudio/internal/agents/application/services"
	agentsInfra "dramastudio/internal/agents/infrastructure/persistence"
	canonApp "dramastudio/internal/canon/application/services"
	canonInfra "dramastudio/internal/canon/infrastructure/persistence"
	charsApp "dramastudio/internal/characters/application/services"
	charsInfra "dramastudio/internal/characters/infrastructure/persistence"
	contApp "dramastudio/internal/continuity/application/services"
	contDomain "dramastudio/internal/continuity/domain"
	contInfra "dramastudio/internal/continuity/infrastructure/persistence"
	identApp "dramastudio/internal/identity/application/services"
	identInfra "dramastudio/internal/identity/infrastructure/persistence"
	mediaApp "dramastudio/internal/media/application/services"
	mediaDomain "dramastudio/internal/media/domain"
	mediaInfra "dramastudio/internal/media/infrastructure/persistence"
	"dramastudio/internal/platform/ai/capabilities"
	"dramastudio/internal/platform/workflow"
	"dramastudio/internal/platform/workflow/activities"
	postApp "dramastudio/internal/postproduction/application/services"
	postDomain "dramastudio/internal/postproduction/domain"
	postInfra "dramastudio/internal/postproduction/infrastructure/persistence"
	"dramastudio/internal/postproduction/infrastructure/ffmpeg"
	prodApp "dramastudio/internal/production/application/services"
	prodDomain "dramastudio/internal/production/domain"
	prodInfra "dramastudio/internal/production/infrastructure/persistence"
	projApp "dramastudio/internal/projects/application/services"
	projDomain "dramastudio/internal/projects/domain"
	projInfra "dramastudio/internal/projects/infrastructure/persistence"
	pubApp "dramastudio/internal/publishing/application/services"
	pubInfra "dramastudio/internal/publishing/infrastructure/persistence"
	storyApp "dramastudio/internal/story/application/services"
	storyInfra "dramastudio/internal/story/infrastructure/persistence"
	worldApp "dramastudio/internal/world/application/services"
	worldInfra "dramastudio/internal/world/infrastructure/persistence"
)

func TestFullDramaStudioWorkflow(t *testing.T) {
	ctx := context.Background()

	// 1. Projects
	projRepo := projInfra.NewInMemoryProjectRepository()
	projSvc := projApp.NewProjectService(projRepo)

	proj, err := projSvc.CreateProject(ctx, "The Last Promise", "Serialized Drama", "Drama/Thriller", "en", projDomain.ProductionModeMonitored)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	bible, err := projSvc.SaveSeriesBible(ctx, proj.ID, "A betraying truth", "Drama", []string{"Betrayal", "Power"}, []string{"No time travel"}, []string{"Action causes consequence"})
	if err != nil {
		t.Fatalf("Failed to save series bible: %v", err)
	}
	if bible.Version != 1 {
		t.Errorf("Expected bible version 1, got %d", bible.Version)
	}

	// 2. Identity
	identRepo := identInfra.NewInMemoryUserRepository()
	identSvc := identApp.NewIdentityService(identRepo)
	user, err := identSvc.CreateUser(ctx, "director@dramastudio.ai", "Lead Director", "Director", "org_1")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 3. Characters
	charRepo := charsInfra.NewInMemoryCharacterRepository()
	charSvc := charsApp.NewCharacterService(charRepo)
	char, err := charSvc.CreateCharacter(ctx, string(proj.ID), "Sarah Johnson", "Protagonist", "Investigative journalist")
	if err != nil {
		t.Fatalf("Failed to create character: %v", err)
	}

	// 4. World
	worldRepo := worldInfra.NewInMemoryWorldRepository()
	worldSvc := worldApp.NewWorldService(worldRepo)
	loc, err := worldSvc.CreateLocation(ctx, string(proj.ID), "Downtown Apartment", "Sarah's living room", "Interior")
	if err != nil {
		t.Fatalf("Failed to create location: %v", err)
	}

	// 5. Story
	storyRepo := storyInfra.NewInMemoryStoryRepository()
	storySvc := storyApp.NewStoryService(storyRepo)

	season, err := storySvc.CreateSeason(ctx, string(proj.ID), "Season 1: Rising Shadows", "Initial arc", 1)
	if err != nil {
		t.Fatalf("Failed to create season: %v", err)
	}

	ep, err := storySvc.CreateEpisode(ctx, season.ID, "Episode 01: Midnight Call", "The initial incident", 1)
	if err != nil {
		t.Fatalf("Failed to create episode: %v", err)
	}

	scene, err := storySvc.CreateScene(ctx, ep.ID, "Scene 01: The Discovery", "Sarah receives an anonymous tip", loc.ID, "Night", 1, []string{string(char.ID)})
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 6. Canon
	canonRepo := canonInfra.NewInMemoryCanonRepository()
	canonSvc := canonApp.NewCanonService(canonRepo)
	fact, err := canonSvc.CreateFact(ctx, char.Name, "received_tip", "confidential_file", ep.ID, ep.ID)
	if err != nil {
		t.Fatalf("Failed to create fact: %v", err)
	}
	if fact.Status != "canonical" {
		t.Errorf("Expected canonical status, got %s", fact.Status)
	}

	// 7. Media Generation (Wan Provider)
	mediaRepo := mediaInfra.NewInMemoryMediaRepository()
	mediaSvc := mediaApp.NewMediaService(mediaRepo, nil)
	asset, err := mediaSvc.GenerateAsset(ctx, string(proj.ID), string(char.ID), scene.ID, "shot_001", "Sarah looking suspiciously at phone in dark apartment", "wan", mediaDomain.MediaTypeVideo)
	if err != nil {
		t.Fatalf("Failed to generate asset: %v", err)
	}
	if asset.Provider != "wan" {
		t.Errorf("Expected provider wan, got %s", asset.Provider)
	}

	// 8. Continuity Check
	contRepo := contInfra.NewInMemoryContinuityRepository()
	contSvc := contApp.NewContinuityService(contRepo)
	issue, err := contSvc.RunCheck(ctx, string(proj.ID), ep.ID, scene.ID, "Wardrobe", contDomain.SeverityWarning, char.Name, "Black blouse", "Red dress", "Missing change event", "Visual mismatch shot 1 vs 2")
	if err != nil {
		t.Fatalf("Failed to run continuity check: %v", err)
	}

	if err := contSvc.ResolveIssue(ctx, issue.ID, "Corrected asset prompt to specify black blouse"); err != nil {
		t.Fatalf("Failed to resolve issue: %v", err)
	}

	// 9. Agents (Lead Director)
	agentRepo := agentsInfra.NewInMemoryAgentRepository()
	agentSvc := agentsApp.NewAgentService(agentRepo)
	dec, err := agentSvc.TriggerLeadDirectorStep(ctx, string(proj.ID), ep.ID)
	if err != nil {
		t.Fatalf("Failed to trigger Lead Director step: %v", err)
	}
	if dec.DecisionMaker != "LEAD_DIRECTOR" {
		t.Errorf("Expected LEAD_DIRECTOR decision maker, got %s", dec.DecisionMaker)
	}

	// 10. Production & Approval
	prodRepo := prodInfra.NewInMemoryProductionRepository()
	prodSvc := prodApp.NewProductionService(prodRepo)

	job, err := prodSvc.StartProductionJob(ctx, "shot_001")
	if err != nil {
		t.Fatalf("Failed to start production job: %v", err)
	}

	app, err := prodSvc.SubmitApproval(ctx, string(proj.ID), ep.ID, "Episode", ep.ID, prodDomain.DecisionApprove, "Episode approved for postproduction", user.Name)
	if err != nil {
		t.Fatalf("Failed to submit approval: %v", err)
	}
	if app.Decision != prodDomain.DecisionApprove {
		t.Errorf("Expected APPROVE, got %s", app.Decision)
	}

	// 11. Postproduction
	postRepo := postInfra.NewInMemoryPostproductionRepository()
	postSvc := postApp.NewPostproductionService(postRepo)
	render, err := postSvc.CreateRender(ctx, ep.ID, "mp4", "1080x1920")
	if err != nil {
		t.Fatalf("Failed to create render: %v", err)
	}

	// 12. Publishing
	pubRepo := pubInfra.NewInMemoryPublishingRepository()
	pubSvc := pubApp.NewPublishingService(pubRepo)
	pub, err := pubSvc.CreatePublication(ctx, ep.ID, "tiktok_channel_1", ep.Title, "Watch Episode 1 now!", []string{"#drama", "#thriller"})
	if err != nil {
		t.Fatalf("Failed to create publication: %v", err)
	}
	if pub.EpisodeID != ep.ID {
		t.Errorf("Expected episode ID %s, got %s", ep.ID, pub.EpisodeID)
	}

	t.Logf("DramaStudio end-to-end integration test passed! Render URL: %s, Publication ID: %s, Job ID: %s", render.OutputURL, pub.ID, job.ID)
}

func TestProduceEpisodeWorkflow(t *testing.T) {
	ctx := context.Background()
	act := activities.NewEpisodeActivities()
	wf := workflow.NewProduceEpisodeWorkflow(act)

	res, err := wf.Execute(ctx, "proj_001", "ep_001")
	if err != nil {
		t.Fatalf("ProduceEpisodeWorkflow failed: %v", err)
	}
	if res.Status != "COMPLETED" {
		t.Errorf("Expected status COMPLETED, got %s", res.Status)
	}
	if len(res.GeneratedShots) != 2 {
		t.Errorf("Expected 2 generated shots, got %d", len(res.GeneratedShots))
	}
	t.Logf("ProduceEpisodeWorkflow passed cleanly! Render URL: %s", res.RenderURL)
}

func TestFFmpegAdapterCommandBuilder(t *testing.T) {
	adapter := ffmpeg.NewFFmpegAdapter("ffmpeg")
	_, args := adapter.BuildConcatCommand(ffmpeg.ConcatOptions{
		VideoURLs:  []string{"shot1.mp4", "shot2.mp4"},
		Resolution: "1080x1920",
		OutputPath: "output.mp4",
	})

	if len(args) == 0 {
		t.Fatalf("Expected non-empty args for concat command")
	}

	ctx := context.Background()
	tl := &postDomain.Timeline{
		ID:        "tl_001",
		EpisodeID: "ep_001",
		VideoTracks: []postDomain.TrackItem{
			{ID: "t1", AssetURL: "shot1.mp4"},
			{ID: "t2", AssetURL: "shot2.mp4"},
		},
	}
	renderURL, err := adapter.RenderTimeline(ctx, tl, "output.mp4")
	if err != nil {
		t.Fatalf("RenderTimeline error: %v", err)
	}
	if renderURL == "" {
		t.Errorf("Expected non-empty render URL")
	}
	t.Logf("FFmpegAdapter test passed! Render URL: %s", renderURL)
}

func TestCapabilityRouter(t *testing.T) {
	ctx := context.Background()
	reg := capabilities.NewModelRegistry()
	router := capabilities.NewCapabilityRouter(reg)

	resp, err := router.Execute(ctx, capabilities.GenerationRequest{
		Capability: capabilities.CapVideoGen,
		ProjectID:  "proj_001",
		EpisodeID:  "ep_001",
		Prompt: capabilities.StructuredPrompt{
			Subject:  "Sarah",
			Action:   "looking at phone",
			ShotType: "medium_close_up",
		},
	})
	if err != nil {
		t.Fatalf("CapabilityRouter Execute failed: %v", err)
	}
	if resp.Provider != "wan" {
		t.Errorf("Expected provider wan for video generation, got %s", resp.Provider)
	}
	t.Logf("CapabilityRouter test passed! Provider: %s, Model: %s, URL: %s", resp.Provider, resp.Model, resp.OutputURL)
}
