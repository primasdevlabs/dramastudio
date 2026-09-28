-- Bounded context: media — versioned generated artifacts (§27).
CREATE SCHEMA IF NOT EXISTS media;

CREATE TABLE IF NOT EXISTS media.assets (
    id            TEXT PRIMARY KEY,
    project_id    TEXT NOT NULL,
    type          TEXT NOT NULL,          -- image|video|voice|music|sfx|subtitle|storyboard|render|character_ref|location_ref|prop
    character_id  TEXT NOT NULL DEFAULT '',
    location_id   TEXT NOT NULL DEFAULT '',
    episode_id    TEXT NOT NULL DEFAULT '',
    scene_id      TEXT NOT NULL DEFAULT '',
    shot_id       TEXT NOT NULL DEFAULT '',
    provider      TEXT NOT NULL DEFAULT '',
    model         TEXT NOT NULL DEFAULT '',
    prompt        TEXT NOT NULL DEFAULT '',
    spec          JSONB NOT NULL DEFAULT '{}',  -- structured generation input (§28)
    reference_assets JSONB NOT NULL DEFAULT '[]',
    parameters    JSONB NOT NULL DEFAULT '{}',
    status        TEXT NOT NULL DEFAULT 'PENDING', -- PENDING|APPROVED|REJECTED|ARCHIVED
    cost          DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_assets_project ON media.assets(project_id);
CREATE INDEX IF NOT EXISTS idx_assets_shot ON media.assets(shot_id);

-- Every generation is an immutable version; assets are never overwritten (§72).
CREATE TABLE IF NOT EXISTS media.asset_versions (
    id         TEXT PRIMARY KEY,
    asset_id   TEXT NOT NULL REFERENCES media.assets(id),
    version    INT NOT NULL,
    object_key TEXT NOT NULL,             -- storage key
    url        TEXT NOT NULL DEFAULT '',
    status     TEXT NOT NULL DEFAULT 'PENDING',
    cost       DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (asset_id, version)
);

CREATE TABLE IF NOT EXISTS media.generation_jobs (
    id              TEXT PRIMARY KEY,
    project_id      TEXT NOT NULL,
    asset_id        TEXT NOT NULL DEFAULT '',
    capability      TEXT NOT NULL DEFAULT '',
    provider        TEXT NOT NULL DEFAULT '',
    model           TEXT NOT NULL DEFAULT '',
    provider_job_id TEXT NOT NULL DEFAULT '',
    input           TEXT NOT NULL DEFAULT '',
    output_url      TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'PENDING',
    attempt         INT NOT NULL DEFAULT 0,
    cost            DOUBLE PRECISION NOT NULL DEFAULT 0,
    error           TEXT NOT NULL DEFAULT '',
    started_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at    TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_genjobs_project ON media.generation_jobs(project_id);
CREATE INDEX IF NOT EXISTS idx_genjobs_provider ON media.generation_jobs(provider_job_id);
