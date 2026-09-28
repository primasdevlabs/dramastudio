package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/platform/security"
	storysvc "dramastudio/internal/story/application/services"
	"dramastudio/internal/story/domain"
	storyinfra "dramastudio/internal/story/infrastructure/persistence"
	storyhttp "dramastudio/internal/story/interfaces/http"
)

func newStoryService(t *testing.T) *storysvc.StoryService {
	t.Helper()
	return storysvc.NewStoryService(storyinfra.NewInMemoryStoryRepository())
}

func TestStoryHierarchy(t *testing.T) {
	svc := newStoryService(t)
	ctx := context.Background()

	series, err := svc.EnsureSeries(ctx, "proj_1", "Series", "")
	if err != nil {
		t.Fatalf("series: %v", err)
	}
	// EnsureSeries is idempotent — same project returns the same series.
	again, err := svc.EnsureSeries(ctx, "proj_1", "Other", "")
	if err != nil || again.ID != series.ID {
		t.Fatalf("EnsureSeries not idempotent: %v %v", err, again)
	}
	season, err := svc.CreateSeason(ctx, series.ID, "S1", "", 1)
	if err != nil {
		t.Fatalf("season: %v", err)
	}
	ep, err := svc.CreateEpisode(ctx, season.ID, "", "Pilot", "", 1)
	if err != nil {
		t.Fatalf("episode: %v", err)
	}
	if ep.Status != domain.EpisodePlanned {
		t.Fatalf("expected PLANNED, got %s", ep.Status)
	}
	if _, err := svc.SetEpisodeScript(ctx, ep.ID, "INT. ROOM — DAY"); err != nil {
		t.Fatalf("script: %v", err)
	}
	got, _ := svc.GetEpisode(ctx, ep.ID)
	if got.Status != domain.EpisodeScripted {
		t.Fatalf("script should advance status to SCRIPTED, got %s", got.Status)
	}
	sc, err := svc.CreateScene(ctx, ep.ID, "Opening", "", "loc_1", "day", 1, nil)
	if err != nil {
		t.Fatalf("scene: %v", err)
	}
	if _, err := svc.CreateBeat(ctx, sc.ID, "She enters.", "", "char_1", 1); err != nil {
		t.Fatalf("beat: %v", err)
	}
}

func TestEpisodeRejectsForeignArc(t *testing.T) {
	svc := newStoryService(t)
	ctx := context.Background()

	series, _ := svc.EnsureSeries(ctx, "proj_1", "S", "")
	s1, _ := svc.CreateSeason(ctx, series.ID, "S1", "", 1)
	s2, _ := svc.CreateSeason(ctx, series.ID, "S2", "", 2)
	arcInS2, _ := svc.CreateArc(ctx, s2.ID, "Arc", 1)

	if _, err := svc.CreateEpisode(ctx, s1.ID, arcInS2.ID, "Ep", "", 1); err != domain.ErrArcNotFound {
		t.Fatalf("arc from another season must be rejected, got %v", err)
	}
}

// TestStoryNestedIDScoping proves /v1/projects/{A}/seasons/{B's season} is
// rejected — the ProjectGuard alone can't catch cross-project nested IDs.
func TestStoryNestedIDScoping(t *testing.T) {
	svc := newStoryService(t)
	ctx := context.Background()

	sA, _ := svc.EnsureSeries(ctx, "proj_a", "A", "")
	seA, _ := svc.CreateSeason(ctx, sA.ID, "S1", "", 1)
	sB, _ := svc.EnsureSeries(ctx, "proj_b", "B", "")
	seB, _ := svc.CreateSeason(ctx, sB.ID, "S1", "", 1)

	mux := http.NewServeMux()
	storyhttp.NewStoryHandler(svc).RegisterRoutes(mux)
	guarded := platformhttp.Chain(mux, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := security.Principal{UserID: "u", OrgID: "org", Roles: []string{"owner"}}
			next.ServeHTTP(w, r.WithContext(security.WithPrincipal(r.Context(), p)))
		})
	})

	call := func(path string) int {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		guarded.ServeHTTP(rec, req)
		return rec.Code
	}

	if code := call("/v1/projects/proj_a/seasons/" + seA.ID); code != 200 {
		t.Fatalf("own season should load, got %d", code)
	}
	if code := call("/v1/projects/proj_a/seasons/" + seB.ID); code != 404 {
		t.Fatalf("foreign season must 404, got %d", code)
	}
	if code := call("/v1/projects/proj_b/seasons/" + seB.ID + "/episodes"); code != 200 {
		t.Fatalf("own season episodes should load, got %d", code)
	}
}
