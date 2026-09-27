package domain

// WorldRule is a canon-level constraint about the story world
// (physics, geography, era technology) that continuity checks enforce.
type WorldRule struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Text      string `json:"text"`
	Category  string `json:"category"`
}
