package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"dramastudio/internal/canon/application/services"
	platformHTTP "dramastudio/internal/platform/http"
)

type CanonHandler struct {
	service *services.CanonService
}

func NewCanonHandler(service *services.CanonService) *CanonHandler {
	return &CanonHandler{service: service}
}

type createFactReq struct {
	Subject    string `json:"subject"`
	Predicate  string `json:"predicate"`
	Object     string `json:"object"`
	Introduced string `json:"introduced"`
	ValidFrom  string `json:"valid_from"`
}

func (h *CanonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/canon")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "facts" {
		switch r.Method {
		case http.MethodGet:
			projectID := r.URL.Query().Get("project_id")
			facts, err := h.service.ListFacts(r.Context(), projectID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"facts": facts})
		case http.MethodPost:
			var req createFactReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			fact, err := h.service.CreateFact(r.Context(), req.Subject, req.Predicate, req.Object, req.Introduced, req.ValidFrom)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, fact)
		default:
			platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
		}
		return
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
