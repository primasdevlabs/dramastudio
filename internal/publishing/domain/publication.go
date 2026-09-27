package domain

import "time"

type Publication struct {
	ID          string          `json:"id"`
	EpisodeID   string          `json:"episode_id"`
	ChannelID   string          `json:"channel_id"`
	Metadata    PublishMetadata `json:"metadata"`
	PublishedAt time.Time       `json:"published_at"`
}
