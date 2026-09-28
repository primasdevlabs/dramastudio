-- Production intelligence layer: execution audit + operator policy layers.

CREATE SCHEMA IF NOT EXISTS intelligence;

-- One row per resolved task execution: the versioned record that makes
-- episode 40 reproducible even after skills/policies evolve.
CREATE TABLE IF NOT EXISTS intelligence.executions (
    id               TEXT PRIMARY KEY,
    project_id       TEXT NOT NULL,
    task_id          TEXT NOT NULL DEFAULT '',
    action           TEXT NOT NULL,
    agent_id         TEXT NOT NULL,
    agent_version    INT NOT NULL DEFAULT 1,
    skill_ids        JSONB NOT NULL DEFAULT '[]',
    skill_versions   JSONB NOT NULL DEFAULT '[]',
    policy_ids       JSONB NOT NULL DEFAULT '[]',
    rule_ids         JSONB NOT NULL DEFAULT '[]',
    evaluator_ids    JSONB NOT NULL DEFAULT '[]',
    provider_id      TEXT NOT NULL DEFAULT '',
    model_id         TEXT NOT NULL DEFAULT '',
    model_version    TEXT NOT NULL DEFAULT '',
    prompt_hash      TEXT NOT NULL DEFAULT '',
    context_hash     TEXT NOT NULL DEFAULT '',
    generation_job_id TEXT NOT NULL DEFAULT '',
    verdict          TEXT NOT NULL,
    violations       JSONB NOT NULL DEFAULT '[]',
    findings         JSONB NOT NULL DEFAULT '[]',
    cost             DOUBLE PRECISION NOT NULL DEFAULT 0,
    error            TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS executions_project_idx ON intelligence.executions (project_id, created_at DESC);
CREATE INDEX IF NOT EXISTS executions_action_idx  ON intelligence.executions (action);

-- Operator policy layers sit on top of catalog defaults. Scope "" means the
-- layer applies wherever its precedence applies.
CREATE TABLE IF NOT EXISTS intelligence.policy_layers (
    policy_id    TEXT NOT NULL,
    layer        TEXT NOT NULL,
    scope_id     TEXT NOT NULL DEFAULT '',
    version      INT NOT NULL DEFAULT 1,
    applies_to   JSONB NOT NULL DEFAULT '[]',
    requirements JSONB NOT NULL DEFAULT '{}',
    protected    JSONB NOT NULL DEFAULT '[]',
    constraints  JSONB NOT NULL DEFAULT '[]',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (policy_id, layer, scope_id)
);
