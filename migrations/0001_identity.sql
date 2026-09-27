-- Bounded context: identity — users, orgs, memberships, credentials.
CREATE SCHEMA IF NOT EXISTS identity;

CREATE TABLE IF NOT EXISTS identity.organizations (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS identity.users (
    id            TEXT PRIMARY KEY,
    org_id        TEXT NOT NULL REFERENCES identity.organizations(id),
    email         TEXT NOT NULL UNIQUE,
    display_name  TEXT NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_users_org ON identity.users(org_id);

CREATE TABLE IF NOT EXISTS identity.memberships (
    user_id    TEXT NOT NULL REFERENCES identity.users(id),
    org_id     TEXT NOT NULL REFERENCES identity.organizations(id),
    role       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, org_id)
);

CREATE TABLE IF NOT EXISTS identity.api_credentials (
    id           TEXT PRIMARY KEY,
    org_id       TEXT NOT NULL REFERENCES identity.organizations(id),
    user_id      TEXT REFERENCES identity.users(id),
    name         TEXT NOT NULL,
    key_digest   TEXT NOT NULL UNIQUE,   -- sha256 of the key; plaintext never stored
    permissions  JSONB NOT NULL DEFAULT '[]',
    service      BOOLEAN NOT NULL DEFAULT false,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
