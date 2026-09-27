package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/world/domain"
)

type InMemoryWorldRepository struct {
	mu        sync.RWMutex
	locations map[string]*domain.Location
}

func NewInMemoryWorldRepository() *InMemoryWorldRepository {
	return &InMemoryWorldRepository{
		locations: make(map[string]*domain.Location),
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
		return nil, fmt.Errorf("location not found: %s", id)
	}
	return loc, nil
}

func (r *InMemoryWorldRepository) ListLocations(ctx context.Context, projectID string) ([]*domain.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Location, 0, len(r.locations))
	for _, l := range r.locations {
		res = append(res, l)
	}
	return res, nil
}
