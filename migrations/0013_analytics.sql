-- Bounded context: analytics — consumption metrics; never mutates canon.
CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.metric_events (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id  TEXT NOT NULL,
    episode_id  TEXT NOT NULL DEFAULT '',
    publication_id TEXT NOT NULL DEFAULT '',
    metric      TEXT NOT NULL,            -- views|watch_time|retention|engagement|shares|comments|follower_growth
    value       DOUBLE PRECISION NOT NULL,
    dimensions  JSONB NOT NULL DEFAULT '{}',
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_metrics_project_metric ON analytics.metric_events(project_id, metric, recorded_at);
