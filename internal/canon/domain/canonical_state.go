package domain

type CanonicalState struct {
	EpisodeID string      `json:"episode_id"`
	Facts     []StoryFact `json:"facts"`
}
