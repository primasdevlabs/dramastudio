package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/media/domain"
	"dramastudio/internal/media/infrastructure/video"
)

type MediaService struct {
	repo          domain.MediaRepository
	videoProvider video.VideoProvider
}

func NewMediaService(repo domain.MediaRepository, vp video.VideoProvider) *MediaService {
	if vp == nil {
		vp = video.NewWanVideoProvider("")
	}
	return &MediaService{
		repo:          repo,
		videoProvider: vp,
	}
}

func (s *MediaService) GenerateAsset(ctx context.Context, projectID, characterID, sceneID, shotID, prompt, providerName string, mediaType domain.MediaType) (*domain.Asset, error) {
	if providerName == "" {
		providerName = "wan"
	}
	url, err := s.videoProvider.GenerateVideo(ctx, prompt)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("asset_%d", time.Now().UnixNano())
	asset := &domain.Asset{
		ID:          id,
		ProjectID:   projectID,
		CharacterID: characterID,
		SceneID:     sceneID,
		ShotID:      shotID,
		Type:        mediaType,
		Provider:    providerName,
		Model:       "wan",
		Prompt:      prompt,
		URL:         url,
		Version:     1,
		Status:      domain.AssetApproved,
		Cost:        0.05,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.SaveAsset(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *MediaService) GetAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.repo.FindAssetByID(ctx, id)
}

func (s *MediaService) ListAssets(ctx context.Context, projectID string) ([]*domain.Asset, error) {
	return s.repo.ListAssetsByProject(ctx, projectID)
}
