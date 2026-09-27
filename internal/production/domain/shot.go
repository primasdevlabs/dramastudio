package domain

type Shot struct {
	ID          string `json:"id"`
	SceneID     string `json:"scene_id"`
	Prompt      string `json:"prompt"`
	DurationSec float64 `json:"duration_sec"`
}
