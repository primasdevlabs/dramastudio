package tests

import (
	"context"
	"path/filepath"
	"testing"

	aidomain "dramastudio/internal/ai/domain"
	aiInfra "dramastudio/internal/ai/infrastructure/registry"
	contdomain "dramastudio/internal/continuity/domain"
	contInfra "dramastudio/internal/continuity/infrastructure/persistence"
	"dramastudio/internal/platform/database/sqlite"
	projdomain "dramastudio/internal/projects/domain"
	projInfra "dramastudio/internal/projects/infrastructure/persistence"
)

// openSQLite migrates a temp database and returns the Querier.
func openSQLite(t *testing.T) *sqlite.DB {
	t.Helper()
	db, err := sqlite.Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("sqlite open: %v", err)
	}
	t.Cleanup(db.Close)
	if err := sqlite.Migrate(context.Background(), db, "../migrations"); err != nil {
		t.Fatalf("sqlite migrate: %v", err)
	}
	return db
}

// TestSQLiteMigrationsApply proves the Postgres migrations translate and run
// under the SQLite driver.
func TestSQLiteMigrationsApply(t *testing.T) {
	db := openSQLite(t)
	var n int
	if err := db.SQL().QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("expected migrations to be applied")
	}
	// schema-qualified table name flattened to story_series
	if err := db.SQL().QueryRow(`SELECT COUNT(*) FROM story_series`).Scan(&n); err != nil {
		t.Fatalf("story_series table missing: %v", err)
	}
}

// TestSQLiteReposEndToEnd exercises the shared SQL repositories over the
// sqlite Querier — the same code that runs on Postgres.
func TestSQLiteReposEndToEnd(t *testing.T) {
	ctx := context.Background()
	db := openSQLite(t)
	projects := projInfra.NewPostgresProjectRepository(db)
	p := projdomain.NewProject("p1", "org-1", "Harbor", "desc", "drama", "en", projdomain.ModeMonitored)
	if err := projects.Save(ctx, p); err != nil {
		t.Fatalf("save project: %v", err)
	}
	got, err := projects.FindByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("find project: %v", err)
	}
	if got.Name != "Harbor" {
		t.Fatalf("find project name: %q", got.Name)
	}
	list, err := projects.ListAll(ctx, "org-1")
	if err != nil || len(list) != 1 {
		t.Fatalf("list projects: %v len=%d", err, len(list))
	}
}

// TestSQLiteDialectQueries covers the two translation edge cases: JSONB
// containment (@> → json_each) and repeated $N placeholder rebinding.
func TestSQLiteDialectQueries(t *testing.T) {
	ctx := context.Background()
	db := openSQLite(t)
	registry := aiInfra.NewPostgresModelRegistry(db)
	if err := registry.SaveProvider(ctx, &aidomain.Provider{ID: "p1", Name: "MockCo", Type: "mock"}); err != nil {
		t.Fatalf("save provider: %v", err)
	}
	if err := registry.SaveModel(ctx, &aidomain.Model{
		ID: "m1", ProviderID: "p1", Name: "mock-vid", Identifier: "mock-1",
		Capabilities: []aidomain.AICapability{"video_generation", "text_to_video"},
	}); err != nil {
		t.Fatalf("save model: %v", err)
	}
	// Exercises the capabilities @> to_jsonb($1) → json_each rewrite.
	models, err := registry.ListModelsByCapability(ctx, aidomain.AICapability("video_generation"))
	if err != nil || len(models) != 1 || models[0].ID != "m1" {
		t.Fatalf("models by capability: %v %v", err, models)
	}

	// Exercises repeated $N placeholders ($2 appears twice per filter clause).
	continuity := contInfra.NewPostgresContinuityRepository(db)
	if err := continuity.SaveIssue(ctx, &contdomain.ContinuityIssue{
		ID: "i1", ProjectID: "p", EpisodeID: "e1",
		Severity: contdomain.SeverityWarning, Status: contdomain.IssueOpen,
		Category: "wardrobe", ActualState: "red coat vs blue",
	}); err != nil {
		t.Fatalf("save issue: %v", err)
	}
	issues, err := continuity.ListIssues(ctx, "p", contdomain.IssueFilter{
		EpisodeID: "e1", Status: contdomain.IssueOpen,
		Severity: contdomain.SeverityWarning, Category: "wardrobe",
	})
	if err != nil || len(issues) != 1 {
		t.Fatalf("list issues filtered: %v len=%d", err, len(issues))
	}
}
