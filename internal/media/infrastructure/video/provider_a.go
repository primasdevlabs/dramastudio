package video

import "context"

type ProviderA struct{}

func (p *ProviderA) GenerateVideo(ctx context.Context, prompt string) (string, error) {
	return "provider_a_result_url", nil
}
