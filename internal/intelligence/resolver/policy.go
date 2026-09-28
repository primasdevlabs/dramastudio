package resolver

import (
	"context"
	"sort"

	"dramastudio/internal/intelligence/domain"
)

// layerRank returns precedence position (lower index = stronger layer).
func layerRank(l domain.PolicyLayer) int {
	order := []domain.PolicyLayer{
		domain.LayerSystem, domain.LayerStudio, domain.LayerProject,
		domain.LayerSeries, domain.LayerSeason, domain.LayerEpisode,
		domain.LayerTask,
	}
	for i, v := range order {
		if v == l {
			return i
		}
	}
	return len(order)
}

// ResolvePolicy merges all applicable layers of one policy into an
// EffectivePolicy. Layers apply strongest-first; a lower layer cannot
// override a requirement key the stronger layer marked protected — those
// overrides are dropped and recorded for auditability.
func (r *Registry) ResolvePolicy(ctx context.Context, policyID string, scope ScopeIDs) (*domain.EffectivePolicy, error) {
	layers := append([]*domain.Policy{}, r.cat.Policies[policyID]...)
	if r.layers != nil {
		stored, err := r.layers.ListPolicyLayers(ctx, policyID)
		if err == nil {
			layers = append(layers, stored...)
		}
	}
	// Keep layers whose scope matches: empty ScopeID = applies to whole
	// layer; otherwise must equal the concrete scope ID for that layer.
	applicable := make([]*domain.Policy, 0, len(layers))
	for _, p := range layers {
		if p.ScopeID == "" || p.ScopeID == scope[p.Layer] {
			applicable = append(applicable, p)
		}
	}
	sort.SliceStable(applicable, func(i, j int) bool {
		ri, rj := layerRank(applicable[i].Layer), layerRank(applicable[j].Layer)
		if ri != rj {
			return ri < rj
		}
		return applicable[i].Version < applicable[j].Version
	})

	return mergeLayers(policyID, applicable), nil
}

// mergeLayers folds strongest-first applicable layers into an
// EffectivePolicy, honoring each layer's protected keys.
func mergeLayers(policyID string, applicable []*domain.Policy) *domain.EffectivePolicy {
	eff := &domain.EffectivePolicy{
		PolicyID: policyID, Requirements: map[string]interface{}{},
		Trace: map[string]string{},
	}
	protected := map[string]bool{}
	for _, p := range applicable {
		mergeRequirements(eff, p, protected)
		markProtected(p, protected)
		eff.Constraints = append(eff.Constraints, p.Constraints...)
	}
	return eff
}

func mergeRequirements(eff *domain.EffectivePolicy, p *domain.Policy, protected map[string]bool) {
	for k, v := range p.Requirements {
		if protected[k] {
			eff.Dropped = append(eff.Dropped, string(p.Layer)+" "+k)
			continue
		}
		eff.Requirements[k] = v
		eff.Trace[k] = string(p.Layer)
	}
}

func markProtected(p *domain.Policy, protected map[string]bool) {
	for _, k := range p.Protected {
		if _, exists := p.Requirements[k]; exists {
			protected[k] = true
		}
	}
}
