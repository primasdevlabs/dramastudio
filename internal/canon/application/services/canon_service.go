package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dramastudio/internal/canon/domain"
	"dramastudio/internal/platform/events"
)

type CanonService struct {
	repo   domain.CanonRepository
	events *events.Bus // may be nil; set via SetEvents
}

func NewCanonService(repo domain.CanonRepository) *CanonService {
	return &CanonService{repo: repo}
}

// SetEvents injects the domain event bus (§48). Nil-safe emitter.
func (s *CanonService) SetEvents(b *events.Bus) {
	s.events = b
}

// EstablishFact records a new canonical fact (§12) with its first version.
func (s *CanonService) EstablishFact(ctx context.Context, f *domain.StoryFact) (*domain.StoryFact, error) {
	if f.ID == "" {
		f.ID = "fact_" + uuid.NewString()
	}
	if f.Status == "" {
		f.Status = domain.FactStatusCanonical
	}
	if f.Confidence == 0 {
		f.Confidence = 1.0
	}
	f.Version = 1
	f.CreatedAt = time.Now().UTC()
	if err := s.repo.SaveFact(ctx, f); err != nil {
		return nil, err
	}
	raw, _ := json.Marshal(f)
	_ = s.repo.SaveFactVersion(ctx, f.ID, f.Version, raw)
	s.events.Emit(ctx, events.StoryFactEstablished, f.ProjectID, f.ID,
		fmt.Sprintf("Fact established: %s %s %s", f.Subject, f.Predicate, f.Object))
	return f, nil
}

// UpdateFact creates a new version of an existing fact (never overwrites
// history — see §72 regeneration/versioning principle).
func (s *CanonService) UpdateFact(ctx context.Context, id string, apply func(*domain.StoryFact) error) (*domain.StoryFact, error) {
	f, err := s.repo.FindFactByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := apply(f); err != nil {
		return nil, err
	}
	f.Version++
	if err := s.repo.SaveFact(ctx, f); err != nil {
		return nil, err
	}
	raw, _ := json.Marshal(f)
	_ = s.repo.SaveFactVersion(ctx, f.ID, f.Version, raw)
	return f, nil
}

// Retcon marks a fact as no longer canonical while preserving history.
// Scoped to projectID like GetFact.
func (s *CanonService) Retcon(ctx context.Context, projectID, id string) (*domain.StoryFact, error) {
	if _, err := s.GetFact(ctx, projectID, id); err != nil {
		return nil, err
	}
	f, err := s.UpdateFact(ctx, id, func(f *domain.StoryFact) error {
		f.Status = domain.FactStatusRetconned
		return nil
	})
	if err == nil {
		s.events.Emit(ctx, events.StoryFactChanged, f.ProjectID, f.ID,
			fmt.Sprintf("Fact retconned: %s %s %s", f.Subject, f.Predicate, f.Object))
	}
	return f, err
}

// GetFact returns the fact only when it belongs to projectID — cross-project
// IDs answer ErrFactNotFound so existence isn't leaked between projects.
func (s *CanonService) GetFact(ctx context.Context, projectID, id string) (*domain.StoryFact, error) {
	f, err := s.repo.FindFactByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f.ProjectID != projectID {
		return nil, domain.ErrFactNotFound
	}
	return f, nil
}

// ListFactVersions returns a fact's immutable change history (§72).
func (s *CanonService) ListFactVersions(ctx context.Context, projectID, id string) ([]*domain.FactVersion, error) {
	if _, err := s.GetFact(ctx, projectID, id); err != nil {
		return nil, err
	}
	return s.repo.ListFactVersions(ctx, id)
}

func (s *CanonService) ListFacts(ctx context.Context, projectID string) ([]*domain.StoryFact, error) {
	return s.repo.ListFacts(ctx, projectID)
}

func (s *CanonService) ListFactsForEntity(ctx context.Context, projectID, entityID string) ([]*domain.StoryFact, error) {
	return s.repo.ListFactsForEntity(ctx, projectID, entityID)
}

// GetKnowledge returns what a character knows at an episode; missing state
// means "knows nothing" — never an error, so callers can't leak canon (§13).
func (s *CanonService) GetKnowledge(ctx context.Context, characterID, episodeID string) (*domain.KnowledgeState, error) {
	ks, err := s.repo.GetKnowledgeState(ctx, characterID, episodeID)
	if err == domain.ErrKnowledgeMissing {
		return &domain.KnowledgeState{
			CharacterID:  characterID,
			EpisodeID:    episodeID,
			KnownFactIDs: []string{},
		}, nil
	}
	return ks, err
}

// GrantKnowledge adds facts to a character's knowledge at an episode. Every
// granted fact must exist and belong to projectID — foreign or dangling fact
// IDs are rejected instead of silently recorded.
func (s *CanonService) GrantKnowledge(ctx context.Context, projectID, characterID, episodeID string, factIDs []string) (*domain.KnowledgeState, error) {
	for _, id := range factIDs {
		if _, err := s.GetFact(ctx, projectID, id); err != nil {
			return nil, err
		}
	}
	ks, err := s.GetKnowledge(ctx, characterID, episodeID)
	if err != nil {
		return nil, err
	}
	known := make(map[string]bool, len(ks.KnownFactIDs))
	for _, id := range ks.KnownFactIDs {
		known[id] = true
	}
	for _, id := range factIDs {
		if !known[id] {
			ks.KnownFactIDs = append(ks.KnownFactIDs, id)
		}
	}
	if err := s.repo.SaveKnowledgeState(ctx, ks); err != nil {
		return nil, err
	}
	return ks, nil
}

// Knows reports whether a character knows a fact at an episode.
func (s *CanonService) Knows(ctx context.Context, characterID, episodeID, factID string) (bool, error) {
	ks, err := s.GetKnowledge(ctx, characterID, episodeID)
	if err != nil {
		return false, err
	}
	for _, id := range ks.KnownFactIDs {
		if id == factID {
			return true, nil
		}
	}
	return false, nil
}

func (s *CanonService) CreateRule(ctx context.Context, projectID, ruleText string, enforced bool) (*domain.CanonRule, error) {
	rule := &domain.CanonRule{
		ID:        "rule_" + uuid.NewString(),
		ProjectID: projectID,
		RuleText:  ruleText,
		Enforced:  enforced,
	}
	if err := s.repo.SaveRule(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *CanonService) ListRules(ctx context.Context, projectID string) ([]*domain.CanonRule, error) {
	return s.repo.ListRules(ctx, projectID)
}
