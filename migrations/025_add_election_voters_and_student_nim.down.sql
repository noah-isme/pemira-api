-- +goose Down
-- Rollback election_voters and restore global NIM uniqueness

DROP TRIGGER IF EXISTS update_election_voters_updated_at ON election_voters;
DROP TABLE IF EXISTS election_voters;
DROP TYPE IF EXISTS election_voter_status;

DROP INDEX IF EXISTS ux_voters_student_nim;
CREATE UNIQUE INDEX IF NOT EXISTS ux_voters_nim
    ON voters (nim);
