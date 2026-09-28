package tests

import (
	"context"
	"testing"

	mediaApp "dramastudio/internal/media/application/services"
	mediaDomain "dramastudio/internal/media/domain"
	mediaInfra "dramastudio/internal/media/infrastructure/persistence"
)

// newMediaSvc builds a media service on in-memory repo + mock resolver.
func newMediaSvc(t *testing.T) *mediaApp.MediaService {
	t.Helper()
	return mediaApp.NewMediaService(mediaInfra.NewInMemoryMediaRepository(), newTestResolver(t))
}

func TestMediaGenerationCreatesPendingAsset(t *testing.T) {
	ctx := context.Background()
	svc := newMediaSvc(t)

	job, asset, err := svc.Generate(ctx, mediaApp.GenerateRequest{
		ProjectID:  "proj_1",
		Capability: "video_generation",
		MediaType:  mediaDomain.MediaTypeVideo,
		SceneID:    "scene_1",
		Prompt:     "wide shot of the apartment at dusk",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if asset.Status != mediaDomain.AssetPending {
		t.Errorf("asset must stay PENDING for human review, got %s", asset.Status)
	}
	if job.ProjectID != "proj_1" || job.AssetID != asset.ID {
		t.Errorf("job not bound to asset/project: %+v", job)
	}
}

// Webhook replays must not double-count: after a job reaches a terminal
// state, repeated callbacks are acknowledged but never re-applied (§61).
func TestMediaWebhookReplayIdempotent(t *testing.T) {
	ctx := context.Background()
	repo := mediaInfra.NewInMemoryMediaRepository()
	svc := mediaApp.NewMediaService(repo, newTestResolver(t))

	job, asset, err := svc.Generate(ctx, mediaApp.GenerateRequest{
		ProjectID:  "proj_1",
		Capability: "video_generation",
		MediaType:  mediaDomain.MediaTypeVideo,
		Prompt:     "close-up on the locket",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if job.Status != mediaDomain.GenerationRunning {
		t.Fatalf("expected async job RUNNING, got %s", job.Status)
	}

	j1, err := svc.HandleProviderCallback(ctx, job.ProviderJobID, "succeeded", "mock://out/1.mp4", "", 0.5)
	if err != nil || j1.Status != mediaDomain.GenerationSucceeded {
		t.Fatalf("first callback: %v status=%s", err, j1.Status)
	}
	a1, _ := svc.GetAsset(ctx, asset.ID)
	if a1.Version != 1 || a1.Cost != 0.5 {
		t.Fatalf("expected version=1 cost=0.5, got v%d cost=%f", a1.Version, a1.Cost)
	}

	// Replay the same delivery — must be a no-op.
	j2, err := svc.HandleProviderCallback(ctx, job.ProviderJobID, "succeeded", "mock://out/2.mp4", "", 0.5)
	if err != nil {
		t.Fatalf("replay callback: %v", err)
	}
	if j2.Status != mediaDomain.GenerationSucceeded {
		t.Errorf("replay changed terminal state: %s", j2.Status)
	}
	a2, _ := svc.GetAsset(ctx, asset.ID)
	if a2.Version != 1 || a2.Cost != 0.5 {
		t.Errorf("replay mutated asset: v%d cost=%f", a2.Version, a2.Cost)
	}
	vs, err := svc.ListVersions(ctx, asset.ID)
	if err != nil || len(vs) != 1 {
		t.Errorf("expected exactly 1 version after replay, got %d err=%v", len(vs), err)
	}
}

func TestAssetApprovalLiftsCurrentVersion(t *testing.T) {
	ctx := context.Background()
	svc := newMediaSvc(t)

	job, asset, err := svc.Generate(ctx, mediaApp.GenerateRequest{
		ProjectID:  "proj_1",
		Capability: "image_generation",
		MediaType:  mediaDomain.MediaTypeImage,
		Prompt:     "character reference sheet",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	// Complete the async job so a version exists to approve.
	if _, err := svc.HandleProviderCallback(ctx, job.ProviderJobID, "succeeded", "mock://out/ref.png", "", 0.1); err != nil {
		t.Fatalf("callback: %v", err)
	}
	a, err := svc.ApproveAsset(ctx, asset.ID)
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if a.Status != mediaDomain.AssetApproved {
		t.Errorf("expected APPROVED, got %s", a.Status)
	}
	if a.URL == "" {
		t.Error("approved asset should expose the approved version URL")
	}
}
