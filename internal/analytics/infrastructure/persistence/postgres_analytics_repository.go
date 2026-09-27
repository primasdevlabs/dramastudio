package persistence

import (
	"context"
	"dramastudio/internal/analytics/domain"
)

type PostgresAnalyticsRepository struct{}

func NewPostgresAnalyticsRepository() *PostgresAnalyticsRepository {
	return &PostgresAnalyticsRepository{}
}

func (r *PostgresAnalyticsRepository) RecordMetric(ctx context.Context, m *domain.Metric) error {
	return nil
}
