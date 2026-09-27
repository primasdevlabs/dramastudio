package http

import (
	"net/http"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/world/application/services"
)

type WorldHandler struct {
	service *services.WorldService
}

func NewWorldHandler(service *services.WorldService) *WorldHandler {
	return &WorldHandler{service: service}
}

type createLocationReq struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
}

func (r *createLocationReq) Validate() error {
	if r.Name == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "name", Message: "required"}}}
	}
	return nil
}

type addVariantReq struct {
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
}

func (r *addVariantReq) Validate() error {
	if r.Name == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "name", Message: "required"}}}
	}
	return nil
}

type createPropReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LocationID  string `json:"location_id"`
}

func (r *createPropReq) Validate() error {
	if r.Name == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "name", Message: "required"}}}
	}
	return nil
}

type createRuleReq struct {
	Text     string `json:"text"`
	Category string `json:"category"`
}

func (r *createRuleReq) Validate() error {
	if r.Text == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "text", Message: "required"}}}
	}
	return nil
}

func (h *WorldHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/world/locations", h.listLocations)
	mux.HandleFunc("POST /v1/projects/{projectId}/world/locations", h.createLocation)
	mux.HandleFunc("GET /v1/projects/{projectId}/world/locations/{locationId}", h.getLocation)
	mux.HandleFunc("POST /v1/projects/{projectId}/world/locations/{locationId}/variants", h.addVariant)
	mux.HandleFunc("GET /v1/projects/{projectId}/world/props", h.listProps)
	mux.HandleFunc("POST /v1/projects/{projectId}/world/props", h.createProp)
	mux.HandleFunc("GET /v1/projects/{projectId}/world/rules", h.listRules)
	mux.HandleFunc("POST /v1/projects/{projectId}/world/rules", h.createRule)
}

func (h *WorldHandler) listLocations(w http.ResponseWriter, r *http.Request) {
	locs, err := h.service.ListLocations(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": locs})
}

func (h *WorldHandler) createLocation(w http.ResponseWriter, r *http.Request) {
	var req createLocationReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	loc, err := h.service.CreateLocation(r.Context(), r.PathValue("projectId"),
		req.Name, req.Kind, req.Description, req.ParentID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, loc)
}

func (h *WorldHandler) getLocation(w http.ResponseWriter, r *http.Request) {
	loc, err := h.service.GetLocation(r.Context(), r.PathValue("locationId"))
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "LOCATION_NOT_FOUND", "Location not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, loc)
}

func (h *WorldHandler) addVariant(w http.ResponseWriter, r *http.Request) {
	var req addVariantReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	loc, err := h.service.AddVariant(r.Context(), r.PathValue("locationId"), req.Name, req.Attributes)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, loc)
}

func (h *WorldHandler) listProps(w http.ResponseWriter, r *http.Request) {
	props, err := h.service.ListProps(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": props})
}

func (h *WorldHandler) createProp(w http.ResponseWriter, r *http.Request) {
	var req createPropReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	p, err := h.service.CreateProp(r.Context(), r.PathValue("projectId"), req.Name, req.Description, req.LocationID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, p)
}

func (h *WorldHandler) listRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.service.ListWorldRules(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": rules})
}

func (h *WorldHandler) createRule(w http.ResponseWriter, r *http.Request) {
	var req createRuleReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	rule, err := h.service.CreateWorldRule(r.Context(), r.PathValue("projectId"), req.Text, req.Category)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, rule)
}
