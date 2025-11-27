-- +goose Down
-- Rollback auth system migration

DROP TABLE IF EXISTS user_sessions;
