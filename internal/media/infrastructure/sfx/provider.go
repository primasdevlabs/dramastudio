package sfx

import "context"

type SFXProvider interface {
	GenerateSFX(ctx context.Context, prompt string) (string, error)
}
