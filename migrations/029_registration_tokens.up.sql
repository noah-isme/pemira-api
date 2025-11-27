-- +goose Up
-- Registration tokens for TPS check-in

CREATE TABLE IF NOT EXISTS registration_tokens (
    id          BIGSERIAL PRIMARY KEY,
    election_id BIGINT NOT NULL REFERENCES elections(id) ON DELETE CASCADE,
    voter_id    BIGINT NOT NULL REFERENCES voters(id) ON DELETE CASCADE,
    tps_id      BIGINT NULL REFERENCES tps(id) ON DELETE SET NULL,
    token       TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NULL,
    used_at     TIMESTAMPTZ NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ux_registration_tokens_token UNIQUE (token)
);

CREATE INDEX IF NOT EXISTS idx_registration_tokens_election ON registration_tokens (election_id);
CREATE INDEX IF NOT EXISTS idx_registration_tokens_tps ON registration_tokens (tps_id);
