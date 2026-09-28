package domain

type Edit struct {
	ID            string `json:"id"`
	EpisodeID     string `json:"episode_id"`
	CompositionID string `json:"composition_id"`
}
