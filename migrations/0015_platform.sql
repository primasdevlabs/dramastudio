-- Platform tables: append-only audit log, domain event log, idempotency.
CREATE SCHEMA IF NOT EXISTS platform;

-- Audit log: append-only from the application perspective (§63).
CREATE TABLE IF NOT EXISTS platform.audit_log (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor       TEXT NOT NULL,            -- user id | agent id | system
    actor_kind  TEXT NOT NULL,            -- USER|LEAD_DIRECTOR|SPECIALIZED_AGENT|SYSTEM
    action      TEXT NOT NULL,
    entity_type TEXT NOT NULL DEFAULT '',
    entity_id   TEXT NOT NULL DEFAULT '',
    project_id  TEXT NOT NULL DEFAULT '',
    detail      JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_audit_project ON platform.audit_log(project_id, created_at);

-- Versioned domain event log (§48).
CREATE TABLE IF NOT EXISTS platform.domain_events (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_type   TEXT NOT NULL,
    event_version INT NOT NULL DEFAULT 1,
    aggregate    TEXT NOT NULL DEFAULT '',
    aggregate_id TEXT NOT NULL DEFAULT '',
    project_id   TEXT NOT NULL DEFAULT '',
    payload      JSONB NOT NULL DEFAULT '{}',
    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_events_type ON platform.domain_events(event_type, occurred_at);

-- Idempotency record store (§49).
CREATE TABLE IF NOT EXISTS platform.idempotency_keys (
    key          TEXT PRIMARY KEY,
    request_hash TEXT NOT NULL,
    response     JSONB NOT NULL DEFAULT '{}',
    status       TEXT NOT NULL DEFAULT 'in_progress', -- in_progress|completed
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at   TIMESTAMPTZ NOT NULL
);

-- Provider webhook dedup (§61).
CREATE TABLE IF NOT EXISTS platform.webhook_deliveries (
    id           TEXT PRIMARY KEY,        -- provider-supplied delivery id
    provider     TEXT NOT NULL,
    received_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    payload      JSONB NOT NULL DEFAULT '{}'
);
