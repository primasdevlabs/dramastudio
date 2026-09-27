package domain

import "context"

type ProjectRepository interface {
	FindByID(ctx context.Context, id ProjectID) (*Project, error)
	ListAll(ctx context.Context) ([]*Project, error)
	Save(ctx context.Context, project *Project) error
	SaveBible(ctx context.Context, bible *SeriesBible) error
	GetLatestBible(ctx context.Context, projectID ProjectID) (*SeriesBible, error)
}
