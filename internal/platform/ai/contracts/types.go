package contracts

import (
	"context"
	"encoding/json"
	"time"
)

// ProviderError is the normalized failure surface every adapter must
// produce (Backend.md §84). Workflows decide retry policy from Kind.
type ProviderError struct {
	Kind      ProviderErrorKind
	Provider  string
	Message   string
	Retryable bool
	RetryAfter time.Duration
}

func (e *ProviderError) Error() string {
	return e.Provider + ": " + e.Kind.String() + ": " + e.Message
}

type ProviderErrorKind string

const (
	ErrRateLimited   ProviderErrorKind = "ProviderRateLimited"
	ErrUnavailable   ProviderErrorKind = "ProviderUnavailable"
	ErrTimeout       ProviderErrorKind = "ProviderTimeout"
	ErrRejected      ProviderErrorKind = "ProviderRejected"
	ErrInvalidRequest ProviderErrorKind = "ProviderInvalidRequest"
	ErrUnknown       ProviderErrorKind = "ProviderError"
)

func (k ProviderErrorKind) String() string { return string(k) }

// GenerationSpec is the provider-neutral input for a generation (§28):
// structured fields, not raw prompt strings, are the source of truth.
type GenerationSpec struct {
	Capability string                 `json:"capability"`
	ProjectID  string                 `json:"project_id"`
	Subject    map[string]string      `json:"subject,omitempty"`
	Location   string                 `json:"location,omitempty"`
	TimeOfDay  string                 `json:"time_of_day,omitempty"`
	Emotion    string                 `json:"emotion,omitempty"`
	Action     string                 `json:"action,omitempty"`
	Camera     *CameraSpec            `json:"camera,omitempty"`
	Style      string                 `json:"style,omitempty"`
	Continuity map[string]interface{} `json:"continuity,omitempty"`
	Prompt     string                 `json:"prompt,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Duration   float64                `json:"duration_sec,omitempty"`
	AspectRatio string                `json:"aspect_ratio,omitempty"`
}

type CameraSpec struct {
	Shot     string `json:"shot,omitempty"`
	Angle    string `json:"angle,omitempty"`
	Movement string `json:"movement,omitempty"`
}

// Result carries the produced artifact metadata back to the platform.
type Result struct {
	Provider        string                 `json:"provider"`
	Model           string                 `json:"model"`
	OutputURI       string                 `json:"output_uri,omitempty"`
	OutputText      string                 `json:"output_text,omitempty"`
	OutputJSON      json.RawMessage        `json:"output_json,omitempty"`
	Cost            float64                `json:"cost"`
	ProviderJobID   string                 `json:"provider_job_id,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// AsyncJob represents a provider-side job for async generation APIs.
type AsyncJob struct {
	ProviderJobID string `json:"provider_job_id"`
	Status        string `json:"status"` // pending|running|succeeded|failed
	Result        *Result
}

// SyncGenerator completes a generation within the request lifetime.
type SyncGenerator interface {
	ProviderName() string
	Generate(ctx context.Context, model string, spec GenerationSpec) (*Result, error)
}

// AsyncGenerator submits a provider job; completion arrives via webhook or
// polling (§61).
type AsyncGenerator interface {
	ProviderName() string
	Submit(ctx context.Context, model string, spec GenerationSpec, callbackURL string) (*AsyncJob, error)
	Poll(ctx context.Context, providerJobID string) (*AsyncJob, error)
	Cancel(ctx context.Context, providerJobID string) error
}
