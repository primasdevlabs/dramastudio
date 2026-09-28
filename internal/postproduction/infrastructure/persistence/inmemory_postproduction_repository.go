package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/postproduction/domain"
)

type InMemoryPostproductionRepository struct {
	mu        sync.RWMutex
	timelines map[string][]*domain.Timeline // episode_id -> versions
	renders   map[string]*domain.RenderTask
}

func NewInMemoryPostproductionRepository() *InMemoryPostproductionRepository {
	return &InMemoryPostproductionRepository{
		timelines: make(map[string][]*domain.Timeline),
		renders:   make(map[string]*domain.RenderTask),
	}
}

func (r *InMemoryPostproductionRepository) SaveTimeline(ctx context.Context, t *domain.Timeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	vs := r.timelines[t.EpisodeID]
	for i, existing := range vs {
		if existing.Version == t.Version {
			vs[i] = t
			r.timelines[t.EpisodeID] = vs
			return nil
		}
	}
	r.timelines[t.EpisodeID] = append(vs, t)
	return nil
}

func (r *InMemoryPostproductionRepository) FindTimelineByEpisode(ctx context.Context, episodeID string) (*domain.Timeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	vs := r.timelines[episodeID]
	if len(vs) == 0 {
		return nil, domain.ErrTimelineNotFound
	}
	latest := vs[0]
	for _, v := range vs {
		if v.Version > latest.Version {
			latest = v
		}
	}
	return latest, nil
}

func (r *InMemoryPostproductionRepository) ListTimelineVersions(ctx context.Context, episodeID string) ([]*domain.Timeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*domain.Timeline{}, r.timelines[episodeID]...), nil
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
		return nil, domain.ErrRenderNotFound
	}
	return rnd, nil
}

func (r *InMemoryPostproductionRepository) ListRendersByEpisode(ctx context.Context, episodeID string) ([]*domain.RenderTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.RenderTask, 0)
	for _, rnd := range r.renders {
		if rnd.EpisodeID == episodeID {
			res = append(res, rnd)
		}
	}
	return res, nil
}
