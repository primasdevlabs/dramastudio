package http

import (
	"context"
	"net/http"

	"dramastudio/internal/agents/application/services"
	"dramastudio/internal/agents/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type AgentsHandler struct {
	service *services.AgentService
}

func NewAgentsHandler(service *services.AgentService) *AgentsHandler {
	return &AgentsHandler{service: service}
}

func (h *AgentsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/agents/tasks", h.listTasks)
	mux.HandleFunc("POST /v1/projects/{projectId}/agents/tasks", h.createTask)
	mux.HandleFunc("GET /v1/projects/{projectId}/agents/tasks/{taskId}", h.getTask)
	mux.HandleFunc("POST /v1/projects/{projectId}/agents/tasks/{taskId}/approve", h.approveTask)
	mux.HandleFunc("POST /v1/projects/{projectId}/agents/tasks/{taskId}/cancel", h.cancelTask)
	mux.HandleFunc("POST /v1/projects/{projectId}/agents/tasks/{taskId}/retry", h.retryTask)
	mux.HandleFunc("POST /v1/projects/{projectId}/agents/tasks/{taskId}/complete", h.completeTask)
	mux.HandleFunc("GET /v1/projects/{projectId}/agents/decisions", h.listDecisions)
	mux.HandleFunc("POST /v1/projects/{projectId}/agents/director/step", h.directorStep)
	mux.HandleFunc("GET /v1/agents/definitions", h.listDefinitions)
	mux.HandleFunc("POST /v1/agents/definitions", h.registerDefinition)
}

// --- requests ---

type createTaskReq struct {
	AgentID         string                 `json:"agent_id"`
	Objective       string                 `json:"objective"`
	Input           map[string]interface{} `json:"input"`
	ExpectedOutput  string                 `json:"expected_output"`
	Constraints     map[string]interface{} `json:"constraints"`
	Budget          *float64               `json:"budget"`
	MaxIterations   int                    `json:"max_iterations"`
	TimeoutSec      int                    `json:"timeout_sec"`
	ApprovalPolicy  domain.ApprovalPolicy  `json:"approval_policy"`
	SuccessCriteria string                 `json:"success_criteria"`
}

func (r *createTaskReq) Validate() error {
	if r.Objective == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "objective", Message: "required"}}}
	}
	return nil
}

type registerDefReq struct {
	Role         domain.AgentRole       `json:"role"`
	Name         string                 `json:"name"`
	Instructions string                 `json:"instructions"`
	Skills       []domain.Skill         `json:"skills"`
	Tools        []string               `json:"tools"`
	Permissions  []string               `json:"permissions"`
	ModelPolicy  map[string]interface{} `json:"model_policy"`
	BudgetLimit  *float64               `json:"budget_limit"`
}

func (r *registerDefReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Role == "" {
		ve.Add("role", "required")
	}
	if r.Name == "" {
		ve.Add("name", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type directorStepReq struct {
	EpisodeID string `json:"episode_id"`
}

func (r *directorStepReq) Validate() error { return nil }

type completeTaskReq struct {
	Result map[string]interface{} `json:"result"`
	Error  string                 `json:"error"`
}

func (r *completeTaskReq) Validate() error { return nil }

// --- handlers ---

func (h *AgentsHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks(r.Context(), r.PathValue("projectId"), domain.TaskStatus(r.URL.Query().Get("status")))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": tasks})
}

func (h *AgentsHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	t, err := h.service.CreateTask(r.Context(), &domain.Task{
		ProjectID:       r.PathValue("projectId"),
		AgentID:         req.AgentID,
		Objective:       req.Objective,
		Input:           req.Input,
		ExpectedOutput:  req.ExpectedOutput,
		Constraints:     req.Constraints,
		Budget:          req.Budget,
		MaxIterations:   req.MaxIterations,
		TimeoutSec:      req.TimeoutSec,
		ApprovalPolicy:  req.ApprovalPolicy,
		SuccessCriteria: req.SuccessCriteria,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, t)
}

// verifyTask confirms taskId belongs to the path project.
func (h *AgentsHandler) verifyTask(w http.ResponseWriter, r *http.Request) bool {
	t, err := h.service.GetTask(r.Context(), r.PathValue("taskId"))
	if err != nil || t.ProjectID != r.PathValue("projectId") {
		platformhttp.WriteError(w, http.StatusNotFound, "TASK_NOT_FOUND", "Agent task not found", platformhttp.RequestIDFrom(r), nil)
		return false
	}
	return true
}

func (h *AgentsHandler) getTask(w http.ResponseWriter, r *http.Request) {
	if !h.verifyTask(w, r) {
		return
	}
	t, _ := h.service.GetTask(r.Context(), r.PathValue("taskId"))
	platformhttp.WriteJSON(w, http.StatusOK, t)
}

func (h *AgentsHandler) taskTransition(w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, id string) (*domain.Task, error)) {
	if !h.verifyTask(w, r) {
		return
	}
	t, err := fn(r.Context(), r.PathValue("taskId"))
	if err != nil {
		if err == domain.ErrInvalidTransition {
			platformhttp.WriteError(w, http.StatusConflict, "INVALID_TRANSITION", err.Error(), platformhttp.RequestIDFrom(r), nil)
			return
		}
		platformhttp.WriteError(w, http.StatusNotFound, "TASK_NOT_FOUND", "Agent task not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, t)
}

func (h *AgentsHandler) approveTask(w http.ResponseWriter, r *http.Request) {
	h.taskTransition(w, r, h.service.ApproveTask)
}

func (h *AgentsHandler) cancelTask(w http.ResponseWriter, r *http.Request) {
	h.taskTransition(w, r, h.service.CancelTask)
}

func (h *AgentsHandler) retryTask(w http.ResponseWriter, r *http.Request) {
	h.taskTransition(w, r, h.service.RetryTask)
}

func (h *AgentsHandler) completeTask(w http.ResponseWriter, r *http.Request) {
	if !h.verifyTask(w, r) {
		return
	}
	var req completeTaskReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	var t *domain.Task
	var err error
	if req.Error != "" {
		t, err = h.service.FailTask(r.Context(), r.PathValue("taskId"), req.Error)
	} else {
		t, err = h.service.CompleteTask(r.Context(), r.PathValue("taskId"), req.Result)
	}
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, t)
}

func (h *AgentsHandler) listDecisions(w http.ResponseWriter, r *http.Request) {
	decs, err := h.service.ListDecisions(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": decs})
}

func (h *AgentsHandler) directorStep(w http.ResponseWriter, r *http.Request) {
	var req directorStepReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	dec, err := h.service.RunDirectorStep(r.Context(), r.PathValue("projectId"), req.EpisodeID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, dec)
}

func (h *AgentsHandler) listDefinitions(w http.ResponseWriter, r *http.Request) {
	defs, err := h.service.ListDefinitions(r.Context())
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": defs})
}

func (h *AgentsHandler) registerDefinition(w http.ResponseWriter, r *http.Request) {
	var req registerDefReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	d, err := h.service.RegisterDefinition(r.Context(), &domain.Definition{
		Role:         req.Role,
		Name:         req.Name,
		Instructions: req.Instructions,
		Skills:       req.Skills,
		Tools:        req.Tools,
		Permissions:  req.Permissions,
		ModelPolicy:  req.ModelPolicy,
		BudgetLimit:  req.BudgetLimit,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, d)
}
