package domain

import "time"

type PublicationStatus string

const (
	PubDraft      PublicationStatus = "draft"
	PubScheduled  PublicationStatus = "scheduled"
	PubPublishing PublicationStatus = "publishing"
	PubPublished  PublicationStatus = "published"
	PubFailed     PublicationStatus = "failed"
)

// Publication is one release of an episode to a channel (§45).
type Publication struct {
	ID               string                 `json:"id"`
	ProjectID        string                 `json:"project_id"`
	EpisodeID        string                 `json:"episode_id"`
	ChannelID        string                 `json:"channel_id"`
	Metadata         PublishMetadata        `json:"metadata"`
	VideoURL         string                 `json:"video_url"`
	ScheduledAt      *time.Time             `json:"scheduled_at,omitempty"`
	Status           PublicationStatus      `json:"status"`
	PlatformResponse map[string]interface{} `json:"platform_response,omitempty"`
	ExternalID       string                 `json:"external_id,omitempty"`
	PublishedAt      *time.Time             `json:"published_at,omitempty"`
	IdempotencyKey   string                 `json:"idempotency_key,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
}

// CanTransitionTo enforces the publication lifecycle.
func (p *Publication) CanTransitionTo(next PublicationStatus) bool {
	switch p.Status {
	case PubDraft:
		return next == PubScheduled || next == PubPublishing
	case PubScheduled:
		return next == PubPublishing || next == PubDraft
	case PubPublishing:
		return next == PubPublished || next == PubFailed
	case PubFailed:
		return next == PubPublishing // retry
	case PubPublished:
		return false
	}
	return false
}
