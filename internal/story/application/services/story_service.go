package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"dramastudio/internal/platform/events"
	"dramastudio/internal/story/domain"
)

type StoryService struct {
	repo   domain.StoryRepository
	events *events.Bus // may be nil; set via SetEvents
}

func NewStoryService(repo domain.StoryRepository) *StoryService {
	return &StoryService{repo: repo}
}

// SetEvents injects the domain event bus (§48). Nil-safe emitter.
func (s *StoryService) SetEvents(b *events.Bus) {
	s.events = b
}

// projectOfSeries resolves the owning project for event attribution.
func (s *StoryService) projectOfSeries(ctx context.Context, seriesID string) string {
	ser, err := s.repo.FindSeriesByID(ctx, seriesID)
	if err != nil {
		return ""
	}
	return ser.ProjectID
}

// projectOfSeason resolves season → series → project.
func (s *StoryService) projectOfSeason(ctx context.Context, seasonID string) string {
	se, err := s.repo.FindSeasonByID(ctx, seasonID)
	if err != nil {
		return ""
	}
	return s.projectOfSeries(ctx, se.SeriesID)
}

// projectOfEpisode resolves episode → season → series → project.
func (s *StoryService) projectOfEpisode(ctx context.Context, episodeID string) string {
	ep, err := s.repo.FindEpisodeByID(ctx, episodeID)
	if err != nil {
		return ""
	}
	return s.projectOfSeason(ctx, ep.SeasonID)
}

// ProjectOfEpisode is the exported resolver for cross-context ownership
// checks (postproduction, media handlers verify episode-scoped params).
func (s *StoryService) ProjectOfEpisode(ctx context.Context, episodeID string) (string, error) {
	if p := s.projectOfEpisode(ctx, episodeID); p != "" {
		return p, nil
	}
	return "", fmt.Errorf("episode %s not found", episodeID)
}

// EnsureSeries returns the project's series, creating it if absent. A
// project has exactly one series per the domain model (§9).
func (s *StoryService) EnsureSeries(ctx context.Context, projectID, title, description string) (*domain.Series, error) {
	if existing, err := s.repo.FindSeriesByProject(ctx, projectID); err == nil {
		return existing, nil
	}
	series := &domain.Series{
		ID:          "series_" + uuid.NewString(),
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		SeasonIDs:   []string{},
	}
	if err := s.repo.SaveSeries(ctx, series); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.SeriesCreated, projectID, series.ID,
		fmt.Sprintf("Series %q created", title))
	return series, nil
}

func (s *StoryService) GetSeries(ctx context.Context, id string) (*domain.Series, error) {
	return s.repo.FindSeriesByID(ctx, id)
}

func (s *StoryService) GetSeriesByProject(ctx context.Context, projectID string) (*domain.Series, error) {
	return s.repo.FindSeriesByProject(ctx, projectID)
}

func (s *StoryService) CreateSeason(ctx context.Context, seriesID, title, summary string, number int) (*domain.Season, error) {
	if _, err := s.repo.FindSeriesByID(ctx, seriesID); err != nil {
		return nil, err
	}
	season := &domain.Season{
		ID:       "season_" + uuid.NewString(),
		SeriesID: seriesID,
		Number:   number,
		Title:    title,
		Summary:  summary,
	}
	if err := s.repo.SaveSeason(ctx, season); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.SeasonCreated, s.projectOfSeries(ctx, seriesID), season.ID,
		fmt.Sprintf("Season %d: %s created", number, title))
	return season, nil
}

func (s *StoryService) ListSeasons(ctx context.Context, seriesID string) ([]*domain.Season, error) {
	return s.repo.ListSeasonsBySeries(ctx, seriesID)
}

func (s *StoryService) GetSeason(ctx context.Context, id string) (*domain.Season, error) {
	return s.repo.FindSeasonByID(ctx, id)
}

func (s *StoryService) CreateArc(ctx context.Context, seasonID, title string, number int) (*domain.StoryArc, error) {
	if _, err := s.repo.FindSeasonByID(ctx, seasonID); err != nil {
		return nil, err
	}
	arc := &domain.StoryArc{ID: "arc_" + uuid.NewString(), SeasonID: seasonID, Title: title, Number: number}
	if err := s.repo.SaveArc(ctx, arc); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.ArcCreated, s.projectOfSeason(ctx, seasonID), arc.ID,
		fmt.Sprintf("Arc %d: %s created", number, title))
	return arc, nil
}

func (s *StoryService) ListArcs(ctx context.Context, seasonID string) ([]*domain.StoryArc, error) {
	return s.repo.ListArcsBySeason(ctx, seasonID)
}

func (s *StoryService) CreateEpisode(ctx context.Context, seasonID, arcID, title, summary string, number int) (*domain.Episode, error) {
	if _, err := s.repo.FindSeasonByID(ctx, seasonID); err != nil {
		return nil, err
	}
	if arcID != "" {
		arc, err := s.repo.FindArcByID(ctx, arcID)
		if err != nil {
			return nil, err
		}
		if arc.SeasonID != seasonID {
			return nil, domain.ErrArcNotFound
		}
	}
	ep := &domain.Episode{
		ID:       "ep_" + uuid.NewString(),
		SeasonID: seasonID,
		ArcID:    arcID,
		Number:   number,
		Title:    title,
		Summary:  summary,
		Status:   domain.EpisodePlanned,
	}
	if err := s.repo.SaveEpisode(ctx, ep); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.EpisodeCreated, s.projectOfSeason(ctx, seasonID), ep.ID,
		fmt.Sprintf("Episode %d: %s created", number, title))
	return ep, nil
}

func (s *StoryService) ListEpisodes(ctx context.Context, seasonID string) ([]*domain.Episode, error) {
	return s.repo.ListEpisodesBySeason(ctx, seasonID)
}

func (s *StoryService) GetEpisode(ctx context.Context, episodeID string) (*domain.Episode, error) {
	return s.repo.FindEpisodeByID(ctx, episodeID)
}

// SetEpisodeScript records generated/approved script content.
func (s *StoryService) SetEpisodeScript(ctx context.Context, episodeID, script string) (*domain.Episode, error) {
	ep, err := s.repo.FindEpisodeByID(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	ep.Script = script
	if ep.Status == domain.EpisodePlanned {
		ep.Status = domain.EpisodeScripted
	}
	if err := s.repo.SaveEpisode(ctx, ep); err != nil {
		return nil, err
	}
	return ep, nil
}

// UpdateEpisodeMeta patches title/summary; nil fields are left unchanged.
func (s *StoryService) UpdateEpisodeMeta(ctx context.Context, episodeID string, title, summary *string) (*domain.Episode, error) {
	ep, err := s.repo.FindEpisodeByID(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	if title != nil {
		ep.Title = *title
	}
	if summary != nil {
		ep.Summary = *summary
	}
	if err := s.repo.SaveEpisode(ctx, ep); err != nil {
		return nil, err
	}
	return ep, nil
}

func (s *StoryService) SetEpisodeStatus(ctx context.Context, episodeID string, status domain.EpisodeStatus) (*domain.Episode, error) {
	ep, err := s.repo.FindEpisodeByID(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	ep.Status = status
	if err := s.repo.SaveEpisode(ctx, ep); err != nil {
		return nil, err
	}
	if status == domain.EpisodeCompleted {
		s.events.Emit(ctx, events.EpisodeCompleted, s.projectOfEpisode(ctx, episodeID), episodeID,
			fmt.Sprintf("Episode %d completed", ep.Number))
	}
	return ep, nil
}

func (s *StoryService) CreateScene(ctx context.Context, episodeID, title, description, locationID, timeOfDay string, number int, characterIDs []string) (*domain.Scene, error) {
	if _, err := s.repo.FindEpisodeByID(ctx, episodeID); err != nil {
		return nil, err
	}
	sc := &domain.Scene{
		ID:           "scene_" + uuid.NewString(),
		EpisodeID:    episodeID,
		Number:       number,
		Title:        title,
		Description:  description,
		LocationID:   locationID,
		TimeOfDay:    timeOfDay,
		CharacterIDs: characterIDs,
	}
	if err := s.repo.SaveScene(ctx, sc); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.SceneCreated, s.projectOfEpisode(ctx, episodeID), sc.ID,
		fmt.Sprintf("Scene %d: %s created", number, title))
	return sc, nil
}

func (s *StoryService) GetScene(ctx context.Context, sceneID string) (*domain.Scene, error) {
	return s.repo.FindSceneByID(ctx, sceneID)
}

func (s *StoryService) ListScenes(ctx context.Context, episodeID string) ([]*domain.Scene, error) {
	return s.repo.ListScenesByEpisode(ctx, episodeID)
}

func (s *StoryService) CreateBeat(ctx context.Context, sceneID, action, dialogue, characterID string, seq int) (*domain.Beat, error) {
	if _, err := s.repo.FindSceneByID(ctx, sceneID); err != nil {
		return nil, err
	}
	b := &domain.Beat{
		ID:          "beat_" + uuid.NewString(),
		SceneID:     sceneID,
		Seq:         seq,
		Action:      action,
		Dialogue:    dialogue,
		CharacterID: characterID,
	}
	if err := s.repo.SaveBeat(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *StoryService) ListBeats(ctx context.Context, sceneID string) ([]*domain.Beat, error) {
	return s.repo.ListBeatsByScene(ctx, sceneID)
}

// AddGraphNode records a narrative event node (§14).
func (s *StoryService) AddGraphNode(ctx context.Context, projectID, nodeType, title, episodeID, sceneID string) (*domain.StoryNode, error) {
	n := &domain.StoryNode{
		ID:        "node_" + uuid.NewString(),
		Type:      nodeType,
		Title:     title,
		EpisodeID: episodeID,
		SceneID:   sceneID,
	}
	if err := s.repo.SaveGraphNode(ctx, projectID, n); err != nil {
		return nil, err
	}
	return n, nil
}

// AddGraphEdge links two nodes with a causal relation.
func (s *StoryService) AddGraphEdge(ctx context.Context, projectID, fromID, toID string, relation domain.EdgeType) (*domain.StoryEdge, error) {
	e := &domain.StoryEdge{FromID: fromID, ToID: toID, Relation: relation}
	if err := s.repo.SaveGraphEdge(ctx, projectID, e); err != nil {
		return nil, err
	}
	return e, nil
}

// GetGraph assembles the project's story graph.
func (s *StoryService) GetGraph(ctx context.Context, projectID string) (*domain.StoryGraph, error) {
	nodes, err := s.repo.ListGraphNodes(ctx, projectID)
	if err != nil {
		return nil, err
	}
	edges, err := s.repo.ListGraphEdges(ctx, projectID)
	if err != nil {
		return nil, err
	}
	g := domain.NewStoryGraph()
	for _, n := range nodes {
		g.Nodes[n.ID] = n
	}
	g.Edges = edges
	return g, nil
}

func (s *StoryService) CreatePlotThread(ctx context.Context, projectID, name, description string) (*domain.PlotThread, error) {
	t := &domain.PlotThread{
		ID:          "thread_" + uuid.NewString(),
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		Status:      "open",
	}
	if err := s.repo.SavePlotThread(ctx, t); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.PlotThreadAdded, projectID, t.ID,
		fmt.Sprintf("Plot thread opened: %s", name))
	return t, nil
}

func (s *StoryService) ListPlotThreads(ctx context.Context, projectID string) ([]*domain.PlotThread, error) {
	return s.repo.ListPlotThreads(ctx, projectID)
}
