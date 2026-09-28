package events

import (
	"context"
	"encoding/json"

	"dramastudio/internal/platform/database/postgres"
)

// EventLogStore appends events to the durable domain event log.
type EventLogStore interface {
	Append(ctx context.Context, e DomainEvent) error
}

// sqlEventLog writes to platform.domain_events on either driver — the
// sqlite adapter satisfies postgres.Querier via dialect translation.
type sqlEventLog struct {
	q postgres.Querier
}

func NewSQLEventLog(q postgres.Querier) EventLogStore {
	return &sqlEventLog{q: q}
}

func (s *sqlEventLog) Append(ctx context.Context, e DomainEvent) error {
	payload := []byte("{}")
	if e.Payload != nil {
		if b, err := json.Marshal(e.Payload); err == nil {
			payload = b
		}
	}
	_, err := s.q.Exec(ctx,
		`INSERT INTO platform.domain_events (event_type, event_version, aggregate, aggregate_id, project_id, payload, occurred_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		e.Type, e.Version, e.Aggregate, e.AggregateID, e.ProjectID, string(payload), e.OccurredAt)
	return err
}

// NoopEventLog discards events (memory store profile / tests).
type NoopEventLog struct{}

func (NoopEventLog) Append(context.Context, DomainEvent) error { return nil }
