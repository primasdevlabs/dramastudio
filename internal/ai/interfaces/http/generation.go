package http

import (
	"encoding/json"
	"net/http"

	"dramastudio/internal/ai/application"
	"dramastudio/internal/ai/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type generateReq struct {
	Capability domain.AICapability `json:"capability"`
	Input      json.RawMessage     `json:"input"`
	Prompt     string              `json:"prompt"`
	Scope      string              `json:"scope"`
}

func (r *generateReq) Validate() error {
	if r.Capability == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "capability", Message: "required"}}}
	}
	return nil
}

func (h *AIHandler) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.ListGenerationJobs(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": jobs})
}

func (h *AIHandler) generate(w http.ResponseWriter, r *http.Request) {
	var req generateReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	scope := req.Scope
	if scope == "" {
		scope = "project:" + r.PathValue("projectId")
	}
	job, err := h.executor.Handle(r.Context(), application.ExecuteGenerationCommand{
		ProjectID:      r.PathValue("projectId"),
		Capability:     req.Capability,
		Scope:          scope,
		Input:          req.Input,
		Prompt:         req.Prompt,
		IdempotencyKey: platformhttp.IdempotencyKey(r),
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusAccepted, job)
}
