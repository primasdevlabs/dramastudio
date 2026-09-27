package domain

type Generation struct {
	ID        string    `json:"id"`
	Prompt    string    `json:"prompt"`
	MediaType MediaType `json:"media_type"`
	Provider  string    `json:"provider"`
}
