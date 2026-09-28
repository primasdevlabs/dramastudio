package domain

import "time"

type RenderTaskStatus string

const (
	RenderQueued    RenderTaskStatus = "queued"
	RenderRendering RenderTaskStatus = "rendering"
	RenderCompleted RenderTaskStatus = "succeeded"
	RenderFailed    RenderTaskStatus = "failed"
)

type RenderTask struct {
	ID         string           `json:"id"`
	TimelineID string           `json:"timeline_id"`
	ProjectID  string           `json:"project_id"`
	EpisodeID  string           `json:"episode_id"`
	Format     string           `json:"format"`     // mp4, webm
	Resolution string           `json:"resolution"` // 1080x1920
	Status     RenderTaskStatus `json:"status"`
	ObjectKey  string           `json:"object_key,omitempty"`
	OutputURL  string           `json:"output_url,omitempty"`
	Error      string           `json:"error,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	FinishedAt *time.Time       `json:"finished_at,omitempty"`
}
