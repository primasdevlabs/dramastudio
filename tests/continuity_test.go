package tests

import (
	"context"
	"testing"

	contApp "dramastudio/internal/continuity/application/services"
	contDomain "dramastudio/internal/continuity/domain"
	contInfra "dramastudio/internal/continuity/infrastructure/persistence"
)

func newContinuitySvc() *contApp.ContinuityService {
	return contApp.NewContinuityService(contInfra.NewInMemoryContinuityRepository())
}

func TestWardrobeCheckSurfacesMissingItems(t *testing.T) {
	ctx := context.Background()
	svc := newContinuitySvc()

	check, issues, err := svc.RunCheck(ctx, contApp.CheckInput{
		ProjectID:         "proj_1",
		EpisodeID:         "ep_1",
		CheckType:         contDomain.CheckWardrobe,
		EntityID:          "char_1",
		CanonicalWardrobe: []string{"red coat", "locket"},
		AssignedWardrobe:  []string{"red coat"},
	})
	if err != nil {
		t.Fatalf("run check: %v", err)
	}
	if check.Status != contDomain.CheckFinished || check.IssueCount != 1 {
		t.Errorf("check %+v", check)
	}
	if len(issues) != 1 || issues[0].Category != string(contDomain.CheckWardrobe) {
		t.Fatalf("expected 1 wardrobe issue, got %+v", issues)
	}
	if issues[0].Status != contDomain.IssueOpen {
		t.Errorf("issue should open as OPEN, got %s", issues[0].Status)
	}
}

func TestTimelineCheckDetectsDoubleBooking(t *testing.T) {
	ctx := context.Background()
	svc := newContinuitySvc()

	// Same participant, same world time, two different scenes.
	for _, ev := range []*contDomain.TimelineEvent{
		{ProjectID: "proj_1", EpisodeID: "ep_1", SceneID: "scene_a", WorldTime: "day1-night", EventOrder: 1, Participants: []string{"char_1"}},
		{ProjectID: "proj_1", EpisodeID: "ep_1", SceneID: "scene_b", WorldTime: "day1-night", EventOrder: 2, Participants: []string{"char_1"}},
	} {
		if _, err := svc.AddTimelineEvent(ctx, ev); err != nil {
			t.Fatalf("add event: %v", err)
		}
	}
	events, _ := svc.ListTimelineEvents(ctx, "proj_1", "ep_1")
	in := contApp.CheckInput{ProjectID: "proj_1", EpisodeID: "ep_1", CheckType: contDomain.CheckTimeline}
	for _, e := range events {
		in.Events = append(in.Events, *e)
	}
	_, issues, err := svc.RunCheck(ctx, in)
	if err != nil {
		t.Fatalf("run check: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected double-booking violation")
	}
	if issues[0].Severity != contDomain.SeverityBlocking {
		t.Errorf("double-booking should be BLOCKING, got %s", issues[0].Severity)
	}
}

func TestIssueResolveLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := newContinuitySvc()

	issue, err := svc.ReportIssue(ctx, &contDomain.ContinuityIssue{
		ProjectID: "proj_1", EpisodeID: "ep_1", Category: "wardrobe",
		Severity: contDomain.SeverityWarning, Entity: "char_1",
		ActualState: "green coat instead of red",
	})
	if err != nil {
		t.Fatalf("report: %v", err)
	}
	resolved, err := svc.ResolveIssue(ctx, issue.ID, "regenerated with red coat", false)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if resolved.Status != contDomain.IssueResolved || resolved.ResolvedAt == nil {
		t.Errorf("issue not resolved: %+v", resolved)
	}
	wf, err := svc.ResolveIssue(ctx, issue.ID, "artistic choice", true)
	if err != nil {
		t.Fatalf("wontfix: %v", err)
	}
	if wf.Status != contDomain.IssueWontFix {
		t.Errorf("expected WONTFIX, got %s", wf.Status)
	}
}
