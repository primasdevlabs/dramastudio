package capability

// Vision covers image understanding for continuity/QA review.
const Vision = "vision"

type VisionSpec struct {
	ImageURIs []string `json:"image_uris"`
	Prompt    string   `json:"prompt"`
}
