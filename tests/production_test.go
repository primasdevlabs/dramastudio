package tests

import (
	"context"
	"testing"

	prodsvc "dramastudio/internal/production/application/services"
	"dramastudio/internal/production/domain"
	prodinfra "dramastudio/internal/production/infrastructure/persistence"
)

func newProductionService(t *testing.T) *prodsvc.ProductionService {
	t.Helper()
	return prodsvc.NewProductionService(prodinfra.NewInMemoryProductionRepository())
}

func TestProductionRunLifecycle(t *testing.T) {
	svc := newProductionService(t)
	ctx := context.Background()

	run, err := svc.StartRun(ctx, "proj_1", "ep_1", 1)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if run.Status != domain.RunStatusPending {
		t.Fatalf("expected pending (no engine), got %s", run.Status)
	}
	// pending → paused is illegal.
	if _, err := svc.PauseRun(ctx, run.ID); err != domain.ErrInvalidTransition {
		t.Fatalf("pending→paused must fail, got %v", err)
	}
	// pending → running → paused → running → stopped is the legal path.
	run2, err := svc.ResumeRun(ctx, run.ID)
	if err != nil || run2.Status != domain.RunStatusRunning {
		t.Fatalf("resume: %v status=%v", err, run2)
	}
	if _, err := svc.PauseRun(ctx, run.ID); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if _, err := svc.ResumeRun(ctx, run.ID); err != nil {
		t.Fatalf("resume2: %v", err)
	}
	stopped, err := svc.StopRun(ctx, run.ID)
	if err != nil || stopped.Status != domain.RunStatusStopped || stopped.CompletedAt == nil {
		t.Fatalf("stop: %v %+v", err, stopped)
	}
	// Stopped is terminal.
	if _, err := svc.ResumeRun(ctx, run.ID); err != domain.ErrInvalidTransition {
		t.Fatalf("stopped→running must fail, got %v", err)
	}
}

func TestProductionJobRetryAndIdempotency(t *testing.T) {
	svc := newProductionService(t)
	ctx := context.Background()

	j, err := svc.CreateJob(ctx, &domain.ProductionJob{
		ProjectID: "proj_1", RunID: "run_1", Kind: "video", IdempotencyKey: "k1",
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}
	// Same idempotency key returns the existing job.
	dup, err := svc.CreateJob(ctx, &domain.ProductionJob{
		ProjectID: "proj_1", RunID: "run_1", Kind: "video", IdempotencyKey: "k1",
	})
	if err != nil || dup.ID != j.ID {
		t.Fatalf("idempotent re-submit returned different job: %v %+v", err, dup)
	}
	// Only failed jobs retry.
	if _, err := svc.RetryJob(ctx, j.ID); err != domain.ErrInvalidTransition {
		t.Fatalf("pending retry must fail, got %v", err)
	}
	if err := svc.FailJob(ctx, j.ID, "provider down"); err != nil {
		t.Fatalf("fail: %v", err)
	}
	retried, err := svc.RetryJob(ctx, j.ID)
	if err != nil || retried.Status != domain.JobStatusPending {
		t.Fatalf("retry: %v %+v", err, retried)
	}
}

func TestApprovalSingleDecision(t *testing.T) {
	svc := newProductionService(t)
	ctx := context.Background()

	a, err := svc.RequestApproval(ctx, "proj_1", "ep_1", "ASSEMBLY", "target_1")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	decided, err := svc.Decide(ctx, a.ID, domain.DecisionApprove, "ok", "user_1")
	if err != nil || decided.Decision != domain.DecisionApprove || decided.DecidedAt == nil {
		t.Fatalf("decide: %v %+v", err, decided)
	}
	// Second decision must be rejected — approvals are one-shot.
	if _, err := svc.Decide(ctx, a.ID, domain.DecisionReject, "changed mind", "user_1"); err != domain.ErrInvalidTransition {
		t.Fatalf("double-decide must fail, got %v", err)
	}
	pending, _ := svc.ListApprovals(ctx, "proj_1", true)
	if len(pending) != 0 {
		t.Fatalf("decided approval must leave pending list, got %d", len(pending))
	}
}
