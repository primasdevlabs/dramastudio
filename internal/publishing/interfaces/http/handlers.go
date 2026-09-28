package http

import (
	"net/http"
	"time"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/publishing/application/services"
	"dramastudio/internal/publishing/domain"
)

type PublishingHandler struct {
	service *services.PublishingService
}

func NewPublishingHandler(service *services.PublishingService) *PublishingHandler {
	return &PublishingHandler{service: service}
}

func (h *PublishingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/publishing/channels", h.listChannels)
	mux.HandleFunc("POST /v1/projects/{projectId}/publishing/channels", h.registerChannel)
	mux.HandleFunc("GET /v1/projects/{projectId}/publishing/publications", h.listPublications)
	mux.HandleFunc("POST /v1/projects/{projectId}/publishing/publications", h.schedulePublication)
	mux.HandleFunc("GET /v1/projects/{projectId}/publishing/publications/{publicationId}", h.getPublication)
	mux.HandleFunc("POST /v1/projects/{projectId}/publishing/publications/{publicationId}/publish", h.publish)
}

type registerChannelReq struct {
	Platform   domain.ChannelType     `json:"platform"`
	Name       string                 `json:"name"`
	AccountRef string                 `json:"account_ref"`
	Config     map[string]interface{} `json:"config"`
}

func (r *registerChannelReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	switch r.Platform {
	case domain.ChannelYouTube, domain.ChannelTikTok, domain.ChannelInstagram:
	default:
		ve.Add("platform", "must be youtube|tiktok|instagram")
	}
	if r.Name == "" {
		ve.Add("name", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type scheduleReq struct {
	EpisodeID   string                 `json:"episode_id"`
	ChannelID   string                 `json:"channel_id"`
	VideoURL    string                 `json:"video_url"`
	Metadata    domain.PublishMetadata `json:"metadata"`
	ScheduledAt *time.Time             `json:"scheduled_at"`
}

func (r *scheduleReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.EpisodeID == "" {
		ve.Add("episode_id", "required")
	}
	if r.ChannelID == "" {
		ve.Add("channel_id", "required")
	}
	if r.VideoURL == "" {
		ve.Add("video_url", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *PublishingHandler) listChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.service.ListChannels(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": channels})
}

func (h *PublishingHandler) registerChannel(w http.ResponseWriter, r *http.Request) {
	var req registerChannelReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	c, err := h.service.RegisterChannel(r.Context(), r.PathValue("projectId"), req.Platform, req.Name, req.AccountRef, req.Config)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, c)
}

func (h *PublishingHandler) listPublications(w http.ResponseWriter, r *http.Request) {
	pubs, err := h.service.ListPublications(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": pubs})
}

func (h *PublishingHandler) schedulePublication(w http.ResponseWriter, r *http.Request) {
	var req scheduleReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	pub, err := h.service.SchedulePublication(r.Context(), r.PathValue("projectId"),
		req.EpisodeID, req.ChannelID, req.VideoURL, req.Metadata, req.ScheduledAt,
		platformhttp.IdempotencyKey(r))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, pub)
}

func (h *PublishingHandler) getPublication(w http.ResponseWriter, r *http.Request) {
	pub, err := h.service.GetPublication(r.Context(), r.PathValue("publicationId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "PUBLICATION_NOT_FOUND", "Publication not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, pub)
}

func (h *PublishingHandler) publish(w http.ResponseWriter, r *http.Request) {
	pub, err := h.service.Publish(r.Context(), r.PathValue("publicationId"))
	if err != nil {
		if err == domain.ErrInvalidTransition {
			platformhttp.WriteError(w, http.StatusConflict, "INVALID_TRANSITION", err.Error(), platformhttp.RequestIDFrom(r), nil)
			return
		}
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, pub)
}
