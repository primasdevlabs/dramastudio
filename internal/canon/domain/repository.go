package domain

import (
	"context"
	"errors"
)

var (
	ErrFactNotFound     = errors.New("story fact not found")
	ErrRuleNotFound     = errors.New("canon rule not found")
	ErrKnowledgeMissing = errors.New("knowledge state not found")
)

type CanonRepository interface {
	SaveFact(ctx context.Context, fact *StoryFact) error
	FindFactByID(ctx context.Context, id string) (*StoryFact, error)
	ListFacts(ctx context.Context, projectID string) ([]*StoryFact, error)
	ListFactsForEntity(ctx context.Context, projectID, entityID string) ([]*StoryFact, error)
	SaveFactVersion(ctx context.Context, factID string, version int, value []byte) error
	// ListFactVersions returns the immutable change history of a fact (§72).
	ListFactVersions(ctx context.Context, factID string) ([]*FactVersion, error)

	GetKnowledgeState(ctx context.Context, characterID, episodeID string) (*KnowledgeState, error)
	SaveKnowledgeState(ctx context.Context, ks *KnowledgeState) error

	SaveRule(ctx context.Context, rule *CanonRule) error
	ListRules(ctx context.Context, projectID string) ([]*CanonRule, error)
}
