package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/story/domain"
)

type InMemoryStoryRepository struct {
	mu       sync.RWMutex
	series   map[string]*domain.Series
	seasons  map[string]*domain.Season
	episodes map[string]*domain.Episode
	scenes   map[string]*domain.Scene
}

func NewInMemoryStoryRepository() *InMemoryStoryRepository {
	return &InMemoryStoryRepository{
		series:   make(map[string]*domain.Series),
		seasons:  make(map[string]*domain.Season),
		episodes: make(map[string]*domain.Episode),
		scenes:   make(map[string]*domain.Scene),
	}
}

func (r *InMemoryStoryRepository) SaveSeries(ctx context.Context, s *domain.Series) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.series[s.ID] = s
	return nil
}

func (r *InMemoryStoryRepository) FindSeriesByID(ctx context.Context, id string) (*domain.Series, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.series[id]
	if !ok {
		return nil, fmt.Errorf("series not found: %s", id)
	}
	return s, nil
}

func (r *InMemoryStoryRepository) FindSeriesByProject(ctx context.Context, projectID string) (*domain.Series, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.series {
		if s.ID == projectID || s.ID == "series_"+projectID {
			return s, nil
		}
	}
	return nil, fmt.Errorf("series not found for project: %s", projectID)
}

func (r *InMemoryStoryRepository) SaveSeason(ctx context.Context, season *domain.Season) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seasons[season.ID] = season
	return nil
}

func (r *InMemoryStoryRepository) FindSeasonByID(ctx context.Context, id string) (*domain.Season, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.seasons[id]
	if !ok {
		return nil, fmt.Errorf("season not found: %s", id)
	}
	return s, nil
}

func (r *InMemoryStoryRepository) ListSeasonsBySeries(ctx context.Context, seriesID string) ([]*domain.Season, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Season, 0)
	for _, s := range r.seasons {
		if s.SeriesID == seriesID {
			res = append(res, s)
		}
	}
	return res, nil
}

func (r *InMemoryStoryRepository) SaveEpisode(ctx context.Context, ep *domain.Episode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.episodes[ep.ID] = ep
	return nil
}

func (r *InMemoryStoryRepository) FindEpisodeByID(ctx context.Context, id string) (*domain.Episode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ep, ok := r.episodes[id]
	if !ok {
		return nil, fmt.Errorf("episode not found: %s", id)
	}
	return ep, nil
}

func (r *InMemoryStoryRepository) ListEpisodesBySeason(ctx context.Context, seasonID string) ([]*domain.Episode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Episode, 0)
	for _, ep := range r.episodes {
		if ep.SeasonID == seasonID {
			res = append(res, ep)
		}
	}
	return res, nil
}

func (r *InMemoryStoryRepository) SaveScene(ctx context.Context, sc *domain.Scene) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scenes[sc.ID] = sc
	return nil
}

func (r *InMemoryStoryRepository) FindSceneByID(ctx context.Context, id string) (*domain.Scene, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sc, ok := r.scenes[id]
	if !ok {
		return nil, fmt.Errorf("scene not found: %s", id)
	}
	return sc, nil
}

func (r *InMemoryStoryRepository) ListScenesByEpisode(ctx context.Context, episodeID string) ([]*domain.Scene, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Scene, 0)
	for _, sc := range r.scenes {
		if sc.EpisodeID == episodeID {
			res = append(res, sc)
		}
	}
	return res, nil
}
