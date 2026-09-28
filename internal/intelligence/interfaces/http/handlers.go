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
	"dramastudio/internal/platform/security"
)

type IntelligenceHandler struct {
	svc        *application.IntelligenceService
	projectOrg platformhttp.ProjectResolver // may be nil; disables org checks
}

// NewIntelligenceHandler wires the facade plus a project→org resolver so
// body-carried project_ids (execute/preview) get the same tenancy checks the
// ProjectGuard applies to path IDs.
func NewIntelligenceHandler(svc *application.IntelligenceService, projectOrg platformhttp.ProjectResolver) *IntelligenceHandler {
	return &IntelligenceHandler{svc: svc, projectOrg: projectOrg}
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

// verifyTaskProject enforces org ownership when the task carries a project
// ID — execute/preview are not under /v1/projects/{id} so the guard cannot
// see them.
func (h *IntelligenceHandler) verifyTaskProject(w http.ResponseWriter, r *http.Request, projectID string) bool {
	if projectID == "" || h.projectOrg == nil {
		return true
	}
	p, _ := security.PrincipalFrom(r.Context())
	orgID, err := h.projectOrg(r.Context(), projectID)
	if err != nil || (orgID != "" && orgID != p.OrgID) {
		platformhttp.WriteError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found", platformhttp.RequestIDFrom(r), nil)
		return false
	}
	return true
}

func (h *IntelligenceHandler) previewContext(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if !platformhttp.DecodeJSON(w, r, &req) {
		return
	}
	if !h.verifyTaskProject(w, r, req.Task.ProjectID) {
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
	if !h.verifyTaskProject(w, r, req.Task.ProjectID) {
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
	// Policy layers override AI routing for the whole org — owner-only.
	principal, _ := security.PrincipalFrom(r.Context())
	if !principal.HasPermission("*") {
		platformhttp.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Only org owners may edit policy layers", platformhttp.RequestIDFrom(r), nil)
		return
	}
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
