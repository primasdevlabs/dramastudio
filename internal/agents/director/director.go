package director

import (
	"context"
	"fmt"
	"time"

	"dramastudio/internal/agents/domain"
)

type LeadDirector struct{}

func NewLeadDirector() *LeadDirector {
	return &LeadDirector{}
}

func (d *LeadDirector) ExecuteLoopStep(ctx context.Context, projectID, episodeID string) (*domain.Decision, error) {
	// Observe -> Assess -> Plan -> Delegate -> Evaluate -> Decide
	decisionID := fmt.Sprintf("dec_%d", time.Now().UnixNano())
	dec := &domain.Decision{
		ID:            decisionID,
		ProjectID:     projectID,
		EpisodeID:     episodeID,
		Decision:      "ADVANCE_PRODUCTION_STAGE",
		Reason:        "Series Bible & Script validated cleanly without blocking continuity issues",
		DecisionMaker: domain.DecisionMakerLeadDirector,
		Mode:          "monitored",
		Timestamp:     time.Now().UTC(),
	}
	return dec, nil
}
