package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/postproduction/domain"
)

type PostproductionService struct {
	repo domain.PostproductionRepository
}

func NewPostproductionService(repo domain.PostproductionRepository) *PostproductionService {
	return &PostproductionService{repo: repo}
}

func (s *PostproductionService) CreateTimeline(ctx context.Context, projectID, episodeID string, videoTracks, audioTracks []domain.TrackItem) (*domain.Timeline, error) {
	id := fmt.Sprintf("tl_%s", episodeID)
	tl := &domain.Timeline{
		ID:          id,
		ProjectID:   projectID,
		EpisodeID:   episodeID,
		VideoTracks: videoTracks,
		AudioTracks: audioTracks,
		Subtitles:   []domain.Subtitle{},
	}
	if err := s.repo.SaveTimeline(ctx, tl); err != nil {
		return nil, err
	}
	return tl, nil
}

func (s *PostproductionService) GetTimeline(ctx context.Context, episodeID string) (*domain.Timeline, error) {
	return s.repo.FindTimelineByEpisode(ctx, episodeID)
}

func (s *PostproductionService) CreateRender(ctx context.Context, episodeID, format, resolution string) (*domain.RenderTask, error) {
	id := fmt.Sprintf("render_%d", time.Now().UnixNano())
	task := &domain.RenderTask{
		ID:         id,
		EpisodeID:  episodeID,
		Format:     format,
		Resolution: resolution,
		Status:     domain.RenderCompleted,
		OutputURL:  fmt.Sprintf("https://storage.dramastudio.ai/renders/%s.%s", id, format),
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.SaveRender(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}
