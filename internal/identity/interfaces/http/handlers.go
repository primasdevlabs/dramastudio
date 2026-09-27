package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"dramastudio/internal/identity/application/services"
	platformHTTP "dramastudio/internal/platform/http"
)

type IdentityHandler struct {
	service *services.IdentityService
}

func NewIdentityHandler(service *services.IdentityService) *IdentityHandler {
	return &IdentityHandler{service: service}
}

type createUserReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	OrgID string `json:"org_id"`
}

func (h *IdentityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/identity")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "users" {
		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				users, err := h.service.ListUsers(r.Context())
				if err != nil {
					platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
					return
				}
				platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"users": users})
			case http.MethodPost:
				var req createUserReq
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
					return
				}
				u, err := h.service.CreateUser(r.Context(), req.Email, req.Name, req.Role, req.OrgID)
				if err != nil {
					platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
					return
				}
				platformHTTP.WriteJSON(w, http.StatusCreated, u)
			default:
				platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
			}
			return
		} else if len(parts) == 2 {
			if r.Method == http.MethodGet {
				u, err := h.service.GetUser(r.Context(), parts[1])
				if err != nil {
					platformHTTP.WriteError(w, http.StatusNotFound, "USER_NOT_FOUND", err.Error(), "", nil)
					return
				}
				platformHTTP.WriteJSON(w, http.StatusOK, u)
				return
			}
		}
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
