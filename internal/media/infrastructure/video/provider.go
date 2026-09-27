package video

import "context"

type VideoProvider interface {
	GenerateVideo(ctx context.Context, prompt string) (string, error)
}
