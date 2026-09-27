package http

import (
	"net/http"

	"dramastudio/internal/identity/application/services"
	"dramastudio/internal/identity/domain"
	platformHTTP "dramastudio/internal/platform/http"
	"dramastudio/internal/platform/security"
)

type IdentityHandler struct {
	service *services.IdentityService
	jwt     *security.JWTService
}

func NewIdentityHandler(service *services.IdentityService, jwt *security.JWTService) *IdentityHandler {
	return &IdentityHandler{service: service, jwt: jwt}
}

func (h *IdentityHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/auth/register", h.register)
	mux.HandleFunc("POST /v1/auth/login", h.login)
	mux.HandleFunc("GET /v1/identity/me", h.me)
	mux.HandleFunc("GET /v1/identity/users", h.listUsers)
	mux.HandleFunc("GET /v1/identity/users/{userId}", h.getUser)
	mux.HandleFunc("POST /v1/identity/api-keys", h.createAPIKey)
	mux.HandleFunc("GET /v1/identity/api-keys", h.listAPIKeys)
	mux.HandleFunc("DELETE /v1/identity/api-keys/{keyId}", h.revokeAPIKey)
}

type registerReq struct {
	OrgName  string `json:"org_name"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *registerReq) Validate() error {
	ve := &platformHTTP.ValidationError{}
	if r.OrgName == "" {
		ve.Add("org_name", "required")
	}
	if r.Email == "" {
		ve.Add("email", "required")
	}
	if r.Name == "" {
		ve.Add("name", "required")
	}
	if len(r.Password) < 8 {
		ve.Add("password", "minimum 8 characters")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *IdentityHandler) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if !platformHTTP.DecodeAndValidate(w, r, &req) {
		return
	}
	u, org, err := h.service.Register(r.Context(), req.OrgName, req.Email, req.Name, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		code := "REGISTER_FAILED"
		if err == domain.ErrEmailTaken {
			status = http.StatusConflict
			code = "EMAIL_TAKEN"
		}
		platformHTTP.WriteError(w, status, code, err.Error(), platformHTTP.RequestIDFrom(r), nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"user":         u,
		"organization": org,
	})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *loginReq) Validate() error {
	ve := &platformHTTP.ValidationError{}
	if r.Email == "" {
		ve.Add("email", "required")
	}
	if r.Password == "" {
		ve.Add("password", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *IdentityHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if !platformHTTP.DecodeAndValidate(w, r, &req) {
		return
	}
	u, err := h.service.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", platformHTTP.RequestIDFrom(r), nil)
		return
	}
	token, err := h.jwt.Issue(security.Principal{
		UserID:      u.ID,
		OrgID:       u.OrgID,
		Roles:       []string{u.Role},
		Permissions: rolePermissions(u.Role),
	})
	if err != nil {
		platformHTTP.WriteError(w, http.StatusInternalServerError, "TOKEN_ISSUE_FAILED", "Could not issue token", platformHTTP.RequestIDFrom(r), nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"user":         u,
	})
}

func (h *IdentityHandler) me(w http.ResponseWriter, r *http.Request) {
	p, _ := security.PrincipalFrom(r.Context())
	platformHTTP.WriteJSON(w, http.StatusOK, p)
}

func (h *IdentityHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	p, _ := security.PrincipalFrom(r.Context())
	users, err := h.service.ListUsers(r.Context(), p.OrgID)
	if err != nil {
		platformHTTP.WriteErrorFrom(w, r, err)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": users})
}

func (h *IdentityHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := platformHTTP.RequirePathValue(w, r, "userId")
	u, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found", platformHTTP.RequestIDFrom(r), nil)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, u)
}

type createKeyReq struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Service     bool     `json:"service"`
}

func (r *createKeyReq) Validate() error {
	if r.Name == "" {
		return &platformHTTP.ValidationError{Fields: []platformHTTP.FieldError{{Field: "name", Message: "required"}}}
	}
	return nil
}

func (h *IdentityHandler) createAPIKey(w http.ResponseWriter, r *http.Request) {
	var req createKeyReq
	if !platformHTTP.DecodeAndValidate(w, r, &req) {
		return
	}
	p, _ := security.PrincipalFrom(r.Context())
	cred, plaintext, err := h.service.CreateAPIKey(r.Context(), p.OrgID, p.UserID, req.Name, req.Permissions, req.Service)
	if err != nil {
		platformHTTP.WriteErrorFrom(w, r, err)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"credential": cred,
		"api_key":    plaintext, // shown once; only the digest is stored
	})
}

func (h *IdentityHandler) listAPIKeys(w http.ResponseWriter, r *http.Request) {
	p, _ := security.PrincipalFrom(r.Context())
	creds, err := h.service.ListCredentials(r.Context(), p.OrgID)
	if err != nil {
		platformHTTP.WriteErrorFrom(w, r, err)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": creds})
}

func (h *IdentityHandler) revokeAPIKey(w http.ResponseWriter, r *http.Request) {
	id, _ := platformHTTP.RequirePathValue(w, r, "keyId")
	if err := h.service.RevokeCredential(r.Context(), id); err != nil {
		platformHTTP.WriteError(w, http.StatusNotFound, "CREDENTIAL_NOT_FOUND", "Credential not found", platformHTTP.RequestIDFrom(r), nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// rolePermissions maps studio roles to permission grants (§8).
func rolePermissions(role string) []string {
	switch role {
	case domain.RoleOwner:
		return []string{"*"}
	case domain.RoleProducer:
		return []string{"project.*", "production.*", "approval.*", "publish.*"}
	case domain.RoleDirector:
		return []string{"story.*", "production.*", "approval.read", "approval.write"}
	case domain.RoleWriter:
		return []string{"story.read", "story.write", "canon.read"}
	case domain.RoleEditor:
		return []string{"postproduction.*", "media.read"}
	case domain.RoleReviewer:
		return []string{"*.read", "approval.write", "continuity.read"}
	case domain.RoleOperator:
		return []string{"production.*", "jobs.*"}
	default: // viewer
		return []string{"*.read"}
	}
}
