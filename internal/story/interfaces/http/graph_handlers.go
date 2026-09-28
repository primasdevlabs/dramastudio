package http

import (
	"net/http"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/story/domain"
)

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
	FromID   string          `json:"from_id"`
	ToID     string          `json:"to_id"`
	Relation domain.EdgeType `json:"relation"`
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
