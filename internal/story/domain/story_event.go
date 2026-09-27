package domain

type StoryEvent struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	EpisodeID   string `json:"episode_id"`
}
