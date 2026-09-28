package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/platform/security"
	projsvc "dramastudio/internal/projects/application/services"
	"dramastudio/internal/projects/domain"
	projinfra "dramastudio/internal/projects/infrastructure/persistence"
)

func newProjectService(t *testing.T) *projsvc.ProjectService {
	t.Helper()
	return projsvc.NewProjectService(projinfra.NewInMemoryProjectRepository())
}

func TestProjectStateMachine(t *testing.T) {
	svc := newProjectService(t)
	ctx := context.Background()

	p, err := svc.CreateProject(ctx, "org-1", "Show", "", "drama", "en", domain.ModeMonitored)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.Status != domain.ProjectStatusCreated {
		t.Fatalf("expected PROJECT_CREATED, got %s", p.Status)
	}
	// Illegal jump must be rejected.
	if _, err := svc.TransitionStatus(ctx, p.ID, domain.ProjectStatusPublished); err == nil {
		t.Fatal("expected invalid transition error")
	}
	// Legal path: created → developing → approved.
	if _, err := svc.TransitionStatus(ctx, p.ID, domain.ProjectStatusDeveloping); err != nil {
		t.Fatalf("developing: %v", err)
	}
	if _, err := svc.TransitionStatus(ctx, p.ID, domain.ProjectStatusApproved); err != nil {
		t.Fatalf("approved: %v", err)
	}
}

func TestSeriesBibleVersioning(t *testing.T) {
	svc := newProjectService(t)
	ctx := context.Background()

	p, _ := svc.CreateProject(ctx, "org-1", "Show", "", "", "en", domain.ModeMonitored)
	b1, err := svc.SaveSeriesBible(ctx, p.ID, "premise v1", "", "", "", "", nil, nil, nil)
	if err != nil || b1.Version != 1 {
		t.Fatalf("v1: %v version=%d", err, b1.Version)
	}
	b2, err := svc.SaveSeriesBible(ctx, p.ID, "premise v2", "", "", "", "", nil, nil, nil)
	if err != nil || b2.Version != 2 {
		t.Fatalf("v2: %v version=%d", err, b2.Version)
	}
	// Older versions stay immutable and retrievable.
	v1, err := svc.GetBibleVersion(ctx, p.ID, 1)
	if err != nil || v1.Premise != "premise v1" {
		t.Fatalf("get v1: %v premise=%q", err, v1.Premise)
	}
	latest, err := svc.GetLatestBible(ctx, p.ID)
	if err != nil || latest.Version != 2 {
		t.Fatalf("latest: %v version=%d", err, latest.Version)
	}
}

func TestProjectGuardScopesByOrg(t *testing.T) {
	resolve := func(_ context.Context, projectID string) (string, error) {
		switch projectID {
		case "proj_a":
			return "org-a", nil
		case "proj_unowned":
			return "", nil
		default:
			return "", errors.New("not found")
		}
	}
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := platformhttp.ProjectGuard(resolve)(inner)

	call := func(path, orgID string) int {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req = req.WithContext(security.WithPrincipal(req.Context(), security.Principal{UserID: "u", OrgID: orgID}))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec.Code
	}

	if code := call("/v1/projects/proj_a/characters", "org-a"); code != 200 {
		t.Fatalf("owner org should pass, got %d", code)
	}
	if code := call("/v1/projects/proj_a/characters", "org-b"); code != 404 {
		t.Fatalf("foreign org should 404, got %d", code)
	}
	if code := call("/v1/projects/missing/characters", "org-a"); code != 404 {
		t.Fatalf("unknown project should 404, got %d", code)
	}
	if code := call("/v1/projects/proj_unowned/characters", "org-b"); code != 200 {
		t.Fatalf("unowned dev project should pass, got %d", code)
	}
	// Non-project and list/create routes bypass the guard.
	if code := call("/v1/projects", "org-b"); code != 200 {
		t.Fatalf("list route should bypass guard, got %d", code)
	}
	if code := call("/v1/identity/me", "org-b"); code != 200 {
		t.Fatalf("non-project route should bypass guard, got %d", code)
	}
}
