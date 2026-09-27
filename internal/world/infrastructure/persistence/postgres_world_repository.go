package persistence

import (
	"context"
	"dramastudio/internal/world/domain"
)

type PostgresWorldRepository struct{}

func NewPostgresWorldRepository() *PostgresWorldRepository {
	return &PostgresWorldRepository{}
}

func (r *PostgresWorldRepository) FindLocationByID(ctx context.Context, id string) (*domain.Location, error) {
	return &domain.Location{ID: id}, nil
}

func (r *PostgresWorldRepository) SaveLocation(ctx context.Context, location *domain.Location) error {
	return nil
}
