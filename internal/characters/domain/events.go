package domain

import "time"

type CharacterDomainEvent interface {
	EventName() string
}

type CharacterCreatedEvent struct {
	CharacterID CharacterID
	CreatedAt   time.Time
}

func (e CharacterCreatedEvent) EventName() string { return "CharacterCreated" }
