package youtube

import "context"

type YouTubeClient struct{}

func (c *YouTubeClient) UploadVideo(ctx context.Context, videoURL string) error {
	return nil
}
