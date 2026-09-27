package services

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/characters/domain"
)

type CharacterService struct {
	repo domain.CharacterRepository
}

func NewCharacterService(repo domain.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

func (s *CharacterService) CreateCharacter(ctx context.Context, projectID, name, role, bio string) (*domain.Character, error) {
	id := domain.CharacterID(fmt.Sprintf("char_%d", time.Now().UnixNano()))
	c := domain.NewCharacter(id, projectID, name, role)
	c.Bio = bio
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, id domain.CharacterID) (*domain.Character, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CharacterService) ListCharacters(ctx context.Context) ([]*domain.Character, error) {
	return s.repo.List(ctx)
}

func (s *CharacterService) UpdateWardrobe(ctx context.Context, id domain.CharacterID, wardrobe domain.Wardrobe) (*domain.Character, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Wardrobe = wardrobe
	c.Version++
	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}
