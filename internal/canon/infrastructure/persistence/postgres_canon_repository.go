package persistence

import (
	"context"
	"dramastudio/internal/canon/domain"
)

type PostgresCanonRepository struct{}

func NewPostgresCanonRepository() *PostgresCanonRepository {
	return &PostgresCanonRepository{}
}

func (r *PostgresCanonRepository) FindFactByID(ctx context.Context, id string) (*domain.StoryFact, error) {
	return &domain.StoryFact{ID: id}, nil
}

func (r *PostgresCanonRepository) GetKnowledgeState(ctx context.Context, characterID, episodeID string) (*domain.KnowledgeState, error) {
	return &domain.KnowledgeState{CharacterID: characterID, EpisodeID: episodeID}, nil
}

func (r *PostgresCanonRepository) SaveFact(ctx context.Context, fact *domain.StoryFact) error {
	return nil
}
