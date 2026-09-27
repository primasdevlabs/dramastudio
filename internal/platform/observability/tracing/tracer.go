package tracing

import "context"

type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, func())
}
