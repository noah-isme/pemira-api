-- +goose Up
-- Auth System Migration
-- Make user_role compatible with earlier migration and add user_sessions

-- Ensure user_role exists and contains roles needed by auth system
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM (
            'ADMIN','PANITIA','KETUA_TPS','OPERATOR_PANEL','VIEWER',
            'STUDENT','TPS_OPERATOR','SUPER_ADMIN'
        );
    ELSE
        IF NOT EXISTS (SELECT 1 FROM pg_enum e JOIN pg_type t ON e.enumtypid = t.oid WHERE t.typname = 'user_role' AND e.enumlabel = 'STUDENT') THEN
            ALTER TYPE user_role ADD VALUE 'STUDENT';
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_enum e JOIN pg_type t ON e.enumtypid = t.oid WHERE t.typname = 'user_role' AND e.enumlabel = 'TPS_OPERATOR') THEN
            ALTER TYPE user_role ADD VALUE 'TPS_OPERATOR';
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_enum e JOIN pg_type t ON e.enumtypid = t.oid WHERE t.typname = 'user_role' AND e.enumlabel = 'SUPER_ADMIN') THEN
            ALTER TYPE user_role ADD VALUE 'SUPER_ADMIN';
        END IF;
    END IF;
END$$;

-- Ensure base user_accounts table exists (created earlier in 004)
CREATE TABLE IF NOT EXISTS user_accounts (
    id              BIGSERIAL PRIMARY KEY,
    username        TEXT NOT NULL,
    email           TEXT NOT NULL,
    password_hash   TEXT NOT NULL,
    full_name       TEXT NOT NULL,
    role            user_role NOT NULL DEFAULT 'VIEWER',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_user_accounts_username ON user_accounts (username);
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_accounts_email ON user_accounts (email);
CREATE INDEX IF NOT EXISTS idx_user_accounts_role ON user_accounts (role);

-- User sessions table for refresh tokens
CREATE TABLE IF NOT EXISTS user_sessions (
    id                   BIGSERIAL PRIMARY KEY,
    user_id              BIGINT NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    refresh_token_hash   TEXT NOT NULL,
    user_agent           TEXT NULL,
    ip_address           INET NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at           TIMESTAMPTZ NOT NULL,
    revoked_at           TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user ON user_sessions (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_user_sessions_token_hash ON user_sessions (refresh_token_hash) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions (expires_at) WHERE revoked_at IS NULL;

COMMENT ON TABLE user_accounts IS 'User accounts for all roles (including students and operators)';
COMMENT ON TABLE user_sessions IS 'User sessions for refresh token management';
COMMENT ON COLUMN user_accounts.username IS 'Username (NIM for students, email/username for others)';
COMMENT ON COLUMN user_sessions.refresh_token_hash IS 'Hashed refresh token (DO NOT store plain token)';
