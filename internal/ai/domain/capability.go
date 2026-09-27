package domain

type AICapability string

const (
	CapabilityBible          AICapability = "story_bible"
	CapabilityStory          AICapability = "story_architecture"
	CapabilitySeason         AICapability = "season_planning"
	CapabilityEpisode        AICapability = "episode_planning"
	CapabilityScript         AICapability = "script_writing"
	CapabilityDialogue       AICapability = "dialogue_writing"
	CapabilityCharacter      AICapability = "character_creation"
	CapabilityLocation       AICapability = "location_creation"
	CapabilityStoryboard     AICapability = "storyboard"
	CapabilityShot           AICapability = "shot_planning"
	CapabilityImage          AICapability = "image_generation"
	CapabilityAnimation      AICapability = "animation"
	CapabilityVoice          AICapability = "voice"
	CapabilityMusic          AICapability = "music"
	CapabilitySFX            AICapability = "sfx_generation"
	CapabilityAssembly       AICapability = "assembly"
	CapabilityCaption        AICapability = "caption_generation"
	CapabilityContinuity     AICapability = "continuity"
	CapabilityQuality        AICapability = "quality_evaluation"
	CapabilityPublishing     AICapability = "publishing_metadata"
)
