package routing

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"dramastudio/internal/platform/ai/provider"
)

// Candidate is one provider+model target inside a resolved policy.
// Identifier is the provider-side model name sent to the API; ModelID is
// the registry row id used for provenance.
type Candidate struct {
	ProviderID      string
	ModelID         string
	ModelIdentifier string
	ModelVersion    string
	CostPerUnit     float64
}

// ResolvedPolicy is the effective model policy after scope override:
// a set of candidates plus the strategy that orders them.
type ResolvedPolicy struct {
	Strategy   string
	MaxCost    float64
	Candidates []Candidate
}

// BoundCandidate is a candidate with its live provider adapter attached.
type BoundCandidate struct {
	Candidate
	Adapter provider.Adapter
}

// ModelResolution is the outcome of capability → policy → provider → model
// resolution (Backend.md §23): an ordered candidate list honoring the
// policy's routing strategy.
type ModelResolution struct {
	Strategy   string
	MaxCost    float64
	Candidates []BoundCandidate
}

// Primary returns the first candidate, or nil when none resolved.
func (r *ModelResolution) Primary() *BoundCandidate {
	if r == nil || len(r.Candidates) == 0 {
		return nil
	}
	return &r.Candidates[0]
}

// Registry is the persistence surface the resolver needs; implemented by
// the internal/ai model registry repository.
type Registry interface {
	ResolvePolicy(ctx context.Context, capability, scopeID string) (*ResolvedPolicy, error)
}

// Resolver maps a capability + scope to concrete adapters and models. It
// never falls back silently: no configured policy means an error (§53).
type Resolver struct {
	registry  Registry
	providers map[string]provider.Adapter
	mu        sync.RWMutex
}

func NewResolver(registry Registry) *Resolver {
	return &Resolver{registry: registry, providers: make(map[string]provider.Adapter)}
}

// RegisterProvider wires a provider implementation under its registry id.
// Called at boot for catalog providers and at runtime when an operator
// adds a provider via the API.
func (r *Resolver) RegisterProvider(id string, adapter provider.Adapter) {
	r.mu.Lock()
	r.providers[id] = adapter
	r.mu.Unlock()
}

// UnregisterProvider removes a provider adapter (provider deleted).
func (r *Resolver) UnregisterProvider(id string) {
	r.mu.Lock()
	delete(r.providers, id)
	r.mu.Unlock()
}

// ProviderAdapter is the legacy name for provider.Adapter.
type ProviderAdapter = provider.Adapter

// AdapterFor returns the live adapter for a provider id (test/discover).
func (r *Resolver) AdapterFor(providerID string) (provider.Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.providers[providerID]
	return a, ok
}

// Resolve returns the ordered adapter+model candidates for a capability
// at a scope ("project:<id>", "episode:<id>", "system", ...).
func (r *Resolver) Resolve(ctx context.Context, capability, scopeID string) (*ModelResolution, error) {
	policy, err := r.registry.ResolvePolicy(ctx, capability, scopeID)
	if err != nil {
		return nil, err
	}
	if policy == nil || len(policy.Candidates) == 0 {
		return nil, fmt.Errorf("no model policy for capability %q at scope %q (no silent fallback)", capability, scopeID)
	}
	return r.bind(policy)
}

// bind attaches live adapters to ordered candidates. Candidates whose
// provider is not registered are skipped; if none bind, it's an error.
func (r *Resolver) bind(policy *ResolvedPolicy) (*ModelResolution, error) {
	ordered := orderCandidates(policy.Strategy, policy.Candidates)
	bound := make([]BoundCandidate, 0, len(ordered))
	r.mu.RLock()
	for _, c := range ordered {
		if a, ok := r.providers[c.ProviderID]; ok {
			bound = append(bound, BoundCandidate{Candidate: c, Adapter: a})
		}
	}
	r.mu.RUnlock()
	if len(bound) == 0 {
		return nil, fmt.Errorf("no registered provider adapter for policy candidates")
	}
	return &ModelResolution{
		Strategy:   policy.Strategy,
		MaxCost:    policy.MaxCost,
		Candidates: bound,
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
