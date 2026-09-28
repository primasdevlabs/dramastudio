package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"dramastudio/internal/platform/events"
	"dramastudio/internal/projects/domain"
)

type ProjectService struct {
	repo   domain.ProjectRepository
	events *events.Bus // may be nil; set via SetEvents
}

func NewProjectService(repo domain.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// SetEvents injects the domain event bus (§48). Nil-safe emitter.
func (s *ProjectService) SetEvents(b *events.Bus) {
	s.events = b
}

func (s *ProjectService) CreateProject(ctx context.Context, orgID, name, description, genre, language string, mode domain.ProductionMode) (*domain.Project, error) {
	id := domain.ProjectID("proj_" + uuid.NewString())
	project := domain.NewProject(id, orgID, name, description, genre, language, mode)
	if err := s.repo.Save(ctx, project); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.ProjectCreated, string(id), string(id),
		fmt.Sprintf("Project %q created", name))
	return project, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context, orgID string) ([]*domain.Project, error) {
	return s.repo.ListAll(ctx, orgID)
}

// UpdateSettings applies a partial settings/policy update.
func (s *ProjectService) UpdateProject(ctx context.Context, id domain.ProjectID, apply func(*domain.Project) error) (*domain.Project, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := apply(p); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// TransitionStatus moves the project through the §32 state machine.
func (s *ProjectService) TransitionStatus(ctx context.Context, id domain.ProjectID, to domain.ProjectStatus) (*domain.Project, error) {
	p, err := s.UpdateProject(ctx, id, func(p *domain.Project) error {
		return p.Transition(to)
	})
	if err == nil {
		s.events.Emit(ctx, events.ProjectUpdated, string(id), string(id),
			fmt.Sprintf("Project status → %s", to))
	}
	return p, err
}

// RecordSpend adds provider cost to the project budget; callers should have
// checked CanSpend before invoking generation (§85).
func (s *ProjectService) RecordSpend(ctx context.Context, id domain.ProjectID, cost float64) error {
	_, err := s.UpdateProject(ctx, id, func(p *domain.Project) error {
		p.Budget.CurrentSpent += cost
		return nil
	})
	return err
}

// BudgetCheck fails when the project budget cannot cover estimated cost.
func (s *ProjectService) BudgetCheck(ctx context.Context, id domain.ProjectID, estimate float64) error {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if !p.Budget.CanSpend(estimate) {
		return fmt.Errorf("budget exceeded: limit %.2f, spent %.2f, requested %.2f",
			p.Budget.TotalBudget, p.Budget.CurrentSpent, estimate)
	}
	return nil
}

// SaveSeriesBible stores a new immutable bible version (§11).
func (s *ProjectService) SaveSeriesBible(ctx context.Context, projectID domain.ProjectID, premise, genre, tone, visualDirection, dialogueStyle string, themes, worldRules, narrativeRules []string) (*domain.SeriesBible, error) {
	if _, err := s.repo.FindByID(ctx, projectID); err != nil {
		return nil, err
	}
	version := 1
	if existing, err := s.repo.GetLatestBible(ctx, projectID); err == nil && existing != nil {
		version = existing.Version + 1
	}
	id := fmt.Sprintf("bible_%s_v%d", projectID, version)
	bible := domain.NewSeriesBible(id, projectID, version, premise)
	bible.Genre = genre
	bible.Tone = tone
	bible.Themes = themes
	bible.WorldRules = worldRules
	bible.NarrativeRules = narrativeRules
	bible.VisualDirection = visualDirection
	bible.DialogueStyle = dialogueStyle
	if err := s.repo.SaveBible(ctx, bible); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.BibleUpdated, string(projectID), id,
		fmt.Sprintf("Series bible v%d saved", version))
	return bible, nil
}

func (s *ProjectService) GetLatestBible(ctx context.Context, projectID domain.ProjectID) (*domain.SeriesBible, error) {
	return s.repo.GetLatestBible(ctx, projectID)
}

func (s *ProjectService) GetBibleVersion(ctx context.Context, projectID domain.ProjectID, version int) (*domain.SeriesBible, error) {
	return s.repo.GetBibleVersion(ctx, projectID, version)
}

func (s *ProjectService) ListBibleVersions(ctx context.Context, projectID domain.ProjectID) ([]*domain.SeriesBible, error) {
	return s.repo.ListBibleVersions(ctx, projectID)
}
