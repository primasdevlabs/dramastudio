package domain

type Season struct {
	ID       string `json:"id"`
	SeriesID string `json:"series_id"`
	Number   int    `json:"number"`
}
