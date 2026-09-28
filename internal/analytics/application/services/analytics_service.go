package services

import (
	"context"
	"time"

	"dramastudio/internal/analytics/domain"
)

type AnalyticsService struct {
	repo domain.AnalyticsRepository
}

func NewAnalyticsService(repo domain.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) RecordMetric(ctx context.Context, m *domain.Metric) (*domain.Metric, error) {
	if m.Timestamp.IsZero() {
		m.Timestamp = time.Now().UTC()
	}
	if err := s.repo.RecordMetric(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AnalyticsService) Summarize(ctx context.Context, projectID, metric string, since, until time.Time) (*domain.MetricSummary, error) {
	if until.IsZero() {
		until = time.Now().UTC()
	}
	if since.IsZero() {
		since = until.Add(-30 * 24 * time.Hour)
	}
	return s.repo.Summarize(ctx, projectID, metric, since, until)
}

func (s *AnalyticsService) ListMetrics(ctx context.Context, projectID string, limit int) ([]*domain.Metric, error) {
	return s.repo.ListMetrics(ctx, projectID, limit)
}
