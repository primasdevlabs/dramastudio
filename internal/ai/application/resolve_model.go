package application

import (
	"context"

	"dramastudio/internal/ai/domain"
	"dramastudio/internal/platform/ai/routing"

	"github.com/google/uuid"
)

// RoutingRegistry adapts the ai repository to platform/ai/routing.Registry,
// walking the scope chain (episode → project → system) until a policy is
// found, then expanding it into ordered candidates joined with model
// records. No silent fallback — if nothing resolves, nil is returned and
// the resolver raises an error (§53).
type RoutingRegistry struct {
	repo domain.ModelRegistryRepository
}

func NewRoutingRegistry(repo domain.ModelRegistryRepository) *RoutingRegistry {
	return &RoutingRegistry{repo: repo}
}

func (a *RoutingRegistry) ResolvePolicy(ctx context.Context, capability, scopeID string) (*routing.ResolvedPolicy, error) {
	for _, scope := range scopeChain(scopeID) {
		row, err := a.repo.ResolvePolicy(ctx, capability, scope)
		if err != nil {
			return nil, err
		}
		if row != nil {
			return a.buildPolicy(ctx, row)
		}
	}
	return nil, nil
}

// buildPolicy expands a policy row into provider+model candidates.
// Targets that fail to resolve to an active, allowed model are skipped
// rather than silently attempted.
func (a *RoutingRegistry) buildPolicy(ctx context.Context, row *domain.PolicyRow) (*routing.ResolvedPolicy, error) {
	targets := append([]domain.ModelTarget{{ProviderID: row.ProviderID, ModelID: row.ModelID}}, row.FallbackModels...)
	candidates := make([]routing.Candidate, 0, len(targets))
	for _, t := range targets {
		m, err := a.repo.FindModelByID(ctx, domain.ModelID(t.ModelID))
		if err != nil || !m.IsActive() {
			continue
		}
		if len(row.AllowedModels) > 0 && !containsString(row.AllowedModels, string(m.ID)) {
			continue
		}
		candidates = append(candidates, routing.Candidate{
			ProviderID:      t.ProviderID,
			ModelID:         string(m.ID),
			ModelIdentifier: m.Identifier,
			ModelVersion:    m.Version,
			CostPerUnit:     m.Pricing.PerUnit,
		})
	}
	return &routing.ResolvedPolicy{
		Strategy:   row.Strategy(),
		MaxCost:    row.MaxCost,
		Candidates: candidates,
	}, nil
}

func containsString(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

// scopeChain expands "episode:e" → [episode:e, system] etc. It preserves
// prefixes that already carry a scope, and always ends at system.
func scopeChain(scopeID string) []string {
	if scopeID == "" || scopeID == "system" {
		return []string{"system"}
	}
	return []string{scopeID, "system"}
}

// PolicyService manages scoped model policies.
type PolicyService struct {
	repo domain.ModelRegistryRepository
}

func NewPolicyService(repo domain.ModelRegistryRepository) *PolicyService {
	return &PolicyService{repo: repo}
}

// SetPolicy upserts a capability policy at a scope. The caller builds the
// full row (primary, fallbacks, constraints, strategy).
func (s *PolicyService) SetPolicy(ctx context.Context, row *domain.PolicyRow) (*domain.PolicyRow, error) {
	if row.ID == "" {
		row.ID = "pol_" + uuid.NewString()
	}
	if err := s.repo.UpsertPolicy(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *PolicyService) ListPolicies(ctx context.Context, scope, scopeID string) ([]*domain.PolicyRow, error) {
	return s.repo.ListPolicies(ctx, scope, scopeID)
}

func (s *PolicyService) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	return s.repo.ListProviders(ctx)
}

func (s *PolicyService) ListModels(ctx context.Context, providerID domain.ProviderID) ([]*domain.Model, error) {
	return s.repo.ListModels(ctx, providerID)
}

func (s *PolicyService) ListModelsByCapability(ctx context.Context, capability domain.AICapability) ([]*domain.Model, error) {
	return s.repo.ListModelsByCapability(ctx, capability)
}
