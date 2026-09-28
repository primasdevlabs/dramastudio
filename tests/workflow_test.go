package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"

	"dramastudio/internal/platform/workflow"
	"dramastudio/internal/platform/workflow/activities"
)

// TestProduceEpisodeWorkflow exercises the durable orchestration: staged
// activities, the human-approval signal gate, and the query surface.
func TestProduceEpisodeWorkflow(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()

	act := &activities.EpisodeActivities{} // method refs only — all calls mocked

	env.OnActivity(act.GenerateEpisodePlan, mock.Anything, "proj_1", "ep_1").Return("plan-text", nil)
	env.OnActivity(act.ValidateContinuity, mock.Anything, "proj_1", "ep_1").Return(0, nil)
	env.OnActivity(act.GenerateScript, mock.Anything, "proj_1", "ep_1").Return("script", nil)
	env.OnActivity(act.GenerateDialogue, mock.Anything, "proj_1", "ep_1").Return("dialogue", nil)
	env.OnActivity(act.GenerateScenePlans, mock.Anything, "ep_1").Return([]string{"sc_1", "sc_2"}, nil)
	env.OnActivity(act.GenerateStoryboard, mock.Anything, "proj_1", "ep_1", []string{"sc_1", "sc_2"}).
		Return([]string{"shot_1", "shot_2"}, nil)
	env.OnActivity(act.GenerateVideoShots, mock.Anything, "proj_1", "ep_1", []string{"shot_1", "shot_2"}).
		Return([]string{"asset_1", "asset_2"}, nil)
	env.OnActivity(act.AssembleEpisode, mock.Anything, "proj_1", "ep_1", []string{"shot_1", "shot_2"}).
		Return("https://files/ep_1.mp4", nil)
	env.OnActivity(act.RequestHumanApproval, mock.Anything, "proj_1", "ep_1", "episode", "ep_1").
		Return("appr_1", nil)

	// The approval gate suspends on the signal channel — approve after 1s.
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(workflow.SignalApproval, "APPROVE")
	}, time.Second)

	env.ExecuteWorkflow(workflow.ProduceEpisode, workflow.ProduceEpisodeInput{
		ProjectID: "proj_1", EpisodeID: "ep_1", RunID: "run_1", BibleVersion: 1,
	})

	if !env.IsWorkflowCompleted() {
		t.Fatal("workflow did not complete")
	}
	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("workflow error: %v", err)
	}
	var res workflow.ProduceEpisodeResult
	if err := env.GetWorkflowResult(&res); err != nil {
		t.Fatalf("get result: %v", err)
	}
	if res.Status != "completed" {
		t.Errorf("expected completed, got %s", res.Status)
	}
	if res.Approval != "APPROVE" {
		t.Errorf("expected APPROVE, got %s", res.Approval)
	}
	if len(res.GeneratedShots) != 2 {
		t.Errorf("expected 2 shot assets, got %d", len(res.GeneratedShots))
	}
	if res.RenderURL == "" {
		t.Error("expected render URL")
	}
}

// TestProduceEpisodeWorkflowRejected verifies a non-approve decision ends
// the run in rejected state without completing production.
func TestProduceEpisodeWorkflowRejected(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	act := &activities.EpisodeActivities{}

	env.OnActivity(act.GenerateEpisodePlan, mock.Anything, mock.Anything, mock.Anything).Return("plan", nil)
	env.OnActivity(act.ValidateContinuity, mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
	env.OnActivity(act.GenerateScript, mock.Anything, mock.Anything, mock.Anything).Return("s", nil)
	env.OnActivity(act.GenerateDialogue, mock.Anything, mock.Anything, mock.Anything).Return("d", nil)
	env.OnActivity(act.GenerateScenePlans, mock.Anything, mock.Anything).Return([]string{"sc_1"}, nil)
	env.OnActivity(act.GenerateStoryboard, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]string{"shot_1"}, nil)
	env.OnActivity(act.GenerateVideoShots, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]string{"asset_1"}, nil)
	env.OnActivity(act.AssembleEpisode, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("u", nil)
	env.OnActivity(act.RequestHumanApproval, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("appr_1", nil)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(workflow.SignalApproval, "REQUEST_REVISION")
	}, time.Second)

	env.ExecuteWorkflow(workflow.ProduceEpisode, workflow.ProduceEpisodeInput{
		ProjectID: "proj_1", EpisodeID: "ep_1", RunID: "run_2",
	})
	if !env.IsWorkflowCompleted() {
		t.Fatal("workflow did not complete")
	}
	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("workflow error: %v", err)
	}
	var res workflow.ProduceEpisodeResult
	if err := env.GetWorkflowResult(&res); err != nil {
		t.Fatalf("get result: %v", err)
	}
	if res.Status != "rejected" {
		t.Errorf("expected rejected, got %s", res.Status)
	}
}
