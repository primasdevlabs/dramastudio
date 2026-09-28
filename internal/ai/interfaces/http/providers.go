package http

import (
	"net/http"

	"dramastudio/internal/ai/application"
	"dramastudio/internal/ai/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type providerReq struct {
	Name      string `json:"name"`
	Type      string `json:"type"` // adapter kind: openai_compatible|anthropic|wan|...
	BaseURL   string `json:"base_url"`
	APIKeyEnv string `json:"api_key_env"`
	APIKey    string `json:"api_key"` // stored server-side only
}

func (r *providerReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Name == "" {
		ve.Add("name", "required")
	}
	if r.Type == "" {
		ve.Add("type", "required (adapter kind, e.g. openai_compatible)")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type providerView struct {
	*domain.Provider
	AdapterAvailable bool `json:"adapter_available"`
}

func (h *AIHandler) listProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.policy.ListProviders(r.Context())
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	items := make([]providerView, 0, len(providers))
	for _, p := range providers {
		items = append(items, providerView{Provider: p, AdapterAvailable: application.AdapterAvailable(p.Type)})
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

func (h *AIHandler) createProvider(w http.ResponseWriter, r *http.Request) {
	if !requireOwner(w, r) {
		return
	}
	var req providerReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	p := &domain.Provider{
		Name:      req.Name,
		Type:      req.Type,
		BaseURL:   req.BaseURL,
		APIKeyEnv: req.APIKeyEnv,
		APIKey:    req.APIKey,
	}
	created, err := h.admin.CreateProvider(r.Context(), p)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, providerView{Provider: created, AdapterAvailable: application.AdapterAvailable(created.Type)})
}

func (h *AIHandler) testProvider(w http.ResponseWriter, r *http.Request) {
	if !requireOwner(w, r) {
		return
	}
	if err := h.admin.TestProvider(r.Context(), r.PathValue("providerId")); err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *AIHandler) discoverModels(w http.ResponseWriter, r *http.Request) {
	if !requireOwner(w, r) {
		return
	}
	models, err := h.admin.DiscoverModels(r.Context(), r.PathValue("providerId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": models, "count": len(models)})
}

func (h *AIHandler) deleteProvider(w http.ResponseWriter, r *http.Request) {
	if !requireOwner(w, r) {
		return
	}
	if err := h.admin.DeleteProvider(r.Context(), r.PathValue("providerId")); err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
