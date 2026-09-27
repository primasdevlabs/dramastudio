package http

import (
	"encoding/json"
	"net/http"
	"strings"

	platformHTTP "dramastudio/internal/platform/http"
	"dramastudio/internal/projects/application/services"
	"dramastudio/internal/projects/domain"
)

type ProjectsHandler struct {
	service *services.ProjectService
}

func NewProjectsHandler(service *services.ProjectService) *ProjectsHandler {
	return &ProjectsHandler{service: service}
}

type createProjectReq struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Genre       string                `json:"genre"`
	Language    string                `json:"language"`
	Mode        domain.ProductionMode `json:"mode"`
}

type saveBibleReq struct {
	Premise        string   `json:"premise"`
	Genre          string   `json:"genre"`
	Themes         []string `json:"themes"`
	WorldRules     []string `json:"world_rules"`
	NarrativeRules []string `json:"narrative_rules"`
}

func (h *ProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/projects")
	path = strings.TrimPrefix(path, "/")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodGet:
			h.listProjects(w, r)
		case http.MethodPost:
			h.createProject(w, r)
		default:
			platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
		}
		return
	}

	parts := strings.Split(path, "/")
	projectID := domain.ProjectID(parts[0])

	if len(parts) == 1 {
		if r.Method == http.MethodGet {
			h.getProject(w, r, projectID)
			return
		}
	} else if len(parts) == 2 && parts[1] == "bible" {
		switch r.Method {
		case http.MethodGet:
			h.getBible(w, r, projectID)
		case http.MethodPost:
			h.saveBible(w, r, projectID)
		default:
			platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
		}
		return
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}

func (h *ProjectsHandler) listProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.ListProjects(r.Context())
	if err != nil {
		platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"projects": projects})
}

func (h *ProjectsHandler) createProject(w http.ResponseWriter, r *http.Request) {
	var req createProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid JSON payload", "", nil)
		return
	}
	if req.Name == "" {
		platformHTTP.WriteError(w, http.StatusBadRequest, "MISSING_NAME", "Project name is required", "", nil)
		return
	}
	if req.Mode == "" {
		req.Mode = domain.ProductionModeMonitored
	}
	p, err := h.service.CreateProject(r.Context(), req.Name, req.Description, req.Genre, req.Language, req.Mode)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusCreated, p)
}

func (h *ProjectsHandler) getProject(w http.ResponseWriter, r *http.Request, id domain.ProjectID) {
	p, err := h.service.GetProject(r.Context(), id)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, p)
}

func (h *ProjectsHandler) getBible(w http.ResponseWriter, r *http.Request, id domain.ProjectID) {
	bible, err := h.service.GetLatestBible(r.Context(), id)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusNotFound, "BIBLE_NOT_FOUND", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, bible)
}

func (h *ProjectsHandler) saveBible(w http.ResponseWriter, r *http.Request, id domain.ProjectID) {
	var req saveBibleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid JSON payload", "", nil)
		return
	}
	bible, err := h.service.SaveSeriesBible(r.Context(), id, req.Premise, req.Genre, req.Themes, req.WorldRules, req.NarrativeRules)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusInternalServerError, "BIBLE_SAVE_FAILED", err.Error(), "", nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusCreated, bible)
}
