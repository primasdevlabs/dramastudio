package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/publishing/domain"
)

type PublishingService struct {
	repo domain.PublishingRepository
}

func NewPublishingService(repo domain.PublishingRepository) *PublishingService {
	return &PublishingService{repo: repo}
}

func (s *PublishingService) CreatePublication(ctx context.Context, episodeID, channelID, title, caption string, tags []string) (*domain.Publication, error) {
	id := fmt.Sprintf("pub_%d", time.Now().UnixNano())
	pub := &domain.Publication{
		ID:        id,
		EpisodeID: episodeID,
		ChannelID: channelID,
		Metadata: domain.PublishMetadata{
			Title:   title,
			Caption: caption,
			Tags:    tags,
		},
		PublishedAt: time.Now().UTC(),
	}
	if err := s.repo.SavePublication(ctx, pub); err != nil {
		return nil, err
	}
	return pub, nil
}

func (s *PublishingService) ListPublications(ctx context.Context, projectID string) ([]*domain.Publication, error) {
	return s.repo.ListPublicationsByProject(ctx, projectID)
}
