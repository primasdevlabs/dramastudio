-- Bounded context: ai — providers, models, scoped policies, generation jobs.
CREATE SCHEMA IF NOT EXISTS ai;

-- Providers are WHERE models are accessed. type selects the adapter
-- (openai_compatible|wan|voicegen|musicgen|mock|...); secrets are
-- referenced by env var name or stored server-side, never committed (§60).
CREATE TABLE IF NOT EXISTS ai.providers (
    id                  TEXT PRIMARY KEY,
    name                TEXT NOT NULL UNIQUE,
    type                TEXT NOT NULL,           -- adapter kind
    base_url            TEXT NOT NULL DEFAULT '',
    api_key_env         TEXT NOT NULL DEFAULT '',
    api_key             TEXT NOT NULL DEFAULT '', -- write-only; never read into API responses
    health_status       TEXT NOT NULL DEFAULT 'unknown', -- unknown|healthy|degraded|down
    last_health_check_at TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Models are data: identifier + observed capabilities + modalities.
-- Capabilities are recorded facts about the model, not assumptions
-- derived from the provider name (tool calling, vision, structured
-- output all vary per provider AND per model).
CREATE TABLE IF NOT EXISTS ai.models (
    id                   TEXT PRIMARY KEY,
    provider_id          TEXT NOT NULL REFERENCES ai.providers(id),
    name                 TEXT NOT NULL,           -- display name
    identifier           TEXT NOT NULL,           -- provider-side model id
    capabilities         JSONB NOT NULL DEFAULT '[]',
    modalities           JSONB NOT NULL DEFAULT '{}', -- {"input":[], "output":[]}
    supported_parameters JSONB NOT NULL DEFAULT '[]',
    limits               JSONB NOT NULL DEFAULT '{}',
    pricing              JSONB NOT NULL DEFAULT '{}', -- {"currency","per_unit","unit"}
    version              TEXT NOT NULL DEFAULT '',
    status               TEXT NOT NULL DEFAULT 'active', -- active|experimental|deprecated|disabled
    metadata             JSONB NOT NULL DEFAULT '{}',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider_id, identifier)
);
CREATE INDEX IF NOT EXISTS idx_models_capability ON ai.models USING GIN (capabilities);

-- Scoped policies: the most specific scope wins (§53). A policy binds a
-- capability to an ordered model set plus constraints; routing_strategy
-- decides candidate order at execution time.
CREATE TABLE IF NOT EXISTS ai.model_policies (
    id                  TEXT PRIMARY KEY,
    scope               TEXT NOT NULL,            -- system|organization|project|series|season|episode|task
    scope_id            TEXT NOT NULL DEFAULT '', -- '' for system
    capability          TEXT NOT NULL,
    provider_id         TEXT NOT NULL,            -- primary
    model_id            TEXT NOT NULL,            -- primary
    fallback_models     JSONB NOT NULL DEFAULT '[]', -- [{provider_id, model_id}]
    allowed_models      JSONB NOT NULL DEFAULT '[]',
    preferred_models    JSONB NOT NULL DEFAULT '[]',
    routing_strategy    TEXT NOT NULL DEFAULT 'automatic',
    quality_requirement TEXT NOT NULL DEFAULT '',
    max_cost            DOUBLE PRECISION NOT NULL DEFAULT 0,
    max_latency_ms      INT NOT NULL DEFAULT 0,
    region              TEXT NOT NULL DEFAULT '',
    concurrency_limit   INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (scope, scope_id, capability)
);

-- Every AI generation is a persistent, auditable job (§26). Provenance
-- records provider + registry model + model version + inputs so outputs
-- are reproducible and portable across providers.
CREATE TABLE IF NOT EXISTS ai.generation_jobs (
    id               TEXT PRIMARY KEY,
    project_id       TEXT NOT NULL,
    capability       TEXT NOT NULL,
    provider_id      TEXT NOT NULL DEFAULT '',
    model_id         TEXT NOT NULL DEFAULT '',
    model_version    TEXT NOT NULL DEFAULT '',
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
