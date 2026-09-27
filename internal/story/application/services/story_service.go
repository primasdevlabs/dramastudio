package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/story/domain"
)

type StoryService struct {
	repo domain.StoryRepository
}

func NewStoryService(repo domain.StoryRepository) *StoryService {
	return &StoryService{repo: repo}
}

func (s *StoryService) GetSeries(ctx context.Context, id string) (*domain.Series, error) {
	return s.repo.FindSeriesByID(ctx, id)
}

func (s *StoryService) CreateSeason(ctx context.Context, seriesID, title, summary string, number int) (*domain.Season, error) {
	id := fmt.Sprintf("season_%d", time.Now().UnixNano())
	season := &domain.Season{
		ID:         id,
		SeriesID:   seriesID,
		Number:     number,
		Title:      title,
		Summary:    summary,
		EpisodeIDs: []string{},
	}
	if err := s.repo.SaveSeason(ctx, season); err != nil {
		return nil, err
	}
	return season, nil
}

func (s *StoryService) ListSeasons(ctx context.Context, seriesID string) ([]*domain.Season, error) {
	return s.repo.ListSeasonsBySeries(ctx, seriesID)
}

func (s *StoryService) CreateEpisode(ctx context.Context, seasonID, title, summary string, number int) (*domain.Episode, error) {
	id := fmt.Sprintf("ep_%d", time.Now().UnixNano())
	ep := &domain.Episode{
		ID:       id,
		SeasonID: seasonID,
		Number:   number,
		Title:    title,
		Summary:  summary,
		Status:   domain.EpisodePlanned,
		SceneIDs: []string{},
	}
	if err := s.repo.SaveEpisode(ctx, ep); err != nil {
		return nil, err
	}
	return ep, nil
}

func (s *StoryService) ListEpisodes(ctx context.Context, seasonID string) ([]*domain.Episode, error) {
	return s.repo.ListEpisodesBySeason(ctx, seasonID)
}

func (s *StoryService) GetEpisode(ctx context.Context, episodeID string) (*domain.Episode, error) {
	return s.repo.FindEpisodeByID(ctx, episodeID)
}

func (s *StoryService) CreateScene(ctx context.Context, episodeID, title, description, locationID, timeOfDay string, number int, characterIDs []string) (*domain.Scene, error) {
	id := fmt.Sprintf("scene_%d", time.Now().UnixNano())
	sc := &domain.Scene{
		ID:           id,
		EpisodeID:    episodeID,
		Number:       number,
		Title:        title,
		Description:  description,
		LocationID:   locationID,
		TimeOfDay:    timeOfDay,
		CharacterIDs: characterIDs,
		BeatIDs:      []string{},
	}
	if err := s.repo.SaveScene(ctx, sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *StoryService) ListScenes(ctx context.Context, episodeID string) ([]*domain.Scene, error) {
	return s.repo.ListScenesByEpisode(ctx, episodeID)
}
