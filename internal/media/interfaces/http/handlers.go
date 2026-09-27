package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"dramastudio/internal/media/application/services"
	"dramastudio/internal/media/domain"
	platformHTTP "dramastudio/internal/platform/http"
)

type MediaHandler struct {
	service *services.MediaService
}

func NewMediaHandler(service *services.MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

type generateAssetReq struct {
	ProjectID   string           `json:"project_id"`
	CharacterID string           `json:"character_id"`
	SceneID     string           `json:"scene_id"`
	ShotID      string           `json:"shot_id"`
	Prompt      string           `json:"prompt"`
	Provider    string           `json:"provider"`
	Type        domain.MediaType `json:"type"`
}

func (h *MediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/media")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "assets" {
		if r.Method == http.MethodGet {
			projectID := r.URL.Query().Get("project_id")
			assets, err := h.service.ListAssets(r.Context(), projectID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"assets": assets})
			return
		}
	} else if parts[0] == "generate" {
		if r.Method == http.MethodPost {
			var req generateAssetReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			asset, err := h.service.GenerateAsset(r.Context(), req.ProjectID, req.CharacterID, req.SceneID, req.ShotID, req.Prompt, req.Provider, req.Type)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "GENERATE_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, asset)
			return
		}
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
