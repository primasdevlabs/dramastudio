package domain

import "context"

// ModelRegistryRepository is the ai context's persistence surface.
// It also satisfies platform/ai/routing.Registry via ResolvePolicy.
type ModelRegistryRepository interface {
	SaveProvider(ctx context.Context, p *Provider) error
	FindProviderByID(ctx context.Context, id ProviderID) (*Provider, error)
	ListProviders(ctx context.Context) ([]*Provider, error)

	SaveModel(ctx context.Context, m *Model) error
	FindModelByID(ctx context.Context, id ModelID) (*Model, error)
	ListModels(ctx context.Context, providerID ProviderID) ([]*Model, error)
	ListModelsByCapability(ctx context.Context, capability AICapability) ([]*Model, error)
	DeleteModel(ctx context.Context, id ModelID) error
	DeleteProvider(ctx context.Context, id ProviderID) error

	UpsertPolicy(ctx context.Context, p *PolicyRow) error
	ResolvePolicy(ctx context.Context, capability, scopeID string) (*PolicyRow, error)
	ListPolicies(ctx context.Context, scope, scopeID string) ([]*PolicyRow, error)

	SaveGenerationJob(ctx context.Context, j *GenerationJob) error
	FindGenerationJobByID(ctx context.Context, id string) (*GenerationJob, error)
	FindGenerationJobByIdempotencyKey(ctx context.Context, key string) (*GenerationJob, error)
	ListGenerationJobs(ctx context.Context, projectID string) ([]*GenerationJob, error)
}
