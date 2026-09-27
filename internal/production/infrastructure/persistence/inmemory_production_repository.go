package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/production/domain"
)

type InMemoryProductionRepository struct {
	mu          sync.RWMutex
	productions map[string]*domain.Production // keyed by project_id
	runs        map[string]*domain.ProductionRun
	jobs        map[string]*domain.ProductionJob
	shots       map[string]*domain.Shot
	approvals   map[string]*domain.ApprovalRequest
}

func NewInMemoryProductionRepository() *InMemoryProductionRepository {
	return &InMemoryProductionRepository{
		productions: make(map[string]*domain.Production),
		runs:        make(map[string]*domain.ProductionRun),
		jobs:        make(map[string]*domain.ProductionJob),
		shots:       make(map[string]*domain.Shot),
		approvals:   make(map[string]*domain.ApprovalRequest),
	}
}

func (r *InMemoryProductionRepository) SaveProduction(ctx context.Context, p *domain.Production) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.productions[p.ProjectID] = p
	return nil
}

func (r *InMemoryProductionRepository) FindProductionByProject(ctx context.Context, projectID string) (*domain.Production, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.productions[projectID]
	if !ok {
		return nil, domain.ErrRunNotFound
	}
	return p, nil
}

func (r *InMemoryProductionRepository) SaveRun(ctx context.Context, run *domain.ProductionRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runs[run.ID] = run
	return nil
}

func (r *InMemoryProductionRepository) FindRunByID(ctx context.Context, id string) (*domain.ProductionRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	run, ok := r.runs[id]
	if !ok {
		return nil, domain.ErrRunNotFound
	}
	return run, nil
}

func (r *InMemoryProductionRepository) ListRunsByProject(ctx context.Context, projectID string) ([]*domain.ProductionRun, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ProductionRun, 0)
	for _, run := range r.runs {
		if run.ProjectID == projectID {
			res = append(res, run)
		}
	}
	return res, nil
}

func (r *InMemoryProductionRepository) SaveJob(ctx context.Context, job *domain.ProductionJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryProductionRepository) FindJobByID(ctx context.Context, id string) (*domain.ProductionJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[id]
	if !ok {
		return nil, domain.ErrJobNotFound
	}
	return j, nil
}

func (r *InMemoryProductionRepository) FindJobByIdempotencyKey(ctx context.Context, key string) (*domain.ProductionJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, j := range r.jobs {
		if j.IdempotencyKey == key && key != "" {
			return j, nil
		}
	}
	return nil, domain.ErrJobNotFound
}

func (r *InMemoryProductionRepository) ListJobsByProject(ctx context.Context, projectID string) ([]*domain.ProductionJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ProductionJob, 0)
	for _, j := range r.jobs {
		if j.ProjectID == projectID {
			res = append(res, j)
		}
	}
	return res, nil
}

func (r *InMemoryProductionRepository) ListJobsByRun(ctx context.Context, runID string) ([]*domain.ProductionJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ProductionJob, 0)
	for _, j := range r.jobs {
		if j.RunID == runID {
			res = append(res, j)
		}
	}
	return res, nil
}

func (r *InMemoryProductionRepository) SaveShot(ctx context.Context, shot *domain.Shot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shots[shot.ID] = shot
	return nil
}

func (r *InMemoryProductionRepository) FindShotByID(ctx context.Context, id string) (*domain.Shot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.shots[id]
	if !ok {
		return nil, domain.ErrShotNotFound
	}
	return s, nil
}

func (r *InMemoryProductionRepository) ListShots(ctx context.Context, episodeID, sceneID string) ([]*domain.Shot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Shot, 0)
	for _, s := range r.shots {
		if sceneID != "" && s.SceneID != sceneID {
			continue
		}
		if sceneID == "" && episodeID != "" && s.EpisodeID != episodeID {
			continue
		}
		res = append(res, s)
	}
	return res, nil
}

func (r *InMemoryProductionRepository) SaveApproval(ctx context.Context, app *domain.ApprovalRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.approvals[app.ID] = app
	return nil
}

func (r *InMemoryProductionRepository) FindApprovalByID(ctx context.Context, id string) (*domain.ApprovalRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.approvals[id]
	if !ok {
		return nil, domain.ErrApprovalNotFound
	}
	return a, nil
}

func (r *InMemoryProductionRepository) ListApprovalsByProject(ctx context.Context, projectID string, pendingOnly bool) ([]*domain.ApprovalRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.ApprovalRequest, 0)
	for _, a := range r.approvals {
		if a.ProjectID != projectID {
			continue
		}
		if pendingOnly && a.Decision != "" {
			continue
		}
		res = append(res, a)
	}
	return res, nil
}
