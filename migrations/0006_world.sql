-- Bounded context: world — locations, variants, props, environmental rules.
CREATE SCHEMA IF NOT EXISTS world;

CREATE TABLE IF NOT EXISTS world.locations (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    name        TEXT NOT NULL,
    kind        TEXT NOT NULL DEFAULT '',   -- building|room|exterior|...
    description TEXT NOT NULL DEFAULT '',
    parent_id   TEXT NOT NULL DEFAULT '',
    visual_ref  TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_locations_project ON world.locations(project_id);

CREATE TABLE IF NOT EXISTS world.location_variants (
    id          TEXT PRIMARY KEY,
    location_id TEXT NOT NULL REFERENCES world.locations(id),
    name        TEXT NOT NULL,              -- day|night|rain|...
    attributes  JSONB NOT NULL DEFAULT '{}',
    UNIQUE (location_id, name)
);

CREATE TABLE IF NOT EXISTS world.props (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    location_id TEXT NOT NULL DEFAULT '',
    visual_ref  TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS world.world_rules (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    rule_text   TEXT NOT NULL,
    category    TEXT NOT NULL DEFAULT ''
);
