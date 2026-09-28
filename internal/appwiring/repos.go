package appwiring

import (
	"context"
	"fmt"

	"dramastudio/configs"
	"dramastudio/internal/platform/database/postgres"
	"dramastudio/internal/platform/database/sqlite"

	agentsInfra "dramastudio/internal/agents/infrastructure/persistence"
	aiInfra "dramastudio/internal/ai/infrastructure/registry"
	analyticsInfra "dramastudio/internal/analytics/infrastructure/persistence"
	canonInfra "dramastudio/internal/canon/infrastructure/persistence"
	charsInfra "dramastudio/internal/characters/infrastructure/persistence"
	contInfra "dramastudio/internal/continuity/infrastructure/persistence"
	identInfra "dramastudio/internal/identity/infrastructure/persistence"
	intelInfra "dramastudio/internal/intelligence/infrastructure/persistence"
	mediaInfra "dramastudio/internal/media/infrastructure/persistence"
	postInfra "dramastudio/internal/postproduction/infrastructure/persistence"
	prodInfra "dramastudio/internal/production/infrastructure/persistence"
	projInfra "dramastudio/internal/projects/infrastructure/persistence"
	pubInfra "dramastudio/internal/publishing/infrastructure/persistence"
	storyInfra "dramastudio/internal/story/infrastructure/persistence"
	worldInfra "dramastudio/internal/world/infrastructure/persistence"

	agentsdomain "dramastudio/internal/agents/domain"
	aidomain "dramastudio/internal/ai/domain"
	analyticsdomain "dramastudio/internal/analytics/domain"
	canondomain "dramastudio/internal/canon/domain"
	charsdomain "dramastudio/internal/characters/domain"
	contdomain "dramastudio/internal/continuity/domain"
	identdomain "dramastudio/internal/identity/domain"
	inteldomain "dramastudio/internal/intelligence/domain"
	mediadomain "dramastudio/internal/media/domain"
	postdomain "dramastudio/internal/postproduction/domain"
	proddomain "dramastudio/internal/production/domain"
	projdomain "dramastudio/internal/projects/domain"
	pubdomain "dramastudio/internal/publishing/domain"
	storydomain "dramastudio/internal/story/domain"
	worlddomain "dramastudio/internal/world/domain"
)

// repos groups the persistence implementations selected by StoreBackend.
type repos struct {
	identity      identdomain.UserRepository
	projects      projdomain.ProjectRepository
	story         storydomain.StoryRepository
	canon         canondomain.CanonRepository
	characters    charsdomain.CharacterRepository
	world         worlddomain.WorldRepository
	agents        agentsdomain.AgentRepository
	production    proddomain.ProductionRepository
	continuity    contdomain.ContinuityRepository
	media         mediadomain.MediaRepository
	postprod      postdomain.PostproductionRepository
	publishing    pubdomain.PublishingRepository
	analytics     analyticsdomain.AnalyticsRepository
	modelRegistry aidomain.ModelRegistryRepository
	intelExec     inteldomain.ExecutionRepository
	intelLayers   inteldomain.PolicyLayerRepository
}

// repoBundle groups the persistence implementations selected by
// StoreBackend, plus the live database handle when a real driver is used.
type repoBundle struct {
	repos
	db    *postgres.DB // postgres store
	sqldb *sqlite.DB   // sqlite store (local dev / tests)
}

// buildRepos returns in-memory repositories by default, or SQL-backed ones
// when Store=postgres|sqlite. Both drivers share the same repository code —
// the sqlite DB satisfies postgres.Querier via dialect translation.
func buildRepos(ctx context.Context, cfg *configs.Config) (repoBundle, error) {
	r := repoBundle{}
	r.identity = identInfra.NewInMemoryUserRepository()
	r.projects = projInfra.NewInMemoryProjectRepository()
	r.story = storyInfra.NewInMemoryStoryRepository()
	r.canon = canonInfra.NewInMemoryCanonRepository()
	r.characters = charsInfra.NewInMemoryCharacterRepository()
	r.world = worldInfra.NewInMemoryWorldRepository()
	r.agents = agentsInfra.NewInMemoryAgentRepository()
	r.production = prodInfra.NewInMemoryProductionRepository()
	r.continuity = contInfra.NewInMemoryContinuityRepository()
	r.media = mediaInfra.NewInMemoryMediaRepository()
	r.postprod = postInfra.NewInMemoryPostproductionRepository()
	r.publishing = pubInfra.NewInMemoryPublishingRepository()
	r.analytics = analyticsInfra.NewInMemoryAnalyticsRepository()
	r.modelRegistry = aiInfra.NewInMemoryModelRegistry()
	intel := intelInfra.NewInMemoryRepository()
	r.intelExec = intel
	r.intelLayers = intel

	switch cfg.Store {
	case configs.StorePostgres:
		db, err := postgres.New(ctx, cfg.Database.URL, cfg.Database.MaxConns, cfg.Database.MinConns)
		if err != nil {
			return r, fmt.Errorf("postgres: %w", err)
		}
		r.db = db
		r.repos = postgresRepos(db.Pool)
	case configs.StoreSQLite:
		sqdb, err := sqlite.Open(ctx, cfg.Database.Path)
		if err != nil {
			return r, err
		}
		if err := sqlite.Migrate(ctx, sqdb, cfg.Database.MigrationsDir); err != nil {
			sqdb.Close()
			return r, err
		}
		r.sqldb = sqdb
		r.repos = postgresRepos(sqdb)
	}
	return r, nil
}

// postgresRepos swaps every in-memory repository for its SQL adapter. Despite
// the name it serves both drivers — sqlite.DB implements postgres.Querier.
func postgresRepos(q postgres.Querier) repos {
	return repos{
		identity:      identInfra.NewPostgresUserRepository(q),
		projects:      projInfra.NewPostgresProjectRepository(q),
		story:         storyInfra.NewPostgresStoryRepository(q),
		canon:         canonInfra.NewPostgresCanonRepository(q),
		characters:    charsInfra.NewPostgresCharacterRepository(q),
		world:         worldInfra.NewPostgresWorldRepository(q),
		agents:        agentsInfra.NewPostgresAgentRepository(q),
		production:    prodInfra.NewPostgresProductionRepository(q),
		continuity:    contInfra.NewPostgresContinuityRepository(q),
		media:         mediaInfra.NewPostgresMediaRepository(q),
		postprod:      postInfra.NewPostgresPostproductionRepository(q),
		publishing:    pubInfra.NewPostgresPublishingRepository(q),
		analytics:     analyticsInfra.NewPostgresAnalyticsRepository(q),
		modelRegistry: aiInfra.NewPostgresModelRegistry(q),
		intelExec:     intelInfra.NewPostgresRepository(q),
		intelLayers:   intelInfra.NewPostgresRepository(q),
	}
}
