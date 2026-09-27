package services

import (
	"context"
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
