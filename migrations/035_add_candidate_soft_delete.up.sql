-- +goose Up
-- Add soft delete support for candidates

ALTER TABLE candidates ADD COLUMN deleted_at TIMESTAMPTZ NULL;
ALTER TABLE candidates ADD COLUMN deleted_by_admin_id BIGINT NULL REFERENCES user_accounts(id) ON DELETE SET NULL;

CREATE INDEX idx_candidates_deleted_at ON candidates(deleted_at);

COMMENT ON COLUMN candidates.deleted_at IS 'Soft delete timestamp. NULL means not deleted.';
COMMENT ON COLUMN candidates.deleted_by_admin_id IS 'Admin who soft deleted this candidate';
