package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/continuity/domain"
)

type InMemoryContinuityRepository struct {
	mu     sync.RWMutex
	issues map[string]*domain.ContinuityIssue
}

func NewInMemoryContinuityRepository() *InMemoryContinuityRepository {
	return &InMemoryContinuityRepository{
		issues: make(map[string]*domain.ContinuityIssue),
	}
}

func (r *InMemoryContinuityRepository) SaveIssue(ctx context.Context, issue *domain.ContinuityIssue) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.issues[issue.ID] = issue
	return nil
}

func (r *InMemoryContinuityRepository) FindIssueByID(ctx context.Context, id string) (*domain.ContinuityIssue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	iss, ok := r.issues[id]
	if !ok {
		return nil, fmt.Errorf("issue not found: %s", id)
	}
	return iss, nil
}

func (r *InMemoryContinuityRepository) ListIssuesByProject(ctx context.Context, projectID string) ([]*domain.ContinuityIssue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ContinuityIssue, 0)
	for _, iss := range r.issues {
		if iss.ProjectID == projectID {
			res = append(res, iss)
		}
	}
	return res, nil
}

func (r *InMemoryContinuityRepository) ResolveIssue(ctx context.Context, id, resolution string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	iss, ok := r.issues[id]
	if !ok {
		return fmt.Errorf("issue not found: %s", id)
	}
	iss.IsResolved = true
	iss.Resolution = resolution
	return nil
}
