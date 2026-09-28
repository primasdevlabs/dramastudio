package http

import (
	"net/http"
	"strconv"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/platform/security"
	"dramastudio/internal/projects/application/services"
	"dramastudio/internal/projects/domain"
)

type ProjectsHandler struct {
	service *services.ProjectService
}

func NewProjectsHandler(service *services.ProjectService) *ProjectsHandler {
	return &ProjectsHandler{service: service}
}

func (h *ProjectsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects", h.listProjects)
	mux.HandleFunc("POST /v1/projects", h.createProject)
	mux.HandleFunc("GET /v1/projects/{projectId}", h.getProject)
	mux.HandleFunc("PATCH /v1/projects/{projectId}", h.updateProject)
	mux.HandleFunc("POST /v1/projects/{projectId}/status", h.transitionStatus)
	mux.HandleFunc("GET /v1/projects/{projectId}/bible", h.getLatestBible)
	mux.HandleFunc("POST /v1/projects/{projectId}/bible", h.saveBible)
	mux.HandleFunc("GET /v1/projects/{projectId}/bible/versions", h.listBibleVersions)
	mux.HandleFunc("GET /v1/projects/{projectId}/bible/versions/{version}", h.getBibleVersion)
}

type createProjectReq struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Genre       string                  `json:"genre"`
	Language    string                  `json:"language"`
	Mode        domain.ProductionMode   `json:"mode"`
	Settings    domain.Settings         `json:"settings"`
	Policy      domain.ProductionPolicy `json:"policy"`
	Budget      domain.Budget           `json:"budget"`
}

func (r *createProjectReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Name == "" {
		ve.Add("name", "required")
	}
	if r.Mode != "" && r.Mode != domain.ModeMonitored && r.Mode != domain.ModeAutonomous {
		ve.Add("mode", "must be monitored|autonomous")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *ProjectsHandler) listProjects(w http.ResponseWriter, r *http.Request) {
	p, _ := security.PrincipalFrom(r.Context())
	projects, err := h.service.ListProjects(r.Context(), p.OrgID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": projects})
}

func (h *ProjectsHandler) createProject(w http.ResponseWriter, r *http.Request) {
	var req createProjectReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	p, _ := security.PrincipalFrom(r.Context())
	if req.Mode == "" {
		req.Mode = domain.ModeMonitored
	}
	proj, err := h.service.CreateProject(r.Context(), p.OrgID, req.Name, req.Description, req.Genre, req.Language, req.Mode)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	proj, err = h.service.UpdateProject(r.Context(), proj.ID, func(p *domain.Project) error {
		p.Settings = req.Settings
		p.Policy = req.Policy
		if req.Budget.TotalBudget > 0 || req.Budget.MaxCostPerGeneration > 0 {
			p.Budget = req.Budget
			if p.Budget.Currency == "" {
				p.Budget.Currency = "USD"
			}
		}
		return nil
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, proj)
}

func (h *ProjectsHandler) getProject(w http.ResponseWriter, r *http.Request) {
	id, _ := platformhttp.RequirePathValue(w, r, "projectId")
	p, err := h.service.GetProject(r.Context(), domain.ProjectID(id))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, p)
}

type updateProjectReq struct {
	Name        *string                  `json:"name"`
	Description *string                  `json:"description"`
	Genre       *string                  `json:"genre"`
	Language    *string                  `json:"language"`
	Mode        *domain.ProductionMode   `json:"mode"`
	Settings    *domain.Settings         `json:"settings"`
	Policy      *domain.ProductionPolicy `json:"policy"`
	Budget      *domain.Budget           `json:"budget"`
}

func (r *updateProjectReq) Validate() error {
	if r.Mode != nil && *r.Mode != domain.ModeMonitored && *r.Mode != domain.ModeAutonomous {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "mode", Message: "must be monitored|autonomous"}}}
	}
	return nil
}

func (h *ProjectsHandler) updateProject(w http.ResponseWriter, r *http.Request) {
	id, _ := platformhttp.RequirePathValue(w, r, "projectId")
	var req updateProjectReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	p, err := h.service.UpdateProject(r.Context(), domain.ProjectID(id), func(p *domain.Project) error {
		if req.Name != nil {
			p.Name = *req.Name
		}
		if req.Description != nil {
			p.Description = *req.Description
		}
		if req.Genre != nil {
			p.Genre = *req.Genre
		}
		if req.Language != nil {
			p.Language = *req.Language
		}
		if req.Mode != nil {
			p.Mode = *req.Mode
		}
		if req.Settings != nil {
			p.Settings = *req.Settings
		}
		if req.Policy != nil {
			p.Policy = *req.Policy
		}
		if req.Budget != nil {
			p.Budget = *req.Budget
		}
		return nil
	})
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, p)
}

type transitionReq struct {
	To domain.ProjectStatus `json:"to"`
}

func (r *transitionReq) Validate() error {
	if r.To == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "to", Message: "required"}}}
	}
	return nil
}

func (h *ProjectsHandler) transitionStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := platformhttp.RequirePathValue(w, r, "projectId")
	var req transitionReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	p, err := h.service.TransitionStatus(r.Context(), domain.ProjectID(id), req.To)
	if err != nil {
		if err == domain.ErrProjectNotFound {
			platformhttp.WriteError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found", platformhttp.RequestIDFrom(r), nil)
			return
		}
		platformhttp.WriteError(w, http.StatusConflict, "INVALID_TRANSITION", err.Error(), platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, p)
}

type saveBibleReq struct {
	Premise         string   `json:"premise"`
	Genre           string   `json:"genre"`
	Tone            string   `json:"tone"`
	Themes          []string `json:"themes"`
	WorldRules      []string `json:"world_rules"`
	NarrativeRules  []string `json:"narrative_rules"`
	VisualDirection string   `json:"visual_direction"`
	DialogueStyle   string   `json:"dialogue_style"`
}

func (r *saveBibleReq) Validate() error {
	if r.Premise == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "premise", Message: "required"}}}
	}
	return nil
}

func (h *ProjectsHandler) getLatestBible(w http.ResponseWriter, r *http.Request) {
	id, _ := platformhttp.RequirePathValue(w, r, "projectId")
	bible, err := h.service.GetLatestBible(r.Context(), domain.ProjectID(id))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "BIBLE_NOT_FOUND", "Series bible not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, bible)
}

func (h *ProjectsHandler) getBibleVersion(w http.ResponseWriter, r *http.Request) {
	id, _ := platformhttp.RequirePathValue(w, r, "projectId")
	v, err := strconv.Atoi(r.PathValue("version"))
	if err != nil || v < 1 {
		platformhttp.WriteError(w, http.StatusBadRequest, "INVALID_VERSION", "Version must be a positive integer", platformhttp.RequestIDFrom(r), nil)
		return
	}
	bible, err := h.service.GetBibleVersion(r.Context(), domain.ProjectID(id), v)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "BIBLE_NOT_FOUND", "Bible version not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, bible)
}

func (h *ProjectsHandler) listBibleVersions(w http.ResponseWriter, r *http.Request) {
	id, _ := platformhttp.RequirePathValue(w, r, "projectId")
	bibles, err := h.service.ListBibleVersions(r.Context(), domain.ProjectID(id))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": bibles})
}

func (h *ProjectsHandler) saveBible(w http.ResponseWriter, r *http.Request) {
	id, _ := platformhttp.RequirePathValue(w, r, "projectId")
	var req saveBibleReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	bible, err := h.service.SaveSeriesBible(r.Context(), domain.ProjectID(id),
		req.Premise, req.Genre, req.Tone, req.VisualDirection, req.DialogueStyle,
		req.Themes, req.WorldRules, req.NarrativeRules)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, bible)
}
