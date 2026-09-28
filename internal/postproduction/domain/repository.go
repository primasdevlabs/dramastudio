package domain

import (
	"context"
	"errors"
)

var (
	ErrTimelineNotFound = errors.New("timeline not found")
	ErrRenderNotFound   = errors.New("render not found")
)

type PostproductionRepository interface {
	SaveTimeline(ctx context.Context, t *Timeline) error
	FindTimelineByEpisode(ctx context.Context, episodeID string) (*Timeline, error)
	ListTimelineVersions(ctx context.Context, episodeID string) ([]*Timeline, error)

	SaveRender(ctx context.Context, r *RenderTask) error
	FindRenderByID(ctx context.Context, id string) (*RenderTask, error)
	ListRendersByEpisode(ctx context.Context, episodeID string) ([]*RenderTask, error)
}
