package domain

import "time"

type GenerationJobStatus string

const (
	GenerationPending   GenerationJobStatus = "PENDING"
	GenerationRunning   GenerationJobStatus = "RUNNING"
	GenerationSucceeded GenerationJobStatus = "SUCCEEDED"
	GenerationFailed    GenerationJobStatus = "FAILED"
	GenerationCancelled GenerationJobStatus = "CANCELLED"
)

type GenerationJob struct {
	ID          string              `json:"id"`
	ProjectID   string              `json:"project_id"`
	Capability  string              `json:"capability"`
	Provider    string              `json:"provider"`
	Model       string              `json:"model"`
	Input       string              `json:"input"`
	OutputURL   string              `json:"output_url"`
	Status      GenerationJobStatus `json:"status"`
	Attempt     int                 `json:"attempt"`
	Cost        float64             `json:"cost"`
	StartedAt   time.Time           `json:"started_at"`
	CompletedAt time.Time           `json:"completed_at"`
	Error       string              `json:"error,omitempty"`
}
