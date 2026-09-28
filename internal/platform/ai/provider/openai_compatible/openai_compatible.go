// Package openai_compatible is the Level-2 provider adapter: one client
// contract that works against OpenAI, OpenRouter, Together, Fireworks,
// vLLM, SGLang, Ollama, gateways, and self-hosted endpoints. What a given
// endpoint actually supports (tool calling, vision, structured output,
// images, audio) varies by provider AND by model — that variance is
// recorded as model capabilities in the registry, not assumed here.
package openai_compatible

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dramastudio/internal/platform/ai/capability"
	"dramastudio/internal/platform/ai/provider"
	"dramastudio/internal/platform/ai/provider/httpjson"
)

// Adapter speaks the OpenAI wire contract against a configurable endpoint.
type Adapter struct {
	name   string
	client *httpjson.Client
}

// New builds an adapter instance for one configured provider row.
func New(name, baseURL, apiKey string) *Adapter {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if name == "" {
		name = "openai_compatible"
	}
	return &Adapter{name: name, client: httpjson.New(name, baseURL, apiKey, httpjson.AuthBearer, 120*time.Second)}
}

func init() {
	factory := func(cfg provider.Config) provider.Adapter {
		return provider.Adapter{Sync: New(cfg.Name, cfg.BaseURL, cfg.APIKey)}
	}
	// Aliases keep provider rows honest about what they are while sharing
	// the same wire adapter (§24).
	for _, kind := range []string{
		"openai_compatible", "openai", "openrouter", "together",
		"fireworks", "vllm", "sglang", "ollama", "custom_gateway",
	} {
		provider.RegisterFactory(kind, factory)
	}
}

func (a *Adapter) ProviderName() string { return a.name }

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// Generate dispatches on capability: text capabilities hit
// /chat/completions; image capabilities hit /images/generations.
func (a *Adapter) Generate(ctx context.Context, model string, spec capability.GenerationSpec) (*capability.Result, error) {
	switch {
	case strings.Contains(spec.Capability, "image"):
		return a.generateImage(ctx, model, spec)
	case strings.Contains(spec.Capability, "speech"), strings.Contains(spec.Capability, "voice"):
		return a.generateSpeech(ctx, model, spec)
	default:
		return a.generateText(ctx, model, spec)
	}
}

func (a *Adapter) generateText(ctx context.Context, model string, spec capability.GenerationSpec) (*capability.Result, error) {
	var out chatResponse
	err := a.client.Do(ctx, "POST", "/chat/completions", chatRequest{
		Model:    model,
		Messages: []chatMessage{{Role: "user", Content: buildPrompt(spec)}},
	}, &out)
	if err != nil {
		return nil, err
	}
	if len(out.Choices) == 0 {
		return nil, &capability.ProviderError{Kind: capability.ErrUnknown, Provider: a.ProviderName(), Message: "empty completion"}
	}
	return &capability.Result{
		Provider:   a.ProviderName(),
		Model:      model,
		OutputText: out.Choices[0].Message.Content,
		Metadata:   map[string]interface{}{"total_tokens": out.Usage.TotalTokens},
	}, nil
}

type imageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size,omitempty"`
}

type imageResponse struct {
	Data []struct {
		URL           string `json:"url"`
		B64JSON       string `json:"b64_json"`
		RevisedPrompt string `json:"revised_prompt"`
	} `json:"data"`
}

func (a *Adapter) generateImage(ctx context.Context, model string, spec capability.GenerationSpec) (*capability.Result, error) {
	var out imageResponse
	err := a.client.Do(ctx, "POST", "/images/generations", imageRequest{
		Model:  model,
		Prompt: buildPrompt(spec),
	}, &out)
	if err != nil {
		return nil, err
	}
	if len(out.Data) == 0 {
		return nil, &capability.ProviderError{Kind: capability.ErrUnknown, Provider: a.ProviderName(), Message: "empty image response"}
	}
	return &capability.Result{
		Provider:  a.ProviderName(),
		Model:     model,
		OutputURI: out.Data[0].URL,
		Metadata:  map[string]interface{}{"revised_prompt": out.Data[0].RevisedPrompt},
	}, nil
}

type speechRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice,omitempty"`
}

// generateSpeech targets the OpenAI-compatible /audio/speech endpoint.
// Bytes arrive base64-wrapped by the httpjson client transport.
func (a *Adapter) generateSpeech(ctx context.Context, model string, spec capability.GenerationSpec) (*capability.Result, error) {
	var out struct {
		Audio string `json:"audio"`
		URL   string `json:"url"`
	}
	err := a.client.Do(ctx, "POST", "/audio/speech", speechRequest{
		Model: model,
		Input: buildPrompt(spec),
	}, &out)
	if err != nil {
		return nil, err
	}
	if out.URL == "" && out.Audio == "" {
		return nil, &capability.ProviderError{Kind: capability.ErrUnknown, Provider: a.ProviderName(), Message: "empty speech response"}
	}
	return &capability.Result{
		Provider:   a.ProviderName(),
		Model:      model,
		OutputURI:  out.URL,
		OutputText: out.Audio,
	}, nil
}

// modelsResponse is the OpenAI GET /models payload shape.
type modelsResponse struct {
	Data []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

// DiscoverModels implements provider.Discoverer via GET /models. The
// endpoint reports identifiers only; capability inference happens in the
// registry layer (name heuristics + operator confirmation), not here.
func (a *Adapter) DiscoverModels(ctx context.Context) ([]provider.DiscoveredModel, error) {
	var out modelsResponse
	if err := a.client.Do(ctx, "GET", "/models", nil, &out); err != nil {
		return nil, err
	}
	models := make([]provider.DiscoveredModel, 0, len(out.Data))
	for _, m := range out.Data {
		models = append(models, provider.DiscoveredModel{
			Identifier: m.ID,
			Metadata:   map[string]string{"owned_by": m.OwnedBy},
		})
	}
	return models, nil
}

// TestConnection implements provider.Tester — GET /models is the cheapest
// authenticated endpoint on every compatible server.
func (a *Adapter) TestConnection(ctx context.Context) error {
	var out modelsResponse
	return a.client.Do(ctx, "GET", "/models", nil, &out)
}

// buildPrompt renders the structured spec into provider text (§29 — prompt
// construction happens at the capability edge, not in domain code).
func buildPrompt(spec capability.GenerationSpec) string {
	if spec.Prompt != "" {
		return spec.Prompt
	}
	var b strings.Builder
	fmt.Fprintf(&b, "capability: %s\n", spec.Capability)
	if subj := spec.Subject["character"]; subj != "" {
		fmt.Fprintf(&b, "character: %s\n", subj)
	}
	for _, pair := range [][2]string{
		{"location", spec.Location},
		{"time", spec.TimeOfDay},
		{"emotion", spec.Emotion},
		{"action", spec.Action},
		{"style", spec.Style},
		{"aspect_ratio", spec.AspectRatio},
	} {
		if pair[1] != "" {
			fmt.Fprintf(&b, "%s: %s\n", pair[0], pair[1])
		}
	}
	if spec.Camera != nil {
		fmt.Fprintf(&b, "camera: %s %s %s\n", spec.Camera.Shot, spec.Camera.Angle, spec.Camera.Movement)
	}
	return b.String()
}
