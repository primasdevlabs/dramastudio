package voicegen

import "dramastudio/internal/platform/ai/provider"

func init() {
	// Generic hosted TTS adapter (ElevenLabs-style HTTP endpoints).
	f := func(cfg provider.Config) provider.Adapter {
		return provider.Adapter{Sync: New(cfg.Name, cfg.BaseURL, cfg.APIKey)}
	}
	provider.RegisterFactory("voicegen", f)
	provider.RegisterFactory("elevenlabs", f)
}
