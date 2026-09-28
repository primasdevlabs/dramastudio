package persistence

import (
	"context"
	"encoding/json"
	"time"

	"dramastudio/internal/analytics/domain"
	"dramastudio/internal/platform/database/postgres"
)

type PostgresAnalyticsRepository struct {
	q postgres.Querier
}

func NewPostgresAnalyticsRepository(q postgres.Querier) *PostgresAnalyticsRepository {
	return &PostgresAnalyticsRepository{q: q}
}

func (r *PostgresAnalyticsRepository) RecordMetric(ctx context.Context, m *domain.Metric) error {
	dims, _ := json.Marshal(m.Dimensions)
	return r.q.QueryRow(ctx, `
		INSERT INTO analytics.metric_events (project_id, episode_id, publication_id, metric, value, dimensions, recorded_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		m.ProjectID, m.EpisodeID, m.PublicationID, m.Name, m.Value, dims, m.Timestamp).
		Scan(&m.ID)
}

func (r *PostgresAnalyticsRepository) Summarize(ctx context.Context, projectID, metric string, since, until time.Time) (*domain.MetricSummary, error) {
	var s domain.MetricSummary
	s.Metric = metric
	err := r.q.QueryRow(ctx, `
		SELECT COALESCE(SUM(value),0), COUNT(*) FROM analytics.metric_events
		WHERE project_id = $1 AND metric = $2 AND recorded_at BETWEEN $3 AND $4`,
		projectID, metric, since, until).Scan(&s.Total, &s.Count)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresAnalyticsRepository) ListMetrics(ctx context.Context, projectID string, limit int) ([]*domain.Metric, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	rows, err := r.q.Query(ctx, `
		SELECT id, project_id, episode_id, publication_id, metric, value, dimensions, recorded_at
		FROM analytics.metric_events WHERE project_id = $1 ORDER BY recorded_at DESC LIMIT $2`,
		projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*domain.Metric{}
	for rows.Next() {
		var m domain.Metric
		var dims []byte
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.EpisodeID, &m.PublicationID,
			&m.Name, &m.Value, &dims, &m.Timestamp); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(dims, &m.Dimensions)
		out = append(out, &m)
	}
	return out, rows.Err()
}
