package services

import (
	"context"
	"dramastudio/internal/production/domain"
)

type ProductionService struct {
	repo domain.ProductionRepository
}

func NewProductionService(repo domain.ProductionRepository) *ProductionService {
	return &ProductionService{repo: repo}
}

func (s *ProductionService) GetProduction(ctx context.Context, id string) (*domain.Production, error) {
	return s.repo.FindByID(ctx, id)
}
