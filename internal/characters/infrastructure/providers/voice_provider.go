package providers

type VoiceProvider interface {
	Synthesize(text string, voiceID string) ([]byte, error)
}
