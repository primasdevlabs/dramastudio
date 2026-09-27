package video

import "context"

type VideoProvider struct{}

func (p *VideoProvider) GenerateVideo(ctx context.Context, prompt string) (string, error) {
	return "https://media.dramastudio.internal/video.mp4", nil
}
