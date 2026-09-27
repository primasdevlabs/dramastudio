package domain

import (
	"context"
	"errors"
)

var ErrLocationNotFound = errors.New("location not found")

type WorldRepository interface {
	SaveLocation(ctx context.Context, loc *Location) error
	FindLocationByID(ctx context.Context, id string) (*Location, error)
	ListLocations(ctx context.Context, projectID string) ([]*Location, error)
	SaveProp(ctx context.Context, prop *Prop) error
	FindPropByID(ctx context.Context, id string) (*Prop, error)
	ListProps(ctx context.Context, projectID string) ([]*Prop, error)
	SaveWorldRule(ctx context.Context, rule *WorldRule) error
	ListWorldRules(ctx context.Context, projectID string) ([]*WorldRule, error)
}
