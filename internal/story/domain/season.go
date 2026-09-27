package domain

type Season struct {
	ID         string   `json:"id"`
	SeriesID   string   `json:"series_id"`
	Number     int      `json:"number"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	EpisodeIDs []string `json:"episode_ids"`
}
