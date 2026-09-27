-- Bounded context: story — narrative hierarchy Series→Season→Arc→Episode→Scene→Beat.
CREATE SCHEMA IF NOT EXISTS story;

CREATE TABLE IF NOT EXISTS story.series (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_series_project ON story.series(project_id);

CREATE TABLE IF NOT EXISTS story.seasons (
    id        TEXT PRIMARY KEY,
    series_id TEXT NOT NULL REFERENCES story.series(id),
    number    INT NOT NULL,
    title     TEXT NOT NULL DEFAULT '',
    summary   TEXT NOT NULL DEFAULT '',
    UNIQUE (series_id, number)
);

CREATE TABLE IF NOT EXISTS story.arcs (
    id        TEXT PRIMARY KEY,
    season_id TEXT NOT NULL REFERENCES story.seasons(id),
    title     TEXT NOT NULL,
    number    INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS story.episodes (
    id         TEXT PRIMARY KEY,
    season_id  TEXT NOT NULL REFERENCES story.seasons(id),
    arc_id     TEXT NOT NULL DEFAULT '',
    number     INT NOT NULL,
    title      TEXT NOT NULL DEFAULT '',
    summary    TEXT NOT NULL DEFAULT '',
    script     TEXT NOT NULL DEFAULT '',
    status     TEXT NOT NULL DEFAULT 'PLANNED',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (season_id, number)
);

CREATE TABLE IF NOT EXISTS story.scenes (
    id           TEXT PRIMARY KEY,
    episode_id   TEXT NOT NULL REFERENCES story.episodes(id),
    number       INT NOT NULL,
    title        TEXT NOT NULL DEFAULT '',
    location_id  TEXT NOT NULL DEFAULT '',
    time_of_day  TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    character_ids JSONB NOT NULL DEFAULT '[]',
    UNIQUE (episode_id, number)
);

CREATE TABLE IF NOT EXISTS story.beats (
    id           TEXT PRIMARY KEY,
    scene_id     TEXT NOT NULL REFERENCES story.scenes(id),
    seq          INT NOT NULL,
    action       TEXT NOT NULL DEFAULT '',
    dialogue     TEXT NOT NULL DEFAULT '',
    character_id TEXT NOT NULL DEFAULT '',
    UNIQUE (scene_id, seq)
);

-- Story graph: nodes are narrative events; edges carry causal relations.
CREATE TABLE IF NOT EXISTS story.graph_nodes (
    id         TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    type       TEXT NOT NULL,          -- Event|Reveal|Conflict|Decision|...
    title      TEXT NOT NULL DEFAULT '',
    episode_id TEXT NOT NULL DEFAULT '',
    scene_id   TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_graph_nodes_project ON story.graph_nodes(project_id);

CREATE TABLE IF NOT EXISTS story.graph_edges (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id TEXT NOT NULL,
    from_id    TEXT NOT NULL,
    to_id      TEXT NOT NULL,
    relation   TEXT NOT NULL,          -- causes|reveals|depends_on|...
    UNIQUE (project_id, from_id, to_id, relation)
);

CREATE TABLE IF NOT EXISTS story.plot_threads (
    id          TEXT PRIMARY KEY,
    project_id  TEXT NOT NULL,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'open'
);
