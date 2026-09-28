// Package persistence provides in-memory and SQL-backed repositories for
// the intelligence layer (execution audit + operator policy layers).
package persistence

import (
	"context"
	"sort"
	"sync"

	"dramastudio/internal/intelligence/domain"
)

// InMemoryRepository implements ExecutionRepository + PolicyLayerRepository.
type InMemoryRepository struct {
	mu         sync.Mutex
	executions map[string]*domain.ExecutionRecord
	layers     []*domain.Policy
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{executions: map[string]*domain.ExecutionRecord{}}
}

func (r *InMemoryRepository) SaveExecution(ctx context.Context, rec *domain.ExecutionRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executions[rec.ID] = rec
	return nil
}

func (r *InMemoryRepository) FindExecution(ctx context.Context, id string) (*domain.ExecutionRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rec, ok := r.executions[id]; ok {
		return rec, nil
	}
	return nil, domain.ErrExecutionNotFound
}

func (r *InMemoryRepository) ListExecutions(ctx context.Context, projectID string) ([]*domain.ExecutionRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*domain.ExecutionRecord
	for _, rec := range r.executions {
		if rec.ProjectID == projectID || projectID == "" {
			out = append(out, rec)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (r *InMemoryRepository) SavePolicyLayer(ctx context.Context, p *domain.Policy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, l := range r.layers {
		if l.ID == p.ID && l.Layer == p.Layer && l.ScopeID == p.ScopeID {
			r.layers[i] = p
			return nil
		}
	}
	r.layers = append(r.layers, p)
	return nil
}

func (r *InMemoryRepository) ListPolicyLayers(ctx context.Context, policyID string) ([]*domain.Policy, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []*domain.Policy
	for _, l := range r.layers {
		if l.ID == policyID {
			out = append(out, l)
		}
	}
	return out, nil
}
