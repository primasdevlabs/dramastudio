package persistence

import (
	"context"
	"fmt"
	"sync"

	"dramastudio/internal/canon/domain"
)

type InMemoryCanonRepository struct {
	mu        sync.RWMutex
	facts     map[string]*domain.StoryFact
	knowledge map[string]*domain.KnowledgeState
}

func NewInMemoryCanonRepository() *InMemoryCanonRepository {
	return &InMemoryCanonRepository{
		facts:     make(map[string]*domain.StoryFact),
		knowledge: make(map[string]*domain.KnowledgeState),
	}
}

func (r *InMemoryCanonRepository) FindFactByID(ctx context.Context, id string) (*domain.StoryFact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fact, ok := r.facts[id]
	if !ok {
		return nil, fmt.Errorf("fact not found: %s", id)
	}
	return fact, nil
}

func (r *InMemoryCanonRepository) ListFacts(ctx context.Context, projectID string) ([]*domain.StoryFact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.StoryFact, 0, len(r.facts))
	for _, f := range r.facts {
		res = append(res, f)
	}
	return res, nil
}

func (r *InMemoryCanonRepository) SaveFact(ctx context.Context, fact *domain.StoryFact) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.facts[fact.ID] = fact
	return nil
}

func (r *InMemoryCanonRepository) GetKnowledgeState(ctx context.Context, characterID, episodeID string) (*domain.KnowledgeState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%s", characterID, episodeID)
	ks, ok := r.knowledge[key]
	if !ok {
		return &domain.KnowledgeState{
			CharacterID: characterID,
			EpisodeID:   episodeID,
			KnownFactIDs: []string{},
		}, nil
	}
	return ks, nil
}

func (r *InMemoryCanonRepository) SaveKnowledgeState(ctx context.Context, ks *domain.KnowledgeState) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%s", ks.CharacterID, ks.EpisodeID)
	r.knowledge[key] = ks
	return nil
}
