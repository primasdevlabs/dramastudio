package domain

import (
	"context"
	"time"
)

type AnalyticsRepository interface {
	RecordMetric(ctx context.Context, m *Metric) error
	Summarize(ctx context.Context, projectID, metric string, since, until time.Time) (*MetricSummary, error)
	ListMetrics(ctx context.Context, projectID string, limit int) ([]*Metric, error)
}
