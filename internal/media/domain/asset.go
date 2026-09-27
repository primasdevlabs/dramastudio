package domain

import "time"

type AssetStatus string

const (
	AssetPending   AssetStatus = "PENDING"
	AssetApproved  AssetStatus = "APPROVED"
	AssetRejected  AssetStatus = "REJECTED"
	AssetArchived  AssetStatus = "ARCHIVED"
)

type Asset struct {
	ID         string         `json:"id"`
	ProjectID  string         `json:"project_id"`
	CharacterID string        `json:"character_id,omitempty"`
	SceneID    string         `json:"scene_id,omitempty"`
	ShotID     string         `json:"shot_id,omitempty"`
	Type       MediaType      `json:"type"`
	Provider   string         `json:"provider"`
	Model      string         `json:"model"`
	Prompt     string         `json:"prompt"`
	URL        string         `json:"url"`
	Version    int            `json:"version"`
	Status     AssetStatus    `json:"status"`
	Cost       float64        `json:"cost"`
	CreatedAt  time.Time      `json:"created_at"`
}
