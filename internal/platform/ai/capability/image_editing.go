package capability

// ImageEditing covers inpainting/outpainting, wardrobe swaps, retouches.
const ImageEditing = "image_editing"

type ImageEditSpec struct {
	SourceURI string  `json:"source_uri"`
	MaskURI   string  `json:"mask_uri,omitempty"`
	Prompt    string  `json:"prompt"`
	Strength  float64 `json:"strength,omitempty"`
}
