package domain

type VoiceProfile struct {
	ProviderVoiceID string  `json:"provider_voice_id"`
	Pitch           float64 `json:"pitch"`
	Speed           float64 `json:"speed"`
	Accent          string  `json:"accent"`
}
