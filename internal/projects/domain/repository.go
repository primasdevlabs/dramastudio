package domain

import (
	"context"
	"errors"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrBibleNotFound   = errors.New("series bible not found")
)

type ProjectRepository interface {
	Save(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id ProjectID) (*Project, error)
	ListAll(ctx context.Context, orgID string) ([]*Project, error)
	SaveBible(ctx context.Context, bible *SeriesBible) error
	GetLatestBible(ctx context.Context, projectID ProjectID) (*SeriesBible, error)
	GetBibleVersion(ctx context.Context, projectID ProjectID, version int) (*SeriesBible, error)
	ListBibleVersions(ctx context.Context, projectID ProjectID) ([]*SeriesBible, error)
}
