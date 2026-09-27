package domain

import "context"

type CanonRepository interface {
	FindFactByID(ctx context.Context, id string) (*StoryFact, error)
	GetKnowledgeState(ctx context.Context, characterID, episodeID string) (*KnowledgeState, error)
	SaveFact(ctx context.Context, fact *StoryFact) error
}
