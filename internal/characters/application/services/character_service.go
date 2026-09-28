package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"dramastudio/internal/characters/domain"
	"dramastudio/internal/platform/events"
)

type CharacterService struct {
	repo   domain.CharacterRepository
	events *events.Bus // may be nil; set via SetEvents
}

func NewCharacterService(repo domain.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

// SetEvents injects the domain event bus (§48). Nil-safe emitter.
func (s *CharacterService) SetEvents(b *events.Bus) {
	s.events = b
}

func (s *CharacterService) CreateCharacter(ctx context.Context, projectID, name, role, bio string) (*domain.Character, error) {
	c := domain.NewCharacter(domain.CharacterID("char_"+uuid.NewString()), projectID, name, role)
	c.Bio = bio
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.CharacterCreated, projectID, string(c.ID),
		fmt.Sprintf("Character %q created (%s)", name, role))
	return c, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, id domain.CharacterID) (*domain.Character, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CharacterService) ListCharacters(ctx context.Context, projectID string) ([]*domain.Character, error) {
	return s.repo.List(ctx, projectID)
}

// UpdateCharacter mutates a character and persists a new immutable version.
// Locked characters reject modification (§15).
func (s *CharacterService) UpdateCharacter(ctx context.Context, id domain.CharacterID, apply func(*domain.Character) error) (*domain.Character, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.IsLocked {
		return nil, domain.ErrCharacterLocked
	}
	if err := apply(c); err != nil {
		return nil, err
	}
	c.Version++
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	s.events.Emit(ctx, events.CharacterUpdated, c.ProjectID, string(c.ID),
		fmt.Sprintf("Character %q updated → v%d", c.Name, c.Version))
	return c, nil
}

func (s *CharacterService) SetAppearance(ctx context.Context, id domain.CharacterID, a domain.Appearance) (*domain.Character, error) {
	return s.UpdateCharacter(ctx, id, func(c *domain.Character) error {
		c.Appearance = a
		return nil
	})
}

func (s *CharacterService) SetPersonality(ctx context.Context, id domain.CharacterID, p domain.Personality) (*domain.Character, error) {
	return s.UpdateCharacter(ctx, id, func(c *domain.Character) error {
		c.Personality = p
		return nil
	})
}

func (s *CharacterService) AssignVoice(ctx context.Context, id domain.CharacterID, v domain.VoiceProfile) (*domain.Character, error) {
	return s.UpdateCharacter(ctx, id, func(c *domain.Character) error {
		c.VoiceProfile = v
		return nil
	})
}

// UpdateWardrobe changes the canonical wardrobe (character-level default).
func (s *CharacterService) UpdateWardrobe(ctx context.Context, id domain.CharacterID, wardrobe domain.Wardrobe) (*domain.Character, error) {
	return s.UpdateCharacter(ctx, id, func(c *domain.Character) error {
		c.Wardrobe = wardrobe
		return nil
	})
}

// AssignWardrobe records wardrobe at an episode/scene scope (§17).
func (s *CharacterService) AssignWardrobe(ctx context.Context, id domain.CharacterID, episodeID, sceneID string, items []domain.WardrobeItem, changeEvent string) (*domain.WardrobeAssignment, error) {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return nil, err
	}
	wa := &domain.WardrobeAssignment{
		ID:          "wa_" + uuid.NewString(),
		CharacterID: id,
		EpisodeID:   episodeID,
		SceneID:     sceneID,
		Items:       items,
		ChangeEvent: changeEvent,
	}
	if err := s.repo.SaveWardrobeAssignment(ctx, wa); err != nil {
		return nil, err
	}
	return wa, nil
}

func (s *CharacterService) ListWardrobe(ctx context.Context, id domain.CharacterID, episodeID string) ([]*domain.WardrobeAssignment, error) {
	return s.repo.ListWardrobeAssignments(ctx, id, episodeID)
}

// Lock freezes the current character version for production use.
func (s *CharacterService) Lock(ctx context.Context, id domain.CharacterID) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.SetLocked(ctx, id, true); err != nil {
		return err
	}
	s.events.Emit(ctx, events.CharacterLocked, c.ProjectID, string(id),
		fmt.Sprintf("Character %q locked at v%d", c.Name, c.Version))
	return nil
}

func (s *CharacterService) Unlock(ctx context.Context, id domain.CharacterID) error {
	return s.repo.SetLocked(ctx, id, false)
}

// AddRelationship records a directional relationship between characters.
// The target must exist in the same project — dangling cross-project links
// are rejected.
func (s *CharacterService) AddRelationship(ctx context.Context, id domain.CharacterID, rel domain.Relationship) (*domain.Character, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.IsLocked {
		return nil, domain.ErrCharacterLocked
	}
	target, err := s.repo.FindByID(ctx, rel.TargetCharacterID)
	if err != nil || target.ProjectID != c.ProjectID {
		return nil, domain.ErrCharacterNotFound
	}
	if err := s.repo.SaveRelationship(ctx, id, &rel); err != nil {
		return nil, err
	}
	c.Relationships = append(c.Relationships, rel)
	return c, nil
}

func (s *CharacterService) GetVersion(ctx context.Context, id domain.CharacterID, version int) (*domain.CharacterVersion, error) {
	return s.repo.FindVersion(ctx, id, version)
}

func (s *CharacterService) ListVersions(ctx context.Context, id domain.CharacterID) ([]*domain.CharacterVersion, error) {
	return s.repo.ListVersions(ctx, id)
}
