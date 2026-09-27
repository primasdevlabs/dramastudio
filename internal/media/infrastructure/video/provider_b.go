package video

import "context"

type ProviderB struct{}

func (p *ProviderB) GenerateVideo(ctx context.Context, prompt string) (string, error) {
	return "provider_b_result_url", nil
}
