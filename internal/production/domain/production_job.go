package domain

import "time"

// ProductionJob is a unit of generated work inside a run (scene render,
// shot video, voice line, music cue, QA check, publish task).
type ProductionJob struct {
	ID             string     `json:"id"`
	ProductionID   string     `json:"production_id,omitempty"`
	ProjectID      string     `json:"project_id"`
	RunID          string     `json:"run_id,omitempty"`
	EpisodeID      string     `json:"episode_id,omitempty"`
	SceneID        string     `json:"scene_id,omitempty"`
	ShotID         string     `json:"shot_id,omitempty"`
	Kind           string     `json:"kind"` // episode|scene|shot|voice|music|render|qa|publish
	Status         JobStatus  `json:"status"`
	Attempt        int        `json:"attempt"`
	ResultURL      string     `json:"result_url,omitempty"`
	Error          string     `json:"error,omitempty"`
	IdempotencyKey string     `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// Retry resets a failed/cancelled job for another attempt.
func (j *ProductionJob) Retry() {
	j.Attempt++
	j.Status = JobStatusPending
	j.Error = ""
	j.ResultURL = ""
	j.CompletedAt = nil
}
