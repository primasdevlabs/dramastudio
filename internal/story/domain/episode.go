package domain

// EpisodeStatus is the episode-level state machine.
type EpisodeStatus string

const (
	EpisodePlanned    EpisodeStatus = "PLANNED"
	EpisodeScripted   EpisodeStatus = "SCRIPTED"
	EpisodeProducing  EpisodeStatus = "PRODUCING"
	EpisodeValidating EpisodeStatus = "VALIDATING"
	EpisodeCompleted  EpisodeStatus = "COMPLETED"
	EpisodePublished  EpisodeStatus = "PUBLISHED"
)

type Episode struct {
	ID          string        `json:"id"`
	SeasonID    string        `json:"season_id"`
	ArcID       string        `json:"arc_id"`
	Number      int           `json:"number"`
	Title       string        `json:"title"`
	Summary     string        `json:"summary"`
	Script      string        `json:"script"`
	Status      EpisodeStatus `json:"status"`
	SceneIDs    []string      `json:"scene_ids"`
}
