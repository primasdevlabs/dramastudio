package domain

type SceneProduction struct {
	SceneID string `json:"scene_id"`
	Shots   []Shot `json:"shots"`
}
