package domain

type ProductionRun struct {
	ID        string          `json:"id"`
	EpisodeID string          `json:"episode_id"`
	Stage     ProductionStage `json:"stage"`
	Jobs      []ProductionJob `json:"jobs"`
}
