package domain

type Prop struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LocationID  string `json:"location_id,omitempty"`
	VisualRef   string `json:"visual_ref,omitempty"`
}
