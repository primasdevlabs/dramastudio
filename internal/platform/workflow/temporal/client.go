package temporal

import (
	"context"
	"crypto/tls"
	"fmt"

	"go.temporal.io/sdk/client"

	"dramastudio/internal/platform/workflow"
)

// Config mirrors configs.TemporalConfig without importing the root config
// package (keeps platform dependency-free).
type Config struct {
	HostPort  string
	Namespace string
	APIKey    string
	UseTLS    bool
}

// Dial connects to a Temporal server.
func Dial(_ context.Context, cfg Config) (client.Client, error) {
	opts := client.Options{
		HostPort:  cfg.HostPort,
		Namespace: cfg.Namespace,
	}
	if cfg.APIKey != "" {
		opts.Credentials = client.NewAPIKeyStaticCredentials(cfg.APIKey)
		opts.ConnectionOptions.TLS = &tls.Config{MinVersion: tls.VersionTLS12}
	} else if cfg.UseTLS {
		opts.ConnectionOptions.TLS = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	c, err := client.Dial(client.Options(opts))
	if err != nil {
		return nil, fmt.Errorf("temporal: dial %s/%s: %w", cfg.HostPort, cfg.Namespace, err)
	}
	return c, nil
}

// Orchestrator is the Temporal adapter for the workflow.Engine port (§83).
// Application code depends on workflow.Engine, never on this type.
type Orchestrator struct {
	client client.Client
}

func NewOrchestrator(c client.Client) *Orchestrator {
	return &Orchestrator{client: c}
}

// Start launches ProduceEpisode on the core queue with the run id in the
// workflow id (idempotent restart-safe, §61).
func (o *Orchestrator) Start(ctx context.Context, in workflow.ProduceEpisodeInput) (string, error) {
	id := "produce-episode-" + in.RunID
	opts := client.StartWorkflowOptions{
		ID:                       id,
		TaskQueue:                workflow.QueueCore,
		WorkflowIDReusePolicy:    3, // WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE
		WorkflowExecutionTimeout: 0, // unbounded; human gates can take days
	}
	if _, err := o.client.ExecuteWorkflow(ctx, opts, workflow.ProduceEpisode, in); err != nil {
		return "", err
	}
	return id, nil
}

func (o *Orchestrator) Signal(ctx context.Context, workflowID, signal string, payload interface{}) error {
	return o.client.SignalWorkflow(ctx, workflowID, "", signal, payload)
}

func (o *Orchestrator) Cancel(ctx context.Context, workflowID string) error {
	return o.client.CancelWorkflow(ctx, workflowID, "")
}

// Status implements workflow.Engine.Status via the Temporal query API.
func (o *Orchestrator) Status(ctx context.Context, workflowID string) (*workflow.WorkflowStatus, error) {
	resp, err := o.client.QueryWorkflow(ctx, workflowID, "", workflow.QueryStatus)
	if err != nil {
		return nil, err
	}
	var st workflow.WorkflowStatus
	if err := resp.Get(&st); err != nil {
		return nil, err
	}
	return &st, nil
}

func (o *Orchestrator) Close() error {
	o.client.Close()
	return nil
}
