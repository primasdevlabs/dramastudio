// Package local provides an in-process workflow.Engine for fully-local
// development ("Go + Postgres + local filesystem" deployments with no
// Temporal or Hatchet server). Runs execute in goroutines, signals are
// delivered over channels, and status is served from memory.
//
// Durability caveat: unlike Temporal, a process restart loses in-flight
// run state. Workflows re-start cleanly because ProductionRun is persisted
// and activities are idempotent, but there is no event-sourced replay.
package local

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/platform/workflow"
	"dramastudio/internal/platform/workflow/activities"
)

// Engine runs production pipelines in-process.
type Engine struct {
	act  *activities.EpisodeActivities
	mu   sync.Mutex
	runs map[string]*run
	done bool
}

func New(act *activities.EpisodeActivities) *Engine {
	return &Engine{act: act, runs: map[string]*run{}}
}

// run mirrors the Temporal workflow's control semantics with channels.
type run struct {
	status    workflow.WorkflowStatus
	approval  chan string
	pause     chan struct{}
	resume    chan struct{}
	stop      chan struct{}
	result    *workflow.ProduceEpisodeResult
	err       error
	done      chan struct{}
	cancel    context.CancelFunc
	statusMux sync.RWMutex
}

func (r *run) snapshot() workflow.WorkflowStatus {
	r.statusMux.RLock()
	defer r.statusMux.RUnlock()
	return r.status
}

func (r *run) setStage(s string) {
	r.statusMux.Lock()
	r.status.Stage = s
	r.statusMux.Unlock()
	r.drain()
}

// drain consumes pending stop/pause signals without blocking.
func (r *run) drain() {
	for {
		select {
		case <-r.stop:
			r.statusMux.Lock()
			r.status.Stopped = true
			r.statusMux.Unlock()
		case <-r.pause:
			r.statusMux.Lock()
			r.status.Paused = true
			r.statusMux.Unlock()
		default:
			return
		}
	}
}

// ok is the control gate: drains signals, blocks while paused, returns
// false when the run was stopped. Same semantics as controlGate.ok().
func (r *run) ok(ctx context.Context) bool {
	r.drain()
	for {
		st := r.snapshot()
		if st.Stopped {
			return false
		}
		if !st.Paused {
			return true
		}
		select {
		case <-r.resume:
			r.statusMux.Lock()
			r.status.Paused = false
			r.statusMux.Unlock()
		case <-r.stop:
			r.statusMux.Lock()
			r.status.Stopped = true
			r.statusMux.Unlock()
			return false
		case <-ctx.Done():
			return false
		}
	}
}

// awaitDecision suspends on the approval signal (or stop/cancel).
func (r *run) awaitDecision(ctx context.Context) string {
	for {
		if r.snapshot().Stopped {
			return ""
		}
		select {
		case d := <-r.approval:
			r.statusMux.Lock()
			r.status.Paused = false
			r.statusMux.Unlock()
			return d
		case <-r.stop:
			r.statusMux.Lock()
			r.status.Stopped = true
			r.statusMux.Unlock()
			return ""
		case <-ctx.Done():
			return ""
		}
	}
}

// Start launches the episode pipeline in a goroutine. The workflow ID is
// deterministic per run ID so duplicate starts are idempotent.
func (e *Engine) Start(ctx context.Context, in workflow.ProduceEpisodeInput) (string, error) {
	id := "produce-episode-" + in.RunID

	e.mu.Lock()
	if existing, ok := e.runs[id]; ok {
		e.mu.Unlock()
		select {
		case <-existing.done:
			return "", fmt.Errorf("workflow %s already completed", id)
		default:
			return id, nil // idempotent: run already in flight
		}
	}
	r := &run{
		status:   workflow.WorkflowStatus{RunID: in.RunID, Stage: "starting"},
		approval: make(chan string, 8),
		pause:    make(chan struct{}, 8),
		resume:   make(chan struct{}, 8),
		stop:     make(chan struct{}, 8),
		done:     make(chan struct{}),
	}
	e.runs[id] = r
	e.mu.Unlock()

	rctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	go e.execute(rctx, r, in)
	return id, nil
}

// execute mirrors ProduceEpisode's stage sequence, calling the same
// activities directly instead of through the Temporal SDK.
func (e *Engine) execute(ctx context.Context, r *run, in workflow.ProduceEpisodeInput) {
	defer close(r.done)
	defer r.cancel()

	res, err := e.stages(ctx, r, in)
	if err != nil {
		r.err = err
		return
	}
	r.result = res
	if r.snapshot().Stopped {
		res.Status = "stopped"
		return
	}

	// Stage 7: human approval gate (§39) — identical to the Temporal path.
	r.setStage("awaiting_approval")
	approvalID, err := e.act.RequestHumanApproval(ctx, in.ProjectID, in.EpisodeID, "episode", in.EpisodeID)
	if err != nil {
		r.err = fmt.Errorf("RequestHumanApproval: %w", err)
		return
	}
	r.statusMux.Lock()
	r.status.ApprovalID = approvalID
	r.status.Paused = true
	r.statusMux.Unlock()

	res.Approval = r.awaitDecision(ctx)
	switch {
	case r.snapshot().Stopped:
		res.Status = "stopped"
	case res.Approval != "APPROVE" && res.Approval != "OVERRIDE":
		res.Status = "rejected"
	default:
		res.Status = "completed"
	}
}

// stages mirrors runProductionStages: plan -> continuity -> writing ->
// planning -> media. The control gate is checked at the same boundaries.
func (e *Engine) stages(ctx context.Context, r *run, in workflow.ProduceEpisodeInput) (*workflow.ProduceEpisodeResult, error) {
	res := &workflow.ProduceEpisodeResult{EpisodeID: in.EpisodeID}

	r.setStage("episode_plan")
	if !r.ok(ctx) {
		return res, nil
	}
	plan, err := e.act.GenerateEpisodePlan(ctx, in.ProjectID, in.EpisodeID)
	if err != nil {
		return nil, fmt.Errorf("GenerateEpisodePlan: %w", err)
	}
	res.Plan = plan

	r.setStage("continuity_check")
	if !r.ok(ctx) {
		return res, nil
	}
	blocking, err := e.act.ValidateContinuity(ctx, in.ProjectID, in.EpisodeID)
	if err != nil {
		return nil, fmt.Errorf("ValidateContinuity: %w", err)
	}
	if blocking > 0 {
		return nil, fmt.Errorf("continuity check found %d blocking issues", blocking)
	}

	if err := e.writing(ctx, r, in); err != nil || !r.ok(ctx) {
		return res, err
	}
	shotIDs, err := e.planning(ctx, r, in)
	if err != nil || r.snapshot().Stopped {
		return res, err
	}
	if err := e.media(ctx, r, in, shotIDs, res); err != nil {
		return res, err
	}
	return res, nil
}

func (e *Engine) writing(ctx context.Context, r *run, in workflow.ProduceEpisodeInput) error {
	r.setStage("script")
	if !r.ok(ctx) {
		return nil
	}
	if _, err := e.act.GenerateScript(ctx, in.ProjectID, in.EpisodeID); err != nil {
		return fmt.Errorf("GenerateScript: %w", err)
	}
	if _, err := e.act.GenerateDialogue(ctx, in.ProjectID, in.EpisodeID); err != nil {
		return fmt.Errorf("GenerateDialogue: %w", err)
	}
	return nil
}

func (e *Engine) planning(ctx context.Context, r *run, in workflow.ProduceEpisodeInput) ([]string, error) {
	r.setStage("scenes")
	if !r.ok(ctx) {
		return nil, nil
	}
	sceneIDs, err := e.act.GenerateScenePlans(ctx, in.EpisodeID)
	if err != nil {
		return nil, fmt.Errorf("GenerateScenePlans: %w", err)
	}
	r.setStage("storyboard")
	shotIDs, err := e.act.GenerateStoryboard(ctx, in.ProjectID, in.EpisodeID, sceneIDs)
	if err != nil {
		return nil, fmt.Errorf("GenerateStoryboard: %w", err)
	}
	return shotIDs, nil
}

func (e *Engine) media(ctx context.Context, r *run, in workflow.ProduceEpisodeInput, shotIDs []string, res *workflow.ProduceEpisodeResult) error {
	r.setStage("video_generation")
	if !r.ok(ctx) {
		return nil
	}
	assetIDs, err := e.act.GenerateVideoShots(ctx, in.ProjectID, in.EpisodeID, shotIDs)
	if err != nil {
		return fmt.Errorf("GenerateVideoShots: %w", err)
	}
	res.GeneratedShots = assetIDs

	r.setStage("assembly")
	if !r.ok(ctx) {
		return nil
	}
	url, err := e.act.AssembleEpisode(ctx, in.ProjectID, in.EpisodeID, shotIDs)
	if err != nil {
		return fmt.Errorf("AssembleEpisode: %w", err)
	}
	res.RenderURL = url
	return nil
}

func (e *Engine) Signal(_ context.Context, workflowID, signal string, payload interface{}) error {
	e.mu.Lock()
	r, ok := e.runs[workflowID]
	e.mu.Unlock()
	if !ok {
		return fmt.Errorf("unknown workflow %s", workflowID)
	}
	switch signal {
	case workflow.SignalApproval:
		r.approval <- fmt.Sprint(payload)
	case workflow.SignalPause:
		r.pause <- struct{}{}
	case workflow.SignalResume:
		r.resume <- struct{}{}
	case workflow.SignalStop:
		r.stop <- struct{}{}
	default:
		return fmt.Errorf("unknown signal %q", signal)
	}
	return nil
}

func (e *Engine) Cancel(_ context.Context, workflowID string) error {
	e.mu.Lock()
	r, ok := e.runs[workflowID]
	e.mu.Unlock()
	if !ok {
		return fmt.Errorf("unknown workflow %s", workflowID)
	}
	r.cancel()
	<-r.done
	return nil
}

func (e *Engine) Status(_ context.Context, workflowID string) (*workflow.WorkflowStatus, error) {
	e.mu.Lock()
	r, ok := e.runs[workflowID]
	e.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("unknown workflow %s", workflowID)
	}
	st := r.snapshot()
	return &st, nil
}

func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.done {
		return nil
	}
	e.done = true
	for _, r := range e.runs {
		r.cancel()
	}
	return nil
}
