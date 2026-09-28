package domain

import "context"

// ChannelAdapter is the provider boundary for a distribution platform.
// Implementations live in infrastructure (youtube/, tiktok/, instagram/).
type ChannelAdapter interface {
	Platform() ChannelType
	// Publish uploads the finished video and returns the platform's
	// external content id plus any response payload worth persisting.
	Publish(ctx context.Context, videoURL string, meta PublishMetadata, cfg map[string]interface{}) (externalID string, response map[string]interface{}, err error)
}
