package services

import (
	"context"

	"github.com/google/uuid"

	"dramastudio/internal/world/domain"
)

type WorldService struct {
	repo domain.WorldRepository
}

func NewWorldService(repo domain.WorldRepository) *WorldService {
	return &WorldService{repo: repo}
}

func (s *WorldService) CreateLocation(ctx context.Context, projectID, name, kind, description, parentID string) (*domain.Location, error) {
	loc := &domain.Location{
		ID:          "loc_" + uuid.NewString(),
		ProjectID:   projectID,
		Name:        name,
		Kind:        kind,
		Description: description,
		ParentID:    parentID,
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

// AddVariant registers a named location state (§19).
func (s *WorldService) AddVariant(ctx context.Context, locationID, name string, attrs map[string]string) (*domain.Location, error) {
	loc, err := s.repo.FindLocationByID(ctx, locationID)
	if err != nil {
		return nil, err
	}
	loc.Variants = append(loc.Variants, domain.LocationVariant{
		ID:         "var_" + uuid.NewString(),
		Name:       name,
		Attributes: attrs,
	})
	if err := s.repo.SaveLocation(ctx, loc); err != nil {
		return nil, err
	}
	return loc, nil
}

func (s *WorldService) CreateProp(ctx context.Context, projectID, name, description, locationID string) (*domain.Prop, error) {
	p := &domain.Prop{
		ID:          "prop_" + uuid.NewString(),
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		LocationID:  locationID,
	}
	if err := s.repo.SaveProp(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *WorldService) ListProps(ctx context.Context, projectID string) ([]*domain.Prop, error) {
	return s.repo.ListProps(ctx, projectID)
}

func (s *WorldService) CreateWorldRule(ctx context.Context, projectID, text, category string) (*domain.WorldRule, error) {
	rule := &domain.WorldRule{
		ID:        "wrule_" + uuid.NewString(),
		ProjectID: projectID,
		Text:      text,
		Category:  category,
	}
	if err := s.repo.SaveWorldRule(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *WorldService) ListWorldRules(ctx context.Context, projectID string) ([]*domain.WorldRule, error) {
	return s.repo.ListWorldRules(ctx, projectID)
}
