package persistence

import (
	"context"
	"sort"
	"sync"

	"dramastudio/internal/story/domain"
)

type InMemoryStoryRepository struct {
	mu       sync.RWMutex
	series   map[string]*domain.Series
	seasons  map[string]*domain.Season
	arcs     map[string]*domain.StoryArc
	episodes map[string]*domain.Episode
	scenes   map[string]*domain.Scene
	beats    map[string]*domain.Beat
	nodes    map[string]*domain.StoryNode
	edges    map[string][]*domain.StoryEdge
	threads  map[string][]*domain.PlotThread
}

func NewInMemoryStoryRepository() *InMemoryStoryRepository {
	return &InMemoryStoryRepository{
		series:   make(map[string]*domain.Series),
		seasons:  make(map[string]*domain.Season),
		arcs:     make(map[string]*domain.StoryArc),
		episodes: make(map[string]*domain.Episode),
		scenes:   make(map[string]*domain.Scene),
		beats:    make(map[string]*domain.Beat),
		nodes:    make(map[string]*domain.StoryNode),
		edges:    make(map[string][]*domain.StoryEdge),
		threads:  make(map[string][]*domain.PlotThread),
	}
}

func (r *InMemoryStoryRepository) SaveSeries(_ context.Context, s *domain.Series) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.series[s.ID] = s
	return nil
}

func (r *InMemoryStoryRepository) FindSeriesByID(_ context.Context, id string) (*domain.Series, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.series[id]; ok {
		return s, nil
	}
	return nil, domain.ErrSeriesNotFound
}

func (r *InMemoryStoryRepository) FindSeriesByProject(_ context.Context, projectID string) (*domain.Series, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.series {
		if s.ProjectID == projectID {
			return s, nil
		}
	}
	return nil, domain.ErrSeriesNotFound
}

func (r *InMemoryStoryRepository) SaveSeason(_ context.Context, s *domain.Season) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seasons[s.ID] = s
	return nil
}

func (r *InMemoryStoryRepository) FindSeasonByID(_ context.Context, id string) (*domain.Season, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.seasons[id]; ok {
		return s, nil
	}
	return nil, domain.ErrSeasonNotFound
}

func (r *InMemoryStoryRepository) ListSeasonsBySeries(_ context.Context, seriesID string) ([]*domain.Season, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Season, 0)
	for _, s := range r.seasons {
		if s.SeriesID == seriesID {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out, nil
}

func (r *InMemoryStoryRepository) SaveArc(_ context.Context, a *domain.StoryArc) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.arcs[a.ID] = a
	return nil
}

func (r *InMemoryStoryRepository) FindArcByID(_ context.Context, id string) (*domain.StoryArc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.arcs[id]; ok {
		return a, nil
	}
	return nil, domain.ErrArcNotFound
}

func (r *InMemoryStoryRepository) ListArcsBySeason(_ context.Context, seasonID string) ([]*domain.StoryArc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.StoryArc, 0)
	for _, a := range r.arcs {
		if a.SeasonID == seasonID {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out, nil
}

func (r *InMemoryStoryRepository) SaveEpisode(_ context.Context, e *domain.Episode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.episodes[e.ID] = e
	return nil
}

func (r *InMemoryStoryRepository) FindEpisodeByID(_ context.Context, id string) (*domain.Episode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if e, ok := r.episodes[id]; ok {
		return e, nil
	}
	return nil, domain.ErrEpisodeNotFound
}

func (r *InMemoryStoryRepository) ListEpisodesBySeason(_ context.Context, seasonID string) ([]*domain.Episode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Episode, 0)
	for _, e := range r.episodes {
		if e.SeasonID == seasonID {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out, nil
}

func (r *InMemoryStoryRepository) SaveScene(_ context.Context, s *domain.Scene) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scenes[s.ID] = s
	return nil
}

func (r *InMemoryStoryRepository) FindSceneByID(_ context.Context, id string) (*domain.Scene, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.scenes[id]; ok {
		return s, nil
	}
	return nil, domain.ErrSceneNotFound
}

func (r *InMemoryStoryRepository) ListScenesByEpisode(_ context.Context, episodeID string) ([]*domain.Scene, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Scene, 0)
	for _, s := range r.scenes {
		if s.EpisodeID == episodeID {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Number < out[j].Number })
	return out, nil
}

func (r *InMemoryStoryRepository) SaveBeat(_ context.Context, b *domain.Beat) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.beats[b.ID] = b
	return nil
}

func (r *InMemoryStoryRepository) FindBeatByID(_ context.Context, id string) (*domain.Beat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if b, ok := r.beats[id]; ok {
		return b, nil
	}
	return nil, domain.ErrBeatNotFound
}

func (r *InMemoryStoryRepository) ListBeatsByScene(_ context.Context, sceneID string) ([]*domain.Beat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Beat, 0)
	for _, b := range r.beats {
		if b.SceneID == sceneID {
			out = append(out, b)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Seq < out[j].Seq })
	return out, nil
}

func (r *InMemoryStoryRepository) SaveGraphNode(_ context.Context, _ string, n *domain.StoryNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[n.ID] = n
	return nil
}

func (r *InMemoryStoryRepository) FindGraphNode(_ context.Context, id string) (*domain.StoryNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if n, ok := r.nodes[id]; ok {
		return n, nil
	}
	return nil, domain.ErrNodeNotFound
}

func (r *InMemoryStoryRepository) ListGraphNodes(_ context.Context, projectID string) ([]*domain.StoryNode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.StoryNode, 0)
	for _, n := range r.nodes {
		out = append(out, n)
	}
	_ = projectID // nodes carry project via episodes; see postgres impl
	return out, nil
}

func (r *InMemoryStoryRepository) SaveGraphEdge(_ context.Context, projectID string, e *domain.StoryEdge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.edges[projectID] = append(r.edges[projectID], e)
	return nil
}

func (r *InMemoryStoryRepository) ListGraphEdges(_ context.Context, projectID string) ([]*domain.StoryEdge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*domain.StoryEdge(nil), r.edges[projectID]...), nil
}

func (r *InMemoryStoryRepository) SavePlotThread(_ context.Context, t *domain.PlotThread) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.threads[t.ProjectID]
	for i, existing := range list {
		if existing.ID == t.ID {
			list[i] = t
			return nil
		}
	}
	r.threads[t.ProjectID] = append(list, t)
	return nil
}

func (r *InMemoryStoryRepository) ListPlotThreads(_ context.Context, projectID string) ([]*domain.PlotThread, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*domain.PlotThread(nil), r.threads[projectID]...), nil
}
