package services

import (
	"context"
	"dramastudio/internal/canon/domain"
)

type CanonService struct {
	repo domain.CanonRepository
}

func NewCanonService(repo domain.CanonRepository) *CanonService {
	return &CanonService{repo: repo}
}

func (s *CanonService) GetFact(ctx context.Context, id string) (*domain.StoryFact, error) {
	return s.repo.FindFactByID(ctx, id)
}
