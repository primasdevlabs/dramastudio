package domain

type TrackItem struct {
	ID        string  `json:"id"`
	StartTime float64 `json:"start_time"`
	Duration  float64 `json:"duration"`
	AssetURL  string  `json:"asset_url"`
}

type Timeline struct {
	VideoTracks [][]TrackItem `json:"video_tracks"`
	AudioMix    AudioMix      `json:"audio_mix"`
	Subtitles   []Subtitle    `json:"subtitles"`
}
