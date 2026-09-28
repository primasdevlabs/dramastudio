package domain

import (
	"context"
	"errors"
)

var (
	ErrAssetNotFound   = errors.New("asset not found")
	ErrVersionNotFound = errors.New("asset version not found")
	ErrJobNotFound     = errors.New("generation job not found")
	ErrInvalidState    = errors.New("invalid asset state transition")
)

type MediaRepository interface {
	SaveAsset(ctx context.Context, asset *Asset) error
	FindAssetByID(ctx context.Context, id string) (*Asset, error)
	ListAssetsByProject(ctx context.Context, projectID string) ([]*Asset, error)

	SaveAssetVersion(ctx context.Context, v *AssetVersion) error
	FindAssetVersion(ctx context.Context, assetID string, version int) (*AssetVersion, error)
	ListAssetVersions(ctx context.Context, assetID string) ([]*AssetVersion, error)

	SaveJob(ctx context.Context, job *GenerationJob) error
	FindJobByID(ctx context.Context, id string) (*GenerationJob, error)
	FindJobByProviderJobID(ctx context.Context, providerJobID string) (*GenerationJob, error)
	ListJobsByProject(ctx context.Context, projectID string) ([]*GenerationJob, error)
}
