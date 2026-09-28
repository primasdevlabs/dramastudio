package registry

import (
	"context"
	"strings"
	"sync"

	"dramastudio/internal/ai/domain"
)

type InMemoryModelRegistry struct {
	mu        sync.RWMutex
	providers map[string]*domain.Provider
	models    map[string]*domain.Model
	policies  map[string]*domain.PolicyRow // key: scope|scopeID|capability
	jobs      map[string]*domain.GenerationJob
}

func NewInMemoryModelRegistry() *InMemoryModelRegistry {
	return &InMemoryModelRegistry{
		providers: make(map[string]*domain.Provider),
		models:    make(map[string]*domain.Model),
		policies:  make(map[string]*domain.PolicyRow),
		jobs:      make(map[string]*domain.GenerationJob),
	}
}

func (r *InMemoryModelRegistry) SaveProvider(ctx context.Context, p *domain.Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[string(p.ID)] = p
	return nil
}

func (r *InMemoryModelRegistry) FindProviderByID(ctx context.Context, id domain.ProviderID) (*domain.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[string(id)]
	if !ok {
		return nil, domain.ErrModelNotFound
	}
	return p, nil
}

func (r *InMemoryModelRegistry) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Provider, 0, len(r.providers))
	for _, p := range r.providers {
		out = append(out, p)
	}
	return out, nil
}

func (r *InMemoryModelRegistry) SaveModel(ctx context.Context, m *domain.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.models[string(m.ID)] = m
	return nil
}

func (r *InMemoryModelRegistry) FindModelByID(ctx context.Context, id domain.ModelID) (*domain.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.models[string(id)]
	if !ok {
		return nil, domain.ErrModelNotFound
	}
	return m, nil
}

func (r *InMemoryModelRegistry) ListModels(ctx context.Context, providerID domain.ProviderID) ([]*domain.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Model, 0)
	for _, m := range r.models {
		if providerID == "" || m.ProviderID == string(providerID) {
			out = append(out, m)
		}
	}
	return out, nil
}

func (r *InMemoryModelRegistry) ListModelsByCapability(ctx context.Context, capability domain.AICapability) ([]*domain.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Model, 0)
	for _, m := range r.models {
		if m.Supports(capability) {
			out = append(out, m)
		}
	}
	return out, nil
}

func (r *InMemoryModelRegistry) DeleteModel(ctx context.Context, id domain.ModelID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.models, string(id))
	return nil
}

func (r *InMemoryModelRegistry) DeleteProvider(ctx context.Context, id domain.ProviderID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.providers, string(id))
	for k, m := range r.models {
		if m.ProviderID == string(id) {
			delete(r.models, k)
		}
	}
	return nil
}

func policyKey(scope, scopeID string, capability domain.AICapability) string {
	return scope + "|" + scopeID + "|" + string(capability)
}

func (r *InMemoryModelRegistry) UpsertPolicy(ctx context.Context, p *domain.PolicyRow) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policies[policyKey(p.Scope, p.ScopeID, p.Capability)] = p
	return nil
}

func (r *InMemoryModelRegistry) ResolvePolicy(ctx context.Context, capability, scopeID string) (*domain.PolicyRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	scope, id := splitScopeID(scopeID)
	p, ok := r.policies[policyKey(scope, id, domain.AICapability(capability))]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func splitScopeID(scopeID string) (scope, id string) {
	if i := strings.Index(scopeID, ":"); i >= 0 {
		return scopeID[:i], scopeID[i+1:]
	}
	if scopeID == "" {
		return domain.ScopeSystem, ""
	}
	return scopeID, ""
}

func (r *InMemoryModelRegistry) ListPolicies(ctx context.Context, scope, scopeID string) ([]*domain.PolicyRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.PolicyRow, 0)
	for _, p := range r.policies {
		if p.Scope == scope && p.ScopeID == scopeID {
			out = append(out, p)
		}
	}
	return out, nil
}

func (r *InMemoryModelRegistry) SaveGenerationJob(ctx context.Context, j *domain.GenerationJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.ID] = j
	return nil
}

func (r *InMemoryModelRegistry) FindGenerationJobByID(ctx context.Context, id string) (*domain.GenerationJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, domain.ErrModelNotFound
	}
	return j, nil
}

func (r *InMemoryModelRegistry) FindGenerationJobByIdempotencyKey(ctx context.Context, key string) (*domain.GenerationJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, j := range r.jobs {
		if j.IdempotencyKey == key && key != "" {
			return j, nil
		}
	}
	return nil, domain.ErrModelNotFound
}

func (r *InMemoryModelRegistry) ListGenerationJobs(ctx context.Context, projectID string) ([]*domain.GenerationJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.GenerationJob, 0)
	for _, j := range r.jobs {
		if j.ProjectID == projectID {
			out = append(out, j)
		}
	}
	return out, nil
}
