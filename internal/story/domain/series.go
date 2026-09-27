package domain

type Series struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	SeasonIDs   []string `json:"season_ids"`
}
