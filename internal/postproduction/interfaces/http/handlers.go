package http

import (
	"context"
	"net/http"
	"strconv"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/postproduction/application/services"
	"dramastudio/internal/postproduction/domain"
)

type PostproductionHandler struct {
	service         *services.PostproductionService
	episodeResolver func(ctx context.Context, episodeID string) (string, error)
}

func NewPostproductionHandler(service *services.PostproductionService) *PostproductionHandler {
	return &PostproductionHandler{service: service}
}

// SetEpisodeResolver injects the episode → project ownership check. When
// unset, episode-scoped params are not verified (tests only).
func (h *PostproductionHandler) SetEpisodeResolver(fn func(ctx context.Context, episodeID string) (string, error)) {
	h.episodeResolver = fn
}

// verifyEpisode confirms the episode_id query/body param belongs to the
// path project. Nested-resource IDs must resolve inside the route's
// hierarchy, not just exist.
func (h *PostproductionHandler) verifyEpisode(w http.ResponseWriter, r *http.Request, episodeID string) bool {
	if h.episodeResolver == nil || episodeID == "" {
		return true
	}
	owner, err := h.episodeResolver(r.Context(), episodeID)
	if err != nil || owner != r.PathValue("projectId") {
		platformhttp.WriteError(w, http.StatusNotFound, "EPISODE_NOT_FOUND", "Episode not found in project", platformhttp.RequestIDFrom(r), nil)
		return false
	}
	return true
}

// verifyRender confirms {renderId} belongs to the path project.
func (h *PostproductionHandler) verifyRender(w http.ResponseWriter, r *http.Request) bool {
	task, err := h.service.GetRender(r.Context(), r.PathValue("renderId"))
	if err != nil || task.ProjectID != r.PathValue("projectId") {
		platformhttp.WriteError(w, http.StatusNotFound, "RENDER_NOT_FOUND", "Render not found", platformhttp.RequestIDFrom(r), nil)
		return false
	}
	return true
}

func (h *PostproductionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/postproduction/timelines", h.getTimeline)
	mux.HandleFunc("POST /v1/projects/{projectId}/postproduction/timelines", h.saveTimeline)
	mux.HandleFunc("GET /v1/projects/{projectId}/postproduction/timelines/versions", h.listTimelineVersions)
	mux.HandleFunc("POST /v1/projects/{projectId}/postproduction/timelines/approve", h.approveTimeline)
	mux.HandleFunc("GET /v1/projects/{projectId}/postproduction/renders", h.listRenders)
	mux.HandleFunc("POST /v1/projects/{projectId}/postproduction/renders", h.queueRender)
	mux.HandleFunc("GET /v1/projects/{projectId}/postproduction/renders/{renderId}", h.getRender)
	mux.HandleFunc("POST /v1/projects/{projectId}/postproduction/renders/{renderId}/execute", h.executeRender)
}

type saveTimelineReq struct {
	EpisodeID   string             `json:"episode_id"`
	VideoTracks []domain.TrackItem `json:"video_tracks"`
	AudioTracks []domain.TrackItem `json:"audio_tracks"`
	Subtitles   []domain.Subtitle  `json:"subtitles"`
}

func (r *saveTimelineReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.EpisodeID == "" {
		ve.Add("episode_id", "required")
	}
	if len(r.VideoTracks) == 0 {
		ve.Add("video_tracks", "at least one video track required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type queueRenderReq struct {
	EpisodeID  string `json:"episode_id"`
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
	Execute    bool   `json:"execute"` // run synchronously (dev/test) instead of via worker
}

func (r *queueRenderReq) Validate() error {
	if r.EpisodeID == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "episode_id", Message: "required"}}}
	}
	return nil
}

func (h *PostproductionHandler) getTimeline(w http.ResponseWriter, r *http.Request) {
	episodeID := r.URL.Query().Get("episode_id")
	if !h.verifyEpisode(w, r, episodeID) {
		return
	}
	if v := r.URL.Query().Get("version"); v != "" {
		// specific version lookup via versions list
		version, _ := strconv.Atoi(v)
		versions, err := h.service.ListTimelineVersions(r.Context(), episodeID)
		if err != nil {
			platformhttp.WriteError(w, http.StatusNotFound, "TIMELINE_NOT_FOUND", "Timeline not found", platformhttp.RequestIDFrom(r), nil)
			return
		}
		for _, t := range versions {
			if t.Version == version {
				platformhttp.WriteJSON(w, http.StatusOK, t)
				return
			}
		}
		platformhttp.WriteError(w, http.StatusNotFound, "TIMELINE_NOT_FOUND", "Timeline version not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	tl, err := h.service.GetTimeline(r.Context(), episodeID)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "TIMELINE_NOT_FOUND", "Timeline not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, tl)
}

func (h *PostproductionHandler) saveTimeline(w http.ResponseWriter, r *http.Request) {
	var req saveTimelineReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	if !h.verifyEpisode(w, r, req.EpisodeID) {
		return
	}
	tl, err := h.service.SaveTimeline(r.Context(), r.PathValue("projectId"), req.EpisodeID,
		req.VideoTracks, req.AudioTracks, req.Subtitles)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, tl)
}

func (h *PostproductionHandler) listTimelineVersions(w http.ResponseWriter, r *http.Request) {
	if !h.verifyEpisode(w, r, r.URL.Query().Get("episode_id")) {
		return
	}
	versions, err := h.service.ListTimelineVersions(r.Context(), r.URL.Query().Get("episode_id"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": versions})
}

func (h *PostproductionHandler) approveTimeline(w http.ResponseWriter, r *http.Request) {
	if !h.verifyEpisode(w, r, r.URL.Query().Get("episode_id")) {
		return
	}
	tl, err := h.service.ApproveTimeline(r.Context(), r.URL.Query().Get("episode_id"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "TIMELINE_NOT_FOUND", "Timeline not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, tl)
}

func (h *PostproductionHandler) listRenders(w http.ResponseWriter, r *http.Request) {
	episodeID := r.URL.Query().Get("episode_id")
	if !h.verifyEpisode(w, r, episodeID) {
		return
	}
	renders, err := h.service.ListRenders(r.Context(), episodeID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	// Defense in depth: rows were created under this path's project, but
	// never let a render carrying a foreign project_id leak back out.
	out := renders[:0]
	for _, t := range renders {
		if t.ProjectID == r.PathValue("projectId") {
			out = append(out, t)
		}
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": out})
}

func (h *PostproductionHandler) queueRender(w http.ResponseWriter, r *http.Request) {
	var req queueRenderReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	if !h.verifyEpisode(w, r, req.EpisodeID) {
		return
	}
	task, err := h.service.QueueRender(r.Context(), r.PathValue("projectId"), req.EpisodeID, req.Format, req.Resolution)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	if req.Execute {
		task, err = h.service.ExecuteRender(r.Context(), task.ID)
		if err != nil {
			platformhttp.WriteErrorFrom(w, r, err)
			return
		}
	}
	platformhttp.WriteJSON(w, http.StatusCreated, task)
}

func (h *PostproductionHandler) getRender(w http.ResponseWriter, r *http.Request) {
	task, err := h.service.GetRender(r.Context(), r.PathValue("renderId"))
	if err != nil || task.ProjectID != r.PathValue("projectId") {
		platformhttp.WriteError(w, http.StatusNotFound, "RENDER_NOT_FOUND", "Render not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, task)
}

func (h *PostproductionHandler) executeRender(w http.ResponseWriter, r *http.Request) {
	if !h.verifyRender(w, r) {
		return
	}
	task, err := h.service.ExecuteRender(r.Context(), r.PathValue("renderId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, task)
}
