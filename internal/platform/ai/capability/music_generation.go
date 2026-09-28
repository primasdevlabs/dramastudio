package capability

// MusicGeneration covers score and source music.
const MusicGeneration = "music_generation"

type MusicSpec struct {
	Prompt       string  `json:"prompt"`
	DurationSec  float64 `json:"duration_sec,omitempty"`
	Genre        string  `json:"genre,omitempty"`
	Mood         string  `json:"mood,omitempty"`
	Instrumental bool    `json:"instrumental,omitempty"`
}
