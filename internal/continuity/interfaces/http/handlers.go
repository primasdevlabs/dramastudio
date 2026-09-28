package http

import (
	"net/http"

	"dramastudio/internal/continuity/application/services"
	"dramastudio/internal/continuity/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type ContinuityHandler struct {
	service *services.ContinuityService
}

func NewContinuityHandler(service *services.ContinuityService) *ContinuityHandler {
	return &ContinuityHandler{service: service}
}

func (h *ContinuityHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/projects/{projectId}/continuity/checks", h.runCheck)
	mux.HandleFunc("GET /v1/projects/{projectId}/continuity/checks", h.listChecks)
	mux.HandleFunc("GET /v1/projects/{projectId}/continuity/issues", h.listIssues)
	mux.HandleFunc("POST /v1/projects/{projectId}/continuity/issues", h.reportIssue)
	mux.HandleFunc("POST /v1/projects/{projectId}/continuity/issues/{issueId}/resolve", h.resolveIssue)
	mux.HandleFunc("GET /v1/projects/{projectId}/continuity/timeline", h.listTimeline)
	mux.HandleFunc("POST /v1/projects/{projectId}/continuity/timeline", h.addTimelineEvent)
}

// --- requests ---

type factView struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

type runCheckReq struct {
	EpisodeID         string            `json:"episode_id"`
	SceneID           string            `json:"scene_id"`
	CheckType         domain.CheckType  `json:"check_type"`
	TargetID          string            `json:"target_id"`
	EntityID          string            `json:"entity_id"`
	CanonicalAttrs    map[string]string `json:"canonical_attrs"`
	ActualAttrs       map[string]string `json:"actual_attrs"`
	CanonicalWardrobe []string          `json:"canonical_wardrobe"`
	AssignedWardrobe  []string          `json:"assigned_wardrobe"`
	Facts             []factView        `json:"facts"`
	Claims            []factView        `json:"claims"`
}

func (r *runCheckReq) Validate() error {
	switch r.CheckType {
	case domain.CheckStory, domain.CheckTimeline, domain.CheckCharacter,
		domain.CheckWardrobe, domain.CheckVisual:
		return nil
	}
	return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "check_type", Message: "must be story|timeline|character|wardrobe|visual"}}}
}

type reportIssueReq struct {
	EpisodeID     string          `json:"episode_id"`
	SceneID       string          `json:"scene_id"`
	Category      string          `json:"category"`
	Severity      domain.Severity `json:"severity"`
	Entity        string          `json:"entity"`
	ExpectedState string          `json:"expected_state"`
	ActualState   string          `json:"actual_state"`
	Cause         string          `json:"cause"`
	Evidence      string          `json:"evidence"`
}

func (r *reportIssueReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Category == "" {
		ve.Add("category", "required")
	}
	if r.Severity == "" {
		ve.Add("severity", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type resolveReq struct {
	Resolution string `json:"resolution"`
	WontFix    bool   `json:"wontfix"`
}

func (r *resolveReq) Validate() error { return nil }

type timelineEventReq struct {
	EpisodeID    string   `json:"episode_id"`
	SceneID      string   `json:"scene_id"`
	WorldTime    string   `json:"world_time"`
	EventOrder   int      `json:"event_order"`
	Participants []string `json:"participants"`
	LocationID   string   `json:"location_id"`
	DurationSec  float64  `json:"duration_sec"`
	Description  string   `json:"description"`
}

func (r *timelineEventReq) Validate() error { return nil }

// --- handlers ---

func (h *ContinuityHandler) runCheck(w http.ResponseWriter, r *http.Request) {
	var req runCheckReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	in := services.CheckInput{
		ProjectID:         r.PathValue("projectId"),
		EpisodeID:         req.EpisodeID,
		SceneID:           req.SceneID,
		CheckType:         req.CheckType,
		TargetID:          req.TargetID,
		EntityID:          req.EntityID,
		CanonicalAttrs:    req.CanonicalAttrs,
		ActualAttrs:       req.ActualAttrs,
		CanonicalWardrobe: req.CanonicalWardrobe,
		AssignedWardrobe:  req.AssignedWardrobe,
	}
	for _, f := range req.Facts {
		in.Facts = append(in.Facts, services.FactView{Subject: f.Subject, Predicate: f.Predicate, Object: f.Object})
	}
	for _, c := range req.Claims {
		in.Claims = append(in.Claims, services.FactView{Subject: c.Subject, Predicate: c.Predicate, Object: c.Object})
	}
	// Timeline checks read events from the repository, not the request body.
	if req.CheckType == domain.CheckTimeline {
		events, err := h.service.ListTimelineEvents(r.Context(), in.ProjectID, req.EpisodeID)
		if err != nil {
			platformhttp.WriteErrorFrom(w, r, err)
			return
		}
		for _, e := range events {
			in.Events = append(in.Events, *e)
		}
	}
	check, issues, err := h.service.RunCheck(r.Context(), in)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, map[string]interface{}{"check": check, "issues": issues})
}

func (h *ContinuityHandler) listChecks(w http.ResponseWriter, r *http.Request) {
	checks, err := h.service.ListChecks(r.Context(), r.PathValue("projectId"), r.URL.Query().Get("episode_id"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": checks})
}

func (h *ContinuityHandler) listIssues(w http.ResponseWriter, r *http.Request) {
	f := domain.IssueFilter{
		EpisodeID: r.URL.Query().Get("episode_id"),
		Status:    domain.IssueStatus(r.URL.Query().Get("status")),
		Severity:  domain.Severity(r.URL.Query().Get("severity")),
		Category:  r.URL.Query().Get("category"),
	}
	issues, err := h.service.ListIssues(r.Context(), r.PathValue("projectId"), f)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": issues})
}

func (h *ContinuityHandler) reportIssue(w http.ResponseWriter, r *http.Request) {
	var req reportIssueReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	issue, err := h.service.ReportIssue(r.Context(), &domain.ContinuityIssue{
		ProjectID:     r.PathValue("projectId"),
		EpisodeID:     req.EpisodeID,
		SceneID:       req.SceneID,
		Category:      req.Category,
		Severity:      req.Severity,
		Entity:        req.Entity,
		ExpectedState: req.ExpectedState,
		ActualState:   req.ActualState,
		Cause:         req.Cause,
		Evidence:      req.Evidence,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, issue)
}

func (h *ContinuityHandler) resolveIssue(w http.ResponseWriter, r *http.Request) {
	var req resolveReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	issue, err := h.service.ResolveIssue(r.Context(), r.PathValue("issueId"), req.Resolution, req.WontFix)
	if err != nil {
		platformhttp.WriteError(w, http.StatusNotFound, "ISSUE_NOT_FOUND", "Continuity issue not found", platformhttp.RequestIDFrom(r), nil)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, issue)
}

func (h *ContinuityHandler) listTimeline(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.ListTimelineEvents(r.Context(), r.PathValue("projectId"), r.URL.Query().Get("episode_id"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": events})
}

func (h *ContinuityHandler) addTimelineEvent(w http.ResponseWriter, r *http.Request) {
	var req timelineEventReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	e, err := h.service.AddTimelineEvent(r.Context(), &domain.TimelineEvent{
		ProjectID:    r.PathValue("projectId"),
		EpisodeID:    req.EpisodeID,
		SceneID:      req.SceneID,
		WorldTime:    req.WorldTime,
		EventOrder:   req.EventOrder,
		Participants: req.Participants,
		LocationID:   req.LocationID,
		DurationSec:  req.DurationSec,
		Description:  req.Description,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, e)
}
