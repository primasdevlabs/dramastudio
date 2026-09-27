package domain

type ModelID string

type Model struct {
	ID           ModelID        `json:"id"`
	Name         string         `json:"name"`
	ProviderID   string         `json:"provider_id"`
	Capabilities []AICapability `json:"capabilities"`
	Version      string         `json:"version"`
	IsActive     bool           `json:"is_active"`
	CostPerUnit  float64        `json:"cost_per_unit"`
}
