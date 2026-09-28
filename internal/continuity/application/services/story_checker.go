package services

import (
	"fmt"

	"dramastudio/internal/continuity/domain"
)

// FactView is a provider-neutral subject/predicate/object triple used by
// the story checker so continuity does not import canon's domain package.
type FactView struct {
	Subject   string
	Predicate string
	Object    string
}

// StoryChecker flags claims that contradict established canon facts —
// same subject+predicate with a different object.
type StoryChecker struct{}

func NewStoryChecker() *StoryChecker { return &StoryChecker{} }

func (c *StoryChecker) Check(facts []FactView, claims []FactView) []domain.Violation {
	violations := []domain.Violation{}
	index := make(map[string]FactView, len(facts))
	for _, f := range facts {
		index[f.Subject+"|"+f.Predicate] = f
	}
	for _, claim := range claims {
		if fact, ok := index[claim.Subject+"|"+claim.Predicate]; ok && fact.Object != claim.Object {
			violations = append(violations, domain.Violation{
				RuleName:    "canon_contradiction",
				Description: fmt.Sprintf("%s %s: canon=%q claim=%q", claim.Subject, claim.Predicate, fact.Object, claim.Object),
				Severity:    domain.SeverityBlocking,
				Entity:      claim.Subject,
				Evidence:    claim.Predicate,
			})
		}
	}
	return violations
}
