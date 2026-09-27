package persistence

import (
	"context"
	"dramastudio/internal/story/domain"
)

type PostgresStoryRepository struct{}

func NewPostgresStoryRepository() *PostgresStoryRepository {
	return &PostgresStoryRepository{}
}

func (r *PostgresStoryRepository) FindSeriesByID(ctx context.Context, id string) (*domain.Series, error) {
	return &domain.Series{ID: id}, nil
}

func (r *PostgresStoryRepository) FindEpisodeByID(ctx context.Context, id string) (*domain.Episode, error) {
	return &domain.Episode{ID: id}, nil
}
