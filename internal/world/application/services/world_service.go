package services

import (
	"context"
	"dramastudio/internal/world/domain"
)

type WorldService struct {
	repo domain.WorldRepository
}

func NewWorldService(repo domain.WorldRepository) *WorldService {
	return &WorldService{repo: repo}
}

func (s *WorldService) GetLocation(ctx context.Context, id string) (*domain.Location, error) {
	return s.repo.FindLocationByID(ctx, id)
}
