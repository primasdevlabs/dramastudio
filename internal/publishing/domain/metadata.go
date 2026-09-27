package domain

type PublishMetadata struct {
	Title       string   `json:"title"`
	Caption     string   `json:"caption"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
