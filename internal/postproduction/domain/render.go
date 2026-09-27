package domain

import "time"

type RenderTaskStatus string

const (
	RenderPending   RenderTaskStatus = "PENDING"
	RenderRendering RenderTaskStatus = "RENDERING"
	RenderCompleted RenderTaskStatus = "COMPLETED"
	RenderFailed    RenderTaskStatus = "FAILED"
)

type RenderTask struct {
	ID        string           `json:"id"`
	EpisodeID string           `json:"episode_id"`
	Format    string           `json:"format"` // mp4, webm
	Resolution string          `json:"resolution"` // 1080x1920
	Status    RenderTaskStatus `json:"status"`
	OutputURL string           `json:"output_url"`
	CreatedAt time.Time        `json:"created_at"`
}
