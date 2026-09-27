package routing

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"dramastudio/internal/platform/ai/contracts"
)

// ProviderAdapter wraps sync or async provider implementations behind one
// lookup surface.
type ProviderAdapter struct {
	Sync  contracts.SyncGenerator
	Async contracts.AsyncGenerator
}

// ModelResolution is the outcome of capability → policy → provider → model
// resolution (Backend.md §23).
type ModelResolution struct {
	Provider string
	Model    string
	Adapter  ProviderAdapter
	Policy   ResolutionPolicy
}

// ResolutionPolicy is the effective model policy after scope override.
type ResolutionPolicy struct {
	ProviderID string
	ModelID    string
	MaxCost    float64
}

// Registry is the persistence surface the resolver needs; implemented by
// the internal/ai model registry repository.
type Registry interface {
	ResolvePolicy(ctx context.Context, capability, scopeID string) (*ResolutionPolicy, error)
}

// Resolver maps a capability + scope to a concrete adapter and model. It
// never falls back silently: no configured policy means an error (§53).
type Resolver struct {
	registry  Registry
	providers map[string]ProviderAdapter
	mu        sync.RWMutex
}

func NewResolver(registry Registry) *Resolver {
	return &Resolver{registry: registry, providers: make(map[string]ProviderAdapter)}
}

// RegisterProvider wires a provider implementation under its registry name.
func (r *Resolver) RegisterProvider(name string, adapter ProviderAdapter) {
	r.mu.Lock()
	r.providers[name] = adapter
	r.mu.Unlock()
}

// Resolve returns the adapter+model for a capability at a scope
// ("project:<id>", "episode:<id>", "system", ...).
func (r *Resolver) Resolve(ctx context.Context, capability, scopeID string) (*ModelResolution, error) {
	policy, err := r.registry.ResolvePolicy(ctx, capability, scopeID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("no model policy for capability %q at scope %q (no silent fallback)", capability, scopeID)
	}
	r.mu.RLock()
	adapter, ok := r.providers[policy.ProviderID]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("provider %q not registered for capability %q", policy.ProviderID, capability)
	}
	return &ModelResolution{
		Provider: policy.ProviderID,
		Model:    policy.ModelID,
		Adapter:  adapter,
		Policy:   *policy,
	}, nil
}

// NormalizeScope builds scope ids used by the policy lookup chain, from
// most to least specific.
func ScopeChain(projectID, episodeID string) []string {
	chain := make([]string, 0, 4)
	if episodeID != "" {
		chain = append(chain, "episode:"+episodeID)
	}
	if projectID != "" {
		chain = append(chain, "project:"+projectID)
	}
	return append(chain, "system")
}

// JoinedScope flattens a scope chain for registry implementations that do
// their own precedence walking.
func JoinedScope(chain []string) string {
	return strings.Join(chain, ",")
}
