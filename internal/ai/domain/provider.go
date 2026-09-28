package domain

import "time"

type ProviderID string

// ProviderType selects the adapter used for this provider. Hosted or
// self-hosted makes no difference — the adapter contract does
// ("openai_compatible" covers OpenAI, OpenRouter, Together, vLLM,
// Ollama, gateways, ...).
const (
	ProviderOpenAICompatible = "openai_compatible"
	ProviderWan              = "wan"
	ProviderVoiceGen         = "voicegen"
	ProviderMusicGen         = "musicgen"
	ProviderMock             = "mock"
)

// Health values reported by TestConnection / health probes.
const (
	HealthUnknown  = "unknown"
	HealthHealthy  = "healthy"
	HealthDegraded = "degraded"
	HealthDown     = "down"
)

// Provider is an external AI service — a where, not a what. Secrets are
// referenced by env var name or stored server-side; APIKey never leaves
// the process (json:"-", §60).
type Provider struct {
	ID                ProviderID `json:"id"`
	Name              string     `json:"name"`
	Type              string     `json:"type"` // adapter kind
	BaseURL           string     `json:"base_url"`
	APIKeyEnv         string     `json:"api_key_env"`
	APIKey            string     `json:"-"` // write-only; never serialized
	HealthStatus      string     `json:"health_status"`
	LastHealthCheckAt *time.Time `json:"last_health_check_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

// IsHealthy reports whether the provider is currently usable.
func (p *Provider) IsHealthy() bool {
	return p.HealthStatus == "" || p.HealthStatus == HealthHealthy || p.HealthStatus == HealthDegraded
}
