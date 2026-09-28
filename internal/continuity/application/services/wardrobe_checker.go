package services

import (
	"fmt"

	"dramastudio/internal/continuity/domain"
)

// WardrobeChecker compares canonical wardrobe items against what a scene
// assignment actually dresses the character in (§17 violations).
type WardrobeChecker struct{}

func NewWardrobeChecker() *WardrobeChecker { return &WardrobeChecker{} }

// Check returns violations for items present in canonical but missing or
// renamed in the assigned set.
func (c *WardrobeChecker) Check(characterID string, canonical, assigned []string) []domain.Violation {
	assignedSet := make(map[string]bool, len(assigned))
	for _, a := range assigned {
		assignedSet[a] = true
	}
	var violations []domain.Violation
	for _, item := range canonical {
		if !assignedSet[item] {
			violations = append(violations, domain.Violation{
				RuleName:    "wardrobe_missing_item",
				Description: fmt.Sprintf("canonical wardrobe item %q not present in assignment", item),
				Severity:    domain.SeverityWarning,
				Entity:      characterID,
				Evidence:    "assigned items: " + fmt.Sprintf("%v", assigned),
			})
		}
	}
	return violations
}
