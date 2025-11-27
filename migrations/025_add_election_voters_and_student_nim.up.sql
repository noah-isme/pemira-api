-- +goose Up
-- Scope NIM uniqueness to students and add election_voters table

-- Limit NIM uniqueness to student voters
DROP INDEX IF EXISTS ux_voters_nim;
CREATE UNIQUE INDEX IF NOT EXISTS ux_voters_student_nim
    ON voters (nim)
    WHERE voter_type = 'STUDENT' AND nim IS NOT NULL;

-- election_voter_status enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'election_voter_status') THEN
        CREATE TYPE election_voter_status AS ENUM ('PENDING','VERIFIED','REJECTED','VOTED','BLOCKED');
    END IF;
END$$;

-- Per-election voter registry to avoid cross-year duplicates
CREATE TABLE IF NOT EXISTS election_voters (
    id BIGSERIAL PRIMARY KEY,

    election_id BIGINT NOT NULL REFERENCES elections(id) ON DELETE CASCADE,
    voter_id    BIGINT NOT NULL REFERENCES voters(id) ON DELETE CASCADE,
    nim         TEXT NOT NULL,

    status        election_voter_status NOT NULL DEFAULT 'PENDING',
    voting_method voting_method NOT NULL DEFAULT 'ONLINE',
    tps_id        BIGINT REFERENCES tps(id) ON DELETE SET NULL,

    checked_in_at TIMESTAMPTZ,
    voted_at      TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT ux_election_voters_election_voter UNIQUE (election_id, voter_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_election_voters_election_nim
    ON election_voters (election_id, nim);

CREATE INDEX IF NOT EXISTS idx_election_voters_election
    ON election_voters (election_id);

CREATE INDEX IF NOT EXISTS idx_election_voters_voter
    ON election_voters (voter_id);

CREATE INDEX IF NOT EXISTS idx_election_voters_status
    ON election_voters (election_id, status);

DROP TRIGGER IF EXISTS update_election_voters_updated_at ON election_voters;
CREATE TRIGGER update_election_voters_updated_at
    BEFORE UPDATE ON election_voters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
