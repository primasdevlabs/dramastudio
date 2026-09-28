package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/media/domain"
)

type InMemoryMediaRepository struct {
	mu       sync.RWMutex
	assets   map[string]*domain.Asset
	versions map[string][]*domain.AssetVersion // keyed by asset_id
	jobs     map[string]*domain.GenerationJob
}

func NewInMemoryMediaRepository() *InMemoryMediaRepository {
	return &InMemoryMediaRepository{
		assets:   make(map[string]*domain.Asset),
		versions: make(map[string][]*domain.AssetVersion),
		jobs:     make(map[string]*domain.GenerationJob),
	}
}

func (r *InMemoryMediaRepository) SaveAsset(ctx context.Context, a *domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assets[a.ID] = a
	return nil
}

func (r *InMemoryMediaRepository) FindAssetByID(ctx context.Context, id string) (*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.assets[id]
	if !ok {
		return nil, domain.ErrAssetNotFound
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

func (r *InMemoryMediaRepository) SaveAssetVersion(ctx context.Context, v *domain.AssetVersion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	vs := r.versions[v.AssetID]
	for i, existing := range vs {
		if existing.Version == v.Version {
			vs[i] = v
			r.versions[v.AssetID] = vs
			return nil
		}
	}
	r.versions[v.AssetID] = append(vs, v)
	return nil
}

func (r *InMemoryMediaRepository) FindAssetVersion(ctx context.Context, assetID string, version int) (*domain.AssetVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.versions[assetID] {
		if v.Version == version {
			return v, nil
		}
	}
	return nil, domain.ErrVersionNotFound
}

func (r *InMemoryMediaRepository) ListAssetVersions(ctx context.Context, assetID string) ([]*domain.AssetVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*domain.AssetVersion{}, r.versions[assetID]...), nil
}

func (r *InMemoryMediaRepository) SaveJob(ctx context.Context, j *domain.GenerationJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.ID] = j
	return nil
}

func (r *InMemoryMediaRepository) FindJobByID(ctx context.Context, id string) (*domain.GenerationJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, domain.ErrJobNotFound
	}
	return j, nil
}

func (r *InMemoryMediaRepository) FindJobByProviderJobID(ctx context.Context, providerJobID string) (*domain.GenerationJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, j := range r.jobs {
		if j.ProviderJobID == providerJobID && providerJobID != "" {
			return j, nil
		}
	}
	return nil, domain.ErrJobNotFound
}

func (r *InMemoryMediaRepository) ListJobsByProject(ctx context.Context, projectID string) ([]*domain.GenerationJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.GenerationJob, 0)
	for _, j := range r.jobs {
		if j.ProjectID == projectID {
			res = append(res, j)
		}
	}
	return res, nil
}
