-- Bounded context: characters — persistent character identity, versioned.
CREATE SCHEMA IF NOT EXISTS characters;

CREATE TABLE IF NOT EXISTS characters.characters (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    name        TEXT NOT NULL,
    role        TEXT NOT NULL DEFAULT '',
    summary     TEXT NOT NULL DEFAULT '',
    locked      BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_characters_project ON characters.characters(project_id);

-- Character versions are immutable once used in production (§15).
CREATE TABLE IF NOT EXISTS characters.character_versions (
    id            TEXT PRIMARY KEY,
    character_id  TEXT NOT NULL REFERENCES characters.characters(id),
    version       INT NOT NULL,
    appearance    JSONB NOT NULL DEFAULT '{}',
    personality   JSONB NOT NULL DEFAULT '{}',
    voice_profile JSONB NOT NULL DEFAULT '{}',
    wardrobe      JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (character_id, version)
);

CREATE TABLE IF NOT EXISTS characters.relationships (
    id              TEXT PRIMARY KEY,
    project_id      TEXT NOT NULL,
    character_id    TEXT NOT NULL REFERENCES characters.characters(id),
    other_character_id TEXT NOT NULL,
    relation        TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    since_episode   TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS characters.wardrobe_assignments (
    id           TEXT PRIMARY KEY,
    character_id TEXT NOT NULL REFERENCES characters.characters(id),
    episode_id   TEXT NOT NULL,
    scene_id     TEXT NOT NULL DEFAULT '',
    items        JSONB NOT NULL DEFAULT '[]',
    change_event TEXT NOT NULL DEFAULT '',  -- justification for the change
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
