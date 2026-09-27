package persistence

import (
	"context"
	"dramastudio/internal/projects/domain"
)

type PostgresProjectRepository struct{}

func NewPostgresProjectRepository() *PostgresProjectRepository {
	return &PostgresProjectRepository{}
}

func (r *PostgresProjectRepository) FindByID(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	return &domain.Project{ID: id}, nil
}

func (r *PostgresProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	return nil
}
