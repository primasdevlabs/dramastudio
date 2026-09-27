package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/production/domain"
)

type InMemoryProductionRepository struct {
	mu        sync.RWMutex
	jobs      map[string]*domain.ProductionJob
	approvals map[string]*domain.ApprovalRequest
}

func NewInMemoryProductionRepository() *InMemoryProductionRepository {
	return &InMemoryProductionRepository{
		jobs:      make(map[string]*domain.ProductionJob),
		approvals: make(map[string]*domain.ApprovalRequest),
	}
}

func (r *InMemoryProductionRepository) SaveJob(ctx context.Context, job *domain.ProductionJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryProductionRepository) FindJobByID(ctx context.Context, id string) (*domain.ProductionJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, fmt.Errorf("production job not found: %s", id)
	}
	return j, nil
}

func (r *InMemoryProductionRepository) ListJobsByProject(ctx context.Context, projectID string) ([]*domain.ProductionJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ProductionJob, 0)
	for _, j := range r.jobs {
		res = append(res, j)
	}
	return res, nil
}

func (r *InMemoryProductionRepository) SaveApproval(ctx context.Context, app *domain.ApprovalRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.approvals[app.ID] = app
	return nil
}

func (r *InMemoryProductionRepository) ListApprovalsByProject(ctx context.Context, projectID string) ([]*domain.ApprovalRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ApprovalRequest, 0)
	for _, a := range r.approvals {
		if a.ProjectID == projectID {
			res = append(res, a)
		}
	}
	return res, nil
}
