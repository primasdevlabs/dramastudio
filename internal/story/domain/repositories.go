package domain

import "context"

type StoryRepository interface {
	FindSeriesByID(ctx context.Context, id string) (*Series, error)
	FindEpisodeByID(ctx context.Context, id string) (*Episode, error)
}
