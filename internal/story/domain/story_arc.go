package domain

type StoryArc struct {
	ID       string `json:"id"`
	SeasonID string `json:"season_id"`
	Title    string `json:"title"`
	Number   int    `json:"number"`
}
