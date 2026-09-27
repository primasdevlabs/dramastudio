package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/world/domain"
)

type WorldService struct {
	repo domain.WorldRepository
}

func NewWorldService(repo domain.WorldRepository) *WorldService {
	return &WorldService{repo: repo}
}

func (s *WorldService) CreateLocation(ctx context.Context, projectID, name, description, locType string) (*domain.Location, error) {
	id := fmt.Sprintf("loc_%d", time.Now().UnixNano())
	loc := &domain.Location{
		ID:          id,
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		Type:        locType,
		Variants:    []domain.LocationVariant{},
	}
	if err := s.repo.SaveLocation(ctx, loc); err != nil {
		return nil, err
	}
	return loc, nil
}

func (s *WorldService) GetLocation(ctx context.Context, id string) (*domain.Location, error) {
	return s.repo.FindLocationByID(ctx, id)
}

func (s *WorldService) ListLocations(ctx context.Context, projectID string) ([]*domain.Location, error) {
	return s.repo.ListLocations(ctx, projectID)
}
