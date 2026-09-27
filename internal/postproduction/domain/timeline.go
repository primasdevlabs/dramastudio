package domain

type TrackItem struct {
	ID        string  `json:"id"`
	ShotID    string  `json:"shot_id"`
	StartTime float64 `json:"start_time"`
	Duration  float64 `json:"duration"`
	AssetURL  string  `json:"asset_url"`
}

type Timeline struct {
	ID          string      `json:"id"`
	ProjectID   string      `json:"project_id"`
	EpisodeID   string      `json:"episode_id"`
	VideoTracks []TrackItem `json:"video_tracks"`
	AudioTracks []TrackItem `json:"audio_tracks"`
	Subtitles   []Subtitle  `json:"subtitles"`
}
