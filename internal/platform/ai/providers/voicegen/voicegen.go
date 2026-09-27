package voicegen

import (
	"context"
	"time"

	"dramastudio/internal/platform/ai/contracts"
	"dramastudio/internal/platform/ai/providers/httpjson"
)

// Adapter is a generic text-to-speech adapter for providers exposing a
// REST JSON → audio-URL surface (e.g. ElevenLabs-style APIs).
type Adapter struct {
	name   string
	client *httpjson.Client
}

func New(providerName, baseURL, apiKey string) *Adapter {
	name := providerName
	if name == "" {
		name = "voicegen"
	}
	return &Adapter{name: name, client: httpjson.New(name, baseURL, apiKey, httpjson.AuthBearer, 90*time.Second)}
}

func (a *Adapter) ProviderName() string { return a.name }

type ttsRequest struct {
	Model string                 `json:"model"`
	Text  string                 `json:"text"`
	Voice string                 `json:"voice,omitempty"`
	Meta  map[string]interface{} `json:"metadata,omitempty"`
}

type ttsResponse struct {
	AudioURL string  `json:"audio_url"`
	Duration float64 `json:"duration_sec"`
	Cost     float64 `json:"cost"`
}

func (a *Adapter) Generate(ctx context.Context, model string, spec contracts.GenerationSpec) (*contracts.Result, error) {
	req := ttsRequest{Model: model, Text: spec.Prompt, Meta: spec.Parameters}
	if v, ok := spec.Parameters["voice"]; ok {
		if s, ok := v.(string); ok {
			req.Voice = s
		}
	}
	var out ttsResponse
	if err := a.client.Do(ctx, "POST", "/tts", req, &out); err != nil {
		return nil, err
	}
	return &contracts.Result{
		Provider:  a.name,
		Model:     model,
		OutputURI: out.AudioURL,
		Cost:      out.Cost,
		Metadata:  map[string]interface{}{"duration_sec": out.Duration},
	}, nil
}
