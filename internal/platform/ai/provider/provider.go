// Package provider is the adapter boundary for external AI services.
// Application code never references provider packages directly; the
// composition root builds Adapter instances from provider rows via
// registered factories, keyed by the row's adapter kind.
package provider

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"dramastudio/internal/platform/ai/capability"
)

// Config is the per-provider-instance settings row projected into adapter
// construction. APIKey is the resolved secret; APIKeyEnv names the env var
// that supplied it (for diagnostics, never the secret itself).
type Config struct {
	Name      string
	BaseURL   string
	APIKey    string
	APIKeyEnv string
	Timeout   time.Duration
}

// Adapter bundles the generation surfaces a provider implements. Either
// field may be nil when the provider only supports one execution style.
// Optional behaviors are discovered via Tester/Discoverer type asserts.
type Adapter struct {
	Sync  capability.SyncGenerator
	Async capability.AsyncGenerator
}

// Tester is implemented by adapters that can verify connectivity
// (provider "Test Connection" button).
type Tester interface {
	TestConnection(ctx context.Context) error
}

// DiscoveredModel is one model reported by a provider endpoint.
type DiscoveredModel struct {
	Identifier   string            `json:"identifier"` // provider-side model id
	Capabilities []string          `json:"capabilities,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Discoverer is implemented by adapters that can enumerate their model
// catalog (e.g. OpenAI-compatible GET /models).
type Discoverer interface {
	DiscoverModels(ctx context.Context) ([]DiscoveredModel, error)
}

// Factory builds an adapter for one provider instance.
type Factory func(cfg Config) Adapter

var (
	factoriesMu sync.RWMutex
	factories   = map[string]Factory{}
)

// RegisterFactory binds an adapter kind ("openai_compatible", "wan",
// "mock", ...) to its constructor. Adapter packages call it from init().
func RegisterFactory(kind string, f Factory) {
	factoriesMu.Lock()
	factories[kind] = f
	factoriesMu.Unlock()
}

// BuildAdapter constructs the adapter for a provider row. Unknown kinds
// fail loudly — a typo'd adapter kind is a config error, not a fallback.
func BuildAdapter(kind string, cfg Config) (Adapter, error) {
	factoriesMu.RLock()
	f, ok := factories[kind]
	factoriesMu.RUnlock()
	if !ok {
		return Adapter{}, fmt.Errorf("unknown adapter kind %q (registered: %v)", kind, KnownKinds())
	}
	return f(cfg), nil
}

// KnownKinds lists registered adapter kinds for diagnostics/UI.
func KnownKinds() []string {
	factoriesMu.RLock()
	out := make([]string, 0, len(factories))
	for k := range factories {
		out = append(out, k)
	}
	factoriesMu.RUnlock()
	sort.Strings(out)
	return out
}

// ResolveAPIKey resolves a provider credential: explicit env var name
// first, then the stored key (self-hosted installs where the operator
// entered it via the settings UI).
func ResolveAPIKey(envName, stored string) string {
	if envName != "" {
		if v := os.Getenv(envName); v != "" {
			return v
		}
	}
	return stored
}
