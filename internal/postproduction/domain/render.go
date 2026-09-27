package domain

type RenderJob struct {
	ID          string `json:"id"`
	EditID      string `json:"edit_id"`
	OutputFormat string `json:"output_format"`
	OutputURL   string `json:"output_url"`
}
