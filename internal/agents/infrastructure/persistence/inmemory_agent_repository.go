package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/agents/domain"
)

type InMemoryAgentRepository struct {
	mu          sync.RWMutex
	definitions map[string]*domain.Definition
	tasks       map[string]*domain.Task
	decisions   map[string]*domain.Decision
}

func NewInMemoryAgentRepository() *InMemoryAgentRepository {
	return &InMemoryAgentRepository{
		definitions: make(map[string]*domain.Definition),
		tasks:       make(map[string]*domain.Task),
		decisions:   make(map[string]*domain.Decision),
	}
}

func (r *InMemoryAgentRepository) SaveDefinition(ctx context.Context, d *domain.Definition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.definitions[d.ID] = d
	return nil
}

func (r *InMemoryAgentRepository) FindDefinitionByID(ctx context.Context, id string) (*domain.Definition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.definitions[id]
	if !ok {
		return nil, domain.ErrDefinitionNotFound
	}
	return d, nil
}

func (r *InMemoryAgentRepository) ListDefinitions(ctx context.Context) ([]*domain.Definition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Definition, 0, len(r.definitions))
	for _, d := range r.definitions {
		res = append(res, d)
	}
	return res, nil
}

func (r *InMemoryAgentRepository) SaveTask(ctx context.Context, t *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID] = t
	return nil
}

func (r *InMemoryAgentRepository) FindTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	return t, nil
}

func (r *InMemoryAgentRepository) ListTasks(ctx context.Context, projectID string, status domain.TaskStatus) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Task, 0)
	for _, t := range r.tasks {
		if t.ProjectID != projectID {
			continue
		}
		if status != "" && t.Status != status {
			continue
		}
		res = append(res, t)
	}
	return res, nil
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
