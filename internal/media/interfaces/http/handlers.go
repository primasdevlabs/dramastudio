package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"dramastudio/internal/media/application/services"
	"dramastudio/internal/media/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type MediaHandler struct {
	service        *services.MediaService
	webhookSecrets map[string][]byte // provider -> HMAC secret
}

// NewMediaHandler wires the handler. webhookSecrets maps provider names to
// the shared secret used to verify X-Provider-Signature (HMAC-SHA256 over
// the raw request body). A missing secret rejects callbacks for that
// provider rather than silently trusting them (§61).
func NewMediaHandler(service *services.MediaService, webhookSecrets map[string]string) *MediaHandler {
	secrets := make(map[string][]byte, len(webhookSecrets))
	for k, v := range webhookSecrets {
		secrets[k] = []byte(v)
	}
	return &MediaHandler{service: service, webhookSecrets: secrets}
}

func (h *MediaHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/media/assets", h.listAssets)
	mux.HandleFunc("POST /v1/projects/{projectId}/media/assets", h.generateAsset)
	mux.HandleFunc("GET /v1/projects/{projectId}/media/assets/{assetId}", h.getAsset)
	mux.HandleFunc("POST /v1/projects/{projectId}/media/assets/{assetId}/approve", h.approveAsset)
	mux.HandleFunc("POST /v1/projects/{projectId}/media/assets/{assetId}/reject", h.rejectAsset)
	mux.HandleFunc("POST /v1/projects/{projectId}/media/assets/{assetId}/archive", h.archiveAsset)
	mux.HandleFunc("POST /v1/projects/{projectId}/media/assets/{assetId}/regenerate", h.regenerateAsset)
	mux.HandleFunc("GET /v1/projects/{projectId}/media/assets/{assetId}/versions", h.listVersions)
	mux.HandleFunc("GET /v1/projects/{projectId}/media/jobs", h.listJobs)
	// Provider webhooks are intentionally NOT project-scoped (§61).
	mux.HandleFunc("POST /v1/webhooks/providers/{provider}", h.providerWebhook)
}

type generateAssetReq struct {
	Capability      string                 `json:"capability"`
	Type            domain.MediaType       `json:"type"`
	CharacterID     string                 `json:"character_id"`
	LocationID      string                 `json:"location_id"`
	EpisodeID       string                 `json:"episode_id"`
	SceneID         string                 `json:"scene_id"`
	ShotID          string                 `json:"shot_id"`
	Prompt          string                 `json:"prompt"`
	Spec            map[string]interface{} `json:"spec"`
	Parameters      map[string]interface{} `json:"parameters"`
	ReferenceAssets []string               `json:"reference_assets"`
	CallbackURL     string                 `json:"callback_url"`
}

func (r *generateAssetReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Type == "" {
		ve.Add("type", "required")
	}
	if r.Prompt == "" && len(r.Spec) == 0 {
		ve.Add("prompt", "prompt or spec required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *MediaHandler) listAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.service.ListAssets(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": assets})
}

func (h *MediaHandler) generateAsset(w http.ResponseWriter, r *http.Request) {
	var req generateAssetReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	capability := req.Capability
	if capability == "" {
		capability = string(req.Type) + "_generation"
	}
	job, asset, err := h.service.Generate(r.Context(), services.GenerateRequest{
		ProjectID:       r.PathValue("projectId"),
		Capability:      capability,
		MediaType:       req.Type,
		CharacterID:     req.CharacterID,
		LocationID:      req.LocationID,
		EpisodeID:       req.EpisodeID,
		SceneID:         req.SceneID,
		ShotID:          req.ShotID,
		Prompt:          req.Prompt,
		Spec:            req.Spec,
		Parameters:      req.Parameters,
		ReferenceAssets: req.ReferenceAssets,
		CallbackURL:     req.CallbackURL,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusAccepted, map[string]interface{}{"job": job, "asset": asset})
}

func (h *MediaHandler) getAsset(w http.ResponseWriter, r *http.Request) {
	a, err := h.service.GetAsset(r.Context(), r.PathValue("assetId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "ASSET_NOT_FOUND", "Asset not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, a)
}

func (h *MediaHandler) approveAsset(w http.ResponseWriter, r *http.Request) {
	h.assetTransition(w, r, h.service.ApproveAsset)
}

func (h *MediaHandler) rejectAsset(w http.ResponseWriter, r *http.Request) {
	h.assetTransition(w, r, h.service.RejectAsset)
}

func (h *MediaHandler) archiveAsset(w http.ResponseWriter, r *http.Request) {
	h.assetTransition(w, r, h.service.ArchiveAsset)
}

func (h *MediaHandler) assetTransition(w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, id string) (*domain.Asset, error)) {
	a, err := fn(r.Context(), r.PathValue("assetId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "ASSET_NOT_FOUND", "Asset not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, a)
}

func (h *MediaHandler) regenerateAsset(w http.ResponseWriter, r *http.Request) {
	job, err := h.service.Regenerate(r.Context(), r.PathValue("assetId"), r.URL.Query().Get("callback_url"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusAccepted, job)
}

func (h *MediaHandler) listVersions(w http.ResponseWriter, r *http.Request) {
	vs, err := h.service.ListVersions(r.Context(), r.PathValue("assetId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": vs})
}

func (h *MediaHandler) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.ListJobs(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": jobs})
}

// --- Provider webhook ---

type providerCallback struct {
	ProviderJobID string  `json:"provider_job_id"`
	Status        string  `json:"status"`
	OutputURL     string  `json:"output_url"`
	Error         string  `json:"error"`
	Cost          float64 `json:"cost"`
}

func (h *MediaHandler) providerWebhook(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	secret, ok := h.webhookSecrets[provider]
	if !ok || len(secret) == 0 {
		platformhttp.WriteError(w, http.StatusNotFound, "UNKNOWN_PROVIDER", "No webhook configured for provider", platformhttp.RequestIDFrom(r), nil)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		platformhttp.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Cannot read body", platformhttp.RequestIDFrom(r), nil)
		return
	}
	if !verifySignature(secret, body, r.Header.Get("X-Provider-Signature")) {
		platformhttp.WriteError(w, http.StatusUnauthorized, "INVALID_SIGNATURE", "Webhook signature verification failed", platformhttp.RequestIDFrom(r), nil)
		return
	}
	var cb providerCallback
	if err := json.Unmarshal(body, &cb); err != nil || cb.ProviderJobID == "" {
		platformhttp.WriteError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Malformed callback payload", platformhttp.RequestIDFrom(r), nil)
		return
	}
	job, err := h.service.HandleProviderCallback(r.Context(), cb.ProviderJobID, cb.Status, cb.OutputURL, cb.Error, cb.Cost)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "JOB_NOT_FOUND", "Generation job not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"job_id": job.ID, "status": job.Status})
}

// verifySignature checks hex(HMAC-SHA256(secret, body)) against the header.
func verifySignature(secret, body []byte, header string) bool {
	if header == "" {
		return false
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(header))
}
