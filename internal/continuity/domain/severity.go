package domain

// TimelineEvent is a narrative-time record (§40): world-time ordering,
// independent of production time.
type TimelineEvent struct {
	ID           string   `json:"id"`
	ProjectID    string   `json:"project_id"`
	EpisodeID    string   `json:"episode_id,omitempty"`
	SceneID      string   `json:"scene_id,omitempty"`
	WorldTime    string   `json:"world_time"` // narrative timestamp label
	EventOrder   int      `json:"event_order"`
	Participants []string `json:"participants"`
	LocationID   string   `json:"location_id,omitempty"`
	DurationSec  float64  `json:"duration_sec"`
	Description  string   `json:"description"`
}
