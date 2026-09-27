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
	s, err := h.service.GetSeason(r.Context(), r.PathValue("seasonId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "SEASON_NOT_FOUND", "Season not found", platformhttp.RequestIDFrom(r), nil)
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
	ep, err := h.service.CreateEpisode(r.Context(), r.PathValue("seasonId"), req.ArcID, req.Title, req.Summary, req.Number)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, ep)
}

func (h *StoryHandler) getEpisode(w http.ResponseWriter, r *http.Request) {
	ep, err := h.service.GetEpisode(r.Context(), r.PathValue("episodeId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "EPISODE_NOT_FOUND", "Episode not found", platformhttp.RequestIDFrom(r), nil)
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

func (r *updateEpisodeReq) Validate() error { return nil }

func (h *StoryHandler) updateEpisode(w http.ResponseWriter, r *http.Request) {
	episodeID := r.PathValue("episodeId")
	var req updateEpisodeReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	ep, err := h.service.GetEpisode(r.Context(), episodeID)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "EPISODE_NOT_FOUND", "Episode not found", platformhttp.RequestIDFrom(r), nil)
		return
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
	sc, err := h.service.CreateScene(r.Context(), r.PathValue("episodeId"),
		req.Title, req.Description, req.LocationID, req.TimeOfDay, req.Number, req.CharacterIDs)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, sc)
}

func (h *StoryHandler) getScene(w http.ResponseWriter, r *http.Request) {
	sc, err := h.service.GetScene(r.Context(), r.PathValue("sceneId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "SCENE_NOT_FOUND", "Scene not found", platformhttp.RequestIDFrom(r), nil)
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
	b, err := h.service.CreateBeat(r.Context(), r.PathValue("sceneId"), req.Action, req.Dialogue, req.CharacterID, req.Seq)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, b)
}

func (h *StoryHandler) getGraph(w http.ResponseWriter, r *http.Request) {
	g, err := h.service.GetGraph(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, g)
}

type addNodeReq struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	EpisodeID string `json:"episode_id"`
	SceneID   string `json:"scene_id"`
}

func (r *addNodeReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Type == "" {
		ve.Add("type", "required (Event|Reveal|Conflict|Decision|Foreshadowing|Resolution)")
	}
	if r.Title == "" {
		ve.Add("title", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *StoryHandler) addNode(w http.ResponseWriter, r *http.Request) {
	var req addNodeReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	n, err := h.service.AddGraphNode(r.Context(), r.PathValue("projectId"), req.Type, req.Title, req.EpisodeID, req.SceneID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, n)
}

type addEdgeReq struct {
	FromID   string           `json:"from_id"`
	ToID     string           `json:"to_id"`
	Relation domain.EdgeType  `json:"relation"`
}

func (r *addEdgeReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.FromID == "" || r.ToID == "" {
		ve.Add("from_id/to_id", "required")
	}
	switch r.Relation {
	case domain.EdgeCauses, domain.EdgeReveals, domain.EdgeDependsOn,
		domain.EdgeContradicts, domain.EdgeResolves, domain.EdgeForeshadows, domain.EdgeFollows:
	default:
		ve.Add("relation", "invalid edge relation")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *StoryHandler) addEdge(w http.ResponseWriter, r *http.Request) {
	var req addEdgeReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	e, err := h.service.AddGraphEdge(r.Context(), r.PathValue("projectId"), req.FromID, req.ToID, req.Relation)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, e)
}

type createThreadReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *createThreadReq) Validate() error {
	if r.Name == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "name", Message: "required"}}}
	}
	return nil
}

func (h *StoryHandler) listThreads(w http.ResponseWriter, r *http.Request) {
	t, err := h.service.ListPlotThreads(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": t})
}

func (h *StoryHandler) createThread(w http.ResponseWriter, r *http.Request) {
	var req createThreadReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	t, err := h.service.CreatePlotThread(r.Context(), r.PathValue("projectId"), req.Name, req.Description)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, t)
}
