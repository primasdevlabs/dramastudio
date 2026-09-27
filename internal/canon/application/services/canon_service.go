package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/canon/domain"
)

type CanonService struct {
	repo domain.CanonRepository
}

func NewCanonService(repo domain.CanonRepository) *CanonService {
	return &CanonService{repo: repo}
}

func (s *CanonService) CreateFact(ctx context.Context, subject, predicate, object, introduced, validFrom string) (*domain.StoryFact, error) {
	id := fmt.Sprintf("fact_%d", time.Now().UnixNano())
	fact := &domain.StoryFact{
		ID:         id,
		Subject:    subject,
		Predicate:  predicate,
		Object:     object,
		Introduced: introduced,
		ValidFrom:  validFrom,
		Status:     domain.FactStatusCanonical,
	}
	if err := s.repo.SaveFact(ctx, fact); err != nil {
		return nil, err
	}
	return fact, nil
}

func (s *CanonService) ListFacts(ctx context.Context, projectID string) ([]*domain.StoryFact, error) {
	return s.repo.ListFacts(ctx, projectID)
}

func (s *CanonService) GetKnowledgeState(ctx context.Context, characterID, episodeID string) (*domain.KnowledgeState, error) {
	return s.repo.GetKnowledgeState(ctx, characterID, episodeID)
}
