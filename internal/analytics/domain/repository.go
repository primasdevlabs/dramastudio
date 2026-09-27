package domain

import "context"

type AnalyticsRepository interface {
	RecordMetric(ctx context.Context, m *Metric) error
}
