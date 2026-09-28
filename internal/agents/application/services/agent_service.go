package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/agents/director"
	"dramastudio/internal/agents/domain"
)

type AgentService struct {
	repo         domain.AgentRepository
	leadDirector *director.LeadDirector
}

func NewAgentService(repo domain.AgentRepository, ld *director.LeadDirector) *AgentService {
	return &AgentService{repo: repo, leadDirector: ld}
}

// --- Definitions ---

func (s *AgentService) RegisterDefinition(ctx context.Context, d *domain.Definition) (*domain.Definition, error) {
	if d.ID == "" {
		d.ID = "agent_" + uuid.NewString()
	}
	d.Active = true
	if err := s.repo.SaveDefinition(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *AgentService) ListDefinitions(ctx context.Context) ([]*domain.Definition, error) {
	return s.repo.ListDefinitions(ctx)
}

// --- Tasks ---

// CreateTask registers a delegated task. Monitored-policy tasks stay
// pending until approved; autonomous tasks are queued for execution.
func (s *AgentService) CreateTask(ctx context.Context, t *domain.Task) (*domain.Task, error) {
	if t.ID == "" {
		t.ID = "task_" + uuid.NewString()
	}
	if t.Status == "" {
		t.Status = domain.TaskPending
	}
	if t.ApprovalPolicy == "" {
		t.ApprovalPolicy = domain.PolicyMonitored
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	if err := s.repo.SaveTask(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *AgentService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.FindTaskByID(ctx, id)
}

func (s *AgentService) ListTasks(ctx context.Context, projectID string, status domain.TaskStatus) ([]*domain.Task, error) {
	return s.repo.ListTasks(ctx, projectID, status)
}

func (s *AgentService) transitionTask(ctx context.Context, id string, next domain.TaskStatus, apply func(*domain.Task)) (*domain.Task, error) {
	t, err := s.repo.FindTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !t.CanTransitionTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	t.Status = next
	if apply != nil {
		apply(t)
	}
	if next == domain.TaskSucceeded || next == domain.TaskFailed || next == domain.TaskCancelled {
		now := time.Now().UTC()
		t.CompletedAt = &now
	}
	if err := s.repo.SaveTask(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// ApproveTask starts a monitored task (human control: approve).
func (s *AgentService) ApproveTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.transitionTask(ctx, id, domain.TaskRunning, nil)
}

func (s *AgentService) CancelTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.transitionTask(ctx, id, domain.TaskCancelled, nil)
}

func (s *AgentService) EscalateTask(ctx context.Context, id, reason string) (*domain.Task, error) {
	return s.transitionTask(ctx, id, domain.TaskEscalated, func(t *domain.Task) {
		t.Error = reason
	})
}

// CompleteTask records a task result.
func (s *AgentService) CompleteTask(ctx context.Context, id string, result map[string]interface{}) (*domain.Task, error) {
	return s.transitionTask(ctx, id, domain.TaskSucceeded, func(t *domain.Task) {
		t.Result = result
	})
}

func (s *AgentService) FailTask(ctx context.Context, id, errMsg string) (*domain.Task, error) {
	return s.transitionTask(ctx, id, domain.TaskFailed, func(t *domain.Task) {
		t.Error = errMsg
	})
}

// RetryTask re-queues a failed or escalated task (human control: retry).
func (s *AgentService) RetryTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.transitionTask(ctx, id, domain.TaskPending, func(t *domain.Task) {
		t.Iteration++
		t.Error = ""
		t.CompletedAt = nil
	})
}

// --- Director ---

// RunDirectorStep executes one Observe→Decide loop and persists the
// decision for audit.
func (s *AgentService) RunDirectorStep(ctx context.Context, projectID, episodeID string) (*domain.Decision, error) {
	dec, err := s.leadDirector.ExecuteLoopStep(ctx, projectID, episodeID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveDecision(ctx, dec); err != nil {
		return nil, err
	}
	return dec, nil
}

// RecordDecision persists a decision made by any actor (USER, SYSTEM, ...).
func (s *AgentService) RecordDecision(ctx context.Context, d *domain.Decision) (*domain.Decision, error) {
	if d.ID == "" {
		d.ID = "dec_" + uuid.NewString()
	}
	if d.Timestamp.IsZero() {
		d.Timestamp = time.Now().UTC()
	}
	if err := s.repo.SaveDecision(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *AgentService) ListDecisions(ctx context.Context, projectID string) ([]*domain.Decision, error) {
	return s.repo.ListDecisionsByProject(ctx, projectID)
}
