package domain

import "errors"

var (
	ErrCharacterNotFound = errors.New("character not found")
	ErrCharacterLocked   = errors.New("character is locked and cannot be modified")
)
