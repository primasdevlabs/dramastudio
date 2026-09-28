package http

import (
	"errors"
	"net/http"

	"dramastudio/internal/identity/application/services"
	"dramastudio/internal/identity/domain"
	"dramastudio/internal/identity/infrastructure/sso"
	platformHTTP "dramastudio/internal/platform/http"
	"dramastudio/internal/platform/security"
)

type IdentityHandler struct {
	service *services.IdentityService
	jwt     *security.JWTService
	sso     *sso.GoogleVerifier
	// exposeResetTokens returns plaintext reset tokens in responses — set
	// only in development where no mailer exists to deliver them.
	exposeResetTokens bool
}

func NewIdentityHandler(service *services.IdentityService, jwt *security.JWTService, ssoVerifier *sso.GoogleVerifier, exposeResetTokens bool) *IdentityHandler {
	return &IdentityHandler{service: service, jwt: jwt, sso: ssoVerifier, exposeResetTokens: exposeResetTokens}
}

func (h *IdentityHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/auth/register", h.register)
	mux.HandleFunc("POST /v1/auth/login", h.login)
	mux.HandleFunc("POST /v1/auth/sso", h.ssoLogin)
	mux.HandleFunc("POST /v1/auth/forgot-password", h.forgotPassword)
	mux.HandleFunc("POST /v1/auth/reset-password", h.resetPassword)
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
	if h.jwt == nil {
		platformHTTP.WriteError(w, http.StatusServiceUnavailable, "AUTH_NOT_CONFIGURED", "Token service is not configured (set JWT_SECRET)", platformHTTP.RequestIDFrom(r), nil)
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

type ssoReq struct {
	Provider string `json:"provider"`
	IDToken  string `json:"id_token"`
	OrgName  string `json:"org_name"`
}

func (r *ssoReq) Validate() error {
	ve := &platformHTTP.ValidationError{}
	if r.Provider != "google" {
		ve.Add("provider", "only google is supported")
	}
	if r.IDToken == "" {
		ve.Add("id_token", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ssoLogin verifies a provider-issued identity token, auto-provisions the
// account on first use, and returns a session token — same shape as login.
func (h *IdentityHandler) ssoLogin(w http.ResponseWriter, r *http.Request) {
	var req ssoReq
	if !platformHTTP.DecodeAndValidate(w, r, &req) {
		return
	}
	if h.sso == nil {
		platformHTTP.WriteError(w, http.StatusServiceUnavailable, "SSO_NOT_CONFIGURED", "Google sign-in is not configured (set GOOGLE_CLIENT_ID)", platformHTTP.RequestIDFrom(r), nil)
		return
	}
	if h.jwt == nil {
		platformHTTP.WriteError(w, http.StatusServiceUnavailable, "AUTH_NOT_CONFIGURED", "Token service is not configured (set JWT_SECRET)", platformHTTP.RequestIDFrom(r), nil)
		return
	}
	profile, err := h.sso.Verify(r.Context(), req.IDToken)
	if err != nil {
		platformHTTP.WriteError(w, http.StatusUnauthorized, "SSO_TOKEN_INVALID", "Could not verify the Google token", platformHTTP.RequestIDFrom(r), nil)
		return
	}
	name := profile.Name
	if name == "" {
		name = profile.Email
	}
	u, _, err := h.service.AuthenticateSSO(r.Context(), profile.Email, name, req.OrgName)
	if err != nil {
		platformHTTP.WriteErrorFrom(w, r, err)
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

type forgotPasswordReq struct {
	Email string `json:"email"`
}

func (r *forgotPasswordReq) Validate() error {
	ve := &platformHTTP.ValidationError{}
	if r.Email == "" {
		ve.Add("email", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

// forgotPassword never reveals whether the email exists (no account
// enumeration). In development the token is returned directly since no
// mailer exists to deliver it.
func (h *IdentityHandler) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordReq
	if !platformHTTP.DecodeAndValidate(w, r, &req) {
		return
	}
	token, err := h.service.RequestPasswordReset(r.Context(), req.Email)
	if err != nil {
		platformHTTP.WriteErrorFrom(w, r, err)
		return
	}
	resp := map[string]interface{}{"ok": true}
	if token != "" && h.exposeResetTokens {
		resp["reset_token"] = token
	}
	platformHTTP.WriteJSON(w, http.StatusOK, resp)
}

type resetPasswordReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (r *resetPasswordReq) Validate() error {
	ve := &platformHTTP.ValidationError{}
	if r.Token == "" {
		ve.Add("token", "required")
	}
	if len(r.Password) < 8 {
		ve.Add("password", "minimum 8 characters")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *IdentityHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordReq
	if !platformHTTP.DecodeAndValidate(w, r, &req) {
		return
	}
	if err := h.service.ResetPassword(r.Context(), req.Token, req.Password); err != nil {
		if errors.Is(err, domain.ErrResetTokenInvalid) {
			platformHTTP.WriteError(w, http.StatusBadRequest, "RESET_TOKEN_INVALID", "Reset token is invalid, expired, or already used", platformHTTP.RequestIDFrom(r), nil)
			return
		}
		platformHTTP.WriteErrorFrom(w, r, err)
		return
	}
	platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"ok": true})
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
	p, _ := security.PrincipalFrom(r.Context())
	u, err := h.service.GetUser(r.Context(), p.OrgID, id)
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
	p, _ := security.PrincipalFrom(r.Context())
	if err := h.service.RevokeCredential(r.Context(), p.OrgID, id); err != nil {
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
