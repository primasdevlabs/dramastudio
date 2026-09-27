package commands

import "dramastudio/internal/characters/domain"

type AssignVoice struct {
	ID           domain.CharacterID
	VoiceProfile domain.VoiceProfile
}
