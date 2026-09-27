package domain

type Personality struct {
	Traits           []string `json:"traits"`
	Motivation       string   `json:"motivation"`
	Flaws            []string `json:"flaws"`
	BackstorySummary string   `json:"backstory_summary"`
}
