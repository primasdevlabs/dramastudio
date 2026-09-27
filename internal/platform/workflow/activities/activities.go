package activities

import "context"

type ActivityRunner interface {
	ExecuteActivity(ctx context.Context, name string, args ...interface{}) error
}
