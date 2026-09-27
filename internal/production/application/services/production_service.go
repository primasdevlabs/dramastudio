package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/production/domain"
)

type ProductionService struct {
	repo domain.ProductionRepository
}

func NewProductionService(repo domain.ProductionRepository) *ProductionService {
	return &ProductionService{repo: repo}
}

func (s *ProductionService) StartProductionJob(ctx context.Context, shotID string) (*domain.ProductionJob, error) {
	id := fmt.Sprintf("job_%d", time.Now().UnixNano())
	job := &domain.ProductionJob{
		ID:     id,
		ShotID: shotID,
		Status: domain.JobStatusRunning,
	}
	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *ProductionService) ListJobs(ctx context.Context, projectID string) ([]*domain.ProductionJob, error) {
	return s.repo.ListJobsByProject(ctx, projectID)
}

func (s *ProductionService) SubmitApproval(ctx context.Context, projectID, episodeID, stage, targetID string, decision domain.ApprovalDecision, notes, decidedBy string) (*domain.ApprovalRequest, error) {
	id := fmt.Sprintf("app_%d", time.Now().UnixNano())
	now := time.Now().UTC()
	req := &domain.ApprovalRequest{
		ID:          id,
		ProjectID:   projectID,
		EpisodeID:   episodeID,
		Stage:       stage,
		TargetID:    targetID,
		Decision:    decision,
		Notes:       notes,
		DecidedBy:   decidedBy,
		SubmittedAt: now,
		DecidedAt:   now,
	}
	if err := s.repo.SaveApproval(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *ProductionService) ListApprovals(ctx context.Context, projectID string) ([]*domain.ApprovalRequest, error) {
	return s.repo.ListApprovalsByProject(ctx, projectID)
}
