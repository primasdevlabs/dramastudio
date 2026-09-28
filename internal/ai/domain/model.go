package domain

import "time"

type ModelID string

// ModelStatus is the lifecycle state of a registered model.
type ModelStatus string

const (
	ModelActive       ModelStatus = "active"
	ModelExperimental ModelStatus = "experimental"
	ModelDeprecated   ModelStatus = "deprecated"
	ModelDisabled     ModelStatus = "disabled"
)

// Modalities record what a model consumes and produces — the registry
// stores observed capability, never assumes it from the provider name.
type Modalities struct {
	Input  []string `json:"input"`  // text|image|audio|video
	Output []string `json:"output"` // text|image|audio|video|json
}

// Pricing is per-unit cost for routing and budget checks.
type Pricing struct {
	Currency string  `json:"currency"` // USD
	PerUnit  float64 `json:"per_unit"` // cost per Unit
	Unit     string  `json:"unit"`     // second|image|1k_tokens|request
}

// Model is a registry record — data, not code. New models are rows,
// not deployments (§24).
type Model struct {
	ID                  ModelID           `json:"id"`
	Name                string            `json:"name"` // display name ("Veo 3.1")
	ProviderID          string            `json:"provider_id"`
	Identifier          string            `json:"identifier"`   // provider-side id ("veo-3.1-generate")
	Capabilities        []AICapability    `json:"capabilities"` // task capabilities it can serve
	Modalities          Modalities        `json:"modalities"`
	SupportedParameters []string          `json:"supported_parameters,omitempty"`
	Limits              map[string]any    `json:"limits,omitempty"` // context_window, max_duration_sec, ...
	Pricing             Pricing           `json:"pricing"`
	Version             string            `json:"version"`
	Status              ModelStatus       `json:"status"`
	Metadata            map[string]string `json:"metadata,omitempty"`
	CreatedAt           time.Time         `json:"created_at"`
}

// IsActive reports whether the model can serve traffic right now.
func (m *Model) IsActive() bool {
	return m.Status == "" || m.Status == ModelActive || m.Status == ModelExperimental
}

// Supports reports whether the model is registered for a capability.
func (m *Model) Supports(cap AICapability) bool {
	for _, c := range m.Capabilities {
		if c == cap {
			return true
		}
	}
	return false
}
