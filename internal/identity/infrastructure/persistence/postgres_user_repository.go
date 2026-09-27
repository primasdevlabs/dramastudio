package persistence

import (
	"context"
	"dramastudio/internal/identity/domain"
)

type PostgresUserRepository struct{}

func NewPostgresUserRepository() *PostgresUserRepository {
	return &PostgresUserRepository{}
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return &domain.User{ID: id}, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	return nil
}
