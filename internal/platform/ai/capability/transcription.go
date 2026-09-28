package capability

// Transcription covers speech-to-text for dailies, captions, QC.
const Transcription = "transcription"

type TranscriptionSpec struct {
	AudioURI string `json:"audio_uri"`
	Language string `json:"language,omitempty"`
	Diarize  bool   `json:"diarize,omitempty"`
}
