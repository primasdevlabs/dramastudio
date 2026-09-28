package persistence

import (
	"context"
	"sync"
	"time"

	"dramastudio/internal/publishing/domain"
)

type InMemoryPublishingRepository struct {
	mu           sync.RWMutex
	channels     map[string]*domain.Channel
	publications map[string]*domain.Publication
}

func NewInMemoryPublishingRepository() *InMemoryPublishingRepository {
	return &InMemoryPublishingRepository{
		channels:     make(map[string]*domain.Channel),
		publications: make(map[string]*domain.Publication),
	}
}

func (r *InMemoryPublishingRepository) SaveChannel(ctx context.Context, c *domain.Channel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.channels[c.ID] = c
	return nil
}

func (r *InMemoryPublishingRepository) FindChannelByID(ctx context.Context, id string) (*domain.Channel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.channels[id]
	if !ok {
		return nil, domain.ErrChannelNotFound
	}
	return c, nil
}

func (r *InMemoryPublishingRepository) ListChannels(ctx context.Context, projectID string) ([]*domain.Channel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Channel, 0)
	for _, c := range r.channels {
		if c.ProjectID == projectID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *InMemoryPublishingRepository) SavePublication(ctx context.Context, pub *domain.Publication) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.publications[pub.ID] = pub
	return nil
}

func (r *InMemoryPublishingRepository) FindPublicationByID(ctx context.Context, id string) (*domain.Publication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.publications[id]
	if !ok {
		return nil, domain.ErrPublicationNotFound
	}
	return p, nil
}

func (r *InMemoryPublishingRepository) FindPublicationByIdempotencyKey(ctx context.Context, key string) (*domain.Publication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.publications {
		if p.IdempotencyKey == key && key != "" {
			return p, nil
		}
	}
	return nil, domain.ErrPublicationNotFound
}

func (r *InMemoryPublishingRepository) ListPublicationsByProject(ctx context.Context, projectID string) ([]*domain.Publication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Publication, 0)
	for _, p := range r.publications {
		if p.ProjectID == projectID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (r *InMemoryPublishingRepository) ListDueScheduled(ctx context.Context, now time.Time) ([]*domain.Publication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Publication, 0)
	for _, p := range r.publications {
		if p.Status == domain.PubScheduled && p.ScheduledAt != nil && !p.ScheduledAt.After(now) {
			res = append(res, p)
		}
	}
	return res, nil
}
