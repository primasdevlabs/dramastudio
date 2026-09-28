package domain

import "time"

type AssetStatus string

const (
	AssetPending  AssetStatus = "PENDING"
	AssetApproved AssetStatus = "APPROVED"
	AssetRejected AssetStatus = "REJECTED"
	AssetArchived AssetStatus = "ARCHIVED"
)

// Asset is a versioned generated artifact. The current row tracks the
// pointer/latest-approved version; all attempts live in asset_versions (§72).
type Asset struct {
	ID              string                 `json:"id"`
	ProjectID       string                 `json:"project_id"`
	Type            MediaType              `json:"type"`
	CharacterID     string                 `json:"character_id,omitempty"`
	LocationID      string                 `json:"location_id,omitempty"`
	EpisodeID       string                 `json:"episode_id,omitempty"`
	SceneID         string                 `json:"scene_id,omitempty"`
	ShotID          string                 `json:"shot_id,omitempty"`
	Provider        string                 `json:"provider"`
	Model           string                 `json:"model"`
	Prompt          string                 `json:"prompt"`
	Spec            map[string]interface{} `json:"spec,omitempty"`
	ReferenceAssets []string               `json:"reference_assets,omitempty"`
	Parameters      map[string]interface{} `json:"parameters,omitempty"`
	URL             string                 `json:"url,omitempty"`
	Version         int                    `json:"version"`
	Status          AssetStatus            `json:"status"`
	Cost            float64                `json:"cost"`
	CreatedAt       time.Time              `json:"created_at"`
}
