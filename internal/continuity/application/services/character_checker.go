package services

import (
	"fmt"

	"dramastudio/internal/continuity/domain"
)

// CharacterChecker diffs a character's canonical attributes against what
// a scene/shot claims (hair color, scars, costumes at a point in time).
type CharacterChecker struct{}

func NewCharacterChecker() *CharacterChecker { return &CharacterChecker{} }

func (c *CharacterChecker) Check(characterID string, canonical, actual map[string]string) []domain.Violation {
	var violations []domain.Violation
	for attr, want := range canonical {
		got, ok := actual[attr]
		if !ok {
			continue // attribute unspecified in scene — not a violation
		}
		if got != want {
			violations = append(violations, domain.Violation{
				RuleName:    "character_attribute_mismatch",
				Description: fmt.Sprintf("attribute %q: canon=%q actual=%q", attr, want, got),
				Severity:    domain.SeverityError,
				Entity:      characterID,
				Evidence:    attr,
			})
		}
	}
	return violations
}
