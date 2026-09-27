package services

import (
	"context"
	"dramastudio/internal/analytics/domain"
)

type AnalyticsService struct {
	repo domain.AnalyticsRepository
}

func NewAnalyticsService(repo domain.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) RecordMetric(ctx context.Context, m *domain.Metric) error {
	return s.repo.RecordMetric(ctx, m)
}
