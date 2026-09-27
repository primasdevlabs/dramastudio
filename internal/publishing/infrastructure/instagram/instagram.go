package instagram

import "context"

type InstagramClient struct{}

func (c *InstagramClient) UploadReel(ctx context.Context, videoURL string) error {
	return nil
}
