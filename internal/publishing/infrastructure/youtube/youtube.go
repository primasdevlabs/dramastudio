package youtube

import (
	"context"
	"time"

	"dramastudio/internal/platform/ai/provider/httpjson"
	"dramastudio/internal/publishing/domain"
)

// Client adapts a YouTube-compatible upload API behind the ChannelAdapter
// boundary. The endpoint shape is config-driven so the same code can point
// at the real Data API or a self-hosted relay.
type Client struct {
	http *httpjson.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{http: httpjson.New("youtube", baseURL, apiKey, httpjson.AuthBearer, 120*time.Second)}
}

func (c *Client) Platform() domain.ChannelType { return domain.ChannelYouTube }

type uploadRequest struct {
	VideoURL    string   `json:"video_url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
	AccountRef  string   `json:"account_ref,omitempty"`
}

type uploadResponse struct {
	VideoID string `json:"video_id"`
}

func (c *Client) Publish(ctx context.Context, videoURL string, meta domain.PublishMetadata, cfg map[string]interface{}) (string, map[string]interface{}, error) {
	accountRef, _ := cfg["account_ref"].(string)
	var out uploadResponse
	if err := c.http.Do(ctx, "POST", "/videos", uploadRequest{
		VideoURL:    videoURL,
		Title:       meta.Title,
		Description: firstNonEmpty(meta.Description, meta.Caption),
		Tags:        meta.Tags,
		AccountRef:  accountRef,
	}, &out); err != nil {
		return "", nil, err
	}
	return out.VideoID, map[string]interface{}{"video_id": out.VideoID}, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
