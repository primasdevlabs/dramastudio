package services

import (
	"fmt"
	"sort"

	"dramastudio/internal/continuity/domain"
)

// TimelineChecker validates narrative-time ordering (§40): the same
// participant cannot appear in two events sharing a world-time label, and
// event_order must be monotonic within an episode.
type TimelineChecker struct{}

func NewTimelineChecker() *TimelineChecker { return &TimelineChecker{} }

func (c *TimelineChecker) Check(events []domain.TimelineEvent) []domain.Violation {
	var violations []domain.Violation

	seen := map[string][]domain.TimelineEvent{} // world_time+participant -> events
	for _, e := range events {
		for _, p := range e.Participants {
			key := e.WorldTime + "|" + p
			for _, prev := range seen[key] {
				if prev.ID != e.ID && prev.SceneID != e.SceneID {
					violations = append(violations, domain.Violation{
						RuleName:    "participant_double_booked",
						Description: fmt.Sprintf("participant %q appears at world time %q in scenes %s and %s", p, e.WorldTime, prev.SceneID, e.SceneID),
						Severity:    domain.SeverityBlocking,
						Entity:      p,
						Evidence:    fmt.Sprintf("events %s and %s share world_time %q", prev.ID, e.ID, e.WorldTime),
					})
				}
			}
			seen[key] = append(seen[key], e)
		}
	}

	sorted := append([]domain.TimelineEvent{}, events...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].EventOrder < sorted[j].EventOrder })
	for i := 1; i < len(sorted); i++ {
		if sorted[i].WorldTime < sorted[i-1].WorldTime {
			violations = append(violations, domain.Violation{
				RuleName:    "world_time_not_monotonic",
				Description: fmt.Sprintf("event %q at %q follows event %q at %q", sorted[i].ID, sorted[i].WorldTime, sorted[i-1].ID, sorted[i-1].WorldTime),
				Severity:    domain.SeverityWarning,
				Entity:      sorted[i].SceneID,
				Evidence:    "event_order conflicts with world_time ordering",
			})
		}
	}
	return violations
}
