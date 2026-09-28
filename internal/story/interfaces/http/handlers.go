package http

import (
	"net/http"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/story/application/services"
	"dramastudio/internal/story/domain"
)

type StoryHandler struct {
	service *services.StoryService
}

func NewStoryHandler(service *services.StoryService) *StoryHandler {
	return &StoryHandler{service: service}
}

func (h *StoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/story", h.getStory)
	mux.HandleFunc("PUT /v1/projects/{projectId}/story", h.ensureSeries)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons", h.listSeasons)
	mux.HandleFunc("POST /v1/projects/{projectId}/seasons", h.createSeason)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons/{seasonId}", h.getSeason)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons/{seasonId}/arcs", h.listArcs)
	mux.HandleFunc("POST /v1/projects/{projectId}/seasons/{seasonId}/arcs", h.createArc)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons/{seasonId}/episodes", h.listEpisodes)
	mux.HandleFunc("POST /v1/projects/{projectId}/seasons/{seasonId}/episodes", h.createEpisode)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons/{seasonId}/episodes/{episodeId}", h.getEpisode)
	mux.HandleFunc("PATCH /v1/projects/{projectId}/seasons/{seasonId}/episodes/{episodeId}", h.updateEpisode)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons/{seasonId}/episodes/{episodeId}/scenes", h.listScenes)
	mux.HandleFunc("POST /v1/projects/{projectId}/seasons/{seasonId}/episodes/{episodeId}/scenes", h.createScene)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons/{seasonId}/episodes/{episodeId}/scenes/{sceneId}", h.getScene)
	mux.HandleFunc("GET /v1/projects/{projectId}/seasons/{seasonId}/episodes/{episodeId}/scenes/{sceneId}/beats", h.listBeats)
	mux.HandleFunc("POST /v1/projects/{projectId}/seasons/{seasonId}/episodes/{episodeId}/scenes/{sceneId}/beats", h.createBeat)
	mux.HandleFunc("GET /v1/projects/{projectId}/story/graph", h.getGraph)
	mux.HandleFunc("POST /v1/projects/{projectId}/story/graph/nodes", h.addNode)
	mux.HandleFunc("POST /v1/projects/{projectId}/story/graph/edges", h.addEdge)
	mux.HandleFunc("GET /v1/projects/{projectId}/story/plot-threads", h.listThreads)
	mux.HandleFunc("POST /v1/projects/{projectId}/story/plot-threads", h.createThread)
}

// verifySeason confirms seasonID lives under projectID (season→series).
// Prevents hierarchy-confusion access where a valid projectId in the path is
// combined with nested IDs from another project/org.
func (h *StoryHandler) verifySeason(w http.ResponseWriter, r *http.Request, projectID, seasonID string) (*domain.Season, bool) {
	se, err := h.service.GetSeason(r.Context(), seasonID)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "SEASON_NOT_FOUND", "Season not found", platformhttp.RequestIDFrom(r), nil)
		return nil, false
	}
	ser, err := h.service.GetSeries(r.Context(), se.SeriesID)
	if err != nil || ser.ProjectID != projectID {
		platformhttp.WriteError(w, http.StatusNotFound, "SEASON_NOT_FOUND", "Season not found", platformhttp.RequestIDFrom(r), nil)
		return nil, false
	}
	return se, true
}

// verifyEpisode confirms episodeID lives under seasonID→projectID.
func (h *StoryHandler) verifyEpisode(w http.ResponseWriter, r *http.Request, projectID, seasonID, episodeID string) (*domain.Episode, bool) {
	if _, ok := h.verifySeason(w, r, projectID, seasonID); !ok {
		return nil, false
	}
	ep, err := h.service.GetEpisode(r.Context(), episodeID)
	if err != nil || ep.SeasonID != seasonID {
		platformhttp.WriteError(w, http.StatusNotFound, "EPISODE_NOT_FOUND", "Episode not found", platformhttp.RequestIDFrom(r), nil)
		return nil, false
	}
	return ep, true
}

// verifyScene confirms sceneID lives under episodeID→seasonID→projectID.
func (h *StoryHandler) verifyScene(w http.ResponseWriter, r *http.Request, projectID, seasonID, episodeID, sceneID string) (*domain.Scene, bool) {
	if _, ok := h.verifyEpisode(w, r, projectID, seasonID, episodeID); !ok {
		return nil, false
	}
	sc, err := h.service.GetScene(r.Context(), sceneID)
	if err != nil || sc.EpisodeID != episodeID {
		platformhttp.WriteError(w, http.StatusNotFound, "SCENE_NOT_FOUND", "Scene not found", platformhttp.RequestIDFrom(r), nil)
		return nil, false
	}
	return sc, true
}

func (h *StoryHandler) getStory(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	series, err := h.service.GetSeriesByProject(r.Context(), projectID)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "SERIES_NOT_FOUND", "No series for project", platformhttp.RequestIDFrom(r), nil)
		return
	}
	seasons, _ := h.service.ListSeasons(r.Context(), series.ID)
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"series":  series,
		"seasons": seasons,
	})
}

type ensureSeriesReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r *ensureSeriesReq) Validate() error {
	if r.Title == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "title", Message: "required"}}}
	}
	return nil
}

func (h *StoryHandler) ensureSeries(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	var req ensureSeriesReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	series, err := h.service.EnsureSeries(r.Context(), projectID, req.Title, req.Description)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, series)
}

func (h *StoryHandler) listSeasons(w http.ResponseWriter, r *http.Request) {
	series, err := h.service.GetSeriesByProject(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "SERIES_NOT_FOUND", "No series for project", platformhttp.RequestIDFrom(r), nil)
		return
	}
	seasons, err := h.service.ListSeasons(r.Context(), series.ID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": seasons})
}

type createSeasonReq struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

func (r *createSeasonReq) Validate() error {
	if r.Number < 1 {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "number", Message: "must be >= 1"}}}
	}
	return nil
}

func (h *StoryHandler) createSeason(w http.ResponseWriter, r *http.Request) {
	var req createSeasonReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	series, err := h.service.GetSeriesByProject(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "SERIES_NOT_FOUND", "No series for project", platformhttp.RequestIDFrom(r), nil)
		return
	}
	s, err := h.service.CreateSeason(r.Context(), series.ID, req.Title, req.Summary, req.Number)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, s)
}

func (h *StoryHandler) getSeason(w http.ResponseWriter, r *http.Request) {
	s, ok := h.verifySeason(w, r, r.PathValue("projectId"), r.PathValue("seasonId"))
	if !ok {
		return
	}
	arcs, _ := h.service.ListArcs(r.Context(), s.ID)
	episodes, _ := h.service.ListEpisodes(r.Context(), s.ID)
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"season":   s,
		"arcs":     arcs,
		"episodes": episodes,
	})
}

type createArcReq struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
}

func (r *createArcReq) Validate() error {
	if r.Title == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "title", Message: "required"}}}
	}
	return nil
}

func (h *StoryHandler) listArcs(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifySeason(w, r, r.PathValue("projectId"), r.PathValue("seasonId")); !ok {
		return
	}
	arcs, err := h.service.ListArcs(r.Context(), r.PathValue("seasonId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": arcs})
}

func (h *StoryHandler) createArc(w http.ResponseWriter, r *http.Request) {
	var req createArcReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	if _, ok := h.verifySeason(w, r, r.PathValue("projectId"), r.PathValue("seasonId")); !ok {
		return
	}
	a, err := h.service.CreateArc(r.Context(), r.PathValue("seasonId"), req.Title, req.Number)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, a)
}

type createEpisodeReq struct {
	ArcID   string `json:"arc_id"`
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

func (r *createEpisodeReq) Validate() error {
	if r.Number < 1 {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "number", Message: "must be >= 1"}}}
	}
	return nil
}

func (h *StoryHandler) listEpisodes(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifySeason(w, r, r.PathValue("projectId"), r.PathValue("seasonId")); !ok {
		return
	}
	eps, err := h.service.ListEpisodes(r.Context(), r.PathValue("seasonId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": eps})
}

func (h *StoryHandler) createEpisode(w http.ResponseWriter, r *http.Request) {
	var req createEpisodeReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	if _, ok := h.verifySeason(w, r, r.PathValue("projectId"), r.PathValue("seasonId")); !ok {
		return
	}
	ep, err := h.service.CreateEpisode(r.Context(), r.PathValue("seasonId"), req.ArcID, req.Title, req.Summary, req.Number)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, ep)
}

func (h *StoryHandler) getEpisode(w http.ResponseWriter, r *http.Request) {
	ep, ok := h.verifyEpisode(w, r, r.PathValue("projectId"), r.PathValue("seasonId"), r.PathValue("episodeId"))
	if !ok {
		return
	}
	scenes, _ := h.service.ListScenes(r.Context(), ep.ID)
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"episode": ep, "scenes": scenes})
}

type updateEpisodeReq struct {
	Title   *string               `json:"title"`
	Summary *string               `json:"summary"`
	Script  *string               `json:"script"`
	Status  *domain.EpisodeStatus `json:"status"`
}

func (r *updateEpisodeReq) Validate() error {
	if r.Status == nil {
		return nil
	}
	switch *r.Status {
	case domain.EpisodePlanned, domain.EpisodeScripted, domain.EpisodeProducing,
		domain.EpisodeValidating, domain.EpisodeCompleted, domain.EpisodePublished:
		return nil
	}
	return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "status", Message: "invalid episode status"}}}
}

func (h *StoryHandler) updateEpisode(w http.ResponseWriter, r *http.Request) {
	episodeID := r.PathValue("episodeId")
	var req updateEpisodeReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	ep, ok := h.verifyEpisode(w, r, r.PathValue("projectId"), r.PathValue("seasonId"), episodeID)
	if !ok {
		return
	}
	var err error
	if req.Title != nil || req.Summary != nil {
		if ep, err = h.service.UpdateEpisodeMeta(r.Context(), episodeID, req.Title, req.Summary); err != nil {
			platformhttp.WriteErrorFrom(w, r, err)
			return
		}
	}
	if req.Script != nil {
		if ep, err = h.service.SetEpisodeScript(r.Context(), episodeID, *req.Script); err != nil {
			platformhttp.WriteErrorFrom(w, r, err)
			return
		}
	}
	if req.Status != nil {
		if ep, err = h.service.SetEpisodeStatus(r.Context(), episodeID, *req.Status); err != nil {
			platformhttp.WriteErrorFrom(w, r, err)
			return
		}
	}
	platformhttp.WriteJSON(w, http.StatusOK, ep)
}

type createSceneReq struct {
	Number       int      `json:"number"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	LocationID   string   `json:"location_id"`
	TimeOfDay    string   `json:"time_of_day"`
	CharacterIDs []string `json:"character_ids"`
}

func (r *createSceneReq) Validate() error {
	if r.Number < 1 {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "number", Message: "must be >= 1"}}}
	}
	return nil
}

func (h *StoryHandler) listScenes(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyEpisode(w, r, r.PathValue("projectId"), r.PathValue("seasonId"), r.PathValue("episodeId")); !ok {
		return
	}
	scenes, err := h.service.ListScenes(r.Context(), r.PathValue("episodeId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": scenes})
}

func (h *StoryHandler) createScene(w http.ResponseWriter, r *http.Request) {
	var req createSceneReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	if _, ok := h.verifyEpisode(w, r, r.PathValue("projectId"), r.PathValue("seasonId"), r.PathValue("episodeId")); !ok {
		return
	}
	sc, err := h.service.CreateScene(r.Context(), r.PathValue("episodeId"),
		req.Title, req.Description, req.LocationID, req.TimeOfDay, req.Number, req.CharacterIDs)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, sc)
}

func (h *StoryHandler) getScene(w http.ResponseWriter, r *http.Request) {
	sc, ok := h.verifyScene(w, r, r.PathValue("projectId"), r.PathValue("seasonId"), r.PathValue("episodeId"), r.PathValue("sceneId"))
	if !ok {
		return
	}
	beats, _ := h.service.ListBeats(r.Context(), sc.ID)
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"scene": sc, "beats": beats})
}

type createBeatReq struct {
	Seq         int    `json:"seq"`
	Action      string `json:"action"`
	Dialogue    string `json:"dialogue"`
	CharacterID string `json:"character_id"`
}

func (r *createBeatReq) Validate() error { return nil }

func (h *StoryHandler) listBeats(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.verifyScene(w, r, r.PathValue("projectId"), r.PathValue("seasonId"), r.PathValue("episodeId"), r.PathValue("sceneId")); !ok {
		return
	}
	beats, err := h.service.ListBeats(r.Context(), r.PathValue("sceneId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": beats})
}

func (h *StoryHandler) createBeat(w http.ResponseWriter, r *http.Request) {
	var req createBeatReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	if _, ok := h.verifyScene(w, r, r.PathValue("projectId"), r.PathValue("seasonId"), r.PathValue("episodeId"), r.PathValue("sceneId")); !ok {
		return
	}
	b, err := h.service.CreateBeat(r.Context(), r.PathValue("sceneId"), req.Action, req.Dialogue, req.CharacterID, req.Seq)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, b)
}
