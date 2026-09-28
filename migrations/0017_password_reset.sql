-- Password reset tokens: sha256 digests only, single-use, time-bounded.
CREATE TABLE IF NOT EXISTS identity.password_reset_tokens (
    token_hash  TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL REFERENCES identity.users(id),
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_pwd_reset_user ON identity.password_reset_tokens(user_id);
