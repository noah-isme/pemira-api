-- +goose Down

ALTER TABLE voter_status
    DROP COLUMN IF EXISTS preferred_method,
    DROP COLUMN IF EXISTS online_allowed,
    DROP COLUMN IF EXISTS tps_allowed;
