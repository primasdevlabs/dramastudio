package tests

import (
	"context"
	"testing"

	agentsvc "dramastudio/internal/agents/application/services"
	"dramastudio/internal/agents/director"
	"dramastudio/internal/agents/domain"
	agentinfra "dramastudio/internal/agents/infrastructure/persistence"
)

type paramObserver struct{ obs *director.Observation }

func (s paramObserver) Observe(_ context.Context, projectID, episodeID string) (*director.Observation, error) {
	o := *s.obs
	o.ProjectID, o.EpisodeID = projectID, episodeID
	return &o, nil
}

func TestAgentTaskLifecycle(t *testing.T) {
	repo := agentinfra.NewInMemoryAgentRepository()
	ld := director.NewLeadDirector(paramObserver{&director.Observation{}}, "monitored")
	svc := agentsvc.NewAgentService(repo, ld)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, &domain.Task{ProjectID: "proj_1", Objective: "draft scene"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if task.Status != domain.TaskPending || task.ApprovalPolicy != domain.PolicyMonitored {
		t.Fatalf("defaults: %+v", task)
	}
	// Running→pending is illegal; only approve moves pending→running.
	if _, err := svc.RetryTask(ctx, task.ID); err != domain.ErrInvalidTransition {
		t.Fatalf("pending→pending retry should fail, got %v", err)
	}
	if _, err := svc.ApproveTask(ctx, task.ID); err != nil {
		t.Fatalf("approve: %v", err)
	}
	got, _ := svc.GetTask(ctx, task.ID)
	if got.Status != domain.TaskRunning {
		t.Fatalf("expected running, got %s", got.Status)
	}
	if _, err := svc.CompleteTask(ctx, task.ID, map[string]interface{}{"out": "ok"}); err != nil {
		t.Fatalf("complete: %v", err)
	}
	got, _ = svc.GetTask(ctx, task.ID)
	if got.Status != domain.TaskSucceeded || got.CompletedAt == nil {
		t.Fatalf("completion: %+v", got)
	}
	// Terminal states are immutable.
	if _, err := svc.CancelTask(ctx, task.ID); err != domain.ErrInvalidTransition {
		t.Fatalf("terminal transition must fail, got %v", err)
	}
}

func TestLeadDirectorDecisionPriorities(t *testing.T) {
	repo := agentinfra.NewInMemoryAgentRepository()
	cases := []struct {
		name string
		obs  director.Observation
		want string
	}{
		{"budget halt", director.Observation{BudgetSpent: 100, BudgetLimit: 100}, director.DecideBudgetExceeded},
		{"blocking issues", director.Observation{BlockingIssues: 2}, director.DecideEscalate},
		{"approvals gate", director.Observation{PendingApprovals: 1}, director.DecideWaitForApproval},
		{"failed jobs", director.Observation{FailedJobs: 3, RunStatus: "running"}, director.DecideRetryFailed},
		{"advance", director.Observation{RunStatus: "running"}, director.DecideAdvanceStage},
	}
	for _, tc := range cases {
		ld := director.NewLeadDirector(paramObserver{&tc.obs}, "autonomous")
		svc := agentsvc.NewAgentService(repo, ld)
		dec, err := svc.RunDirectorStep(context.Background(), "proj_1", "ep_1")
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if dec.Decision != tc.want {
			t.Fatalf("%s: expected %s, got %s", tc.name, tc.want, dec.Decision)
		}
	}
}
