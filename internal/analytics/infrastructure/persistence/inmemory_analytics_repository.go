package persistence

import (
	"context"
	"sort"
	"sync"
	"time"

	"dramastudio/internal/analytics/domain"
)

type InMemoryAnalyticsRepository struct {
	mu      sync.RWMutex
	nextID  int64
	metrics []*domain.Metric
}

func NewInMemoryAnalyticsRepository() *InMemoryAnalyticsRepository {
	return &InMemoryAnalyticsRepository{}
}

func (r *InMemoryAnalyticsRepository) RecordMetric(ctx context.Context, m *domain.Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	m.ID = r.nextID
	r.metrics = append(r.metrics, m)
	return nil
}

func (r *InMemoryAnalyticsRepository) Summarize(ctx context.Context, projectID, metric string, since, until time.Time) (*domain.MetricSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s := &domain.MetricSummary{Metric: metric}
	for _, m := range r.metrics {
		if m.ProjectID == projectID && m.Name == metric &&
			!m.Timestamp.Before(since) && !m.Timestamp.After(until) {
			s.Total += m.Value
			s.Count++
		}
	}
	return s, nil
}

func (r *InMemoryAnalyticsRepository) ListMetrics(ctx context.Context, projectID string, limit int) ([]*domain.Metric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	out := make([]*domain.Metric, 0)
	for _, m := range r.metrics {
		if m.ProjectID == projectID {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp.After(out[j].Timestamp) })
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
