package services

import (
	"fmt"

	"dramastudio/internal/continuity/domain"
)

// VisualChecker diffs expected visual attributes (palette, lighting,
// location dressing) against generated output metadata.
type VisualChecker struct{}

func NewVisualChecker() *VisualChecker { return &VisualChecker{} }

func (c *VisualChecker) Check(entityID string, expected, actual map[string]string) []domain.Violation {
	var violations []domain.Violation
	for attr, want := range expected {
		got, ok := actual[attr]
		if !ok {
			continue
		}
		if got != want {
			violations = append(violations, domain.Violation{
				RuleName:    "visual_attribute_mismatch",
				Description: fmt.Sprintf("visual attribute %q: expected=%q actual=%q", attr, want, got),
				Severity:    domain.SeverityWarning,
				Entity:      entityID,
				Evidence:    attr,
			})
		}
	}
	return violations
}
