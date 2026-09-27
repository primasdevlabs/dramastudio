package domain

import "context"

type StoryRepository interface {
	SaveSeries(ctx context.Context, series *Series) error
	FindSeriesByID(ctx context.Context, id string) (*Series, error)
	FindSeriesByProject(ctx context.Context, projectID string) (*Series, error)

	SaveSeason(ctx context.Context, season *Season) error
	FindSeasonByID(ctx context.Context, id string) (*Season, error)
	ListSeasonsBySeries(ctx context.Context, seriesID string) ([]*Season, error)

	SaveEpisode(ctx context.Context, episode *Episode) error
	FindEpisodeByID(ctx context.Context, id string) (*Episode, error)
	ListEpisodesBySeason(ctx context.Context, seasonID string) ([]*Episode, error)

	SaveScene(ctx context.Context, scene *Scene) error
	FindSceneByID(ctx context.Context, id string) (*Scene, error)
	ListScenesByEpisode(ctx context.Context, episodeID string) ([]*Scene, error)
}
