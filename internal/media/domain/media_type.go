package domain

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeVoice MediaType = "voice"
	MediaTypeMusic MediaType = "music"
	MediaTypeSFX   MediaType = "sfx"
)
