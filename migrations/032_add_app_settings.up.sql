-- +goose Up
-- Add app settings table (migrated from timestamped file)

CREATE TABLE IF NOT EXISTS app_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by BIGINT REFERENCES user_accounts(id)
);

INSERT INTO app_settings (key, value, description)
VALUES 
    ('active_election_id', '1', 'ID election yang aktif saat ini untuk admin dashboard'),
    ('default_election_id', '1', 'ID election default untuk voter')
ON CONFLICT (key) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_app_settings_key ON app_settings(key);
