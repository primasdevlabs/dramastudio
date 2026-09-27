package mock

import (
	"context"
	"fmt"
	"sync/atomic"

	"dramastudio/internal/platform/ai/contracts"
)

// Sync is a deterministic SyncGenerator for tests and offline development
// (Backend.md §80 — mocks satisfy the same contract as real providers).
type Sync struct {
	Name    string
	Counter int64
	// FailWith, when set, makes every call return this error.
	FailWith error
}

func NewSync(name string) *Sync {
	return &Sync{Name: name}
}

func (m *Sync) ProviderName() string { return m.Name }

func (m *Sync) Generate(ctx context.Context, model string, spec contracts.GenerationSpec) (*contracts.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if m.FailWith != nil {
		return nil, m.FailWith
	}
	n := atomic.AddInt64(&m.Counter, 1)
	res := &contracts.Result{
		Provider: m.Name,
		Model:    model,
		Cost:     0.001,
		Metadata: map[string]interface{}{"mock": true, "call": n},
	}
	if isTextCapability(spec.Capability) {
		res.OutputText = fmt.Sprintf("mock-%s-output-%d", spec.Capability, n)
	} else {
		res.OutputURI = fmt.Sprintf("mock://%s/%s/%s-%d", m.Name, spec.Capability, model, n)
	}
	return res, nil
}

func isTextCapability(cap string) bool {
	switch cap {
	case "story_bible", "story_architecture", "season_planning",
		"episode_planning", "script_writing", "dialogue_writing",
		"shot_planning", "continuity_analysis", "quality_evaluation",
		"storyboard_generation":
		return true
	}
	return false
}

// Async is a deterministic AsyncGenerator: Submit returns a job that
// succeeds on the next Poll unless FailWith is set.
type Async struct {
	Name    string
	Counter int64
	FailWith error
	jobs    chan *contracts.AsyncJob
}

func NewAsync(name string) *Async {
	return &Async{Name: name, jobs: make(chan *contracts.AsyncJob, 64)}
}

func (m *Async) ProviderName() string { return m.Name }

func (m *Async) Submit(ctx context.Context, model string, spec contracts.GenerationSpec, _ string) (*contracts.AsyncJob, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if m.FailWith != nil {
		return nil, m.FailWith
	}
	n := atomic.AddInt64(&m.Counter, 1)
	job := &contracts.AsyncJob{
		ProviderJobID: fmt.Sprintf("mockjob_%s_%d", m.Name, n),
		Status:        "pending",
	}
	job.Result = &contracts.Result{
		Provider:      m.Name,
		Model:         model,
		OutputURI:     fmt.Sprintf("mock://%s/%s-%d.mp4", m.Name, spec.Capability, n),
		Cost:          0.01,
		ProviderJobID: job.ProviderJobID,
	}
	m.jobs <- job
	return job, nil
}

func (m *Async) Poll(_ context.Context, providerJobID string) (*contracts.AsyncJob, error) {
	if m.FailWith != nil {
		return nil, m.FailWith
	}
	select {
	case job := <-m.jobs:
		if job.ProviderJobID == providerJobID {
			job.Status = "succeeded"
			return job, nil
		}
		m.jobs <- job
	default:
	}
	return &contracts.AsyncJob{ProviderJobID: providerJobID, Status: "running"}, nil
}

func (m *Async) Cancel(_ context.Context, providerJobID string) error {
	return nil
}
