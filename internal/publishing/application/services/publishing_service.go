package services

import (
	"context"
	"dramastudio/internal/publishing/domain"
)

type PublishingService struct {
	repo domain.PublishingRepository
}

func NewPublishingService(repo domain.PublishingRepository) *PublishingService {
	return &PublishingService{repo: repo}
}

func (s *PublishingService) Publish(ctx context.Context, pub *domain.Publication) error {
	return s.repo.SavePublication(ctx, pub)
}
