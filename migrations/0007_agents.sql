-- Bounded context: agents — virtual production team, tasks, decisions.
CREATE SCHEMA IF NOT EXISTS agents;

CREATE TABLE IF NOT EXISTS agents.definitions (
    id             TEXT PRIMARY KEY,
    role           TEXT NOT NULL,          -- lead_director|story_director|...
    name           TEXT NOT NULL,
    instructions   TEXT NOT NULL DEFAULT '',
    skills         JSONB NOT NULL DEFAULT '[]',
    tools          JSONB NOT NULL DEFAULT '[]',
    permissions    JSONB NOT NULL DEFAULT '[]',
    model_policy   JSONB NOT NULL DEFAULT '{}',
    budget_limit   DOUBLE PRECISION,
    active         BOOLEAN NOT NULL DEFAULT true,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agents.tasks (
    id               TEXT PRIMARY KEY,
    project_id       TEXT NOT NULL,
    agent_id         TEXT NOT NULL DEFAULT '',
    objective        TEXT NOT NULL,
    input            JSONB NOT NULL DEFAULT '{}',
    expected_output  TEXT NOT NULL DEFAULT '',
    constraints      JSONB NOT NULL DEFAULT '{}',
    budget           DOUBLE PRECISION,
    max_iterations   INT NOT NULL DEFAULT 3,
    iteration        INT NOT NULL DEFAULT 0,
    timeout_sec      INT NOT NULL DEFAULT 3600,
    approval_policy  TEXT NOT NULL DEFAULT 'monitored',
    success_criteria TEXT NOT NULL DEFAULT '',
    status           TEXT NOT NULL DEFAULT 'pending', -- pending|running|succeeded|failed|escalated|cancelled
    result           JSONB,
    error            TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at     TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_agent_tasks_project ON agents.tasks(project_id);

CREATE TABLE IF NOT EXISTS agents.decisions (
    id             TEXT PRIMARY KEY,
    project_id     TEXT NOT NULL,
    episode_id     TEXT NOT NULL DEFAULT '',
    decision       TEXT NOT NULL,
    reason         TEXT NOT NULL DEFAULT '',
    decision_maker TEXT NOT NULL,        -- USER|LEAD_DIRECTOR|SPECIALIZED_AGENT|SYSTEM
    mode           TEXT NOT NULL DEFAULT '', -- monitored|autonomous
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_agent_decisions_project ON agents.decisions(project_id);
