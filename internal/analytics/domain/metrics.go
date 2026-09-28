package domain

import "time"

// Metric is one consumption event; analytics never mutates canon (§47).
type Metric struct {
	ID            int64                  `json:"id"`
	ProjectID     string                 `json:"project_id"`
	EpisodeID     string                 `json:"episode_id,omitempty"`
	PublicationID string                 `json:"publication_id,omitempty"`
	Name          string                 `json:"metric"` // views|watch_time|retention|engagement|shares|comments|follower_growth
	Value         float64                `json:"value"`
	Dimensions    map[string]interface{} `json:"dimensions,omitempty"`
	Timestamp     time.Time              `json:"recorded_at"`
}

// MetricSummary is an aggregate over a metric within a window.
type MetricSummary struct {
	Metric string  `json:"metric"`
	Total  float64 `json:"total"`
	Count  int64   `json:"count"`
}
