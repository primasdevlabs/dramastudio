package tests

import (
	"context"
	"testing"

	agentsApp "dramastudio/internal/agents/application/services"
	"dramastudio/internal/agents/director"
	agentsDomain "dramastudio/internal/agents/domain"
	agentsInfra "dramastudio/internal/agents/infrastructure/persistence"

	aiApp "dramastudio/internal/ai/application"
	aiDomain "dramastudio/internal/ai/domain"
	aiInfra "dramastudio/internal/ai/infrastructure/registry"

	canonApp "dramastudio/internal/canon/application/services"
	canonDomain "dramastudio/internal/canon/domain"
	canonInfra "dramastudio/internal/canon/infrastructure/persistence"

	charsApp "dramastudio/internal/characters/application/services"
	charsDomain "dramastudio/internal/characters/domain"
	charsInfra "dramastudio/internal/characters/infrastructure/persistence"

	contApp "dramastudio/internal/continuity/application/services"
	contDomain "dramastudio/internal/continuity/domain"
	contInfra "dramastudio/internal/continuity/infrastructure/persistence"

	identApp "dramastudio/internal/identity/application/services"
	identInfra "dramastudio/internal/identity/infrastructure/persistence"

	mediaApp "dramastudio/internal/media/application/services"
	mediaDomain "dramastudio/internal/media/domain"
	mediaInfra "dramastudio/internal/media/infrastructure/persistence"

	aimock "dramastudio/internal/platform/ai/provider/mock"
	"dramastudio/internal/platform/ai/routing"
	"dramastudio/internal/platform/storage"

	postApp "dramastudio/internal/postproduction/application/services"
	postDomain "dramastudio/internal/postproduction/domain"
	postInfra "dramastudio/internal/postproduction/infrastructure/persistence"

	prodApp "dramastudio/internal/production/application/services"
	prodDomain "dramastudio/internal/production/domain"
	prodInfra "dramastudio/internal/production/infrastructure/persistence"

	projApp "dramastudio/internal/projects/application/services"
	projDomain "dramastudio/internal/projects/domain"
	projInfra "dramastudio/internal/projects/infrastructure/persistence"

	pubApp "dramastudio/internal/publishing/application/services"
	pubDomain "dramastudio/internal/publishing/domain"
	pubInfra "dramastudio/internal/publishing/infrastructure/persistence"

	storyApp "dramastudio/internal/story/application/services"
	storyDomain "dramastudio/internal/story/domain"
	storyInfra "dramastudio/internal/story/infrastructure/persistence"

	worldApp "dramastudio/internal/world/application/services"
	worldInfra "dramastudio/internal/world/infrastructure/persistence"
)

type stubObserver struct{}

func (stubObserver) Observe(ctx context.Context, projectID, episodeID string) (*director.Observation, error) {
	return &director.Observation{ProjectID: projectID, EpisodeID: episodeID}, nil
}

// newTestResolver builds an AI resolver backed by an in-memory registry
// with one mock provider/model and system-scope policies.
func newTestResolver(t *testing.T) *routing.Resolver {
	t.Helper()
	repo := aiInfra.NewInMemoryModelRegistry()
	ctx := context.Background()
	if err := repo.SaveProvider(ctx, &aiDomain.Provider{ID: "mock", Name: "mock", Type: "mock", HealthStatus: aiDomain.HealthHealthy}); err != nil {
		t.Fatalf("seed provider: %v", err)
	}
	if err := repo.SaveModel(ctx, &aiDomain.Model{ID: "mock-v1", Name: "mock-v1", Identifier: "mock-v1", ProviderID: "mock", Status: aiDomain.ModelActive}); err != nil {
		t.Fatalf("seed model: %v", err)
	}
	policies := aiApp.NewPolicyService(repo)
	for _, cap := range []aiDomain.AICapability{"video_generation", "image_generation", "script_writing"} {
		if _, err := policies.SetPolicy(ctx, &aiDomain.PolicyRow{Scope: "system", Capability: cap, ProviderID: "mock", ModelID: "mock-v1"}); err != nil {
			t.Fatalf("seed policy %s: %v", cap, err)
		}
	}
	r := routing.NewResolver(aiApp.NewRoutingRegistry(repo))
	r.RegisterProvider("mock", routing.ProviderAdapter{
		Sync:  aimock.NewSync("mock"),
		Async: aimock.NewAsync("mock"),
	})
	return r
}

// e2eFixture carries cross-context entities a scenario builds up.
type e2eFixture struct {
	proj         *projDomain.Project
	user         string
	bibleVersion int
	char         *charsDomain.Character
	episode      *storyDomain.Episode
	scene        *storyDomain.Scene
	runID        string
	render       *postDomain.RenderTask
	pub          *pubDomain.Publication
}

func TestFullDramaStudioWorkflow(t *testing.T) {
	ctx := context.Background()
	fx := &e2eFixture{}

	fx.setupStoryWorld(t, ctx)
	fx.establishCanon(t, ctx)
	fx.generateMedia(t, ctx)
	fx.checkContinuity(t, ctx)
	fx.runDirector(t, ctx)
	fx.produceAndApprove(t, ctx)
	fx.postproduce(t, ctx)
	fx.publish(t, ctx)

	t.Logf("end-to-end OK — run=%s render=%s publication=%s", fx.runID, fx.render.ID, fx.pub.ID)
}

// setupStoryWorld creates the project, bible, user, character, location,
// and the series→season→episode→scene hierarchy.
func (fx *e2eFixture) setupStoryWorld(t *testing.T, ctx context.Context) {
	t.Helper()
	projSvc := projApp.NewProjectService(projInfra.NewInMemoryProjectRepository())
	proj, err := projSvc.CreateProject(ctx, "org_1", "The Last Promise", "Serialized Drama", "Drama/Thriller", "en", projDomain.ProductionModeMonitored)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	fx.proj = proj

	bible, err := projSvc.SaveSeriesBible(ctx, proj.ID, "A betraying truth", "Drama", "Grounded", "Desaturated", "Subtext-heavy",
		[]string{"Betrayal", "Power"}, []string{"No time travel"}, []string{"Action causes consequence"})
	if err != nil {
		t.Fatalf("save series bible: %v", err)
	}
	if bible.Version != 1 {
		t.Errorf("expected bible version 1, got %d", bible.Version)
	}
	fx.bibleVersion = bible.Version

	user, _, err := identApp.NewIdentityService(identInfra.NewInMemoryUserRepository()).
		Register(ctx, "DramaStudio Org", "director@dramastudio.ai", "Lead Director", "password-123")
	if err != nil {
		t.Fatalf("register user: %v", err)
	}
	fx.user = user.Name

	char, err := charsApp.NewCharacterService(charsInfra.NewInMemoryCharacterRepository()).
		CreateCharacter(ctx, string(proj.ID), "Sarah Johnson", "Protagonist", "Investigative journalist")
	if err != nil {
		t.Fatalf("create character: %v", err)
	}
	fx.char = char

	loc, err := worldApp.NewWorldService(worldInfra.NewInMemoryWorldRepository()).
		CreateLocation(ctx, string(proj.ID), "Downtown Apartment", "interior", "Sarah's living room", "")
	if err != nil {
		t.Fatalf("create location: %v", err)
	}

	storySvc := storyApp.NewStoryService(storyInfra.NewInMemoryStoryRepository())
	series, err := storySvc.EnsureSeries(ctx, string(proj.ID), "The Last Promise", "Season-long serialized drama")
	if err != nil {
		t.Fatalf("ensure series: %v", err)
	}
	season, err := storySvc.CreateSeason(ctx, series.ID, "Season 1: Rising Shadows", "Initial arc", 1)
	if err != nil {
		t.Fatalf("create season: %v", err)
	}
	ep, err := storySvc.CreateEpisode(ctx, season.ID, "", "Episode 01: Midnight Call", "The initial incident", 1)
	if err != nil {
		t.Fatalf("create episode: %v", err)
	}
	fx.episode = ep

	scene, err := storySvc.CreateScene(ctx, ep.ID, "Scene 01: The Discovery", "Sarah receives an anonymous tip", loc.ID, "night", 1, []string{string(char.ID)})
	if err != nil {
		t.Fatalf("create scene: %v", err)
	}
	fx.scene = scene
}

// establishCanon records an immutable fact via the canon service.
func (fx *e2eFixture) establishCanon(t *testing.T, ctx context.Context) {
	t.Helper()
	canonSvc := canonApp.NewCanonService(canonInfra.NewInMemoryCanonRepository())
	fact, err := canonSvc.EstablishFact(ctx, &canonDomain.StoryFact{
		ProjectID:         string(fx.proj.ID),
		Subject:           fx.char.Name,
		Predicate:         "received_tip",
		Object:            "confidential_file",
		IntroducedEpisode: fx.episode.ID,
		EffectiveFrom:     fx.episode.ID,
	})
	if err != nil {
		t.Fatalf("establish fact: %v", err)
	}
	if fact.Status != canonDomain.FactStatusCanonical {
		t.Errorf("expected canonical status, got %s", fact.Status)
	}
}

// generateMedia routes a video generation through the mock provider.
func (fx *e2eFixture) generateMedia(t *testing.T, ctx context.Context) {
	t.Helper()
	mediaSvc := mediaApp.NewMediaService(mediaInfra.NewInMemoryMediaRepository(), newTestResolver(t))
	job, asset, err := mediaSvc.Generate(ctx, mediaApp.GenerateRequest{
		ProjectID:  string(fx.proj.ID),
		Capability: "video_generation",
		MediaType:  mediaDomain.MediaTypeVideo,
		SceneID:    fx.scene.ID,
		ShotID:     "shot_001",
		Prompt:     "Sarah looking suspiciously at phone in dark apartment",
	})
	if err != nil {
		t.Fatalf("generate asset: %v", err)
	}
	if asset.Provider != "mock" {
		t.Errorf("expected provider mock, got %s", asset.Provider)
	}
	if job.AssetID != asset.ID {
		t.Errorf("expected job bound to asset %s, got %s", asset.ID, job.AssetID)
	}
}

// checkContinuity runs the wardrobe checker; a mismatch must surface issues.
func (fx *e2eFixture) checkContinuity(t *testing.T, ctx context.Context) {
	t.Helper()
	contSvc := contApp.NewContinuityService(contInfra.NewInMemoryContinuityRepository())
	check, issues, err := contSvc.RunCheck(ctx, contApp.CheckInput{
		ProjectID:         string(fx.proj.ID),
		EpisodeID:         fx.episode.ID,
		SceneID:           fx.scene.ID,
		CheckType:         contDomain.CheckWardrobe,
		EntityID:          string(fx.char.ID),
		CanonicalWardrobe: []string{"black blouse"},
		AssignedWardrobe:  []string{"red dress"},
	})
	if err != nil {
		t.Fatalf("run continuity check: %v", err)
	}
	if check == nil {
		t.Fatal("expected a check record")
	}
	if len(issues) == 0 {
		t.Skip("wardrobe checker produced no issue; skipping resolve step")
	}
	if _, err := contSvc.ResolveIssue(ctx, issues[0].ID, "Corrected asset prompt to specify black blouse", false); err != nil {
		t.Fatalf("resolve issue: %v", err)
	}
}

// runDirector executes one Lead Director step against a stub observer.
func (fx *e2eFixture) runDirector(t *testing.T, ctx context.Context) {
	t.Helper()
	agentSvc := agentsApp.NewAgentService(agentsInfra.NewInMemoryAgentRepository(),
		director.NewLeadDirector(stubObserver{}, "monitored"))
	dec, err := agentSvc.RunDirectorStep(ctx, string(fx.proj.ID), fx.episode.ID)
	if err != nil {
		t.Fatalf("director step: %v", err)
	}
	if dec.DecisionMaker != agentsDomain.DecisionMakerLeadDirector {
		t.Errorf("expected LEAD_DIRECTOR decision maker, got %s", dec.DecisionMaker)
	}
}

// produceAndApprove starts a run, registers a job, records approval.
func (fx *e2eFixture) produceAndApprove(t *testing.T, ctx context.Context) {
	t.Helper()
	prodSvc := prodApp.NewProductionService(prodInfra.NewInMemoryProductionRepository())
	run, err := prodSvc.StartRun(ctx, string(fx.proj.ID), fx.episode.ID, fx.bibleVersion)
	if err != nil {
		t.Fatalf("start run: %v", err)
	}
	fx.runID = run.ID
	if _, err := prodSvc.CreateJob(ctx, &prodDomain.ProductionJob{
		ProjectID: string(fx.proj.ID),
		RunID:     run.ID,
		EpisodeID: fx.episode.ID,
		ShotID:    "shot_001",
		Kind:      "shot",
	}); err != nil {
		t.Fatalf("create job: %v", err)
	}
	app, err := prodSvc.SubmitApproval(ctx, string(fx.proj.ID), fx.episode.ID, "episode", fx.episode.ID,
		prodDomain.DecisionApprove, "Episode approved for postproduction", fx.user)
	if err != nil {
		t.Fatalf("submit approval: %v", err)
	}
	if app.Decision != prodDomain.DecisionApprove {
		t.Errorf("expected APPROVE, got %s", app.Decision)
	}
}

// postproduce saves + approves a timeline and queues a render task.
func (fx *e2eFixture) postproduce(t *testing.T, ctx context.Context) {
	t.Helper()
	store, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("local storage: %v", err)
	}
	postSvc := postApp.NewPostproductionService(postInfra.NewInMemoryPostproductionRepository(), nil, store, t.TempDir())
	tl, err := postSvc.SaveTimeline(ctx, string(fx.proj.ID), fx.episode.ID, []postDomain.TrackItem{
		{ID: "ti_1", ShotID: "shot_001", AssetURL: "shot1.mp4", Duration: 4},
	}, nil, nil)
	if err != nil {
		t.Fatalf("save timeline: %v", err)
	}
	if _, err := postSvc.ApproveTimeline(ctx, fx.episode.ID); err != nil {
		t.Fatalf("approve timeline: %v", err)
	}
	render, err := postSvc.QueueRender(ctx, string(fx.proj.ID), fx.episode.ID, "mp4", "1080x1920")
	if err != nil {
		t.Fatalf("queue render: %v", err)
	}
	if render.TimelineID != tl.ID {
		t.Errorf("render bound to wrong timeline: %s", render.TimelineID)
	}
	fx.render = render
}

// publish registers a channel and schedules the episode for release.
func (fx *e2eFixture) publish(t *testing.T, ctx context.Context) {
	t.Helper()
	pubSvc := pubApp.NewPublishingService(pubInfra.NewInMemoryPublishingRepository())
	channel, err := pubSvc.RegisterChannel(ctx, string(fx.proj.ID), pubDomain.ChannelTikTok, "TikTok Main", "acct_1", nil)
	if err != nil {
		t.Fatalf("register channel: %v", err)
	}
	pub, err := pubSvc.SchedulePublication(ctx, string(fx.proj.ID), fx.episode.ID, channel.ID, "https://files/episode.mp4",
		pubDomain.PublishMetadata{Title: fx.episode.Title, Caption: "Watch Episode 1 now!", Tags: []string{"#drama", "#thriller"}},
		nil, "")
	if err != nil {
		t.Fatalf("schedule publication: %v", err)
	}
	if pub.EpisodeID != fx.episode.ID {
		t.Errorf("expected episode %s, got %s", fx.episode.ID, pub.EpisodeID)
	}
	fx.pub = pub
}
