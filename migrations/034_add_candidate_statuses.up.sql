-- +goose Up
-- Add new candidate statuses

ALTER TYPE myschema.candidate_status ADD VALUE IF NOT EXISTS 'PUBLISHED';
ALTER TYPE myschema.candidate_status ADD VALUE IF NOT EXISTS 'DRAFT';
ALTER TYPE myschema.candidate_status ADD VALUE IF NOT EXISTS 'HIDDEN';
ALTER TYPE myschema.candidate_status ADD VALUE IF NOT EXISTS 'ARCHIVED';

-- Update default from PENDING to DRAFT for new candidates
ALTER TABLE candidates ALTER COLUMN status SET DEFAULT 'DRAFT';

COMMENT ON TYPE myschema.candidate_status IS 'Candidate publication status: DRAFT (not submitted), PENDING (under review), APPROVED/PUBLISHED (visible to public), HIDDEN (temporarily hidden), REJECTED (rejected by admin), WITHDRAWN (withdrawn by candidate), ARCHIVED (archived)';
