package tests

import (
	"context"
	"testing"
	"time"

	pubApp "dramastudio/internal/publishing/application/services"
	pubDomain "dramastudio/internal/publishing/domain"
	pubInfra "dramastudio/internal/publishing/infrastructure/persistence"
)

type stubAdapter struct{ platform pubDomain.ChannelType }

func (a stubAdapter) Platform() pubDomain.ChannelType { return a.platform }
func (a stubAdapter) Publish(_ context.Context, _ string, _ pubDomain.PublishMetadata, _ map[string]interface{}) (string, map[string]interface{}, error) {
	return "ext_123", map[string]interface{}{"ok": true}, nil
}

func newPubSvc() *pubApp.PublishingService {
	svc := pubApp.NewPublishingService(pubInfra.NewInMemoryPublishingRepository())
	svc.RegisterAdapter(stubAdapter{platform: pubDomain.ChannelYouTube})
	return svc
}

// A publication must not target another project's channel — that would
// push content to a foreign platform account.
func TestScheduleRejectsForeignChannel(t *testing.T) {
	ctx := context.Background()
	svc := newPubSvc()

	foreign, err := svc.RegisterChannel(ctx, "proj_B", pubDomain.ChannelYouTube, "B's channel", "acct_b", nil)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := svc.SchedulePublication(ctx, "proj_A", "ep_1", foreign.ID, "https://v/1.mp4",
		pubDomain.PublishMetadata{Title: "t"}, nil, ""); err == nil {
		t.Fatal("scheduled against foreign channel")
	}
}

func TestPublishLifecycleAndDueSweep(t *testing.T) {
	ctx := context.Background()
	svc := newPubSvc()

	ch, err := svc.RegisterChannel(ctx, "proj_1", pubDomain.ChannelYouTube, "YT", "acct_1", nil)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	// Schedule in the past so the sweep picks it up.
	past := time.Now().Add(-time.Hour)
	pub, err := svc.SchedulePublication(ctx, "proj_1", "ep_1", ch.ID, "https://v/1.mp4",
		pubDomain.PublishMetadata{Title: "t"}, &past, "")
	if err != nil || pub.Status != pubDomain.PubScheduled {
		t.Fatalf("schedule: %v", err)
	}
	published, failed, err := svc.PublishDue(ctx, time.Now())
	if err != nil || published != 1 || failed != 0 {
		t.Fatalf("sweep: published=%d failed=%d err=%v", published, failed, err)
	}
	got, _ := svc.GetPublication(ctx, pub.ID)
	if got.Status != pubDomain.PubPublished || got.ExternalID != "ext_123" {
		t.Errorf("publication not published: %+v", got)
	}

	// Terminal publications cannot republish.
	if _, err := svc.Publish(ctx, pub.ID); err != pubDomain.ErrInvalidTransition {
		t.Errorf("republish should fail with ErrInvalidTransition, got %v", err)
	}
}

func TestSchedulePublicationIdempotent(t *testing.T) {
	ctx := context.Background()
	svc := newPubSvc()

	ch, _ := svc.RegisterChannel(ctx, "proj_1", pubDomain.ChannelYouTube, "YT", "a", nil)
	p1, err := svc.SchedulePublication(ctx, "proj_1", "ep_1", ch.ID, "https://v/1.mp4",
		pubDomain.PublishMetadata{Title: "t"}, nil, "idem-1")
	if err != nil {
		t.Fatalf("schedule: %v", err)
	}
	p2, err := svc.SchedulePublication(ctx, "proj_1", "ep_1", ch.ID, "https://v/1.mp4",
		pubDomain.PublishMetadata{Title: "t"}, nil, "idem-1")
	if err != nil || p2.ID != p1.ID {
		t.Fatalf("idempotent retry must return same publication, got %s vs %s", p2.ID, p1.ID)
	}
}
