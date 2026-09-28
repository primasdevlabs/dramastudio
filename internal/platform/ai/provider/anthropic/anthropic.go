// Package anthropic is the native adapter for the Anthropic Messages API.
// It exists because the contract genuinely differs from the
// OpenAI-compatible shape — that is the rule for when a native adapter is
// warranted (§24).
package anthropic

import (
	"context"
	"strings"
	"time"

	"dramastudio/internal/platform/ai/capability"
	"dramastudio/internal/platform/ai/provider"
	"dramastudio/internal/platform/ai/provider/httpjson"
)

type Adapter struct {
	name   string
	client *httpjson.Client
}

func New(name, baseURL, apiKey string) *Adapter {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	c := httpjson.New(name, baseURL, apiKey, httpjson.AuthXAPIKey, 120*time.Second)
	c.ExtraHeaders = map[string]string{"anthropic-version": "2023-06-01"}
	return &Adapter{name: name, client: c}
}

func init() {
	provider.RegisterFactory("anthropic", func(cfg provider.Config) provider.Adapter {
		return provider.Adapter{Sync: New(cfg.Name, cfg.BaseURL, cfg.APIKey)}
	})
}

func (a *Adapter) ProviderName() string { return a.name }

type messagesRequest struct {
	Model     string        `json:"model"`
	MaxTokens int           `json:"max_tokens"`
	Messages  []chatMessage `json:"messages"`
	System    string        `json:"system,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messagesResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// Generate handles text capabilities via /messages.
func (a *Adapter) Generate(ctx context.Context, model string, spec capability.GenerationSpec) (*capability.Result, error) {
	var out messagesResponse
	err := a.client.Do(ctx, "POST", "/messages", messagesRequest{
		Model:     model,
		MaxTokens: 8192,
		Messages:  []chatMessage{{Role: "user", Content: buildPrompt(spec)}},
	}, &out)
	if err != nil {
		return nil, err
	}
	var text strings.Builder
	for _, c := range out.Content {
		if c.Type == "text" {
			text.WriteString(c.Text)
		}
	}
	if text.Len() == 0 {
		return nil, &capability.ProviderError{Kind: capability.ErrUnknown, Provider: a.ProviderName(), Message: "empty completion"}
	}
	return &capability.Result{
		Provider:   a.ProviderName(),
		Model:      model,
		OutputText: text.String(),
		Metadata: map[string]interface{}{
			"input_tokens":  out.Usage.InputTokens,
			"output_tokens": out.Usage.OutputTokens,
		},
	}, nil
}

// TestConnection implements provider.Tester — a minimal /messages call.
func (a *Adapter) TestConnection(ctx context.Context) error {
	var out messagesResponse
	return a.client.Do(ctx, "POST", "/messages", messagesRequest{
		Model:     "claude-haiku-4-5",
		MaxTokens: 1,
		Messages:  []chatMessage{{Role: "user", Content: "ping"}},
	}, &out)
}

func buildPrompt(spec capability.GenerationSpec) string {
	if spec.Prompt != "" {
		return spec.Prompt
	}
	var b strings.Builder
	b.WriteString("capability: " + spec.Capability + "\n")
	if spec.Action != "" {
		b.WriteString("action: " + spec.Action + "\n")
	}
	if spec.Style != "" {
		b.WriteString("style: " + spec.Style + "\n")
	}
	return b.String()
}
