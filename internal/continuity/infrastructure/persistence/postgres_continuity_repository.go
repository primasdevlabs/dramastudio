package persistence

import (
	"context"
	"dramastudio/internal/continuity/domain"
)

type PostgresContinuityRepository struct{}

func NewPostgresContinuityRepository() *PostgresContinuityRepository {
	return &PostgresContinuityRepository{}
}

func (r *PostgresContinuityRepository) FindIssueByID(ctx context.Context, id string) (*domain.ContinuityIssue, error) {
	return &domain.ContinuityIssue{ID: id}, nil
}

func (r *PostgresContinuityRepository) ListIssuesByEpisode(ctx context.Context, episodeID string) ([]*domain.ContinuityIssue, error) {
	return []*domain.ContinuityIssue{}, nil
}

func (r *PostgresContinuityRepository) SaveIssue(ctx context.Context, issue *domain.ContinuityIssue) error {
	return nil
}
