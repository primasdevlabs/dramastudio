package domain

type AudioTrack struct {
	Type   string  `json:"type"` // voice, music, sfx
	URL    string  `json:"url"`
	Volume float64 `json:"volume"`
}

type AudioMix struct {
	Tracks []AudioTrack `json:"tracks"`
}
