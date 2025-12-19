-- Rollback: Remove digital_signature_url column
ALTER TABLE voter_status DROP COLUMN IF EXISTS digital_signature_url;
