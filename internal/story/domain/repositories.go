package domain

import (
	"context"
	"errors"
)

var (
	ErrSeriesNotFound  = errors.New("series not found")
	ErrSeasonNotFound  = errors.New("season not found")
	ErrArcNotFound     = errors.New("arc not found")
	ErrEpisodeNotFound = errors.New("episode not found")
	ErrSceneNotFound   = errors.New("scene not found")
	ErrBeatNotFound    = errors.New("beat not found")
	ErrNodeNotFound    = errors.New("story graph node not found")
)

type StoryRepository interface {
	SaveSeries(ctx context.Context, series *Series) error
	FindSeriesByID(ctx context.Context, id string) (*Series, error)
	FindSeriesByProject(ctx context.Context, projectID string) (*Series, error)

	SaveSeason(ctx context.Context, season *Season) error
	FindSeasonByID(ctx context.Context, id string) (*Season, error)
	ListSeasonsBySeries(ctx context.Context, seriesID string) ([]*Season, error)

	SaveArc(ctx context.Context, arc *StoryArc) error
	FindArcByID(ctx context.Context, id string) (*StoryArc, error)
	ListArcsBySeason(ctx context.Context, seasonID string) ([]*StoryArc, error)

	SaveEpisode(ctx context.Context, episode *Episode) error
	FindEpisodeByID(ctx context.Context, id string) (*Episode, error)
	ListEpisodesBySeason(ctx context.Context, seasonID string) ([]*Episode, error)

	SaveScene(ctx context.Context, scene *Scene) error
	FindSceneByID(ctx context.Context, id string) (*Scene, error)
	ListScenesByEpisode(ctx context.Context, episodeID string) ([]*Scene, error)

	SaveBeat(ctx context.Context, beat *Beat) error
	FindBeatByID(ctx context.Context, id string) (*Beat, error)
	ListBeatsByScene(ctx context.Context, sceneID string) ([]*Beat, error)

	SaveGraphNode(ctx context.Context, projectID string, node *StoryNode) error
	FindGraphNode(ctx context.Context, id string) (*StoryNode, error)
	ListGraphNodes(ctx context.Context, projectID string) ([]*StoryNode, error)
	SaveGraphEdge(ctx context.Context, projectID string, edge *StoryEdge) error
	ListGraphEdges(ctx context.Context, projectID string) ([]*StoryEdge, error)

	SavePlotThread(ctx context.Context, t *PlotThread) error
	ListPlotThreads(ctx context.Context, projectID string) ([]*PlotThread, error)
}
