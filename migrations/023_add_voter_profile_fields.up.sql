-- Migration: Add profile fields to voters and users tables
-- Date: 2025-11-25

-- Add profile fields to voters table
ALTER TABLE voters
ADD COLUMN IF NOT EXISTS email VARCHAR(255),
ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
ADD COLUMN IF NOT EXISTS photo_url TEXT,
ADD COLUMN IF NOT EXISTS bio TEXT,
ADD COLUMN IF NOT EXISTS voting_method_preference voting_method DEFAULT 'ONLINE',
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- Add index for email lookups
CREATE INDEX IF NOT EXISTS idx_voters_email ON voters(email);
CREATE INDEX IF NOT EXISTS idx_voters_updated_at ON voters(updated_at);

-- Add login tracking fields to user_accounts table
ALTER TABLE user_accounts
ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS login_count INTEGER DEFAULT 0;

-- Add index for last_login_at
CREATE INDEX IF NOT EXISTS idx_user_accounts_last_login ON user_accounts(last_login_at);

-- Update existing records to have default voting_method_preference
UPDATE voters 
SET voting_method_preference = 'ONLINE' 
WHERE voting_method_preference IS NULL;

COMMENT ON COLUMN voters.email IS 'Voter email address (optional, can be updated)';
COMMENT ON COLUMN voters.phone IS 'Voter phone number (optional)';
COMMENT ON COLUMN voters.photo_url IS 'Profile photo URL from storage';
COMMENT ON COLUMN voters.bio IS 'Short biography or description';
COMMENT ON COLUMN voters.voting_method_preference IS 'Preferred voting method: ONLINE or TPS';
COMMENT ON COLUMN user_accounts.last_login_at IS 'Last successful login timestamp';
COMMENT ON COLUMN user_accounts.login_count IS 'Total number of successful logins';
