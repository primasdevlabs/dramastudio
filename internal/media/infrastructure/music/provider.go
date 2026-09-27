package music

import "context"

type MusicProvider interface {
	GenerateMusic(ctx context.Context, prompt string) (string, error)
}
