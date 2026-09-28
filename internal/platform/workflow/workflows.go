package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"dramastudio/internal/platform/workflow/activities"
)

// controlGate owns the pause/resume/stop signal channels and the shared
// status snapshot between workflow stages (§39).
type controlGate struct {
	ctx        workflow.Context
	st         *WorkflowStatus
	approvalCh workflow.ReceiveChannel
	pauseCh    workflow.ReceiveChannel
	resumeCh   workflow.ReceiveChannel
	stopCh     workflow.ReceiveChannel
}

func newControlGate(ctx workflow.Context, st *WorkflowStatus) *controlGate {
	return &controlGate{
		ctx:        ctx,
		st:         st,
		approvalCh: workflow.GetSignalChannel(ctx, SignalApproval),
		pauseCh:    workflow.GetSignalChannel(ctx, SignalPause),
		resumeCh:   workflow.GetSignalChannel(ctx, SignalResume),
		stopCh:     workflow.GetSignalChannel(ctx, SignalStop),
	}
}

// drain consumes pending stop/pause signals without blocking.
func (g *controlGate) drain() {
	var sig string
	for g.stopCh.ReceiveAsync(&sig) {
		g.st.Stopped = true
	}
	for g.pauseCh.ReceiveAsync(&sig) {
		g.st.Paused = true
	}
}

// waitWhilePaused blocks on resume/stop while the run is paused.
func (g *controlGate) waitWhilePaused() {
	for g.st.Paused && !g.st.Stopped {
		sel := workflow.NewSelector(g.ctx)
		sel.AddReceive(g.resumeCh, func(c workflow.ReceiveChannel, more bool) {
			var sig string
			c.Receive(g.ctx, &sig)
			g.st.Paused = false
		})
		sel.AddReceive(g.stopCh, func(c workflow.ReceiveChannel, more bool) {
			var sig string
			c.Receive(g.ctx, &sig)
			g.st.Stopped = true
		})
		sel.Select(g.ctx)
	}
}

// ok runs the control gate; false means the run was stopped.
func (g *controlGate) ok() bool {
	g.drain()
	g.waitWhilePaused()
	return !g.st.Stopped
}

// setStage records the stage and drains pending control signals.
func (g *controlGate) setStage(s string) { g.st.Stage = s; g.drain() }

// awaitDecision suspends on the approval signal (or stop) and returns the
// human decision; empty string means stopped.
func (g *controlGate) awaitDecision() string {
	var decision string
	for decision == "" && !g.st.Stopped {
		sel := workflow.NewSelector(g.ctx)
		sel.AddReceive(g.approvalCh, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(g.ctx, &decision)
			g.st.Paused = false
		})
		sel.AddReceive(g.stopCh, func(c workflow.ReceiveChannel, more bool) {
			var sig string
			c.Receive(g.ctx, &sig)
			g.st.Stopped = true
		})
		sel.Select(g.ctx)
	}
	return decision
}

func episodeActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
}

// ProduceEpisode is the durable episode production workflow. It is
// deterministic: all side effects run inside activities; control arrives
// only via signals; state is queryable (§41, §83).
func ProduceEpisode(ctx workflow.Context, in ProduceEpisodeInput) (*ProduceEpisodeResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ProduceEpisode started", "project", in.ProjectID, "episode", in.EpisodeID)

	st := WorkflowStatus{RunID: in.RunID, Stage: "starting"}
	if err := registerQueries(ctx, &st); err != nil {
		return nil, err
	}
	ctx = workflow.WithActivityOptions(ctx, episodeActivityOptions())
	g := newControlGate(ctx, &st)
	var act *activities.EpisodeActivities

	res, err := runProductionStages(ctx, g, act, in)
	if err != nil {
		return nil, err
	}
	if g.st.Stopped {
		res.Status = "stopped"
		return res, nil
	}

	// Stage 7: human approval gate — suspends on the signal channel until
	// a decision arrives (§39).
	g.setStage("awaiting_approval")
	var approvalID string
	if err := workflow.ExecuteActivity(ctx, act.RequestHumanApproval, in.ProjectID, in.EpisodeID, "episode", in.EpisodeID).Get(ctx, &approvalID); err != nil {
		return nil, fmt.Errorf("RequestHumanApproval: %w", err)
	}
	st.ApprovalID = approvalID
	st.Paused = true

	res.Approval = g.awaitDecision()
	if g.st.Stopped {
		res.Status = "stopped"
		return res, nil
	}
	if res.Approval != "APPROVE" && res.Approval != "OVERRIDE" {
		res.Status = "rejected"
		return res, nil
	}

	g.setStage("completed")
	res.Status = "completed"
	logger.Info("ProduceEpisode completed", "episode", in.EpisodeID)
	return res, nil
}

// runProductionStages executes stages 1-6 (plan → continuity → script →
// storyboard → video → assembly), filling res as it goes.
func runProductionStages(ctx workflow.Context, g *controlGate, act *activities.EpisodeActivities, in ProduceEpisodeInput) (*ProduceEpisodeResult, error) {
	res := &ProduceEpisodeResult{EpisodeID: in.EpisodeID}

	var plan string
	if g.st.Stopped || !g.ok() {
		return res, nil
	}
	g.setStage("episode_plan")
	if err := workflow.ExecuteActivity(ctx, act.GenerateEpisodePlan, in.ProjectID, in.EpisodeID).Get(ctx, &plan); err != nil {
		return nil, fmt.Errorf("GenerateEpisodePlan: %w", err)
	}
	res.Plan = plan

	g.setStage("continuity_check")
	if !g.ok() {
		return res, nil
	}
	var blocking int
	if err := workflow.ExecuteActivity(ctx, act.ValidateContinuity, in.ProjectID, in.EpisodeID).Get(ctx, &blocking); err != nil {
		return nil, fmt.Errorf("ValidateContinuity: %w", err)
	}
	if blocking > 0 {
		return nil, fmt.Errorf("continuity check found %d blocking issues", blocking)
	}

	if err := runWritingStage(ctx, g, act, in); err != nil || g.st.Stopped {
		return res, err
	}
	shotIDs, err := runPlanningStages(ctx, g, act, in)
	if err != nil || g.st.Stopped {
		return res, err
	}
	if err := runMediaStages(ctx, g, act, in, shotIDs, res); err != nil {
		return res, err
	}
	return res, nil
}

func runWritingStage(ctx workflow.Context, g *controlGate, act *activities.EpisodeActivities, in ProduceEpisodeInput) error {
	g.setStage("script")
	if !g.ok() {
		return nil
	}
	var script string
	if err := workflow.ExecuteActivity(ctx, act.GenerateScript, in.ProjectID, in.EpisodeID).Get(ctx, &script); err != nil {
		return fmt.Errorf("GenerateScript: %w", err)
	}
	var dialogue string
	if err := workflow.ExecuteActivity(ctx, act.GenerateDialogue, in.ProjectID, in.EpisodeID).Get(ctx, &dialogue); err != nil {
		return fmt.Errorf("GenerateDialogue: %w", err)
	}
	return nil
}

func runPlanningStages(ctx workflow.Context, g *controlGate, act *activities.EpisodeActivities, in ProduceEpisodeInput) ([]string, error) {
	g.setStage("scenes")
	if !g.ok() {
		return nil, nil
	}
	var sceneIDs []string
	if err := workflow.ExecuteActivity(ctx, act.GenerateScenePlans, in.EpisodeID).Get(ctx, &sceneIDs); err != nil {
		return nil, fmt.Errorf("GenerateScenePlans: %w", err)
	}
	g.setStage("storyboard")
	var shotIDs []string
	if err := workflow.ExecuteActivity(ctx, act.GenerateStoryboard, in.ProjectID, in.EpisodeID, sceneIDs).Get(ctx, &shotIDs); err != nil {
		return nil, fmt.Errorf("GenerateStoryboard: %w", err)
	}
	return shotIDs, nil
}

// runMediaStages submits shot generation and assembles the render.
func runMediaStages(ctx workflow.Context, g *controlGate, act *activities.EpisodeActivities, in ProduceEpisodeInput, shotIDs []string, res *ProduceEpisodeResult) error {
	g.setStage("video_generation")
	if !g.ok() {
		return nil
	}
	var assetIDs []string
	if err := workflow.ExecuteActivity(ctx, act.GenerateVideoShots, in.ProjectID, in.EpisodeID, shotIDs).Get(ctx, &assetIDs); err != nil {
		return fmt.Errorf("GenerateVideoShots: %w", err)
	}
	res.GeneratedShots = assetIDs

	g.setStage("assembly")
	if !g.ok() {
		return nil
	}
	if err := workflow.ExecuteActivity(ctx, act.AssembleEpisode, in.ProjectID, in.EpisodeID, shotIDs).Get(ctx, &res.RenderURL); err != nil {
		return fmt.Errorf("AssembleEpisode: %w", err)
	}
	return nil
}

func registerQueries(ctx workflow.Context, st *WorkflowStatus) error {
	if err := workflow.SetQueryHandler(ctx, QueryStatus, func() (WorkflowStatus, error) { return *st, nil }); err != nil {
		return err
	}
	return workflow.SetQueryHandler(ctx, QueryStage, func() (string, error) { return st.Stage, nil })
}
