package http

import (
	"encoding/json"
	"net/http"
	"strings"

	platformHTTP "dramastudio/internal/platform/http"
	"dramastudio/internal/story/application/services"
)

type StoryHandler struct {
	service *services.StoryService
}

func NewStoryHandler(service *services.StoryService) *StoryHandler {
	return &StoryHandler{service: service}
}

type createSeasonReq struct {
	SeriesID string `json:"series_id"`
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
}

type createEpisodeReq struct {
	SeasonID string `json:"season_id"`
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
}

type createSceneReq struct {
	EpisodeID    string   `json:"episode_id"`
	Number       int      `json:"number"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	LocationID   string   `json:"location_id"`
	TimeOfDay    string   `json:"time_of_day"`
	CharacterIDs []string `json:"character_ids"`
}

func (h *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/story")
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
		return
	}

	switch parts[0] {
	case "seasons":
		h.handleSeasons(w, r)
	case "episodes":
		h.handleEpisodes(w, r, parts[1:])
	case "scenes":
		h.handleScenes(w, r)
	default:
		platformHTTP.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "", nil)
	}
}

func (h *StoryHandler) handleSeasons(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		seriesID := r.URL.Query().Get("series_id")
		seasons, err := h.service.ListSeasons(r.Context(), seriesID)
		if err != nil {
			platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
			return
		}
		platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"seasons": seasons})
	case http.MethodPost:
		var req createSeasonReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
			return
		}
		s, err := h.service.CreateSeason(r.Context(), req.SeriesID, req.Title, req.Summary, req.Number)
		if err != nil {
			platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
			return
		}
		platformHTTP.WriteJSON(w, http.StatusCreated, s)
	default:
		platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
	}
}

func (h *StoryHandler) handleEpisodes(w http.ResponseWriter, r *http.Request, rest []string) {
	if len(rest) > 0 && rest[0] != "" {
		epID := rest[0]
		ep, err := h.service.GetEpisode(r.Context(), epID)
		if err != nil {
			platformHTTP.WriteError(w, http.StatusNotFound, "EPISODE_NOT_FOUND", err.Error(), "", nil)
			return
		}
		platformHTTP.WriteJSON(w, http.StatusOK, ep)
		return
	}

	switch r.Method {
	case http.MethodGet:
		seasonID := r.URL.Query().Get("season_id")
		episodes, err := h.service.ListEpisodes(r.Context(), seasonID)
		if err != nil {
			platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
			return
		}
		platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"episodes": episodes})
	case http.MethodPost:
		var req createEpisodeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
			return
		}
		ep, err := h.service.CreateEpisode(r.Context(), req.SeasonID, req.Title, req.Summary, req.Number)
		if err != nil {
			platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
			return
		}
		platformHTTP.WriteJSON(w, http.StatusCreated, ep)
	default:
		platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
	}
}

func (h *StoryHandler) handleScenes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		episodeID := r.URL.Query().Get("episode_id")
		scenes, err := h.service.ListScenes(r.Context(), episodeID)
		if err != nil {
			platformHTTP.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), "", nil)
			return
		}
		platformHTTP.WriteJSON(w, http.StatusOK, map[string]interface{}{"scenes": scenes})
	case http.MethodPost:
		var req createSceneReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			platformHTTP.WriteError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid payload", "", nil)
			return
		}
		sc, err := h.service.CreateScene(r.Context(), req.EpisodeID, req.Title, req.Description, req.LocationID, req.TimeOfDay, req.Number, req.CharacterIDs)
		if err != nil {
			platformHTTP.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "", nil)
			return
		}
		platformHTTP.WriteJSON(w, http.StatusCreated, sc)
	default:
		platformHTTP.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "", nil)
	}
}
