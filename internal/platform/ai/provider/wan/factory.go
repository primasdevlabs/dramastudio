package wan

import "dramastudio/internal/platform/ai/provider"

func init() {
	// Native async video adapter (job submit/poll per §61).
	f := func(cfg provider.Config) provider.Adapter {
		return provider.Adapter{Async: New(cfg.Name, cfg.BaseURL, cfg.APIKey)}
	}
	provider.RegisterFactory("wan", f)
	provider.RegisterFactory("alibaba", f)
}
