-- Bounded context: postproduction — edit timelines, renders, subtitles.
CREATE SCHEMA IF NOT EXISTS postproduction;

-- Structured edit decision list; the render source of truth (§42).
CREATE TABLE IF NOT EXISTS postproduction.timelines (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    episode_id  TEXT NOT NULL,
    version     INT NOT NULL DEFAULT 1,
    tracks      JSONB NOT NULL DEFAULT '{}', -- video/dialogue/music/sfx/subtitle tracks
    status      TEXT NOT NULL DEFAULT 'draft', -- draft|approved|rendered
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (episode_id, version)
);

CREATE TABLE IF NOT EXISTS postproduction.renders (
    id          TEXT PRIMARY KEY,
    timeline_id TEXT NOT NULL REFERENCES postproduction.timelines(id),
    project_id  TEXT NOT NULL,
    episode_id  TEXT NOT NULL,
    object_key  TEXT NOT NULL DEFAULT '',
    url         TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'queued', -- queued|rendering|succeeded|failed
    error       TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS postproduction.subtitles (
    id          TEXT PRIMARY KEY,
    timeline_id TEXT NOT NULL,
    language    TEXT NOT NULL DEFAULT 'en',
    cues        JSONB NOT NULL DEFAULT '[]'
);
