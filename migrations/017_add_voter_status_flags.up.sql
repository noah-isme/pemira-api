-- +goose Up
-- Add missing preferred/allowed flags to voter_status to support registration + monitoring

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'voting_method') THEN
        CREATE TYPE voting_method AS ENUM ('ONLINE','TPS');
    END IF;
END$$;

ALTER TABLE voter_status
    ADD COLUMN IF NOT EXISTS preferred_method voting_method NULL,
    ADD COLUMN IF NOT EXISTS online_allowed BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS tps_allowed BOOLEAN NOT NULL DEFAULT TRUE;
