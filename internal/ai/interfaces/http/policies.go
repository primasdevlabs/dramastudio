package http

import (
	"net/http"

	"dramastudio/internal/ai/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type policyReq struct {
	Scope      string              `json:"scope"`    // system|organization|project|series|season|episode|task
	ScopeID    string              `json:"scope_id"` // "" for system
	Capability domain.AICapability `json:"capability"`

	ProviderID      string               `json:"provider_id"` // primary target
	ModelID         string               `json:"model_id"`
	FallbackModels  []domain.ModelTarget `json:"fallback_models"`
	AllowedModels   []string             `json:"allowed_models"`
	PreferredModels []string             `json:"preferred_models"`

	RoutingStrategy    string  `json:"routing_strategy"` // primary_only|primary_fallback|cost_optimized|quality_optimized|automatic
	QualityRequirement string  `json:"quality_requirement"`
	MaxCost            float64 `json:"max_cost"`
	MaxLatencyMs       int     `json:"max_latency_ms"`
	Region             string  `json:"region"`
	ConcurrencyLimit   int     `json:"concurrency_limit"`
}

func (r *policyReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	switch r.Scope {
	case domain.ScopeSystem, domain.ScopeOrganization, domain.ScopeProject,
		domain.ScopeSeries, domain.ScopeSeason, domain.ScopeEpisode, domain.ScopeTask:
	default:
		ve.Add("scope", "invalid scope")
	}
	if r.Capability == "" {
		ve.Add("capability", "required")
	}
	if r.ProviderID == "" {
		ve.Add("provider_id", "required (primary provider)")
	}
	if r.ModelID == "" {
		ve.Add("model_id", "required (primary model)")
	}
	switch r.RoutingStrategy {
	case "", domain.StrategyPrimaryOnly, domain.StrategyPrimaryFallback,
		domain.StrategyCostOptimized, domain.StrategyQualityOptimized,
		domain.StrategyAutomatic:
	default:
		ve.Add("routing_strategy", "invalid strategy")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *AIHandler) setPolicy(w http.ResponseWriter, r *http.Request) {
	var req policyReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	row := &domain.PolicyRow{
		Scope:              req.Scope,
		ScopeID:            req.ScopeID,
		Capability:         req.Capability,
		ProviderID:         req.ProviderID,
		ModelID:            req.ModelID,
		FallbackModels:     req.FallbackModels,
		AllowedModels:      req.AllowedModels,
		PreferredModels:    req.PreferredModels,
		RoutingStrategy:    req.RoutingStrategy,
		QualityRequirement: req.QualityRequirement,
		MaxCost:            req.MaxCost,
		MaxLatencyMs:       req.MaxLatencyMs,
		Region:             req.Region,
		ConcurrencyLimit:   req.ConcurrencyLimit,
	}
	out, err := h.policy.SetPolicy(r.Context(), row)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, out)
}

func (h *AIHandler) listPolicies(w http.ResponseWriter, r *http.Request) {
	scope := r.URL.Query().Get("scope")
	if scope == "" {
		scope = domain.ScopeSystem
	}
	rows, err := h.policy.ListPolicies(r.Context(), scope, r.URL.Query().Get("scope_id"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": rows})
}

func (h *AIHandler) listProjectPolicies(w http.ResponseWriter, r *http.Request) {
	rows, err := h.policy.ListPolicies(r.Context(), domain.ScopeProject, r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": rows})
}

func (h *AIHandler) listCapabilities(w http.ResponseWriter, r *http.Request) {
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": domain.AllCapabilities})
}
