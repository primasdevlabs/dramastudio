package routing

import "sort"

// Routing strategy constants (mirror internal/ai/domain without importing
// it — routing is platform-level and must stay context-agnostic).
const (
	StrategyPrimaryOnly      = "primary_only"
	StrategyPrimaryFallback  = "primary_fallback"
	StrategyCostOptimized    = "cost_optimized"
	StrategyQualityOptimized = "quality_optimized"
	StrategyAutomatic        = "automatic"
)

// orderCandidates applies the policy's routing strategy to its candidate
// list. Candidates arrive in declared order: [primary, ...fallbacks].
func orderCandidates(strategy string, in []Candidate) []Candidate {
	out := make([]Candidate, len(in))
	copy(out, in)
	switch strategy {
	case StrategyPrimaryOnly:
		if len(out) > 1 {
			out = out[:1]
		}
	case StrategyCostOptimized:
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].CostPerUnit < out[j].CostPerUnit
		})
	case StrategyQualityOptimized:
		// Declared order already encodes preference; keep as-is.
	case StrategyPrimaryFallback, StrategyAutomatic:
		// Primary first, then declared fallbacks — the default.
	}
	return out
}
