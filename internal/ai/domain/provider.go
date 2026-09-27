package domain

type ProviderID string

type Provider struct {
	ID          ProviderID `json:"id"`
	Name        string     `json:"name"`
	APIKeyEnv   string     `json:"api_key_env"`
	BaseURL     string     `json:"base_url"`
	IsHealthy   bool       `json:"is_healthy"`
}
