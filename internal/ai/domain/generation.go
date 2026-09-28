package domain

import (
	"encoding/json"
	"time"
)

type GenerationStatus string

const (
	StatusPending   GenerationStatus = "PENDING"
	StatusRunning   GenerationStatus = "RUNNING"
	StatusSucceeded GenerationStatus = "SUCCEEDED"
	StatusFailed    GenerationStatus = "FAILED"
	StatusCancelled GenerationStatus = "CANCELLED"
)

// DialogueContextPackage is the structured prompt payload for dialogue
// generation (§30).
type DialogueContextPackage struct {
	Task                  string              `json:"task"`
	Episode               int                 `json:"episode"`
	Scene                 int                 `json:"scene"`
	Characters            []string            `json:"characters"`
	SceneObjective        string              `json:"scene_objective"`
	StoryContext          string              `json:"story_context"`
	CharacterState        map[string]string   `json:"character_state"`
	KnowledgeState        map[string][]string `json:"knowledge_state"`
	Relationships         []string            `json:"relationships"`
	Tone                  string              `json:"tone"`
	CanonicalLanguage     string              `json:"canonical_language"`
	ContinuityConstraints []string            `json:"continuity_constraints"`
}

type StoryboardContextPackage struct {
	Episode               int      `json:"episode"`
	Scene                 int      `json:"scene"`
	ApprovedScript        string   `json:"approved_script"`
	ShotRequirements      []string `json:"shot_requirements"`
	CharacterReferences   []string `json:"character_references"`
	LocationReferences    []string `json:"location_references"`
	VisualBible           string   `json:"visual_bible"`
	CameraLanguage        string   `json:"camera_language"`
	ContinuityConstraints []string `json:"continuity_constraints"`
}

type AnimationContextPackage struct {
	Shot                       string   `json:"shot"`
	ApprovedStoryboard         string   `json:"approved_storyboard"`
	CharacterReference         string   `json:"character_reference"`
	LocationReference          string   `json:"location_reference"`
	MotionSpecification        string   `json:"motion_specification"`
	CameraMovement             string   `json:"camera_movement"`
	DurationSec                float64  `json:"duration_sec"`
	PreviousNextShotReferences []string `json:"previous_next_shot_references"`
}

// GenerationJob is a persistent, auditable AI generation record (§26).
type GenerationJob struct {
	ID             string           `json:"id"`
	ProjectID      string           `json:"project_id"`
	Capability     AICapability     `json:"capability"`
	ProviderID     string           `json:"provider_id"`
	ModelID        ModelID          `json:"model_id"`
	ModelVersion   string           `json:"model_version,omitempty"` // provenance: model.Version at run time
	Input          json.RawMessage  `json:"input"`
	Output         json.RawMessage  `json:"output"`
	Status         GenerationStatus `json:"status"`
	Attempt        int              `json:"attempt"`
	Cost           float64          `json:"cost"`
	ProviderJobID  string           `json:"provider_job_id,omitempty"`
	IdempotencyKey string           `json:"idempotency_key,omitempty"`
	Error          string           `json:"error,omitempty"`
	StartedAt      *time.Time       `json:"started_at,omitempty"`
	CompletedAt    *time.Time       `json:"completed_at,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
}
