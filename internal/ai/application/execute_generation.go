package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/ai/domain"
	"dramastudio/internal/platform/ai/capability"
	"dramastudio/internal/platform/ai/routing"
)

// ExecuteGenerationCommand records and (for sync providers) executes one
// AI generation as a durable, auditable job (§26).
type ExecuteGenerationCommand struct {
	ProjectID      string
	Capability     domain.AICapability
	Scope          string // "project:<id>", "episode:<id>", "system"
	Input          json.RawMessage
	Prompt         string
	IdempotencyKey string
}

type ExecuteGenerationHandler struct {
	repo     domain.ModelRegistryRepository
	resolver *routing.Resolver
}

func NewExecuteGenerationHandler(repo domain.ModelRegistryRepository, resolver *routing.Resolver) *ExecuteGenerationHandler {
	return &ExecuteGenerationHandler{repo: repo, resolver: resolver}
}

func (h *ExecuteGenerationHandler) Handle(ctx context.Context, cmd ExecuteGenerationCommand) (*domain.GenerationJob, error) {
	// Idempotent submission (§61).
	if cmd.IdempotencyKey != "" {
		if existing, err := h.repo.FindGenerationJobByIdempotencyKey(ctx, cmd.IdempotencyKey); err == nil {
			return existing, nil
		}
	}

	job := &domain.GenerationJob{
		ID:             "genj_" + uuid.NewString(),
		ProjectID:      cmd.ProjectID,
		Capability:     cmd.Capability,
		Input:          cmd.Input,
		Status:         domain.StatusPending,
		Attempt:        1,
		IdempotencyKey: cmd.IdempotencyKey,
		CreatedAt:      time.Now().UTC(),
	}

	res, err := h.resolver.Resolve(ctx, string(cmd.Capability), cmd.Scope)
	if err != nil {
		h.fail(ctx, job, err)
		return nil, err
	}

	spec := capability.GenerationSpec{
		Capability: string(cmd.Capability),
		ProjectID:  cmd.ProjectID,
		Prompt:     cmd.Prompt,
	}

	// Try candidates in the policy's routing order. A retryable provider
	// failure moves to the next candidate; permanent failures stop.
	var lastErr error
	for i, cand := range res.Candidates {
		// Provenance: every attempt records the concrete provider + model +
		// version so outputs are reproducible (§26).
		job.ProviderID = cand.ProviderID
		job.ModelID = domain.ModelID(cand.ModelID)
		job.ModelVersion = cand.ModelVersion
		job.Attempt = i + 1

		now := time.Now().UTC()
		job.StartedAt = &now
		job.Status = domain.StatusRunning

		dispatchErr := h.dispatch(ctx, job, cand, spec)
		if dispatchErr == nil {
			if err := h.repo.SaveGenerationJob(ctx, job); err != nil {
				return nil, err
			}
			return job, nil
		}
		lastErr = dispatchErr
		if !isRetryable(dispatchErr) {
			break
		}
	}

	h.fail(ctx, job, fmt.Errorf("all %d candidates failed, last: %w", len(res.Candidates), lastErr))
	return nil, lastErr
}

// dispatch routes to the sync or async surface of one candidate.
func (h *ExecuteGenerationHandler) dispatch(ctx context.Context, job *domain.GenerationJob, cand routing.BoundCandidate, spec capability.GenerationSpec) error {
	switch {
	case cand.Adapter.Sync != nil:
		return h.dispatchSync(ctx, job, cand, spec)
	case cand.Adapter.Async != nil:
		return h.dispatchAsync(ctx, job, cand, spec)
	default:
		return &capability.ProviderError{Kind: capability.ErrUnavailable, Provider: cand.ProviderID, Message: "resolved provider has no registered adapter"}
	}
}

// dispatchSync runs a synchronous generator and records its output inline.
func (h *ExecuteGenerationHandler) dispatchSync(ctx context.Context, job *domain.GenerationJob, cand routing.BoundCandidate, spec capability.GenerationSpec) error {
	result, err := cand.Adapter.Sync.Generate(ctx, cand.ModelIdentifier, spec)
	if err != nil {
		return err
	}
	out := map[string]interface{}{
		"text":   result.OutputText,
		"uri":    result.OutputURI,
		"result": json.RawMessage(result.OutputJSON),
	}
	job.Output, _ = json.Marshal(out)
	job.Cost = result.Cost
	job.Status = domain.StatusSucceeded
	done := time.Now().UTC()
	job.CompletedAt = &done
	return nil
}

// dispatchAsync submits to an async provider; completion arrives via poll
// or webhook and is reconciled elsewhere (§40).
func (h *ExecuteGenerationHandler) dispatchAsync(ctx context.Context, job *domain.GenerationJob, cand routing.BoundCandidate, spec capability.GenerationSpec) error {
	aj, err := cand.Adapter.Async.Submit(ctx, cand.ModelIdentifier, spec, "")
	if err != nil {
		return err
	}
	job.ProviderJobID = aj.ProviderJobID
	return nil
}

// isRetryable reports whether a dispatch failure should advance to the
// next policy candidate.
func isRetryable(err error) bool {
	var pe *capability.ProviderError
	if errors.As(err, &pe) {
		return pe.Retryable ||
			pe.Kind == capability.ErrRateLimited ||
			pe.Kind == capability.ErrUnavailable ||
			pe.Kind == capability.ErrTimeout
	}
	return false
}

// fail records the failed state best-effort before returning the error.
func (h *ExecuteGenerationHandler) fail(ctx context.Context, job *domain.GenerationJob, err error) {
	job.Status = domain.StatusFailed
	job.Error = err.Error()
	done := time.Now().UTC()
	job.CompletedAt = &done
	_ = h.repo.SaveGenerationJob(ctx, job)
}
