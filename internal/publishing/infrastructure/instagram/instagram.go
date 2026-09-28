package instagram

import (
	"context"
	"time"

	"dramastudio/internal/platform/ai/provider/httpjson"
	"dramastudio/internal/publishing/domain"
)

// Client adapts an Instagram-compatible reels API behind ChannelAdapter.
type Client struct {
	http *httpjson.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{http: httpjson.New("instagram", baseURL, apiKey, httpjson.AuthBearer, 120*time.Second)}
}

func (c *Client) Platform() domain.ChannelType { return domain.ChannelInstagram }

func (c *Client) Publish(ctx context.Context, videoURL string, meta domain.PublishMetadata, cfg map[string]interface{}) (string, map[string]interface{}, error) {
	body := map[string]interface{}{
		"video_url":   videoURL,
		"caption":     meta.Caption,
		"account_ref": cfg["account_ref"],
	}
	var out struct {
		MediaID string `json:"media_id"`
	}
	if err := c.http.Do(ctx, "POST", "/media", body, &out); err != nil {
		return "", nil, err
	}
	return out.MediaID, map[string]interface{}{"media_id": out.MediaID}, nil
}
