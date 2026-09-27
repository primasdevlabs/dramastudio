package domain

// LocationVariant is a named state of a location (day/night/rain) with its
// own reference attributes and generated asset.
type LocationVariant struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"` // Day, Night, Rain, Fog
	Attributes map[string]string `json:"attributes"`
	AssetURL   string            `json:"asset_url"`
}

type Location struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Name        string            `json:"name"`
	Kind        string            `json:"kind"` // Interior, Exterior
	Description string            `json:"description"`
	ParentID    string            `json:"parent_id,omitempty"`
	VisualRef   string            `json:"visual_ref,omitempty"`
	Variants    []LocationVariant `json:"variants"`
}
