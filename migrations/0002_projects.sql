-- Bounded context: projects — productions, settings, bibles, budgets.
CREATE SCHEMA IF NOT EXISTS projects;

CREATE TABLE IF NOT EXISTS projects.projects (
    id              TEXT PRIMARY KEY,
    org_id          TEXT NOT NULL DEFAULT '',
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    genre           TEXT NOT NULL DEFAULT '',
    language        TEXT NOT NULL DEFAULT '',
    mode            TEXT NOT NULL DEFAULT 'monitored',
    status          TEXT NOT NULL DEFAULT 'PROJECT_CREATED',
    settings        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS projects.series_bibles (
    id              TEXT PRIMARY KEY,
    project_id      TEXT NOT NULL REFERENCES projects.projects(id),
    version         INT NOT NULL,
    premise         TEXT NOT NULL DEFAULT '',
    genre           TEXT NOT NULL DEFAULT '',
    themes          JSONB NOT NULL DEFAULT '[]',
    tone            TEXT NOT NULL DEFAULT '',
    world_rules     JSONB NOT NULL DEFAULT '[]',
    narrative_rules JSONB NOT NULL DEFAULT '[]',
    visual_style    TEXT NOT NULL DEFAULT '',
    dialogue_style  TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, version)
);

CREATE TABLE IF NOT EXISTS projects.budgets (
    id           TEXT PRIMARY KEY,
    project_id   TEXT NOT NULL REFERENCES projects.projects(id),
    scope        TEXT NOT NULL,          -- project|season|episode|task
    scope_id     TEXT NOT NULL DEFAULT '',
    limit_amount DOUBLE PRECISION NOT NULL,
    spent        DOUBLE PRECISION NOT NULL DEFAULT 0,
    currency     TEXT NOT NULL DEFAULT 'USD',
    UNIQUE (project_id, scope, scope_id)
);
