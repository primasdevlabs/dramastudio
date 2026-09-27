package domain

type Scene struct {
	ID        string   `json:"id"`
	EpisodeID string   `json:"episode_id"`
	Number    int      `json:"number"`
	Location  string   `json:"location"`
	BeatIDs   []string `json:"beat_ids"`
}
