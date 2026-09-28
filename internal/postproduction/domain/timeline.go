package domain

import "time"

type TimelineStatus string

const (
	TimelineDraft    TimelineStatus = "draft"
	TimelineApproved TimelineStatus = "approved"
	TimelineRendered TimelineStatus = "rendered"
)

type TrackItem struct {
	ID        string  `json:"id"`
	ShotID    string  `json:"shot_id"`
	StartTime float64 `json:"start_time"`
	Duration  float64 `json:"duration"`
	AssetURL  string  `json:"asset_url"`
}

// Timeline is the structured edit decision list — the render source of
// truth (§42). Versions are immutable per episode.
type Timeline struct {
	ID          string         `json:"id"`
	ProjectID   string         `json:"project_id"`
	EpisodeID   string         `json:"episode_id"`
	Version     int            `json:"version"`
	Status      TimelineStatus `json:"status"`
	VideoTracks []TrackItem    `json:"video_tracks"`
	AudioTracks []TrackItem    `json:"audio_tracks"`
	Subtitles   []Subtitle     `json:"subtitles"`
	CreatedAt   time.Time      `json:"created_at"`
}
