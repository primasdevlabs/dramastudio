package domain

type Episode struct {
	ID       string `json:"id"`
	ArcID    string `json:"arc_id"`
	Number   int    `json:"number"`
	Title    string `json:"title"`
	SceneIDs []string `json:"scene_ids"`
}
