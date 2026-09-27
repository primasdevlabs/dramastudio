package tiktok

import "context"

type TikTokClient struct{}

func (c *TikTokClient) UploadVideo(ctx context.Context, videoURL string) error {
	return nil
}
