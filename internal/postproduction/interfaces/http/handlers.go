package http

import (
	"encoding/json"
	"net/http"
	"strings"

	platformHTTP "dramastudio/internal/platform/http"
	"dramastudio/internal/postproduction/application/services"
	"dramastudio/internal/postproduction/domain"
)

type PostproductionHandler struct {
	service *services.PostproductionService
}

func NewPostproductionHandler(service *services.PostproductionService) *PostproductionHandler {
	return &PostproductionHandler{service: service}
}

type createTimelineReq struct {
	ProjectID   string             `json:"project_id"`
	EpisodeID   string             `json:"episode_id"`
	VideoTracks []domain.TrackItem `json:"video_tracks"`
	AudioTracks []domain.TrackItem `json:"audio_tracks"`
}

type createRenderReq struct {
	EpisodeID  string `json:"episode_id"`
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
}

func (h *PostproductionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/postproduction")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if parts[0] == "timelines" {
		switch r.Method {
		case http.MethodGet:
			episodeID := r.URL.Query().Get("episode_id")
			tl, err := h.service.GetTimeline(r.Context(), episodeID)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusNotFound, "TIMELINE_NOT_FOUND", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusOK, tl)
		case http.MethodPost:
			var req createTimelineReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			tl, err := h.service.CreateTimeline(r.Context(), req.ProjectID, req.EpisodeID, req.VideoTracks, req.AudioTracks)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, tl)
		default:
			platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
		}
		return
	} else if parts[0] == "renders" {
		if r.Method == http.MethodPost {
			var req createRenderReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
				return
			}
			rnd, err := h.service.CreateRender(r.Context(), req.EpisodeID, req.Format, req.Resolution)
			if err != nil {
				platformHTTP.WriteError(w, http.StatusInternalServerError, "RENDER_FAILED", err.Error(), "", nil)
				return
			}
			platformHTTP.WriteJSON(w, http.StatusCreated, rnd)
			return
		}
	}

	platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
}
