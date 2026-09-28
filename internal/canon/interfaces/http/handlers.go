package http

import (
	"net/http"

	"dramastudio/internal/canon/application/services"
	"dramastudio/internal/canon/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type CanonHandler struct {
	service *services.CanonService
}

func NewCanonHandler(service *services.CanonService) *CanonHandler {
	return &CanonHandler{service: service}
}

func (h *CanonHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/canon/facts", h.listFacts)
	mux.HandleFunc("POST /v1/projects/{projectId}/canon/facts", h.createFact)
	mux.HandleFunc("GET /v1/projects/{projectId}/canon/facts/{factId}", h.getFact)
	mux.HandleFunc("GET /v1/projects/{projectId}/canon/facts/{factId}/versions", h.factVersions)
	mux.HandleFunc("POST /v1/projects/{projectId}/canon/facts/{factId}/retcon", h.retconFact)
	mux.HandleFunc("GET /v1/projects/{projectId}/canon/rules", h.listRules)
	mux.HandleFunc("POST /v1/projects/{projectId}/canon/rules", h.createRule)
	mux.HandleFunc("GET /v1/projects/{projectId}/canon/knowledge/{characterId}", h.getKnowledge)
	mux.HandleFunc("POST /v1/projects/{projectId}/canon/knowledge/{characterId}", h.grantKnowledge)
}

type createFactReq struct {
	EntityID          string  `json:"entity_id"`
	Subject           string  `json:"subject"`
	Predicate         string  `json:"predicate"`
	Object            string  `json:"object"`
	Type              string  `json:"type"`
	IntroducedEpisode string  `json:"introduced_episode"`
	EffectiveFrom     string  `json:"effective_from"`
	EffectiveUntil    string  `json:"effective_until"`
	Source            string  `json:"source"`
	Confidence        float64 `json:"confidence"`
}

func (r *createFactReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Subject == "" {
		ve.Add("subject", "required")
	}
	if r.Predicate == "" {
		ve.Add("predicate", "required")
	}
	if r.Object == "" {
		ve.Add("object", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *CanonHandler) listFacts(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	entity := r.URL.Query().Get("entity_id")
	var facts []*domain.StoryFact
	var err error
	if entity != "" {
		facts, err = h.service.ListFactsForEntity(r.Context(), projectID, entity)
	} else {
		facts, err = h.service.ListFacts(r.Context(), projectID)
	}
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": facts})
}

func (h *CanonHandler) createFact(w http.ResponseWriter, r *http.Request) {
	var req createFactReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	fact, err := h.service.EstablishFact(r.Context(), &domain.StoryFact{
		ProjectID:         r.PathValue("projectId"),
		EntityID:          req.EntityID,
		Subject:           req.Subject,
		Predicate:         req.Predicate,
		Object:            req.Object,
		Type:              req.Type,
		IntroducedEpisode: req.IntroducedEpisode,
		EffectiveFrom:     req.EffectiveFrom,
		EffectiveUntil:    req.EffectiveUntil,
		Source:            req.Source,
		Confidence:        req.Confidence,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, fact)
}

func (h *CanonHandler) getFact(w http.ResponseWriter, r *http.Request) {
	f, err := h.service.GetFact(r.Context(), r.PathValue("projectId"), r.PathValue("factId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "FACT_NOT_FOUND", "Fact not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, f)
}

func (h *CanonHandler) factVersions(w http.ResponseWriter, r *http.Request) {
	vers, err := h.service.ListFactVersions(r.Context(), r.PathValue("projectId"), r.PathValue("factId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "FACT_NOT_FOUND", "Fact not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": vers})
}

func (h *CanonHandler) retconFact(w http.ResponseWriter, r *http.Request) {
	f, err := h.service.Retcon(r.Context(), r.PathValue("projectId"), r.PathValue("factId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "FACT_NOT_FOUND", "Fact not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, f)
}

type createRuleReq struct {
	RuleText string `json:"rule_text"`
	Enforced bool   `json:"enforced"`
}

func (r *createRuleReq) Validate() error {
	if r.RuleText == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "rule_text", Message: "required"}}}
	}
	return nil
}

func (h *CanonHandler) listRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.service.ListRules(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": rules})
}

func (h *CanonHandler) createRule(w http.ResponseWriter, r *http.Request) {
	var req createRuleReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	rule, err := h.service.CreateRule(r.Context(), r.PathValue("projectId"), req.RuleText, req.Enforced)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, rule)
}

func (h *CanonHandler) getKnowledge(w http.ResponseWriter, r *http.Request) {
	episodeID := r.URL.Query().Get("episode_id")
	ks, err := h.service.GetKnowledge(r.Context(), r.PathValue("characterId"), episodeID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, ks)
}

type grantKnowledgeReq struct {
	EpisodeID string   `json:"episode_id"`
	FactIDs   []string `json:"fact_ids"`
}

func (r *grantKnowledgeReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.EpisodeID == "" {
		ve.Add("episode_id", "required")
	}
	if len(r.FactIDs) == 0 {
		ve.Add("fact_ids", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *CanonHandler) grantKnowledge(w http.ResponseWriter, r *http.Request) {
	var req grantKnowledgeReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	ks, err := h.service.GrantKnowledge(r.Context(), r.PathValue("projectId"), r.PathValue("characterId"), req.EpisodeID, req.FactIDs)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, ks)
}
