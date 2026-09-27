package structuredoutput

import "context"

type Generator interface {
	GenerateStructured(ctx context.Context, prompt string, target interface{}) error
}
