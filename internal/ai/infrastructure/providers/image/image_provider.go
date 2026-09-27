package image

import "context"

type ImageProvider struct{}

func (p *ImageProvider) GenerateImage(ctx context.Context, prompt string) (string, error) {
	return "https://media.dramastudio.internal/image.png", nil
}
