-- Migration: Add digital_signature_url column to voter_status
-- Date: 2025-12-19
-- Description: Store URL to signature image in Supabase Storage
--              Format: {SUPABASE_URL}/storage/v1/object/public/pemira/signatures/{election_id}/{voter_id}.png

ALTER TABLE voter_status ADD COLUMN IF NOT EXISTS digital_signature_url TEXT;

-- Add comment
COMMENT ON COLUMN voter_status.digital_signature_url IS 'URL to signature image stored in Supabase Storage: pemira/signatures/{election_id}/{voter_id}.png';
