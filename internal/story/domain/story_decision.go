package domain

type StoryDecision struct {
	ID        string `json:"id"`
	Rationale string `json:"rationale"`
	Outcome   string `json:"outcome"`
}
