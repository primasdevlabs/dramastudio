package voice

import "context"

type VoiceProvider struct{}

func (p *VoiceProvider) GenerateVoice(ctx context.Context, text, voiceID string) (string, error) {
	return "https://media.dramastudio.internal/voice.wav", nil
}
