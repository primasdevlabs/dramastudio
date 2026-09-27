package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/publishing/domain"
)

type InMemoryPublishingRepository struct {
	mu           sync.RWMutex
	publications map[string]*domain.Publication
}

func NewInMemoryPublishingRepository() *InMemoryPublishingRepository {
	return &InMemoryPublishingRepository{
		publications: make(map[string]*domain.Publication),
	}
}

func (r *InMemoryPublishingRepository) SavePublication(ctx context.Context, pub *domain.Publication) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.publications[pub.ID] = pub
	return nil
}

func (r *InMemoryPublishingRepository) ListPublicationsByProject(ctx context.Context, projectID string) ([]*domain.Publication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Publication, 0)
	for _, p := range r.publications {
		res = append(res, p)
	}
	return res, nil
}
