// Package audit writes append-only records of important actions (§63):
// user approvals, overrides, agent decisions, credential changes,
// publishing actions. Entries land in platform.audit_log.
package audit

import (
	"context"
	"encoding/json"

	"dramastudio/internal/platform/database/postgres"
)

// Actor kinds recorded on every entry.
const (
	ActorUser     = "USER"
	ActorDirector = "LEAD_DIRECTOR"
	ActorAgent    = "SPECIALIZED_AGENT"
	ActorSystem   = "SYSTEM"
)

// Entry is one immutable audit record.
type Entry struct {
	Actor      string                 `json:"actor"`
	ActorKind  string                 `json:"actor_kind"`
	Action     string                 `json:"action"`
	EntityType string                 `json:"entity_type,omitempty"`
	EntityID   string                 `json:"entity_id,omitempty"`
	ProjectID  string                 `json:"project_id,omitempty"`
	Detail     map[string]interface{} `json:"detail,omitempty"`
}

// Logger records audit entries.
type Logger interface {
	Record(ctx context.Context, e Entry) error
}

type sqlLogger struct {
	q postgres.Querier
}

// NewLogger returns a SQL-backed audit logger, or a no-op when q is nil
// (memory store profile / tests).
func NewLogger(q postgres.Querier) Logger {
	if q == nil {
		return noopLogger{}
	}
	return &sqlLogger{q: q}
}

func (l *sqlLogger) Record(ctx context.Context, e Entry) error {
	detail := []byte("{}")
	if e.Detail != nil {
		if b, err := json.Marshal(e.Detail); err == nil {
			detail = b
		}
	}
	_, err := l.q.Exec(ctx,
		`INSERT INTO platform.audit_log (actor, actor_kind, action, entity_type, entity_id, project_id, detail)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		e.Actor, e.ActorKind, e.Action, e.EntityType, e.EntityID, e.ProjectID, string(detail))
	return err
}

type noopLogger struct{}

func (noopLogger) Record(context.Context, Entry) error { return nil }

// Record is the nil-safe helper services call; audit failures must never
// break the audited operation.
func Record(l Logger, ctx context.Context, e Entry) {
	if l == nil {
		return
	}
	_ = l.Record(ctx, e)
}
