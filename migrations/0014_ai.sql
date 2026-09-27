-- Bounded context: ai — providers, models, scoped policies, generation jobs.
CREATE SCHEMA IF NOT EXISTS ai;

-- Provider credentials live server-side only (§60); secrets are referenced
-- by env name, never stored in plaintext.
CREATE TABLE IF NOT EXISTS ai.providers (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    kind        TEXT NOT NULL,            -- llm|image|video|voice|music|sfx
    base_url    TEXT NOT NULL DEFAULT '',
    api_key_env TEXT NOT NULL DEFAULT '',
    is_healthy  BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ai.models (
    id           TEXT PRIMARY KEY,
    provider_id  TEXT NOT NULL REFERENCES ai.providers(id),
    name         TEXT NOT NULL,           -- provider-side model name
    capabilities JSONB NOT NULL DEFAULT '[]',
    version      TEXT NOT NULL DEFAULT '',
    cost_per_unit DOUBLE PRECISION NOT NULL DEFAULT 0,
    is_active    BOOLEAN NOT NULL DEFAULT true,
    UNIQUE (provider_id, name)
);

-- Scoped policies: the most specific scope wins (§53).
CREATE TABLE IF NOT EXISTS ai.model_policies (
    id          TEXT PRIMARY KEY,
    scope       TEXT NOT NULL,            -- system|organization|project|season|episode|task
    scope_id    TEXT NOT NULL DEFAULT '', -- '' for system
    capability  TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    model_id    TEXT NOT NULL,
    max_cost    DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (scope, scope_id, capability)
);

-- Every AI generation is a persistent, auditable job (§26).
CREATE TABLE IF NOT EXISTS ai.generation_jobs (
    id               TEXT PRIMARY KEY,
    project_id       TEXT NOT NULL,
    capability       TEXT NOT NULL,
    provider_id      TEXT NOT NULL DEFAULT '',
    model_id         TEXT NOT NULL DEFAULT '',
    input            JSONB NOT NULL DEFAULT '{}',
    output           JSONB NOT NULL DEFAULT '{}',
    status           TEXT NOT NULL DEFAULT 'PENDING', -- PENDING|RUNNING|SUCCEEDED|FAILED|CANCELLED
    attempt          INT NOT NULL DEFAULT 0,
    cost             DOUBLE PRECISION NOT NULL DEFAULT 0,
    provider_job_id  TEXT NOT NULL DEFAULT '',
    idempotency_key  TEXT UNIQUE,
    error            TEXT NOT NULL DEFAULT '',
    started_at       TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_genjobs_project ON ai.generation_jobs(project_id, capability);
