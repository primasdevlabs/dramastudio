package llm

import "context"

type LLMProvider struct{}

func (p *LLMProvider) GenerateText(ctx context.Context, prompt string) (string, error) {
	return "generated text response", nil
}
