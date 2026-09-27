package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/media/domain"
)

type InMemoryMediaRepository struct {
	mu     sync.RWMutex
	assets map[string]*domain.Asset
	jobs   map[string]*domain.GenerationJob
}

func NewInMemoryMediaRepository() *InMemoryMediaRepository {
	return &InMemoryMediaRepository{
		assets: make(map[string]*domain.Asset),
		jobs:   make(map[string]*domain.GenerationJob),
	}
}

func (r *InMemoryMediaRepository) SaveAsset(ctx context.Context, asset *domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assets[asset.ID] = asset
	return nil
}

func (r *InMemoryMediaRepository) FindAssetByID(ctx context.Context, id string) (*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.assets[id]
	if !ok {
		return nil, fmt.Errorf("asset not found: %s", id)
	}
	return a, nil
}

func (r *InMemoryMediaRepository) ListAssetsByProject(ctx context.Context, projectID string) ([]*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Asset, 0)
	for _, a := range r.assets {
		if a.ProjectID == projectID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (r *InMemoryMediaRepository) SaveJob(ctx context.Context, job *domain.GenerationJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryMediaRepository) FindJobByID(ctx context.Context, id string) (*domain.GenerationJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, fmt.Errorf("generation job not found: %s", id)
	}
	return j, nil
}
