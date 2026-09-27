package domain

type LocationVariant struct {
	ID        string `json:"id"`
	TimeOfDay string `json:"time_of_day"` // Day, Night, Rain, Fog
	AssetURL  string `json:"asset_url"`
}

type Location struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"` // Interior, Exterior
	Variants    []LocationVariant `json:"variants"`
}
