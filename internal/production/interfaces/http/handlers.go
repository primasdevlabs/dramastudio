package http

import (
	"encoding/json"
	"net/http"
	"strings"

	platformHTTP "dramastudio/internal/platform/http"
	"dramastudio/internal/production/application/services"
	"dramastudio/internal/production/domain"
)

type ProductionHandler struct {
	service *services.ProductionService
}

func NewProductionHandler(service *services.ProductionService) *ProductionHandler {
	return &ProductionHandler{service: service}
}

type startJobReq struct {
	ShotID string `json:"shot_id"`
}

type approvalReq struct {
	ProjectID string                  `json:"project_id"`
	EpisodeID string                  `json:"episode_id"`
	Stage     string                  `json:"stage"`
	TargetID  string                  `json:"target_id"`
	Decision  domain.ApprovalDecision `json:"decision"`
	Notes     string                  `json:"notes"`
	DecidedBy string                  `json:"decided_by"`
}

func (h *ProductionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/production")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "jobs" {
		if r.Method == http.MethodGet {
			projectID := r.URL.Query().Get("project_id")
			jobs, err := h.service.ListJobs(r.Context(), projectID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"jobs": jobs})
			return
		}
	} else if parts[0] == "start" {
		if r.Method == http.MethodPost {
			var req startJobReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			job, err := h.service.StartProductionJob(r.Context(), req.ShotID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "START_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, job)
			return
		}
	} else if parts[0] == "approval" {
		switch r.Method {
		case http.MethodGet:
			projectID := r.URL.Query().Get("project_id")
			apps, err := h.service.ListApprovals(r.Context(), projectID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"approvals": apps})
		case http.MethodPost:
			var req approvalReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			app, err := h.service.SubmitApproval(r.Context(), req.ProjectID, req.EpisodeID, req.Stage, req.TargetID, req.Decision, req.Notes, req.DecidedBy)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "APPROVAL_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, app)
		default:
			platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
		}
		return
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
