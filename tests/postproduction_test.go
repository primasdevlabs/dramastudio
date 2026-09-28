package tests

import (
	"context"
	"strings"
	"testing"

	postApp "dramastudio/internal/postproduction/application/services"
	postDomain "dramastudio/internal/postproduction/domain"
	postInfra "dramastudio/internal/postproduction/infrastructure/persistence"
)

func newPostSvc() *postApp.PostproductionService {
	return postApp.NewPostproductionService(postInfra.NewInMemoryPostproductionRepository(), nil, nil, "")
}

// Timelines are immutable versions: every save bumps version, and a
// draft timeline cannot be rendered (human must approve first).
func TestTimelineVersioningAndDraftGate(t *testing.T) {
	ctx := context.Background()
	svc := newPostSvc()

	tracks := []postDomain.TrackItem{{AssetURL: "mem://v1.mp4"}}
	tl1, err := svc.SaveTimeline(ctx, "proj_1", "ep_1", tracks, nil, nil)
	if err != nil || tl1.Version != 1 {
		t.Fatalf("save v1: %v", err)
	}
	tl2, err := svc.SaveTimeline(ctx, "proj_1", "ep_1", tracks, nil, nil)
	if err != nil || tl2.Version != 2 {
		t.Fatalf("save v2: %v", err)
	}

	if _, err := svc.QueueRender(ctx, "proj_1", "ep_1", "mp4", "1080x1920"); err == nil ||
		!strings.Contains(err.Error(), "approve") {
		t.Fatalf("draft timeline must not render, got %v", err)
	}

	ap, err := svc.ApproveTimeline(ctx, "ep_1")
	if err != nil || ap.Status != postDomain.TimelineApproved {
		t.Fatalf("approve: %v", err)
	}
	render, err := svc.QueueRender(ctx, "proj_1", "ep_1", "mp4", "1080x1920")
	if err != nil {
		t.Fatalf("queue render: %v", err)
	}
	if render.Status != postDomain.RenderQueued || render.TimelineID != tl2.ID {
		t.Errorf("render %+v", render)
	}

	vs, err := svc.ListTimelineVersions(ctx, "ep_1")
	if err != nil || len(vs) != 2 {
		t.Errorf("expected 2 timeline versions, got %d", len(vs))
	}
}
