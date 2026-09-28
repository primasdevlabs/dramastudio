package http

import (
	"net/http"
	"time"

	"dramastudio/internal/ai/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type modelReq struct {
	ProviderID          string                `json:"provider_id"`
	Name                string                `json:"name"`
	Identifier          string                `json:"identifier"` // provider-side model id
	Capabilities        []domain.AICapability `json:"capabilities"`
	Modalities          domain.Modalities     `json:"modalities"`
	SupportedParameters []string              `json:"supported_parameters"`
	Limits              map[string]any        `json:"limits"`
	Pricing             domain.Pricing        `json:"pricing"`
	Version             string                `json:"version"`
	Status              domain.ModelStatus    `json:"status"`
	Metadata            map[string]string     `json:"metadata"`
}

func (r *modelReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.ProviderID == "" {
		ve.Add("provider_id", "required")
	}
	if r.Identifier == "" {
		ve.Add("identifier", "required (provider-side model id)")
	}
	if len(r.Capabilities) == 0 {
		ve.Add("capabilities", "at least one capability required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *AIHandler) listModels(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var (
		models []*domain.Model
		err    error
	)
	switch cap := q.Get("capability"); {
	case cap != "":
		models, err = h.policy.ListModelsByCapability(r.Context(), domain.AICapability(cap))
	default:
		models, err = h.policy.ListModels(r.Context(), domain.ProviderID(q.Get("provider_id")))
	}
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": models})
}

func (h *AIHandler) createModel(w http.ResponseWriter, r *http.Request) {
	var req modelReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	m := &domain.Model{
		ID:                  domain.ModelID(req.ProviderID + "." + req.Identifier),
		Name:                req.Name,
		ProviderID:          req.ProviderID,
		Identifier:          req.Identifier,
		Capabilities:        req.Capabilities,
		Modalities:          req.Modalities,
		SupportedParameters: req.SupportedParameters,
		Limits:              req.Limits,
		Pricing:             req.Pricing,
		Version:             req.Version,
		Status:              req.Status,
		Metadata:            req.Metadata,
		CreatedAt:           time.Now().UTC(),
	}
	if m.Name == "" {
		m.Name = req.Identifier
	}
	if err := h.repo.SaveModel(r.Context(), m); err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, m)
}

func (h *AIHandler) deleteModel(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteModel(r.Context(), domain.ModelID(r.PathValue("modelId"))); err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
