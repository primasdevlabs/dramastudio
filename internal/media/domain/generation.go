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

// GenerationJob tracks a provider-side generation attempt. Completion
// arrives via verified webhook or polling (§61).
type GenerationJob struct {
	ID            string              `json:"id"`
	ProjectID     string              `json:"project_id"`
	AssetID       string              `json:"asset_id,omitempty"`
	Capability    string              `json:"capability"`
	Provider      string              `json:"provider"`
	Model         string              `json:"model"`
	ProviderJobID string              `json:"provider_job_id,omitempty"`
	Input         string              `json:"input,omitempty"`
	OutputURL     string              `json:"output_url,omitempty"`
	Status        GenerationJobStatus `json:"status"`
	Attempt       int                 `json:"attempt"`
	Cost          float64             `json:"cost"`
	StartedAt     time.Time           `json:"started_at"`
	CompletedAt   *time.Time          `json:"completed_at,omitempty"`
	Error         string              `json:"error,omitempty"`
}
