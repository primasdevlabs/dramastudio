package http

import (
	"net/http"
	"strconv"
	"time"

	"dramastudio/internal/analytics/application/services"
	"dramastudio/internal/analytics/domain"
	platformhttp "dramastudio/internal/platform/http"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/analytics/metrics", h.listMetrics)
	mux.HandleFunc("POST /v1/projects/{projectId}/analytics/metrics", h.recordMetric)
	mux.HandleFunc("GET /v1/projects/{projectId}/analytics/summary", h.summary)
}

type recordMetricReq struct {
	EpisodeID     string                 `json:"episode_id"`
	PublicationID string                 `json:"publication_id"`
	Metric        string                 `json:"metric"`
	Value         float64                `json:"value"`
	Dimensions    map[string]interface{} `json:"dimensions"`
}

func (r *recordMetricReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Metric == "" {
		ve.Add("metric", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func (h *AnalyticsHandler) listMetrics(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	ms, err := h.service.ListMetrics(r.Context(), r.PathValue("projectId"), limit)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": ms})
}

func (h *AnalyticsHandler) recordMetric(w http.ResponseWriter, r *http.Request) {
	var req recordMetricReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	m, err := h.service.RecordMetric(r.Context(), &domain.Metric{
		ProjectID:     r.PathValue("projectId"),
		EpisodeID:     req.EpisodeID,
		PublicationID: req.PublicationID,
		Name:          req.Metric,
		Value:         req.Value,
		Dimensions:    req.Dimensions,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, m)
}

func (h *AnalyticsHandler) summary(w http.ResponseWriter, r *http.Request) {
	metric := r.URL.Query().Get("metric")
	since, _ := time.Parse(time.RFC3339, r.URL.Query().Get("since"))
	until, _ := time.Parse(time.RFC3339, r.URL.Query().Get("until"))
	s, err := h.service.Summarize(r.Context(), r.PathValue("projectId"), metric, since, until)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, s)
}
