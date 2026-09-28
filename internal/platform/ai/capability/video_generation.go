package capability

// VideoGeneration covers text_to_video and image_to_video shot synthesis.
const VideoGeneration = "video_generation"

type VideoSpec struct {
	Prompt         string   `json:"prompt"`
	ImageReference []string `json:"image_reference,omitempty"`
	DurationSec    float64  `json:"duration_sec,omitempty"`
	AspectRatio    string   `json:"aspect_ratio,omitempty"`
	CameraMotion   string   `json:"camera_motion,omitempty"`
	FPS            int      `json:"fps,omitempty"`
}
