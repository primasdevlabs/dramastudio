package domain

type Composition struct {
	ID       string   `json:"id"`
	Timeline Timeline `json:"timeline"`
}
