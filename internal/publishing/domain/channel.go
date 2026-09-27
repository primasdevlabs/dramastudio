package domain

type ChannelType string

const (
	ChannelTikTok    ChannelType = "tiktok"
	ChannelYouTube   ChannelType = "youtube"
	ChannelInstagram ChannelType = "instagram"
)

type Channel struct {
	ID   string      `json:"id"`
	Type ChannelType `json:"type"`
	Name string      `json:"name"`
}
