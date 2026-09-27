package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/agents/domain"
)

type InMemoryAgentRepository struct {
	mu        sync.RWMutex
	decisions map[string]*domain.Decision
}

func NewInMemoryAgentRepository() *InMemoryAgentRepository {
	return &InMemoryAgentRepository{
		decisions: make(map[string]*domain.Decision),
	}
}

func (r *InMemoryAgentRepository) SaveDecision(ctx context.Context, d *domain.Decision) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.decisions[d.ID] = d
	return nil
}

func (r *InMemoryAgentRepository) ListDecisionsByProject(ctx context.Context, projectID string) ([]*domain.Decision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Decision, 0)
	for _, d := range r.decisions {
		if d.ProjectID == projectID {
			res = append(res, d)
		}
	}
	return res, nil
}
