package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/world/domain"
)

type InMemoryWorldRepository struct {
	mu        sync.RWMutex
	locations map[string]*domain.Location
	props     map[string]*domain.Prop
	rules     map[string]*domain.WorldRule
}

func NewInMemoryWorldRepository() *InMemoryWorldRepository {
	return &InMemoryWorldRepository{
		locations: make(map[string]*domain.Location),
		props:     make(map[string]*domain.Prop),
		rules:     make(map[string]*domain.WorldRule),
	}
}

func (r *InMemoryWorldRepository) SaveLocation(ctx context.Context, loc *domain.Location) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.locations[loc.ID] = loc
	return nil
}

func (r *InMemoryWorldRepository) FindLocationByID(ctx context.Context, id string) (*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	loc, ok := r.locations[id]
	if !ok {
		return nil, domain.ErrLocationNotFound
	}
	return loc, nil
}

func (r *InMemoryWorldRepository) ListLocations(ctx context.Context, projectID string) ([]*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Location, 0, len(r.locations))
	for _, l := range r.locations {
		if l.ProjectID == projectID {
			res = append(res, l)
		}
	}
	return res, nil
}

func (r *InMemoryWorldRepository) SaveProp(ctx context.Context, p *domain.Prop) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.props[p.ID] = p
	return nil
}

func (r *InMemoryWorldRepository) FindPropByID(ctx context.Context, id string) (*domain.Prop, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.props[id]
	if !ok {
		return nil, domain.ErrLocationNotFound
	}
	return p, nil
}

func (r *InMemoryWorldRepository) ListProps(ctx context.Context, projectID string) ([]*domain.Prop, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Prop, 0)
	for _, p := range r.props {
		if p.ProjectID == projectID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (r *InMemoryWorldRepository) SaveWorldRule(ctx context.Context, rule *domain.WorldRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryWorldRepository) ListWorldRules(ctx context.Context, projectID string) ([]*domain.WorldRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.WorldRule, 0)
	for _, rl := range r.rules {
		if rl.ProjectID == projectID {
			res = append(res, rl)
		}
	}
	return res, nil
}
