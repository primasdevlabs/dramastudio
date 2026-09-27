package domain

type ContinuityIssue struct {
	ID          string      `json:"id"`
	EpisodeID   string      `json:"episode_id"`
	CheckType   CheckType   `json:"check_type"`
	Violation   Violation   `json:"violation"`
	IsResolved  bool        `json:"is_resolved"`
}
