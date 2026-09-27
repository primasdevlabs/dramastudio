package musicgen

import (
	"context"
	"time"

	"dramastudio/internal/platform/ai/contracts"
	"dramastudio/internal/platform/ai/providers/httpjson"
)

// Adapter is a generic music/SFX generation adapter for providers exposing
// a REST JSON → audio-URL surface.
type Adapter struct {
	name   string
	client *httpjson.Client
}

func New(providerName, baseURL, apiKey string) *Adapter {
	name := providerName
	if name == "" {
		name = "musicgen"
	}
	return &Adapter{name: name, client: httpjson.New(name, baseURL, apiKey, httpjson.AuthBearer, 120*time.Second)}
}

func (a *Adapter) ProviderName() string { return a.name }

type musicRequest struct {
	Model    string  `json:"model"`
	Prompt   string  `json:"prompt"`
	Duration float64 `json:"duration_sec,omitempty"`
}

type musicResponse struct {
	AudioURL string  `json:"audio_url"`
	Cost     float64 `json:"cost"`
}

func (a *Adapter) Generate(ctx context.Context, model string, spec contracts.GenerationSpec) (*contracts.Result, error) {
	var out musicResponse
	err := a.client.Do(ctx, "POST", "/generate", musicRequest{
		Model:    model,
		Prompt:   spec.Prompt,
		Duration: spec.Duration,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &contracts.Result{
		Provider:  a.name,
		Model:     model,
		OutputURI: out.AudioURL,
		Cost:      out.Cost,
	}, nil
}
