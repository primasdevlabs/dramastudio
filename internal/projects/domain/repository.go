package domain

import "context"

type ProjectRepository interface {
	FindByID(ctx context.Context, id ProjectID) (*Project, error)
	Save(ctx context.Context, project *Project) error
}
