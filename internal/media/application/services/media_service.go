package services

import (
	"context"
	"dramastudio/internal/media/domain"
)

type MediaService struct {
	repo domain.MediaRepository
}

func NewMediaService(repo domain.MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

func (s *MediaService) GetAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.repo.FindAssetByID(ctx, id)
}
