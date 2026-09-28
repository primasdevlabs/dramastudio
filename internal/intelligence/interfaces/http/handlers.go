// Package http exposes the intelligence subsystem: definitions for the
// settings/admin UI, execution-context preview, task execution, and the
// audit trail.
package http

import (
	"encoding/json"
	"net/http"

	"dramastudio/internal/intelligence/application"
	"dramastudio/internal/intelligence/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type IntelligenceHandler struct {
	svc *application.IntelligenceService
}

func NewIntelligenceHandler(svc *application.IntelligenceService) *IntelligenceHandler {
	return &IntelligenceHandler{svc: svc}
}

func (h *IntelligenceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/intelligence/agents", h.listAgents)
	mux.HandleFunc("GET /v1/intelligence/skills", h.listSkills)
	mux.HandleFunc("GET /v1/intelligence/policies", h.listPolicies)
	mux.HandleFunc("GET /v1/intelligence/rules", h.listRules)
	mux.HandleFunc("GET /v1/intelligence/evaluators", h.listEvaluators)
	mux.HandleFunc("POST /v1/intelligence/context", h.previewContext)
	mux.HandleFunc("POST /v1/intelligence/execute", h.execute)
	mux.HandleFunc("POST /v1/intelligence/policies/{id}/layers", h.savePolicyLayer)
	mux.HandleFunc("GET /v1/projects/{projectId}/intelligence/executions", h.listExecutions)
}

func (h *IntelligenceHandler) listAgents(w http.ResponseWriter, r *http.Request) {
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"agents": h.svc.Agents()})
}

func (h *IntelligenceHandler) listSkills(w http.ResponseWriter, r *http.Request) {
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"skills": h.svc.Skills()})
}

func (h *IntelligenceHandler) listPolicies(w http.ResponseWriter, r *http.Request) {
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"policies": h.svc.Policies()})
}

func (h *IntelligenceHandler) listRules(w http.ResponseWriter, r *http.Request) {
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"rules": h.svc.Rules()})
}

func (h *IntelligenceHandler) listEvaluators(w http.ResponseWriter, r *http.Request) {
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"evaluators": h.svc.Evaluators()})
}

func (h *IntelligenceHandler) listExecutions(w http.ResponseWriter, r *http.Request) {
	projectID, ok := platformhttp.RequirePathValue(w, r, "projectId")
	if !ok {
		return
	}
	recs, err := h.svc.ListExecutions(r.Context(), projectID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"executions": recs})
}

type taskRequest struct {
	AgentID string             `json:"agent_id"`
	Task    domain.TaskContext `json:"task"`
}

func (h *IntelligenceHandler) previewContext(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if !platformhttp.DecodeJSON(w, r, &req) {
		return
	}
	ec, err := h.svc.Preview(r.Context(), req.AgentID, req.Task)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, ec)
}

func (h *IntelligenceHandler) execute(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if !platformhttp.DecodeJSON(w, r, &req) {
		return
	}
	rec, res, err := h.svc.Execute(r.Context(), req.AgentID, req.Task)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"record": rec, "result": res,
	})
}

func (h *IntelligenceHandler) savePolicyLayer(w http.ResponseWriter, r *http.Request) {
	policyID, ok := platformhttp.RequirePathValue(w, r, "id")
	if !ok {
		return
	}
	var p domain.Policy
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		platformhttp.WriteErrorFrom(w, r, platformhttp.ErrBadRequest("invalid_body", "invalid JSON body"))
		return
	}
	p.ID = policyID
	if err := h.svc.SavePolicyLayer(r.Context(), &p); err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, p)
}
