package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/continuity/domain"
	"dramastudio/internal/platform/events"
)

// CheckInput carries the data a checker needs for one run. Callers
// assemble it from the owning bounded context — continuity never imports
// other contexts' domain packages.
type CheckInput struct {
	ProjectID string
	EpisodeID string
	SceneID   string
	CheckType domain.CheckType
	TargetID  string

	// checker inputs — populate only what the check type needs
	Events            []domain.TimelineEvent // timeline
	Facts             []FactView             // story
	Claims            []FactView             // story
	CanonicalAttrs    map[string]string      // character / visual
	ActualAttrs       map[string]string      // character / visual
	EntityID          string                 // character / visual
	CanonicalWardrobe []string               // wardrobe
	AssignedWardrobe  []string               // wardrobe
}

type ContinuityService struct {
	repo      domain.ContinuityRepository
	timeline  *TimelineChecker
	story     *StoryChecker
	character *CharacterChecker
	visual    *VisualChecker
	wardrobe  *WardrobeChecker
	events    *events.Bus // may be nil; set via SetEvents
}

// SetEvents injects the domain event bus (§48). Nil-safe emitter.
func (s *ContinuityService) SetEvents(b *events.Bus) {
	s.events = b
}

func NewContinuityService(repo domain.ContinuityRepository) *ContinuityService {
	return &ContinuityService{
		repo:      repo,
		timeline:  NewTimelineChecker(),
		story:     NewStoryChecker(),
		character: NewCharacterChecker(),
		visual:    NewVisualChecker(),
		wardrobe:  NewWardrobeChecker(),
	}
}

// RunCheck executes a check suite and persists the findings (§38).
func (s *ContinuityService) RunCheck(ctx context.Context, in CheckInput) (*domain.ContinuityCheck, []*domain.ContinuityIssue, error) {
	check := &domain.ContinuityCheck{
		ID:        "chk_" + uuid.NewString(),
		ProjectID: in.ProjectID,
		EpisodeID: in.EpisodeID,
		CheckType: in.CheckType,
		TargetID:  in.TargetID,
		Status:    domain.CheckRunning,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.SaveCheck(ctx, check); err != nil {
		return nil, nil, err
	}

	var violations []domain.Violation
	switch in.CheckType {
	case domain.CheckTimeline:
		violations = s.timeline.Check(in.Events)
	case domain.CheckStory:
		violations = s.story.Check(in.Facts, in.Claims)
	case domain.CheckCharacter:
		violations = s.character.Check(in.EntityID, in.CanonicalAttrs, in.ActualAttrs)
	case domain.CheckVisual:
		violations = s.visual.Check(in.EntityID, in.CanonicalAttrs, in.ActualAttrs)
	case domain.CheckWardrobe:
		violations = s.wardrobe.Check(in.EntityID, in.CanonicalWardrobe, in.AssignedWardrobe)
	}

	issues := make([]*domain.ContinuityIssue, 0, len(violations))
	for _, v := range violations {
		issue := &domain.ContinuityIssue{
			ID:            "iss_" + uuid.NewString(),
			CheckID:       check.ID,
			ProjectID:     in.ProjectID,
			EpisodeID:     in.EpisodeID,
			SceneID:       in.SceneID,
			Category:      string(in.CheckType),
			Severity:      v.Severity,
			Entity:        v.Entity,
			Cause:         v.RuleName,
			Evidence:      v.Evidence,
			ExpectedState: "",
			ActualState:   v.Description,
			Status:        domain.IssueOpen,
			CreatedAt:     time.Now().UTC(),
		}
		if err := s.repo.SaveIssue(ctx, issue); err != nil {
			return nil, nil, err
		}
		s.events.Emit(ctx, events.ContinuityIssueDetected, in.ProjectID, issue.ID,
			fmt.Sprintf("Continuity %s: %s (%s)", v.Severity, v.Description, v.Entity))
		issues = append(issues, issue)
	}

	now := time.Now().UTC()
	check.Status = domain.CheckFinished
	check.FinishedAt = &now
	check.IssueCount = len(issues)
	if err := s.repo.SaveCheck(ctx, check); err != nil {
		return nil, nil, err
	}
	return check, issues, nil
}

// ReportIssue records a manually-raised finding.
func (s *ContinuityService) ReportIssue(ctx context.Context, issue *domain.ContinuityIssue) (*domain.ContinuityIssue, error) {
	if issue.ID == "" {
		issue.ID = "iss_" + uuid.NewString()
	}
	if issue.Status == "" {
		issue.Status = domain.IssueOpen
	}
	if issue.CreatedAt.IsZero() {
		issue.CreatedAt = time.Now().UTC()
	}
	if err := s.repo.SaveIssue(ctx, issue); err != nil {
		return nil, err
	}
	return issue, nil
}

func (s *ContinuityService) ListIssues(ctx context.Context, projectID string, f domain.IssueFilter) ([]*domain.ContinuityIssue, error) {
	return s.repo.ListIssues(ctx, projectID, f)
}

func (s *ContinuityService) ListChecks(ctx context.Context, projectID, episodeID string) ([]*domain.ContinuityCheck, error) {
	return s.repo.ListChecks(ctx, projectID, episodeID)
}

func (s *ContinuityService) GetIssue(ctx context.Context, id string) (*domain.ContinuityIssue, error) {
	return s.repo.FindIssueByID(ctx, id)
}

// ResolveIssue marks an issue resolved or wontfix (human control).
func (s *ContinuityService) ResolveIssue(ctx context.Context, id, resolution string, wontFix bool) (*domain.ContinuityIssue, error) {
	i, err := s.repo.FindIssueByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	i.Resolution = resolution
	i.ResolvedAt = &now
	if wontFix {
		i.Status = domain.IssueWontFix
	} else {
		i.Status = domain.IssueResolved
	}
	if err := s.repo.SaveIssue(ctx, i); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.ContinuityIssueResolved, i.ProjectID, i.ID,
		fmt.Sprintf("Continuity issue resolved: %s", resolution))
	return i, nil
}

// --- Timeline ---

func (s *ContinuityService) AddTimelineEvent(ctx context.Context, e *domain.TimelineEvent) (*domain.TimelineEvent, error) {
	if e.ID == "" {
		e.ID = "evt_" + uuid.NewString()
	}
	if err := s.repo.SaveTimelineEvent(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *ContinuityService) ListTimelineEvents(ctx context.Context, projectID, episodeID string) ([]*domain.TimelineEvent, error) {
	return s.repo.ListTimelineEvents(ctx, projectID, episodeID)
}
