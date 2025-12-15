-- +goose Down
-- Revert candidate status changes

-- Note: PostgreSQL does not support removing enum values once added
-- We can only revert the default value

ALTER TABLE candidates ALTER COLUMN status SET DEFAULT 'PENDING';
