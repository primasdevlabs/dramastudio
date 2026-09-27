package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/projects/domain"
)

type InMemoryProjectRepository struct {
	mu       sync.RWMutex
	projects map[domain.ProjectID]*domain.Project
	bibles   map[domain.ProjectID][]*domain.SeriesBible
}

func NewInMemoryProjectRepository() *InMemoryProjectRepository {
	return &InMemoryProjectRepository{
		projects: make(map[domain.ProjectID]*domain.Project),
		bibles:   make(map[domain.ProjectID][]*domain.SeriesBible),
	}
}

func (r *InMemoryProjectRepository) FindByID(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.projects[id]
	if !ok {
		return nil, fmt.Errorf("project not found: %s", id)
	}
	return p, nil
}

func (r *InMemoryProjectRepository) ListAll(ctx context.Context) ([]*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Project, 0, len(r.projects))
	for _, p := range r.projects {
		res = append(res, p)
	}
	return res, nil
}

func (r *InMemoryProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[project.ID] = project
	return nil
}

func (r *InMemoryProjectRepository) SaveBible(ctx context.Context, bible *domain.SeriesBible) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bibles[bible.ProjectID] = append(r.bibles[bible.ProjectID], bible)
	return nil
}

func (r *InMemoryProjectRepository) GetLatestBible(ctx context.Context, projectID domain.ProjectID) (*domain.SeriesBible, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, ok := r.bibles[projectID]
	if !ok || len(list) == 0 {
		return nil, fmt.Errorf("series bible not found for project: %s", projectID)
	}
	return list[len(list)-1], nil
}
