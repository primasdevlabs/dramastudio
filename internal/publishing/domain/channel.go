package domain

import "time"

type ChannelType string

const (
	ChannelTikTok    ChannelType = "tiktok"
	ChannelYouTube   ChannelType = "youtube"
	ChannelInstagram ChannelType = "instagram"
)

// Channel is a distribution target for a project.
type Channel struct {
	ID         string                 `json:"id"`
	ProjectID  string                 `json:"project_id"`
	Platform   ChannelType            `json:"platform"`
	Name       string                 `json:"name"`
	AccountRef string                 `json:"account_ref,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Enabled    bool                   `json:"enabled"`
	CreatedAt  time.Time              `json:"created_at"`
}
