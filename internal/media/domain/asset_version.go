package domain

import "time"

// AssetVersion is one immutable generation attempt for an asset.
type AssetVersion struct {
	ID        string      `json:"id"`
	AssetID   string      `json:"asset_id"`
	Version   int         `json:"version"`
	ObjectKey string      `json:"object_key"`
	URL       string      `json:"url"`
	Status    AssetStatus `json:"status"`
	Cost      float64     `json:"cost"`
	CreatedAt time.Time   `json:"created_at"`
}
