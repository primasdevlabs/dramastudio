package persistence

import (
	"context"
	"sync"

	"dramastudio/internal/characters/domain"
)

type InMemoryCharacterRepository struct {
	mu         sync.RWMutex
	characters map[domain.CharacterID]*domain.Character
	versions   map[domain.CharacterID][]*domain.CharacterVersion
	wardrobe   map[string][]*domain.WardrobeAssignment
}

func NewInMemoryCharacterRepository() *InMemoryCharacterRepository {
	return &InMemoryCharacterRepository{
		characters: make(map[domain.CharacterID]*domain.Character),
		versions:   make(map[domain.CharacterID][]*domain.CharacterVersion),
		wardrobe:   make(map[string][]*domain.WardrobeAssignment),
	}
}

func (r *InMemoryCharacterRepository) Save(_ context.Context, c *domain.Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.characters[c.ID] = c
	return nil
}

func (r *InMemoryCharacterRepository) FindByID(_ context.Context, id domain.CharacterID) (*domain.Character, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.characters[id]
	if !ok {
		return nil, domain.ErrCharacterNotFound
	}
	return c, nil
}

func (r *InMemoryCharacterRepository) List(_ context.Context, projectID string) ([]*domain.Character, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Character, 0)
	for _, c := range r.characters {
		if projectID == "" || c.ProjectID == projectID {
			out = append(out, c)
		}
	}
	return out, nil
}

func (r *InMemoryCharacterRepository) SetLocked(_ context.Context, id domain.CharacterID, locked bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.characters[id]
	if !ok {
		return domain.ErrCharacterNotFound
	}
	c.IsLocked = locked
	return nil
}

func (r *InMemoryCharacterRepository) SaveVersion(_ context.Context, v *domain.CharacterVersion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.versions[v.CharacterID]
	for i, existing := range list {
		if existing.Version == v.Version {
			list[i] = v
			return nil
		}
	}
	r.versions[v.CharacterID] = append(list, v)
	return nil
}

func (r *InMemoryCharacterRepository) FindVersion(_ context.Context, characterID domain.CharacterID, version int) (*domain.CharacterVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.versions[characterID] {
		if v.Version == version {
			return v, nil
		}
	}
	return nil, domain.ErrCharacterNotFound
}

func (r *InMemoryCharacterRepository) ListVersions(_ context.Context, characterID domain.CharacterID) ([]*domain.CharacterVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*domain.CharacterVersion(nil), r.versions[characterID]...), nil
}

func (r *InMemoryCharacterRepository) SaveRelationship(_ context.Context, characterID domain.CharacterID, rel *domain.Relationship) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.characters[characterID]
	if !ok {
		return domain.ErrCharacterNotFound
	}
	c.Relationships = append(c.Relationships, *rel)
	return nil
}

func wardrobeKey(characterID domain.CharacterID, episodeID string) string {
	return string(characterID) + "|" + episodeID
}

func (r *InMemoryCharacterRepository) SaveWardrobeAssignment(_ context.Context, wa *domain.WardrobeAssignment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := wardrobeKey(wa.CharacterID, wa.EpisodeID)
	r.wardrobe[k] = append(r.wardrobe[k], wa)
	return nil
}

func (r *InMemoryCharacterRepository) ListWardrobeAssignments(_ context.Context, characterID domain.CharacterID, episodeID string) ([]*domain.WardrobeAssignment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]*domain.WardrobeAssignment(nil), r.wardrobe[wardrobeKey(characterID, episodeID)]...), nil
}
