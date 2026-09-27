package video

import "context"

type ProviderC struct{}

func (p *ProviderC) GenerateVideo(ctx context.Context, prompt string) (string, error) {
	return "provider_c_result_url", nil
}
