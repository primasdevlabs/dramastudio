package persistence

import (
	"context"
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

func (r *InMemoryProjectRepository) FindByID(_ context.Context, id domain.ProjectID) (*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.projects[id]
	if !ok {
		return nil, domain.ErrProjectNotFound
	}
	return p, nil
}

func (r *InMemoryProjectRepository) ListAll(_ context.Context, orgID string) ([]*domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Project, 0, len(r.projects))
	for _, p := range r.projects {
		if orgID == "" || p.OrgID == orgID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (r *InMemoryProjectRepository) Save(_ context.Context, project *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[project.ID] = project
	return nil
}

func (r *InMemoryProjectRepository) SaveBible(_ context.Context, bible *domain.SeriesBible) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.bibles[bible.ProjectID]
	for i, b := range list {
		if b.Version == bible.Version {
			list[i] = bible
			return nil
		}
	}
	r.bibles[bible.ProjectID] = append(list, bible)
	return nil
}

func (r *InMemoryProjectRepository) GetLatestBible(_ context.Context, projectID domain.ProjectID) (*domain.SeriesBible, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.bibles[projectID]
	if len(list) == 0 {
		return nil, domain.ErrBibleNotFound
	}
	latest := list[0]
	for _, b := range list {
		if b.Version > latest.Version {
			latest = b
		}
	}
	return latest, nil
}

func (r *InMemoryProjectRepository) GetBibleVersion(_ context.Context, projectID domain.ProjectID, version int) (*domain.SeriesBible, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, b := range r.bibles[projectID] {
		if b.Version == version {
			return b, nil
		}
	}
	return nil, domain.ErrBibleNotFound
}

func (r *InMemoryProjectRepository) ListBibleVersions(_ context.Context, projectID domain.ProjectID) ([]*domain.SeriesBible, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*domain.SeriesBible(nil), r.bibles[projectID]...), nil
}
