package persistence

import (
	"context"
	"dramastudio/internal/characters/domain"
)

type PostgresCharacterRepository struct{}

func NewPostgresCharacterRepository() *PostgresCharacterRepository {
	return &PostgresCharacterRepository{}
}

func (r *PostgresCharacterRepository) FindByID(ctx context.Context, id domain.CharacterID) (*domain.Character, error) {
	return domain.NewCharacter(id, "Sample Character"), nil
}

func (r *PostgresCharacterRepository) List(ctx context.Context) ([]*domain.Character, error) {
	return []*domain.Character{}, nil
}

func (r *PostgresCharacterRepository) Save(ctx context.Context, character *domain.Character) error {
	return nil
}
