package mock

import "dramastudio/internal/platform/ai/provider"

func init() {
	provider.RegisterFactory("mock", func(cfg provider.Config) provider.Adapter {
		return provider.Adapter{Sync: NewSync(cfg.Name), Async: NewAsync(cfg.Name)}
	})
}
