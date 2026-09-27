package http

import (
	"encoding/json"
	"net/http"
	"strings"

	platformHTTP "dramastudio/internal/platform/http"
	"dramastudio/internal/publishing/application/services"
)

type PublishingHandler struct {
	service *services.PublishingService
}

func NewPublishingHandler(service *services.PublishingService) *PublishingHandler {
	return &PublishingHandler{service: service}
}

type publishReq struct {
	EpisodeID string   `json:"episode_id"`
	ChannelID string   `json:"channel_id"`
	Title     string   `json:"title"`
	Caption   string   `json:"caption"`
	Tags      []string `json:"tags"`
}

func (h *PublishingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/publishing")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "publications" {
		if r.Method == http.MethodGet {
			projectID := r.URL.Query().Get("project_id")
			pubs, err := h.service.ListPublications(r.Context(), projectID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"publications": pubs})
			return
		}
	} else if parts[0] == "publish" {
		if r.Method == http.MethodPost {
			var req publishReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			pub, err := h.service.CreatePublication(r.Context(), req.EpisodeID, req.ChannelID, req.Title, req.Caption, req.Tags)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "PUBLISH_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, pub)
			return
		}
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
