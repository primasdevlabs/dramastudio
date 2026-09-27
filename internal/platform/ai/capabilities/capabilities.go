package capabilities

import "context"

type CapabilityType string

const (
	CapStoryBible          CapabilityType = "story_bible"
	CapStoryArchitecture   CapabilityType = "story_architecture"
	CapSeasonPlanning      CapabilityType = "season_planning"
	CapEpisodePlanning     CapabilityType = "episode_planning"
	CapScriptWriting       CapabilityType = "script_writing"
	CapDialogueWriting     CapabilityType = "dialogue_writing"
	CapCharacterCreation   CapabilityType = "character_creation"
	CapLocationCreation    CapabilityType = "location_creation"
	CapStoryboardGen       CapabilityType = "storyboard_generation"
	CapShotPlanning        CapabilityType = "shot_planning"
	CapImageGen            CapabilityType = "image_generation"
	CapVideoGen            CapabilityType = "video_generation"
	CapVoiceGen            CapabilityType = "voice_generation"
	CapMusicGen            CapabilityType = "music_generation"
	CapSfxGen              CapabilityType = "sfx_generation"
	CapContinuityAnalysis  CapabilityType = "continuity_analysis"
	CapQualityEvaluation   CapabilityType = "quality_evaluation"
)

type StructuredPrompt struct {
	Subject               string            `json:"subject"`
	Location              string            `json:"location"`
	TimeOfDay             string            `json:"time_of_day"`
	Emotion               string            `json:"emotion"`
	Action                string            `json:"action"`
	ShotType              string            `json:"shot_type"` // medium_close_up, wide, etc.
	Angle                 string            `json:"angle"`     // eye_level, low_angle, etc.
	AspectRatio           string            `json:"aspect_ratio"` // 9:16, 16:9
	Style                 string            `json:"style"`
	ContinuityConstraints map[string]string `json:"continuity_constraints"`
}

type GenerationRequest struct {
	Capability CapabilityType   `json:"capability"`
	ProjectID  string           `json:"project_id"`
	EpisodeID  string           `json:"episode_id"`
	Prompt     StructuredPrompt `json:"prompt"`
}

type GenerationResponse struct {
	OutputURL string  `json:"output_url"`
	Provider  string  `json:"provider"`
	Model     string  `json:"model"`
	Cost      float64 `json:"cost"`
}

type CapabilityExecutor interface {
	Execute(ctx context.Context, req GenerationRequest) (*GenerationResponse, error)
}
