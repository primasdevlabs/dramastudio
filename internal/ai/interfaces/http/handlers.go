package http

import (
	"net/http"

	"dramastudio/internal/ai/application"
	"dramastudio/internal/ai/domain"
)

type AIHandler struct {
	policy   *application.PolicyService
	admin    *application.ProviderAdminService
	executor *application.ExecuteGenerationHandler
	repo     domain.ModelRegistryRepository
}

func NewAIHandler(policy *application.PolicyService, admin *application.ProviderAdminService, executor *application.ExecuteGenerationHandler, repo domain.ModelRegistryRepository) *AIHandler {
	return &AIHandler{policy: policy, admin: admin, executor: executor, repo: repo}
}

func (h *AIHandler) RegisterRoutes(mux *http.ServeMux) {
	// Providers (adapter config + credentials + health)
	mux.HandleFunc("GET /v1/ai/providers", h.listProviders)
	mux.HandleFunc("POST /v1/ai/providers", h.createProvider)
	mux.HandleFunc("POST /v1/ai/providers/{providerId}/test", h.testProvider)
	mux.HandleFunc("POST /v1/ai/providers/{providerId}/discover", h.discoverModels)
	mux.HandleFunc("DELETE /v1/ai/providers/{providerId}", h.deleteProvider)

	// Models (registry data; models are rows, not code)
	mux.HandleFunc("GET /v1/ai/models", h.listModels)
	mux.HandleFunc("POST /v1/ai/models", h.createModel)
	mux.HandleFunc("DELETE /v1/ai/models/{modelId}", h.deleteModel)

	// Capabilities (task-level routable types for the settings UI)
	mux.HandleFunc("GET /v1/ai/capabilities", h.listCapabilities)

	// Policies (scoped capability → model bindings)
	mux.HandleFunc("PUT /v1/ai/policies", h.setPolicy)
	mux.HandleFunc("GET /v1/ai/policies", h.listPolicies)
	mux.HandleFunc("GET /v1/projects/{projectId}/ai/policies", h.listProjectPolicies)

	// Generation jobs
	mux.HandleFunc("GET /v1/projects/{projectId}/ai/jobs", h.listJobs)
	mux.HandleFunc("POST /v1/projects/{projectId}/ai/generate", h.generate)
}
