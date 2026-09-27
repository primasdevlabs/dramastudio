package domain

import "context"

type PostproductionRepository interface {
	SaveTimeline(ctx context.Context, t *Timeline) error
	FindTimelineByEpisode(ctx context.Context, episodeID string) (*Timeline, error)
	SaveRender(ctx context.Context, r *RenderTask) error
	FindRenderByID(ctx context.Context, id string) (*RenderTask, error)
}
