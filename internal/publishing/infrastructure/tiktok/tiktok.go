package tiktok

import (
	"context"
	"time"

	"dramastudio/internal/platform/ai/provider/httpjson"
	"dramastudio/internal/publishing/domain"
)

// Client adapts a TikTok-compatible posting API behind ChannelAdapter.
type Client struct {
	http *httpjson.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{http: httpjson.New("tiktok", baseURL, apiKey, httpjson.AuthBearer, 120*time.Second)}
}

func (c *Client) Platform() domain.ChannelType { return domain.ChannelTikTok }

func (c *Client) Publish(ctx context.Context, videoURL string, meta domain.PublishMetadata, cfg map[string]interface{}) (string, map[string]interface{}, error) {
	body := map[string]interface{}{
		"video_url":   videoURL,
		"caption":     meta.Caption,
		"title":       meta.Title,
		"account_ref": cfg["account_ref"],
	}
	var out struct {
		ShareID string `json:"share_id"`
	}
	if err := c.http.Do(ctx, "POST", "/posts", body, &out); err != nil {
		return "", nil, err
	}
	return out.ShareID, map[string]interface{}{"share_id": out.ShareID}, nil
}
