package domain

type Beat struct {
	ID          string `json:"id"`
	SceneID     string `json:"scene_id"`
	Seq         int    `json:"seq"`
	Action      string `json:"action"`
	Dialogue    string `json:"dialogue"`
	CharacterID string `json:"character_id"`
}
