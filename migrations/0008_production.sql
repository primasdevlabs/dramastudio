-- Bounded context: production — executable production state.
CREATE SCHEMA IF NOT EXISTS production;

CREATE TABLE IF NOT EXISTS production.productions (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL UNIQUE,
    status      TEXT NOT NULL DEFAULT 'PROJECT_CREATED',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS production.runs (
    id             TEXT PRIMARY KEY,
    production_id  TEXT NOT NULL REFERENCES production.productions(id),
    project_id     TEXT NOT NULL,
    episode_id     TEXT NOT NULL,
    stage          TEXT NOT NULL DEFAULT 'pre_production',
    status         TEXT NOT NULL DEFAULT 'pending',
    bible_version  INT,
    workflow_id    TEXT NOT NULL DEFAULT '',   -- temporal workflow handle
    snapshot       JSONB NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_runs_episode ON production.runs(episode_id);

CREATE TABLE IF NOT EXISTS production.jobs (
    id              TEXT PRIMARY KEY,
    production_id   TEXT NOT NULL DEFAULT '',
    project_id      TEXT NOT NULL,
    run_id          TEXT NOT NULL DEFAULT '',
    episode_id      TEXT NOT NULL DEFAULT '',
    scene_id        TEXT NOT NULL DEFAULT '',
    shot_id         TEXT NOT NULL DEFAULT '',
    kind            TEXT NOT NULL DEFAULT '',  -- episode|scene|shot|render|...
    status          TEXT NOT NULL DEFAULT 'pending',
    attempt         INT NOT NULL DEFAULT 0,
    result_url      TEXT NOT NULL DEFAULT '',
    error           TEXT NOT NULL DEFAULT '',
    idempotency_key TEXT UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at    TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_jobs_project ON production.jobs(project_id);
CREATE INDEX IF NOT EXISTS idx_jobs_episode ON production.jobs(episode_id);

CREATE TABLE IF NOT EXISTS production.shots (
    id             TEXT PRIMARY KEY,
    project_id     TEXT NOT NULL,
    episode_id     TEXT NOT NULL,
    scene_id       TEXT NOT NULL,
    seq            INT NOT NULL,
    description    TEXT NOT NULL DEFAULT '',
    camera         JSONB NOT NULL DEFAULT '{}',
    characters     JSONB NOT NULL DEFAULT '[]',
    location_id    TEXT NOT NULL DEFAULT '',
    duration_sec   DOUBLE PRECISION NOT NULL DEFAULT 0,
    status         TEXT NOT NULL DEFAULT 'planned',
    approved_asset TEXT NOT NULL DEFAULT '',
    UNIQUE (scene_id, seq)
);

CREATE TABLE IF NOT EXISTS production.approvals (
    id           TEXT PRIMARY KEY,
    project_id   TEXT NOT NULL,
    episode_id   TEXT NOT NULL DEFAULT '',
    stage        TEXT NOT NULL,
    target_id    TEXT NOT NULL,
    decision     TEXT NOT NULL DEFAULT '', -- APPROVE|REJECT|REQUEST_REVISION|REGENERATE|PAUSE|STOP|OVERRIDE
    notes        TEXT NOT NULL DEFAULT '',
    decided_by   TEXT NOT NULL DEFAULT '',
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    decided_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_approvals_project ON production.approvals(project_id);
