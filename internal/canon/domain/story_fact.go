package domain

type FactStatus string

const (
	FactStatusCanonical FactStatus = "canonical"
	FactStatusDisputed  FactStatus = "disputed"
	FactStatusRetconned FactStatus = "retconned"
)

type StoryFact struct {
	ID          string     `json:"id"`
	Subject     string     `json:"subject"`
	Predicate   string     `json:"predicate"`
	Object      string     `json:"object"`
	Introduced  string     `json:"introduced"`
	ValidFrom   string     `json:"valid_from"`
	Status      FactStatus `json:"status"`
}
