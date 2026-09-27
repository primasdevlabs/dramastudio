package voice

import "context"

type VoiceProvider interface {
	GenerateVoice(ctx context.Context, text, voiceID string) (string, error)
}
