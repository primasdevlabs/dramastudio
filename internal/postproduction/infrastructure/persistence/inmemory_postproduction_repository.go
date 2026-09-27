package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/postproduction/domain"
)

type InMemoryPostproductionRepository struct {
	mu        sync.RWMutex
	timelines map[string]*domain.Timeline
	renders   map[string]*domain.RenderTask
}

func NewInMemoryPostproductionRepository() *InMemoryPostproductionRepository {
	return &InMemoryPostproductionRepository{
		timelines: make(map[string]*domain.Timeline),
		renders:   make(map[string]*domain.RenderTask),
	}
}

func (r *InMemoryPostproductionRepository) SaveTimeline(ctx context.Context, t *domain.Timeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.timelines[t.EpisodeID] = t
	return nil
}

func (r *InMemoryPostproductionRepository) FindTimelineByEpisode(ctx context.Context, episodeID string) (*domain.Timeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.timelines[episodeID]
	if !ok {
		return nil, fmt.Errorf("timeline not found for episode: %s", episodeID)
	}
	return t, nil
}

func (r *InMemoryPostproductionRepository) SaveRender(ctx context.Context, render *domain.RenderTask) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.renders[render.ID] = render
	return nil
}

func (r *InMemoryPostproductionRepository) FindRenderByID(ctx context.Context, id string) (*domain.RenderTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rnd, ok := r.renders[id]
	if !ok {
		return nil, fmt.Errorf("render task not found: %s", id)
	}
	return rnd, nil
}
