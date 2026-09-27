package openai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dramastudio/internal/platform/ai/contracts"
	"dramastudio/internal/platform/ai/providers/httpjson"
)

// Adapter speaks to OpenAI-compatible endpoints (chat completions for text,
// images API for image generation). Works against OpenAI and compatible
// providers by configuring BaseURL.
type Adapter struct {
	client *httpjson.Client
}

func New(baseURL, apiKey string) *Adapter {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &Adapter{client: httpjson.New("openai", baseURL, apiKey, httpjson.AuthBearer, 120*time.Second)}
}

func (a *Adapter) ProviderName() string { return "openai" }

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

// Generate handles text-oriented capabilities (scripts, dialogue, bibles).
func (a *Adapter) Generate(ctx context.Context, model string, spec contracts.GenerationSpec) (*contracts.Result, error) {
	if strings.Contains(string(spec.Capability), "image") {
		return a.generateImage(ctx, model, spec)
	}
	prompt := buildPrompt(spec)
	var out chatResponse
	err := a.client.Do(ctx, "POST", "/chat/completions", chatRequest{
		Model:    model,
		Messages: []chatMessage{{Role: "user", Content: prompt}},
	}, &out)
	if err != nil {
		return nil, err
	}
	if len(out.Choices) == 0 {
		return nil, &contracts.ProviderError{Kind: contracts.ErrUnknown, Provider: a.ProviderName(), Message: "empty completion"}
	}
	return &contracts.Result{
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
		URL string `json:"url"`
	} `json:"data"`
}

func (a *Adapter) generateImage(ctx context.Context, model string, spec contracts.GenerationSpec) (*contracts.Result, error) {
	var out imageResponse
	err := a.client.Do(ctx, "POST", "/images/generations", imageRequest{
		Model:  model,
		Prompt: buildPrompt(spec),
	}, &out)
	if err != nil {
		return nil, err
	}
	if len(out.Data) == 0 {
		return nil, &contracts.ProviderError{Kind: contracts.ErrUnknown, Provider: a.ProviderName(), Message: "empty image response"}
	}
	return &contracts.Result{
		Provider:  a.ProviderName(),
		Model:     model,
		OutputURI: out.Data[0].URL,
	}, nil
}

// buildPrompt renders the structured spec into provider text (§29 — prompt
// construction happens at the capability edge, not in domain code).
func buildPrompt(spec contracts.GenerationSpec) string {
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
