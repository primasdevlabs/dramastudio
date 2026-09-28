package capability

// SpeechGeneration covers TTS / voice lines.
const SpeechGeneration = "speech_generation"

type SpeechSpec struct {
	Text     string  `json:"text"`
	VoiceID  string  `json:"voice_id,omitempty"`
	Language string  `json:"language,omitempty"`
	Speed    float64 `json:"speed,omitempty"`
	Emotion  string  `json:"emotion,omitempty"`
}
