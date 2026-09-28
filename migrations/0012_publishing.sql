-- Bounded context: publishing — channels, accounts, scheduled releases.
CREATE SCHEMA IF NOT EXISTS publishing;

CREATE TABLE IF NOT EXISTS publishing.channels (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    platform    TEXT NOT NULL,             -- youtube|tiktok|instagram|...
    account_ref TEXT NOT NULL DEFAULT '',
    config      JSONB NOT NULL DEFAULT '{}',
    enabled     BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS publishing.publications (
    id                TEXT PRIMARY KEY,
    project_id        TEXT NOT NULL,
    episode_id        TEXT NOT NULL,
    channel_id        TEXT NOT NULL REFERENCES publishing.channels(id),
    metadata          JSONB NOT NULL DEFAULT '{}', -- title|description|captions|thumbnail
    video_url         TEXT NOT NULL DEFAULT '',
    scheduled_at      TIMESTAMPTZ,
    status            TEXT NOT NULL DEFAULT 'draft', -- draft|scheduled|publishing|published|failed
    platform_response JSONB NOT NULL DEFAULT '{}',
    external_id       TEXT NOT NULL DEFAULT '',
    published_at      TIMESTAMPTZ,
    idempotency_key   TEXT UNIQUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_publications_project ON publishing.publications(project_id);
