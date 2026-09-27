package domain

type Asset struct {
	ID        string         `json:"id"`
	MediaType MediaType      `json:"media_type"`
	Versions  []AssetVersion `json:"versions"`
}
