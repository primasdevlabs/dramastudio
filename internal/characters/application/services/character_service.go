package services

import (
	"context"

	"github.com/google/uuid"

	"dramastudio/internal/characters/domain"
)

type CharacterService struct {
	repo domain.CharacterRepository
}

func NewCharacterService(repo domain.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

func (s *CharacterService) CreateCharacter(ctx context.Context, projectID, name, role, bio string) (*domain.Character, error) {
	c := domain.NewCharacter(domain.CharacterID("char_"+uuid.NewString()), projectID, name, role)
	c.Bio = bio
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
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
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.SetLocked(ctx, id, true)
}

func (s *CharacterService) Unlock(ctx context.Context, id domain.CharacterID) error {
	return s.repo.SetLocked(ctx, id, false)
}

// AddRelationship records a directional relationship between characters.
func (s *CharacterService) AddRelationship(ctx context.Context, id domain.CharacterID, rel domain.Relationship) (*domain.Character, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.IsLocked {
		return nil, domain.ErrCharacterLocked
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
