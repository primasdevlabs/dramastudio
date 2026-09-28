package http

import (
	"context"
	"net/http"

	platformhttp "dramastudio/internal/platform/http"
	"dramastudio/internal/platform/security"
	"dramastudio/internal/production/application/services"
	"dramastudio/internal/production/domain"
)

type ProductionHandler struct {
	service *services.ProductionService
}

func NewProductionHandler(service *services.ProductionService) *ProductionHandler {
	return &ProductionHandler{service: service}
}

func (h *ProductionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/projects/{projectId}/production/runs", h.listRuns)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/runs", h.startRun)
	mux.HandleFunc("GET /v1/projects/{projectId}/production/runs/{runId}", h.getRun)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/runs/{runId}/pause", h.pauseRun)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/runs/{runId}/resume", h.resumeRun)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/runs/{runId}/stop", h.stopRun)
	mux.HandleFunc("GET /v1/projects/{projectId}/production/runs/{runId}/jobs", h.listRunJobs)
	mux.HandleFunc("GET /v1/projects/{projectId}/production/jobs", h.listJobs)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/jobs", h.createJob)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/jobs/{jobId}/retry", h.retryJob)
	mux.HandleFunc("GET /v1/projects/{projectId}/production/shots", h.listShots)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/shots", h.createShot)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/shots/{shotId}/approve", h.approveShot)
	mux.HandleFunc("GET /v1/projects/{projectId}/production/approvals", h.listApprovals)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/approvals", h.requestApproval)
	mux.HandleFunc("POST /v1/projects/{projectId}/production/approvals/{approvalId}/decide", h.decideApproval)
}

// --- requests ---

type startRunReq struct {
	EpisodeID    string `json:"episode_id"`
	BibleVersion int    `json:"bible_version"`
}

func (r *startRunReq) Validate() error {
	if r.EpisodeID == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "episode_id", Message: "required"}}}
	}
	return nil
}

type createJobReq struct {
	RunID          string `json:"run_id"`
	EpisodeID      string `json:"episode_id"`
	SceneID        string `json:"scene_id"`
	ShotID         string `json:"shot_id"`
	Kind           string `json:"kind"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (r *createJobReq) Validate() error {
	if r.Kind == "" {
		return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "kind", Message: "required"}}}
	}
	return nil
}

type createShotReq struct {
	EpisodeID   string                 `json:"episode_id"`
	SceneID     string                 `json:"scene_id"`
	Seq         int                    `json:"seq"`
	Description string                 `json:"description"`
	Camera      map[string]interface{} `json:"camera"`
	Characters  []string               `json:"characters"`
	LocationID  string                 `json:"location_id"`
	DurationSec float64                `json:"duration_sec"`
}

func (r *createShotReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.EpisodeID == "" {
		ve.Add("episode_id", "required")
	}
	if r.SceneID == "" {
		ve.Add("scene_id", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type approveShotReq struct {
	AssetURL string `json:"asset_url"`
}

func (r *approveShotReq) Validate() error { return nil }

type approvalReq struct {
	EpisodeID string `json:"episode_id"`
	Stage     string `json:"stage"`
	TargetID  string `json:"target_id"`
}

func (r *approvalReq) Validate() error {
	ve := &platformhttp.ValidationError{}
	if r.Stage == "" {
		ve.Add("stage", "required")
	}
	if r.TargetID == "" {
		ve.Add("target_id", "required")
	}
	if ve.HasErrors() {
		return ve
	}
	return nil
}

type decideReq struct {
	Decision  domain.ApprovalDecision `json:"decision"`
	Notes     string                  `json:"notes"`
	DecidedBy string                  `json:"decided_by"`
}

func (r *decideReq) Validate() error {
	switch r.Decision {
	case domain.DecisionApprove, domain.DecisionReject, domain.DecisionRequestRevision,
		domain.DecisionRegenerate, domain.DecisionPause, domain.DecisionStop,
		domain.DecisionOverride:
		return nil
	}
	return &platformhttp.ValidationError{Fields: []platformhttp.FieldError{{Field: "decision", Message: "must be APPROVE|REJECT|REQUEST_REVISION|REGENERATE|PAUSE|STOP|OVERRIDE"}}}
}

// --- handlers ---

// verifyOwned returns true when the entity's project_id matches the path
// project; writes 404 otherwise (no cross-project existence leak).
func (h *ProductionHandler) verifyOwned(w http.ResponseWriter, r *http.Request, entityProjectID string) bool {
	if entityProjectID != "" && entityProjectID != r.PathValue("projectId") {
		platformhttp.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Resource not found", platformhttp.RequestIDFrom(r), nil)
		return false
	}
	return true
}

func (h *ProductionHandler) listRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := h.service.ListRuns(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": runs})
}

func (h *ProductionHandler) startRun(w http.ResponseWriter, r *http.Request) {
	var req startRunReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	run, err := h.service.StartRun(r.Context(), r.PathValue("projectId"), req.EpisodeID, req.BibleVersion)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, run)
}

func (h *ProductionHandler) getRun(w http.ResponseWriter, r *http.Request) {
	run, err := h.service.GetRun(r.Context(), r.PathValue("runId"))
	if err != nil {
		h.writeNotFound(w, r, "RUN_NOT_FOUND", "Production run not found")
		return
	}
	if !h.verifyOwned(w, r, run.ProjectID) {
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, run)
}

func (h *ProductionHandler) pauseRun(w http.ResponseWriter, r *http.Request) {
	h.runTransition(w, r, h.service.PauseRun)
}

func (h *ProductionHandler) resumeRun(w http.ResponseWriter, r *http.Request) {
	h.runTransition(w, r, h.service.ResumeRun)
}

func (h *ProductionHandler) stopRun(w http.ResponseWriter, r *http.Request) {
	h.runTransition(w, r, h.service.StopRun)
}

func (h *ProductionHandler) runTransition(w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, id string) (*domain.ProductionRun, error)) {
	existing, err := h.service.GetRun(r.Context(), r.PathValue("runId"))
	if err != nil {
		h.writeNotFound(w, r, "RUN_NOT_FOUND", "Production run not found")
		return
	}
	if !h.verifyOwned(w, r, existing.ProjectID) {
		return
	}
	run, err := fn(r.Context(), r.PathValue("runId"))
	if err != nil {
		if err == domain.ErrInvalidTransition {
			platformhttp.WriteError(w, http.StatusConflict, "INVALID_TRANSITION", err.Error(), platformhttp.RequestIDFrom(r), nil)
			return
		}
		h.writeNotFound(w, r, "RUN_NOT_FOUND", "Production run not found")
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, run)
}

func (h *ProductionHandler) listRunJobs(w http.ResponseWriter, r *http.Request) {
	run, err := h.service.GetRun(r.Context(), r.PathValue("runId"))
	if err != nil {
		h.writeNotFound(w, r, "RUN_NOT_FOUND", "Production run not found")
		return
	}
	if !h.verifyOwned(w, r, run.ProjectID) {
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": run.Jobs})
}

func (h *ProductionHandler) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.ListJobs(r.Context(), r.PathValue("projectId"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": jobs})
}

func (h *ProductionHandler) createJob(w http.ResponseWriter, r *http.Request) {
	var req createJobReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	key := req.IdempotencyKey
	if k := platformhttp.IdempotencyKey(r); k != "" {
		key = k
	}
	job, err := h.service.CreateJob(r.Context(), &domain.ProductionJob{
		ProjectID:      r.PathValue("projectId"),
		RunID:          req.RunID,
		EpisodeID:      req.EpisodeID,
		SceneID:        req.SceneID,
		ShotID:         req.ShotID,
		Kind:           req.Kind,
		IdempotencyKey: key,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, job)
}

func (h *ProductionHandler) retryJob(w http.ResponseWriter, r *http.Request) {
	existing, err := h.service.GetJob(r.Context(), r.PathValue("jobId"))
	if err != nil {
		h.writeNotFound(w, r, "JOB_NOT_FOUND", "Job not found")
		return
	}
	if !h.verifyOwned(w, r, existing.ProjectID) {
		return
	}
	j, err := h.service.RetryJob(r.Context(), r.PathValue("jobId"))
	if err != nil {
		if err == domain.ErrInvalidTransition {
			platformhttp.WriteError(w, http.StatusConflict, "INVALID_TRANSITION", "Only failed jobs can be retried", platformhttp.RequestIDFrom(r), nil)
			return
		}
		h.writeNotFound(w, r, "JOB_NOT_FOUND", "Job not found")
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, j)
}

func (h *ProductionHandler) listShots(w http.ResponseWriter, r *http.Request) {
	shots, err := h.service.ListShots(r.Context(), r.PathValue("projectId"), r.URL.Query().Get("episode_id"), r.URL.Query().Get("scene_id"))
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": shots})
}

func (h *ProductionHandler) createShot(w http.ResponseWriter, r *http.Request) {
	var req createShotReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	shot, err := h.service.CreateShot(r.Context(), &domain.Shot{
		ProjectID:   r.PathValue("projectId"),
		EpisodeID:   req.EpisodeID,
		SceneID:     req.SceneID,
		Seq:         req.Seq,
		Description: req.Description,
		Camera:      req.Camera,
		Characters:  req.Characters,
		LocationID:  req.LocationID,
		DurationSec: req.DurationSec,
	})
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, shot)
}

func (h *ProductionHandler) approveShot(w http.ResponseWriter, r *http.Request) {
	var req approveShotReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	shot, err := h.service.GetShot(r.Context(), r.PathValue("shotId"))
	if err != nil {
		h.writeNotFound(w, r, "SHOT_NOT_FOUND", "Shot not found")
		return
	}
	if !h.verifyOwned(w, r, shot.ProjectID) {
		return
	}
	shot, err = h.service.ApproveShot(r.Context(), r.PathValue("shotId"), req.AssetURL)
	if err != nil {
		h.writeNotFound(w, r, "SHOT_NOT_FOUND", "Shot not found")
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, shot)
}

func (h *ProductionHandler) listApprovals(w http.ResponseWriter, r *http.Request) {
	pendingOnly := r.URL.Query().Get("pending") == "true"
	apps, err := h.service.ListApprovals(r.Context(), r.PathValue("projectId"), pendingOnly)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{"items": apps})
}

func (h *ProductionHandler) requestApproval(w http.ResponseWriter, r *http.Request) {
	var req approvalReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	a, err := h.service.RequestApproval(r.Context(), r.PathValue("projectId"), req.EpisodeID, req.Stage, req.TargetID)
	if err != nil {
		platformhttp.WriteErrorFrom(w, r, err)
		return
	}
	platformhttp.WriteJSON(w, http.StatusCreated, a)
}

func (h *ProductionHandler) decideApproval(w http.ResponseWriter, r *http.Request) {
	var req decideReq
	if !platformhttp.DecodeAndValidate(w, r, &req) {
		return
	}
	a, err := h.service.GetApproval(r.Context(), r.PathValue("approvalId"))
	if err != nil {
		h.writeNotFound(w, r, "APPROVAL_NOT_FOUND", "Approval request not found")
		return
	}
	if !h.verifyOwned(w, r, a.ProjectID) {
		return
	}
	// decided_by is the authenticated principal — never client-supplied
	// (audit integrity, §63).
	decidedBy := req.DecidedBy
	if p, ok := security.PrincipalFrom(r.Context()); ok && p.UserID != "" {
		decidedBy = p.UserID
	}
	a, err = h.service.Decide(r.Context(), r.PathValue("approvalId"), req.Decision, req.Notes, decidedBy)
	if err != nil {
		if err == domain.ErrInvalidTransition {
			platformhttp.WriteError(w, http.StatusConflict, "ALREADY_DECIDED", "Approval already decided", platformhttp.RequestIDFrom(r), nil)
			return
		}
		h.writeNotFound(w, r, "APPROVAL_NOT_FOUND", "Approval request not found")
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, a)
}

func (h *ProductionHandler) writeNotFound(w http.ResponseWriter, r *http.Request, code, msg string) {
	platformhttp.WriteError(w, http.StatusNotFound, code, msg, platformhttp.RequestIDFrom(r), nil)
}
