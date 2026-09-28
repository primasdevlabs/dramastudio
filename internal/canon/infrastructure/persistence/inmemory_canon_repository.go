package persistence

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"dramastudio/internal/canon/domain"
)

type InMemoryCanonRepository struct {
	mu        sync.RWMutex
	facts     map[string]*domain.StoryFact
	knowledge map[string]*domain.KnowledgeState
	rules     map[string]*domain.CanonRule
	versions  map[string]map[int]*domain.FactVersion
}

func NewInMemoryCanonRepository() *InMemoryCanonRepository {
	return &InMemoryCanonRepository{
		facts:     make(map[string]*domain.StoryFact),
		knowledge: make(map[string]*domain.KnowledgeState),
		rules:     make(map[string]*domain.CanonRule),
		versions:  make(map[string]map[int]*domain.FactVersion),
	}
}

func (r *InMemoryCanonRepository) SaveFact(_ context.Context, f *domain.StoryFact) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.facts[f.ID] = f
	return nil
}

func (r *InMemoryCanonRepository) FindFactByID(_ context.Context, id string) (*domain.StoryFact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if f, ok := r.facts[id]; ok {
		return f, nil
	}
	return nil, domain.ErrFactNotFound
}

func (r *InMemoryCanonRepository) ListFacts(_ context.Context, projectID string) ([]*domain.StoryFact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.StoryFact, 0)
	for _, f := range r.facts {
		if f.ProjectID == projectID {
			out = append(out, f)
		}
	}
	return out, nil
}

func (r *InMemoryCanonRepository) ListFactsForEntity(_ context.Context, projectID, entityID string) ([]*domain.StoryFact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.StoryFact, 0)
	for _, f := range r.facts {
		if f.ProjectID == projectID && (f.EntityID == entityID || f.Subject == entityID) {
			out = append(out, f)
		}
	}
	return out, nil
}

func (r *InMemoryCanonRepository) SaveFactVersion(_ context.Context, factID string, version int, value []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.versions[factID] == nil {
		r.versions[factID] = make(map[int]*domain.FactVersion)
	}
	r.versions[factID][version] = &domain.FactVersion{
		FactID:    factID,
		Version:   version,
		Value:     append(json.RawMessage(nil), value...),
		UpdatedAt: time.Now().UTC(),
	}
	return nil
}

func (r *InMemoryCanonRepository) ListFactVersions(_ context.Context, factID string) ([]*domain.FactVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	vers := r.versions[factID]
	out := make([]*domain.FactVersion, 0, len(vers))
	nums := make([]int, 0, len(vers))
	for n := range vers {
		nums = append(nums, n)
	}
	sort.Ints(nums)
	for _, n := range nums {
		out = append(out, vers[n])
	}
	return out, nil
}

func knowledgeKey(characterID, episodeID string) string {
	return characterID + "|" + episodeID
}

func (r *InMemoryCanonRepository) GetKnowledgeState(_ context.Context, characterID, episodeID string) (*domain.KnowledgeState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ks, ok := r.knowledge[knowledgeKey(characterID, episodeID)]; ok {
		return ks, nil
	}
	return nil, domain.ErrKnowledgeMissing
}

func (r *InMemoryCanonRepository) SaveKnowledgeState(_ context.Context, ks *domain.KnowledgeState) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.knowledge[knowledgeKey(ks.CharacterID, ks.EpisodeID)] = ks
	return nil
}

func (r *InMemoryCanonRepository) SaveRule(_ context.Context, rule *domain.CanonRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryCanonRepository) ListRules(_ context.Context, projectID string) ([]*domain.CanonRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.CanonRule, 0)
	for _, rule := range r.rules {
		if rule.ProjectID == projectID {
			out = append(out, rule)
		}
	}
	return out, nil
}
