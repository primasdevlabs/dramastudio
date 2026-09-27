package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/characters/domain"
)

type InMemoryCharacterRepository struct {
	mu         sync.RWMutex
	characters map[domain.CharacterID]*domain.Character
}

func NewInMemoryCharacterRepository() *InMemoryCharacterRepository {
	return &InMemoryCharacterRepository{
		characters: make(map[domain.CharacterID]*domain.Character),
	}
}

func (r *InMemoryCharacterRepository) FindByID(ctx context.Context, id domain.CharacterID) (*domain.Character, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.characters[id]
	if !ok {
		return nil, fmt.Errorf("character not found: %s", id)
	}
	return c, nil
}

func (r *InMemoryCharacterRepository) List(ctx context.Context) ([]*domain.Character, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Character, 0, len(r.characters))
	for _, c := range r.characters {
		res = append(res, c)
	}
	return res, nil
}

func (r *InMemoryCharacterRepository) Save(ctx context.Context, character *domain.Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.characters[character.ID] = character
	return nil
}
