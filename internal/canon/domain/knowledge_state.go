package domain

type KnowledgeState struct {
	CharacterID  string   `json:"character_id"`
	EpisodeID    string   `json:"episode_id"`
	KnownFactIDs []string `json:"known_fact_ids"`
}
