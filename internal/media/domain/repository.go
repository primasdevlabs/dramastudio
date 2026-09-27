package domain

import "context"

type MediaRepository interface {
	SaveAsset(ctx context.Context, asset *Asset) error
	FindAssetByID(ctx context.Context, id string) (*Asset, error)
	ListAssetsByProject(ctx context.Context, projectID string) ([]*Asset, error)
	SaveJob(ctx context.Context, job *GenerationJob) error
	FindJobByID(ctx context.Context, id string) (*GenerationJob, error)
}
