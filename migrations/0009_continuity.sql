-- Bounded context: continuity — validation findings + narrative timeline.
CREATE SCHEMA IF NOT EXISTS continuity;

CREATE TABLE IF NOT EXISTS continuity.checks (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    episode_id  TEXT NOT NULL DEFAULT '',
    check_type  TEXT NOT NULL,        -- story|character|visual|timeline|production
    status      TEXT NOT NULL DEFAULT 'running',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_checks_project ON continuity.checks(project_id);

CREATE TABLE IF NOT EXISTS continuity.issues (
    id             TEXT PRIMARY KEY,
    check_id       TEXT NOT NULL DEFAULT '',
    project_id     TEXT NOT NULL,
    episode_id     TEXT NOT NULL DEFAULT '',
    category       TEXT NOT NULL,     -- wardrobe|knowledge|timeline|...
    severity       TEXT NOT NULL,     -- INFO|WARNING|ERROR|BLOCKING
    entity         TEXT NOT NULL DEFAULT '',
    evidence       TEXT NOT NULL DEFAULT '',
    expected_state TEXT NOT NULL DEFAULT '',
    actual_state   TEXT NOT NULL DEFAULT '',
    resolution     TEXT NOT NULL DEFAULT '',
    status         TEXT NOT NULL DEFAULT 'open', -- open|resolved|wontfix
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at    TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_issues_project ON continuity.issues(project_id);

-- Narrative-time engine (§40): world-time ordering, independent of
-- production time.
CREATE TABLE IF NOT EXISTS continuity.timeline_events (
    id           TEXT PRIMARY KEY,
    project_id   TEXT NOT NULL,
    episode_id   TEXT NOT NULL DEFAULT '',
    scene_id     TEXT NOT NULL DEFAULT '',
    world_time   TEXT NOT NULL DEFAULT '',  -- narrative timestamp label
    event_order  INT NOT NULL DEFAULT 0,
    participants JSONB NOT NULL DEFAULT '[]',
    location_id  TEXT NOT NULL DEFAULT '',
    duration_sec DOUBLE PRECISION NOT NULL DEFAULT 0,
    description  TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_timeline_project ON continuity.timeline_events(project_id, episode_id);
