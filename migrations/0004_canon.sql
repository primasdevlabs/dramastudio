-- Bounded context: canon — authoritative story truth + character knowledge.
CREATE SCHEMA IF NOT EXISTS canon;

CREATE TABLE IF NOT EXISTS canon.story_facts (
    id                TEXT PRIMARY KEY,
    project_id        TEXT NOT NULL,
    entity_id         TEXT NOT NULL DEFAULT '',
    subject           TEXT NOT NULL,
    predicate         TEXT NOT NULL,
    object            TEXT NOT NULL,
    fact_type         TEXT NOT NULL DEFAULT '',
    introduced_episode TEXT NOT NULL DEFAULT '',
    effective_from    TEXT NOT NULL DEFAULT '',
    effective_until   TEXT NOT NULL DEFAULT '',
    source            TEXT NOT NULL DEFAULT '',
    confidence        DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    status            TEXT NOT NULL DEFAULT 'canonical', -- canonical|disputed|retconned
    version           INT NOT NULL DEFAULT 1,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_facts_project ON canon.story_facts(project_id);
CREATE INDEX IF NOT EXISTS idx_facts_entity ON canon.story_facts(project_id, entity_id);

CREATE TABLE IF NOT EXISTS canon.fact_versions (
    fact_id     TEXT NOT NULL,
    version     INT NOT NULL,
    value       JSONB NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (fact_id, version)
);

-- Knowledge state: what a character knows — separate from objective canon (§13).
CREATE TABLE IF NOT EXISTS canon.knowledge_states (
    character_id  TEXT NOT NULL,
    episode_id    TEXT NOT NULL,
    known_fact_ids JSONB NOT NULL DEFAULT '[]',
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (character_id, episode_id)
);

CREATE TABLE IF NOT EXISTS canon.canon_rules (
    id         TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    rule_text  TEXT NOT NULL,
    enforced   BOOLEAN NOT NULL DEFAULT true
);
