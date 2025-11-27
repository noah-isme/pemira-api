-- Rollback: Remove profile fields from voters and users tables

-- Drop indexes
DROP INDEX IF EXISTS idx_voters_email;
DROP INDEX IF EXISTS idx_voters_updated_at;
DROP INDEX IF EXISTS idx_user_accounts_last_login;

-- Remove columns from user_accounts table
ALTER TABLE user_accounts
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS login_count;

-- Remove columns from voters table
ALTER TABLE voters
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS photo_url,
DROP COLUMN IF EXISTS bio,
DROP COLUMN IF EXISTS voting_method_preference;
