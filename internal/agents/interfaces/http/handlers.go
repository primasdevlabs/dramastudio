package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"dramastudio/internal/agents/application/services"
	platformHTTP "dramastudio/internal/platform/http"
)

type AgentsHandler struct {
	service *services.AgentService
}

func NewAgentsHandler(service *services.AgentService) *AgentsHandler {
	return &AgentsHandler{service: service}
}

type directorStepReq struct {
	ProjectID string `json:"project_id"`
	EpisodeID string `json:"episode_id"`
}

func (h *AgentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/agents")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "decisions" {
		if r.Method == http.MethodGet {
			projectID := r.URL.Query().Get("project_id")
			decs, err := h.service.ListDecisions(r.Context(), projectID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"decisions": decs})
			return
		}
	} else if len(parts) >= 2 && parts[0] == "lead-director" && parts[1] == "step" {
		if r.Method == http.MethodPost {
			var req directorStepReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			dec, err := h.service.TriggerLeadDirectorStep(r.Context(), req.ProjectID, req.EpisodeID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "STEP_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, dec)
			return
		}
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
