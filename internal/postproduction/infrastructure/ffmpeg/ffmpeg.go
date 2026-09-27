package ffmpeg

import "context"

type FFmpegAdapter struct{}

func (f *FFmpegAdapter) Render(ctx context.Context, editID string) (string, error) {
	return "rendered_video.mp4", nil
}
