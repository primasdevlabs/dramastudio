package domain

// ShotStatus is the lifecycle of a planned shot.
type ShotStatus string

const (
	ShotStatusPlanned    ShotStatus = "planned"
	ShotStatusGenerating ShotStatus = "generating"
	ShotStatusGenerated  ShotStatus = "generated"
	ShotStatusApproved   ShotStatus = "approved"
	ShotStatusRejected   ShotStatus = "rejected"
)

// Shot is a planned camera unit within a scene.
type Shot struct {
	ID            string                 `json:"id"`
	ProjectID     string                 `json:"project_id"`
	EpisodeID     string                 `json:"episode_id"`
	SceneID       string                 `json:"scene_id"`
	Seq           int                    `json:"seq"`
	Description   string                 `json:"description"`
	Camera        map[string]interface{} `json:"camera,omitempty"`
	Characters    []string               `json:"characters,omitempty"`
	LocationID    string                 `json:"location_id,omitempty"`
	DurationSec   float64                `json:"duration_sec"`
	Status        ShotStatus             `json:"status"`
	ApprovedAsset string                 `json:"approved_asset,omitempty"`
}
