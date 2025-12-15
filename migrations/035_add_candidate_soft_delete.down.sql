-- +goose Down
-- Remove soft delete support

DROP INDEX IF EXISTS idx_candidates_deleted_at;
ALTER TABLE candidates DROP COLUMN IF EXISTS deleted_by_admin_id;
ALTER TABLE candidates DROP COLUMN IF EXISTS deleted_at;
