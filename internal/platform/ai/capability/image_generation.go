package capability

// ImageGeneration covers stills: character refs, storyboard frames,
// location plates, posters.
const ImageGeneration = "image_generation"

type ImageSpec struct {
	Prompt      string   `json:"prompt"`
	Size        string   `json:"size,omitempty"`
	AspectRatio string   `json:"aspect_ratio,omitempty"`
	Style       string   `json:"style,omitempty"`
	References  []string `json:"references,omitempty"`
	Count       int      `json:"count,omitempty"`
}
