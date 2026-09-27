package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"dramastudio/internal/continuity/application/services"
	"dramastudio/internal/continuity/domain"
	platformHTTP "dramastudio/internal/platform/http"
)

type ContinuityHandler struct {
	service *services.ContinuityService
}

func NewContinuityHandler(service *services.ContinuityService) *ContinuityHandler {
	return &ContinuityHandler{service: service}
}

type runCheckReq struct {
	ProjectID string          `json:"project_id"`
	EpisodeID string          `json:"episode_id"`
	SceneID   string          `json:"scene_id"`
	Category  string          `json:"category"`
	Severity  domain.Severity `json:"severity"`
	Entity    string          `json:"entity"`
	Expected  string          `json:"expected_state"`
	Actual    string          `json:"actual_state"`
	Cause     string          `json:"cause"`
	Evidence  string          `json:"evidence"`
}

type resolveReq struct {
	Resolution string `json:"resolution"`
}

func (h *ContinuityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/continuity")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "issues" {
		if len(parts) == 1 {
			if r.Method == http.MethodGet {
				projectID := r.URL.Query().Get("project_id")
				issues, err := h.service.ListIssues(r.Context(), projectID)
				if err != nil {
					platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
					return
				}
				platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"issues": issues})
				return
			}
		} else if len(parts) == 3 && parts[2] == "resolve" {
			if r.Method == http.MethodPost {
				issueID := parts[1]
				var req resolveReq
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
					return
				}
				if err := h.service.ResolveIssue(r.Context(), issueID, req.Resolution); err != nil {
					platformHTTP.WriteError(w, http.StatusInternalServerError, "RESOLVE_FAILED", err.Error(), "", nil)
					return
				}
				platformHTTP.WriteJSON(w, http.StatusOK, map[string]string{"status": "resolved"})
				return
			}
		}
	} else if parts[0] == "check" {
		if r.Method == http.MethodPost {
			var req runCheckReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			issue, err := h.service.RunCheck(r.Context(), req.ProjectID, req.EpisodeID, req.SceneID, req.Category, req.Severity, req.Entity, req.Expected, req.Actual, req.Cause, req.Evidence)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "CHECK_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, issue)
			return
		}
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
