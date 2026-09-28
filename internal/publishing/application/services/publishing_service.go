package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/publishing/domain"
)

type PublishingService struct {
	repo     domain.PublishingRepository
	adapters map[domain.ChannelType]domain.ChannelAdapter
}

func NewPublishingService(repo domain.PublishingRepository) *PublishingService {
	return &PublishingService{repo: repo, adapters: map[domain.ChannelType]domain.ChannelAdapter{}}
}

// RegisterAdapter wires a platform adapter (youtube/tiktok/instagram).
func (s *PublishingService) RegisterAdapter(a domain.ChannelAdapter) {
	s.adapters[a.Platform()] = a
}

// --- Channels ---

func (s *PublishingService) RegisterChannel(ctx context.Context, projectID string, platform domain.ChannelType, name, accountRef string, config map[string]interface{}) (*domain.Channel, error) {
	c := &domain.Channel{
		ID:         "ch_" + uuid.NewString(),
		ProjectID:  projectID,
		Platform:   platform,
		Name:       name,
		AccountRef: accountRef,
		Config:     config,
		Enabled:    true,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.SaveChannel(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *PublishingService) ListChannels(ctx context.Context, projectID string) ([]*domain.Channel, error) {
	return s.repo.ListChannels(ctx, projectID)
}

// --- Publications ---

// SchedulePublication creates a publication in draft or scheduled state.
// Idempotent on IdempotencyKey.
func (s *PublishingService) SchedulePublication(ctx context.Context, projectID, episodeID, channelID, videoURL string, meta domain.PublishMetadata, scheduledAt *time.Time, idempotencyKey string) (*domain.Publication, error) {
	if idempotencyKey != "" {
		if existing, err := s.repo.FindPublicationByIdempotencyKey(ctx, idempotencyKey); err == nil {
			return existing, nil
		}
	}
	status := domain.PubDraft
	if scheduledAt != nil {
		status = domain.PubScheduled
	}
	pub := &domain.Publication{
		ID:             "pub_" + uuid.NewString(),
		ProjectID:      projectID,
		EpisodeID:      episodeID,
		ChannelID:      channelID,
		VideoURL:       videoURL,
		Metadata:       meta,
		ScheduledAt:    scheduledAt,
		Status:         status,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now().UTC(),
	}
	if err := s.repo.SavePublication(ctx, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

// Publish executes a publication against the channel's platform adapter.
func (s *PublishingService) Publish(ctx context.Context, publicationID string) (*domain.Publication, error) {
	pub, err := s.repo.FindPublicationByID(ctx, publicationID)
	if err != nil {
		return nil, err
	}
	if !pub.CanTransitionTo(domain.PubPublishing) {
		return nil, domain.ErrInvalidTransition
	}
	channel, err := s.repo.FindChannelByID(ctx, pub.ChannelID)
	if err != nil {
		return nil, err
	}
	if !channel.Enabled {
		return nil, fmt.Errorf("channel %s is disabled", channel.ID)
	}
	adapter, ok := s.adapters[channel.Platform]
	if !ok {
		return nil, fmt.Errorf("no adapter registered for platform %q", channel.Platform)
	}

	pub.Status = domain.PubPublishing
	_ = s.repo.SavePublication(ctx, pub)

	externalID, resp, err := adapter.Publish(ctx, pub.VideoURL, pub.Metadata, channel.Config)
	if err != nil {
		pub.Status = domain.PubFailed
		pub.PlatformResponse = map[string]interface{}{"error": err.Error()}
		_ = s.repo.SavePublication(ctx, pub)
		return nil, err
	}
	now := time.Now().UTC()
	pub.Status = domain.PubPublished
	pub.ExternalID = externalID
	pub.PlatformResponse = resp
	pub.PublishedAt = &now
	if err := s.repo.SavePublication(ctx, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

func (s *PublishingService) GetPublication(ctx context.Context, id string) (*domain.Publication, error) {
	return s.repo.FindPublicationByID(ctx, id)
}

func (s *PublishingService) ListPublications(ctx context.Context, projectID string) ([]*domain.Publication, error) {
	return s.repo.ListPublicationsByProject(ctx, projectID)
}

// PublishDue sweeps scheduled publications whose time has arrived and
// publishes each through its channel adapter (§50). Per-item failures are
// logged by callers; each publication is independent.
func (s *PublishingService) PublishDue(ctx context.Context, now time.Time) (published, failed int, err error) {
	due, err := s.repo.ListDueScheduled(ctx, now)
	if err != nil {
		return 0, 0, err
	}
	for _, p := range due {
		if _, perr := s.Publish(ctx, p.ID); perr != nil {
			failed++
			continue
		}
		published++
	}
	return published, failed, nil
}
