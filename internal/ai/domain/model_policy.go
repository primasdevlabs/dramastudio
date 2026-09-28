package domain

import "time"

// Scope levels for model policies; the most specific scope wins (§53).
const (
	ScopeSystem       = "system"
	ScopeOrganization = "organization"
	ScopeProject      = "project"
	ScopeSeason       = "season"
	ScopeSeries       = "series"
	ScopeEpisode      = "episode"
	ScopeTask         = "task"
)

// Routing strategies decide how a policy's candidates are ordered.
const (
	StrategyPrimaryOnly      = "primary_only"      // primary or fail
	StrategyPrimaryFallback  = "primary_fallback"  // primary → fallbacks in order
	StrategyCostOptimized    = "cost_optimized"    // cheapest allowed first
	StrategyQualityOptimized = "quality_optimized" // preferred order, quality first
	StrategyAutomatic        = "automatic"         // primary → fallbacks (default)
)

// ModelTarget is a provider+model candidate inside a policy.
type ModelTarget struct {
	ProviderID string `json:"provider_id"`
	ModelID    string `json:"model_id"`
}

// PolicyRow is one capability→model binding at a scope: primary model,
// ordered fallbacks, constraints, and a routing strategy (§24, §53).
type PolicyRow struct {
	ID         string       `json:"id"`
	Scope      string       `json:"scope"`
	ScopeID    string       `json:"scope_id"` // "" for system scope
	Capability AICapability `json:"capability"`

	// Primary target. Kept as top-level fields for the common case.
	ProviderID string `json:"provider_id"`
	ModelID    string `json:"model_id"`

	// Fallback candidates tried in order when primary fails (retryable).
	FallbackModels []ModelTarget `json:"fallback_models,omitempty"`
	// AllowedModels restricts which registered models may serve this
	// capability (empty = no restriction).
	AllowedModels []string `json:"allowed_models,omitempty"`
	// PreferredModels orders preference for quality-optimized routing.
	PreferredModels []string `json:"preferred_models,omitempty"`

	RoutingStrategy    string  `json:"routing_strategy"`
	QualityRequirement string  `json:"quality_requirement,omitempty"` // draft|standard|high
	MaxCost            float64 `json:"max_cost"`
	MaxLatencyMs       int     `json:"max_latency_ms,omitempty"`
	Region             string  `json:"region,omitempty"`
	ConcurrencyLimit   int     `json:"concurrency_limit,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// Strategy returns the effective routing strategy (automatic default).
func (p *PolicyRow) Strategy() string {
	if p.RoutingStrategy == "" {
		return StrategyAutomatic
	}
	return p.RoutingStrategy
}
