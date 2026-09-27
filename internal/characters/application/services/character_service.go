package services

import (
	"context"
	"dramastudio/internal/characters/domain"
)

type CharacterService struct {
	repo domain.CharacterRepository
}

func NewCharacterService(repo domain.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

func (s *CharacterService) GetCharacter(ctx context.Context, id domain.CharacterID) (*domain.Character, error) {
	return s.repo.FindByID(ctx, id)
}
