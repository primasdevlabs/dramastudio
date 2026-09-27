package services

import (
	"context"
	"dramastudio/internal/postproduction/domain"
)

type PostProductionService struct {
	repo domain.PostProductionRepository
}

func NewPostProductionService(repo domain.PostProductionRepository) *PostProductionService {
	return &PostProductionService{repo: repo}
}

func (s *PostProductionService) GetEdit(ctx context.Context, id string) (*domain.Edit, error) {
	return s.repo.FindEditByID(ctx, id)
}
