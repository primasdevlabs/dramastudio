package domain

import "context"

type MediaRepository interface {
	FindAssetByID(ctx context.Context, id string) (*Asset, error)
	SaveAsset(ctx context.Context, asset *Asset) error
}
