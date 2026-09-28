package domain

type AICapability string

const (
	CapabilityBible      AICapability = "story_bible"
	CapabilityStory      AICapability = "story_architecture"
	CapabilitySeason     AICapability = "season_planning"
	CapabilityEpisode    AICapability = "episode_planning"
	CapabilityScript     AICapability = "script_writing"
	CapabilityDialogue   AICapability = "dialogue_writing"
	CapabilityCharacter  AICapability = "character_creation"
	CapabilityLocation   AICapability = "location_creation"
	CapabilityStoryboard AICapability = "storyboard"
	CapabilityShot       AICapability = "shot_planning"
	CapabilityImage      AICapability = "image_generation"
	CapabilityAnimation  AICapability = "animation"
	CapabilityVoice      AICapability = "voice"
	CapabilityMusic      AICapability = "music"
	CapabilitySFX        AICapability = "sfx_generation"
	CapabilityAssembly   AICapability = "assembly"
	CapabilityCaption    AICapability = "caption_generation"
	CapabilityContinuity AICapability = "continuity"
	CapabilityQuality    AICapability = "quality_evaluation"
	CapabilityPublishing AICapability = "publishing_metadata"

	// Matrix-aligned additions: reasoning, video, image editing,
	// speech-to-text, and vision/QA are distinct routable capabilities.
	CapabilityLeadDirector AICapability = "lead_director"
	CapabilityVideo        AICapability = "video_generation"
	CapabilityImageEdit    AICapability = "image_editing"
	CapabilitySpeechToText AICapability = "speech_to_text"
	CapabilityVision       AICapability = "vision"
)

// AllCapabilities enumerates the routable task capabilities for the
// settings UI and catalog validation.
var AllCapabilities = []AICapability{
	CapabilityLeadDirector,
	CapabilityBible,
	CapabilityStory,
	CapabilitySeason,
	CapabilityEpisode,
	CapabilityScript,
	CapabilityDialogue,
	CapabilityContinuity,
	CapabilityQuality,
	CapabilityCharacter,
	CapabilityLocation,
	CapabilityStoryboard,
	CapabilityShot,
	CapabilityImage,
	CapabilityImageEdit,
	CapabilityAnimation,
	CapabilityVideo,
	CapabilityVoice,
	CapabilityMusic,
	CapabilitySFX,
	CapabilitySpeechToText,
	CapabilityVision,
	CapabilityAssembly,
	CapabilityCaption,
	CapabilityPublishing,
}
