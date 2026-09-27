package domain

type Scene struct {
	ID           string   `json:"id"`
	EpisodeID    string   `json:"episode_id"`
	Number       int      `json:"number"`
	Title        string   `json:"title"`
	LocationID   string   `json:"location_id"`
	TimeOfDay    string   `json:"time_of_day"`
	Description  string   `json:"description"`
	CharacterIDs []string `json:"character_ids"`
	BeatIDs      []string `json:"beat_ids"`
}
