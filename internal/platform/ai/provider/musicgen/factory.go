package musicgen

import "dramastudio/internal/platform/ai/provider"

func init() {
	// Generic hosted music-generation adapter.
	f := func(cfg provider.Config) provider.Adapter {
		return provider.Adapter{Sync: New(cfg.Name, cfg.BaseURL, cfg.APIKey)}
	}
	provider.RegisterFactory("musicgen", f)
	provider.RegisterFactory("suno", f)
}
