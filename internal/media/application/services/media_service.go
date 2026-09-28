package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/media/domain"
	aicontracts "dramastudio/internal/platform/ai/capability"
	"dramastudio/internal/platform/ai/routing"
)

// GenerateRequest is the application-layer input for a generation.
type GenerateRequest struct {
	ProjectID       string
	Capability      string // e.g. "image_generation", "video_generation"
	MediaType       domain.MediaType
	CharacterID     string
	LocationID      string
	EpisodeID       string
	SceneID         string
	ShotID          string
	Prompt          string
	Spec            map[string]interface{}
	Parameters      map[string]interface{}
	ReferenceAssets []string
	CallbackURL     string
}

type MediaService struct {
	repo     domain.MediaRepository
	resolver *routing.Resolver
}

func NewMediaService(repo domain.MediaRepository, resolver *routing.Resolver) *MediaService {
	return &MediaService{repo: repo, resolver: resolver}
}

// Generate submits a generation to the resolved provider. Sync providers
// complete inline; async providers return a running job finished later via
// webhook/poll (§61). The asset is created PENDING — never auto-approved.
func (s *MediaService) Generate(ctx context.Context, req GenerateRequest) (*domain.GenerationJob, *domain.Asset, error) {
	res, err := s.resolver.Resolve(ctx, req.Capability, "project:"+req.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	primary := res.Primary()
	if primary == nil {
		return nil, nil, fmt.Errorf("no bound provider for capability %q", req.Capability)
	}

	asset := &domain.Asset{
		ID:              "asset_" + uuid.NewString(),
		ProjectID:       req.ProjectID,
		Type:            req.MediaType,
		CharacterID:     req.CharacterID,
		LocationID:      req.LocationID,
		EpisodeID:       req.EpisodeID,
		SceneID:         req.SceneID,
		ShotID:          req.ShotID,
		Provider:        primary.ProviderID,
		Model:           primary.ModelID,
		Prompt:          req.Prompt,
		Spec:            req.Spec,
		Parameters:      req.Parameters,
		ReferenceAssets: req.ReferenceAssets,
		Status:          domain.AssetPending,
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.repo.SaveAsset(ctx, asset); err != nil {
		return nil, nil, err
	}

	job, err := s.submitJob(ctx, asset, req.Capability, req.Prompt, req.Parameters, req.Spec, res, req.CallbackURL)
	if err != nil {
		return nil, nil, err
	}
	return job, asset, nil
}

// submitJob creates a generation job bound to an asset and dispatches it
// through the policy's candidate chain. A retryable provider failure moves
// to the next candidate; the winning candidate's identity is recorded on
// both job and asset for provenance (§26).
func (s *MediaService) submitJob(ctx context.Context, asset *domain.Asset, capability, prompt string, params, spec map[string]interface{}, res *routing.ModelResolution, callbackURL string) (*domain.GenerationJob, error) {
	job := &domain.GenerationJob{
		ID:         "gen_" + uuid.NewString(),
		ProjectID:  asset.ProjectID,
		AssetID:    asset.ID,
		Capability: capability,
		Input:      prompt,
		Status:     domain.GenerationPending,
		Attempt:    1,
		StartedAt:  time.Now().UTC(),
	}

	genSpec := aicontracts.GenerationSpec{
		Capability:  capability,
		ProjectID:   asset.ProjectID,
		Prompt:      prompt,
		Parameters:  params,
		AspectRatio: strParam(spec, "aspect_ratio"),
		Continuity:  mapParam(spec, "continuity"),
		Subject:     strMapParam(spec, "subject"),
	}

	fail := func(err error) (*domain.GenerationJob, error) {
		job.Status = domain.GenerationFailed
		job.Error = err.Error()
		now := time.Now().UTC()
		job.CompletedAt = &now
		_ = s.repo.SaveJob(ctx, job)
		return nil, err
	}

	var lastErr error
	for i, cand := range res.Candidates {
		job.Provider = cand.ProviderID
		job.Model = cand.ModelID
		job.Attempt = i + 1
		asset.Provider = cand.ProviderID
		asset.Model = cand.ModelID
		switch {
		case cand.Adapter.Async != nil:
			aj, err := cand.Adapter.Async.Submit(ctx, cand.ModelIdentifier, genSpec, callbackURL)
			if err != nil {
				lastErr = err
				if !isRetryable(err) {
					return fail(err)
				}
				continue
			}
			job.ProviderJobID = aj.ProviderJobID
			job.Status = domain.GenerationRunning
			if aj.Status == "succeeded" && aj.Result != nil {
				s.completeJob(ctx, job, asset, aj.Result)
			}
		case cand.Adapter.Sync != nil:
			job.Status = domain.GenerationRunning
			result, err := cand.Adapter.Sync.Generate(ctx, cand.ModelIdentifier, genSpec)
			if err != nil {
				lastErr = err
				if !isRetryable(err) {
					return fail(err)
				}
				continue
			}
			s.completeJob(ctx, job, asset, result)
		default:
			lastErr = fmt.Errorf("provider %q has no generator registered", cand.ProviderID)
			continue
		}
		if err := s.repo.SaveJob(ctx, job); err != nil {
			return nil, err
		}
		return job, nil
	}
	return fail(lastErr)
}

// isRetryable reports whether a provider failure should advance to the
// next policy candidate (rate limits, timeouts, unavailability).
func isRetryable(err error) bool {
	var pe *aicontracts.ProviderError
	if errors.As(err, &pe) {
		return pe.Retryable ||
			pe.Kind == aicontracts.ErrRateLimited ||
			pe.Kind == aicontracts.ErrUnavailable ||
			pe.Kind == aicontracts.ErrTimeout
	}
	return false
}

func strParam(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func mapParam(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}

func strMapParam(m map[string]interface{}, key string) map[string]string {
	if m == nil {
		return nil
	}
	raw, ok := m[key].(map[string]interface{})
	if !ok {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}

// completeJob records a finished provider result: new immutable asset
// version + job success. The asset itself stays PENDING for human review.
func (s *MediaService) completeJob(ctx context.Context, job *domain.GenerationJob, asset *domain.Asset, res *aicontracts.Result) {
	now := time.Now().UTC()
	job.Status = domain.GenerationSucceeded
	job.OutputURL = res.OutputURI
	job.Cost = res.Cost
	job.CompletedAt = &now

	asset.Version++
	asset.Cost += res.Cost
	v := &domain.AssetVersion{
		ID:        "av_" + uuid.NewString(),
		AssetID:   asset.ID,
		Version:   asset.Version,
		ObjectKey: res.OutputURI,
		URL:       res.OutputURI,
		Status:    domain.AssetPending,
		Cost:      res.Cost,
		CreatedAt: now,
	}
	_ = s.repo.SaveAssetVersion(ctx, v)
	_ = s.repo.SaveAsset(ctx, asset)
}

// HandleProviderCallback applies a verified provider webhook/poll result
// to the tracked generation job and its asset (§61).
func (s *MediaService) HandleProviderCallback(ctx context.Context, providerJobID, status, outputURL, errMsg string, cost float64) (*domain.GenerationJob, error) {
	job, err := s.repo.FindJobByProviderJobID(ctx, providerJobID)
	if err != nil {
		return nil, err
	}
	asset, err := s.repo.FindAssetByID(ctx, job.AssetID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	switch status {
	case "succeeded", "completed":
		s.completeJob(ctx, job, asset, &aicontracts.Result{
			Provider:  job.Provider,
			Model:     job.Model,
			OutputURI: outputURL,
			Cost:      cost,
		})
	case "failed", "error":
		job.Status = domain.GenerationFailed
		job.Error = errMsg
		job.CompletedAt = &now
	case "cancelled":
		job.Status = domain.GenerationCancelled
		job.CompletedAt = &now
	default:
		job.Status = domain.GenerationRunning
	}
	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

// setAssetStatus applies a human decision to an asset.
func (s *MediaService) setAssetStatus(ctx context.Context, id string, next domain.AssetStatus) (*domain.Asset, error) {
	a, err := s.repo.FindAssetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	a.Status = next
	if next == domain.AssetApproved {
		if v, err := s.repo.FindAssetVersion(ctx, id, a.Version); err == nil {
			v.Status = domain.AssetApproved
			_ = s.repo.SaveAssetVersion(ctx, v)
			a.URL = v.URL
		}
	}
	if err := s.repo.SaveAsset(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *MediaService) ApproveAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.setAssetStatus(ctx, id, domain.AssetApproved)
}

func (s *MediaService) RejectAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.setAssetStatus(ctx, id, domain.AssetRejected)
}

func (s *MediaService) ArchiveAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.setAssetStatus(ctx, id, domain.AssetArchived)
}

// Regenerate issues a new generation producing a new version on the same
// asset (human control: regenerate).
func (s *MediaService) Regenerate(ctx context.Context, assetID, callbackURL string) (*domain.GenerationJob, error) {
	asset, err := s.repo.FindAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	capability := capabilityForType(asset.Type)
	res, err := s.resolver.Resolve(ctx, capability, "project:"+asset.ProjectID)
	if err != nil {
		return nil, err
	}
	asset.Status = domain.AssetPending
	if err := s.repo.SaveAsset(ctx, asset); err != nil {
		return nil, err
	}
	return s.submitJob(ctx, asset, capability, asset.Prompt, asset.Parameters, asset.Spec, res, callbackURL)
}

func capabilityForType(t domain.MediaType) string {
	switch t {
	case domain.MediaTypeImage:
		return "image_generation"
	case domain.MediaTypeVideo:
		return "video_generation"
	case domain.MediaTypeVoice:
		return "voice"
	case domain.MediaTypeMusic:
		return "music_generation"
	case domain.MediaTypeSFX:
		return "sfx_generation"
	}
	return string(t)
}

func (s *MediaService) GetAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.repo.FindAssetByID(ctx, id)
}

func (s *MediaService) ListAssets(ctx context.Context, projectID string) ([]*domain.Asset, error) {
	return s.repo.ListAssetsByProject(ctx, projectID)
}

func (s *MediaService) ListVersions(ctx context.Context, assetID string) ([]*domain.AssetVersion, error) {
	return s.repo.ListAssetVersions(ctx, assetID)
}

func (s *MediaService) ListJobs(ctx context.Context, projectID string) ([]*domain.GenerationJob, error) {
	return s.repo.ListJobsByProject(ctx, projectID)
}

func (s *MediaService) GetJob(ctx context.Context, id string) (*domain.GenerationJob, error) {
	return s.repo.FindJobByID(ctx, id)
}
