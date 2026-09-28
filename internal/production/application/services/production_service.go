package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/platform/workflow/contracts"
	"dramastudio/internal/production/domain"
)

type ProductionService struct {
	repo   domain.ProductionRepository
	engine contracts.Engine // may be nil; set via SetEngine
}

func NewProductionService(repo domain.ProductionRepository) *ProductionService {
	return &ProductionService{repo: repo}
}

// SetEngine injects the workflow engine port (§83). Called by the
// composition root after construction to break the wiring cycle
// (the engine's activities depend on this service).
func (s *ProductionService) SetEngine(e contracts.Engine) {
	s.engine = e
}

// signalRun delivers a control signal to the run's workflow when an
// engine is configured. Best-effort: signal failures are logged but do
// not fail the control operation itself.
func (s *ProductionService) signalRun(ctx context.Context, runID, signal string, payload interface{}) {
	if s.engine == nil {
		return
	}
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil || run.WorkflowID == "" {
		return
	}
	_ = s.engine.Signal(ctx, run.WorkflowID, signal, payload)
}

// ensureProduction creates the productions row for a project if absent.
func (s *ProductionService) ensureProduction(ctx context.Context, projectID string) (*domain.Production, error) {
	p, err := s.repo.FindProductionByProject(ctx, projectID)
	if err == nil {
		return p, nil
	}
	p = &domain.Production{ID: "prod_" + uuid.NewString(), ProjectID: projectID}
	if err := s.repo.SaveProduction(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// StartRun begins a production run for an episode (§37).
func (s *ProductionService) StartRun(ctx context.Context, projectID, episodeID string, bibleVersion int) (*domain.ProductionRun, error) {
	p, err := s.ensureProduction(ctx, projectID)
	if err != nil {
		return nil, err
	}
	run := &domain.ProductionRun{
		ID:           "run_" + uuid.NewString(),
		ProductionID: p.ID,
		ProjectID:    projectID,
		EpisodeID:    episodeID,
		Stage:        domain.StagePreProduction,
		Status:       domain.RunStatusPending,
		BibleVersion: bibleVersion,
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.repo.SaveRun(ctx, run); err != nil {
		return nil, err
	}
	// Hand off to the workflow engine when configured: the pipeline runs
	// durably there while the run row remains the production record.
	if s.engine != nil {
		wfID, err := s.engine.Start(ctx, contracts.ProduceEpisodeInput{
			ProjectID:    projectID,
			EpisodeID:    episodeID,
			RunID:        run.ID,
			BibleVersion: bibleVersion,
		})
		if err != nil {
			return nil, err
		}
		run.WorkflowID = wfID
		run.Status = domain.RunStatusRunning
		if err := s.repo.SaveRun(ctx, run); err != nil {
			return nil, err
		}
	}
	return run, nil
}

func (s *ProductionService) GetRun(ctx context.Context, runID string) (*domain.ProductionRun, error) {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		return nil, err
	}
	jobs, err := s.repo.ListJobsByRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	run.Jobs = derefJobs(jobs)
	return run, nil
}

func derefJobs(js []*domain.ProductionJob) []domain.ProductionJob {
	out := make([]domain.ProductionJob, 0, len(js))
	for _, j := range js {
		out = append(out, *j)
	}
	return out
}

func (s *ProductionService) ListRuns(ctx context.Context, projectID string) ([]*domain.ProductionRun, error) {
	return s.repo.ListRunsByProject(ctx, projectID)
}

// transitionRun applies a run state transition guarded by the state machine.
func (s *ProductionService) transitionRun(ctx context.Context, runID string, next domain.RunStatus) (*domain.ProductionRun, error) {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		return nil, err
	}
	if !run.CanTransitionTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	run.Status = next
	if next == domain.RunStatusCompleted || next == domain.RunStatusStopped {
		now := time.Now().UTC()
		run.CompletedAt = &now
	}
	if err := s.repo.SaveRun(ctx, run); err != nil {
		return nil, err
	}
	return run, nil
}

func (s *ProductionService) PauseRun(ctx context.Context, runID string) (*domain.ProductionRun, error) {
	run, err := s.transitionRun(ctx, runID, domain.RunStatusPaused)
	if err == nil {
		s.signalRun(ctx, runID, contracts.SignalPause, "")
	}
	return run, err
}

func (s *ProductionService) ResumeRun(ctx context.Context, runID string) (*domain.ProductionRun, error) {
	run, err := s.transitionRun(ctx, runID, domain.RunStatusRunning)
	if err == nil {
		s.signalRun(ctx, runID, contracts.SignalResume, "")
	}
	return run, err
}

func (s *ProductionService) StopRun(ctx context.Context, runID string) (*domain.ProductionRun, error) {
	run, err := s.transitionRun(ctx, runID, domain.RunStatusStopped)
	if err == nil {
		s.signalRun(ctx, runID, contracts.SignalStop, "")
	}
	return run, err
}

// --- Jobs ---

func (s *ProductionService) CreateJob(ctx context.Context, job *domain.ProductionJob) (*domain.ProductionJob, error) {
	if job.IdempotencyKey != "" {
		if existing, err := s.repo.FindJobByIdempotencyKey(ctx, job.IdempotencyKey); err == nil {
			return existing, nil
		}
	}
	if job.ID == "" {
		job.ID = "job_" + uuid.NewString()
	}
	if job.Status == "" {
		job.Status = domain.JobStatusPending
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now().UTC()
	}
	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *ProductionService) GetJob(ctx context.Context, id string) (*domain.ProductionJob, error) {
	return s.repo.FindJobByID(ctx, id)
}

func (s *ProductionService) ListJobs(ctx context.Context, projectID string) ([]*domain.ProductionJob, error) {
	return s.repo.ListJobsByProject(ctx, projectID)
}

// CompleteJob records a job's result URL (called by workers/webhooks).
func (s *ProductionService) CompleteJob(ctx context.Context, id, resultURL string) error {
	j, err := s.repo.FindJobByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	j.Status = domain.JobStatusCompleted
	j.ResultURL = resultURL
	j.CompletedAt = &now
	return s.repo.SaveJob(ctx, j)
}

func (s *ProductionService) FailJob(ctx context.Context, id, errMsg string) error {
	j, err := s.repo.FindJobByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	j.Status = domain.JobStatusFailed
	j.Error = errMsg
	j.CompletedAt = &now
	return s.repo.SaveJob(ctx, j)
}

// RetryJob re-queues a failed job (human control: retry).
func (s *ProductionService) RetryJob(ctx context.Context, id string) (*domain.ProductionJob, error) {
	j, err := s.repo.FindJobByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if j.Status != domain.JobStatusFailed {
		return nil, domain.ErrInvalidTransition
	}
	j.Retry()
	if err := s.repo.SaveJob(ctx, j); err != nil {
		return nil, err
	}
	return j, nil
}

// --- Shots ---

func (s *ProductionService) CreateShot(ctx context.Context, shot *domain.Shot) (*domain.Shot, error) {
	if shot.ID == "" {
		shot.ID = "shot_" + uuid.NewString()
	}
	if shot.Status == "" {
		shot.Status = domain.ShotStatusPlanned
	}
	if err := s.repo.SaveShot(ctx, shot); err != nil {
		return nil, err
	}
	return shot, nil
}

func (s *ProductionService) GetShot(ctx context.Context, id string) (*domain.Shot, error) {
	return s.repo.FindShotByID(ctx, id)
}

func (s *ProductionService) UpdateShot(ctx context.Context, shot *domain.Shot) (*domain.Shot, error) {
	if err := s.repo.SaveShot(ctx, shot); err != nil {
		return nil, err
	}
	return shot, nil
}

func (s *ProductionService) ListShots(ctx context.Context, episodeID, sceneID string) ([]*domain.Shot, error) {
	return s.repo.ListShots(ctx, episodeID, sceneID)
}

func (s *ProductionService) SetShotStatus(ctx context.Context, id string, status domain.ShotStatus) (*domain.Shot, error) {
	shot, err := s.repo.FindShotByID(ctx, id)
	if err != nil {
		return nil, err
	}
	shot.Status = status
	if err := s.repo.SaveShot(ctx, shot); err != nil {
		return nil, err
	}
	return shot, nil
}

// ApproveShot marks a shot approved with the chosen asset (human control: approve).
func (s *ProductionService) ApproveShot(ctx context.Context, id, assetURL string) (*domain.Shot, error) {
	shot, err := s.repo.FindShotByID(ctx, id)
	if err != nil {
		return nil, err
	}
	shot.Status = domain.ShotStatusApproved
	shot.ApprovedAsset = assetURL
	if err := s.repo.SaveShot(ctx, shot); err != nil {
		return nil, err
	}
	return shot, nil
}

// --- Approvals ---

// RequestApproval creates a pending approval gate for a stage/target.
func (s *ProductionService) RequestApproval(ctx context.Context, projectID, episodeID, stage, targetID string) (*domain.ApprovalRequest, error) {
	req := &domain.ApprovalRequest{
		ID:          "app_" + uuid.NewString(),
		ProjectID:   projectID,
		EpisodeID:   episodeID,
		Stage:       stage,
		TargetID:    targetID,
		SubmittedAt: time.Now().UTC(),
	}
	if err := s.repo.SaveApproval(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

// Decide records a human decision on an approval request and applies the
// corresponding run transition for control decisions (PAUSE/STOP).
func (s *ProductionService) Decide(ctx context.Context, approvalID string, decision domain.ApprovalDecision, notes, decidedBy string) (*domain.ApprovalRequest, error) {
	a, err := s.repo.FindApprovalByID(ctx, approvalID)
	if err != nil {
		return nil, err
	}
	if a.DecidedAt != nil {
		return nil, domain.ErrInvalidTransition // already decided
	}
	now := time.Now().UTC()
	a.Decision = decision
	a.Notes = notes
	a.DecidedBy = decidedBy
	a.DecidedAt = &now
	if err := s.repo.SaveApproval(ctx, a); err != nil {
		return nil, err
	}
	// Forward the decision to any workflow waiting at an approval gate.
	s.signalApproval(ctx, a, decision)
	return a, nil
}

// signalApproval finds the run for an approval's episode and forwards the
// decision to its workflow gate. Best-effort: approvals may target assets
// or shots with no running workflow behind them.
func (s *ProductionService) signalApproval(ctx context.Context, a *domain.ApprovalRequest, decision domain.ApprovalDecision) {
	if s.engine == nil {
		return
	}
	runs, err := s.repo.ListRunsByProject(ctx, a.ProjectID)
	if err != nil {
		return
	}
	var latest *domain.ProductionRun
	for _, r := range runs {
		if r.EpisodeID == a.EpisodeID && r.WorkflowID != "" &&
			(latest == nil || r.CreatedAt.After(latest.CreatedAt)) {
			latest = r
		}
	}
	if latest == nil {
		return
	}
	_ = s.engine.Signal(ctx, latest.WorkflowID, contracts.SignalApproval, string(decision))
}

// SubmitApproval is a convenience combining RequestApproval + Decide (scripted flows).
func (s *ProductionService) SubmitApproval(ctx context.Context, projectID, episodeID, stage, targetID string, decision domain.ApprovalDecision, notes, decidedBy string) (*domain.ApprovalRequest, error) {
	a, err := s.RequestApproval(ctx, projectID, episodeID, stage, targetID)
	if err != nil {
		return nil, err
	}
	return s.Decide(ctx, a.ID, decision, notes, decidedBy)
}

func (s *ProductionService) ListApprovals(ctx context.Context, projectID string, pendingOnly bool) ([]*domain.ApprovalRequest, error) {
	return s.repo.ListApprovalsByProject(ctx, projectID, pendingOnly)
}
