package music

import "context"

type MusicProvider struct{}

func (p *MusicProvider) GenerateMusic(ctx context.Context, prompt string) (string, error) {
	return "https://media.dramastudio.internal/music.mp3", nil
}
