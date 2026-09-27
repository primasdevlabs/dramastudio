package image

import "context"

type ImageProvider interface {
	GenerateImage(ctx context.Context, prompt string) (string, error)
}
