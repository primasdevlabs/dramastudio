package domain

import "context"

type CanonRepository interface {
	FindFactByID(ctx context.Context, id string) (*StoryFact, error)
	ListFacts(ctx context.Context, projectID string) ([]*StoryFact, error)
	SaveFact(ctx context.Context, fact *StoryFact) error
	GetKnowledgeState(ctx context.Context, characterID, episodeID string) (*KnowledgeState, error)
	SaveKnowledgeState(ctx context.Context, ks *KnowledgeState) error
}
