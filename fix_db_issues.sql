-- ====================================================================
-- FIX DATABASE ISSUES - PEMIRA API
-- Date: 2025-11-26
-- Description: Fix critical and medium priority database issues
-- ====================================================================

BEGIN;

-- ====================================================================
-- CRITICAL FIX #1: Fix app_settings.updated_by type mismatch
-- ====================================================================
DO $$ BEGIN
    RAISE NOTICE 'Fixing app_settings.updated_by type mismatch...';
END $$;

-- Drop FK constraint temporarily
ALTER TABLE app_settings 
DROP CONSTRAINT IF EXISTS app_settings_updated_by_fkey;

-- Change column type from INT to BIGINT
ALTER TABLE app_settings 
ALTER COLUMN updated_by TYPE BIGINT;

-- Recreate FK constraint
ALTER TABLE app_settings 
ADD CONSTRAINT app_settings_updated_by_fkey 
FOREIGN KEY (updated_by) REFERENCES user_accounts(id);

-- ====================================================================
-- CRITICAL FIX #2: Resolve voting_method duplication
-- ====================================================================
DO $$ BEGIN
    RAISE NOTICE 'Syncing voting_method_preference to voting_method...';
END $$;

-- Sync data from preference to method for valid values
UPDATE voters 
SET voting_method = voting_method_preference::voting_method
WHERE voting_method_preference IS NOT NULL 
AND voting_method_preference IN ('ONLINE', 'TPS')
AND (voting_method IS NULL OR voting_method::text != voting_method_preference);

-- Show affected rows
DO $$ 
DECLARE
    affected_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO affected_count
    FROM voters 
    WHERE voting_method_preference IS NOT NULL 
    AND voting_method::text != voting_method_preference;
    
    RAISE NOTICE 'Synced % voter records', affected_count;
END $$;

-- Drop redundant column
DO $$ BEGIN
    RAISE NOTICE 'Dropping redundant voting_method_preference column...';
END $$;

ALTER TABLE voters 
DROP COLUMN IF EXISTS voting_method_preference;

-- ====================================================================
-- MEDIUM FIX #3: Standardize timestamp types
-- ====================================================================
DO $$ BEGIN
    RAISE NOTICE 'Standardizing timestamp types to TIMESTAMPTZ...';
END $$;

-- Fix app_settings.updated_at
ALTER TABLE app_settings 
ALTER COLUMN updated_at TYPE TIMESTAMPTZ 
USING updated_at AT TIME ZONE 'UTC';

-- Fix user_accounts.last_login_at
ALTER TABLE user_accounts 
ALTER COLUMN last_login_at TYPE TIMESTAMPTZ 
USING last_login_at AT TIME ZONE 'UTC';

-- ====================================================================
-- MEDIUM FIX #4: Clean up duplicate FK constraints
-- ====================================================================
DO $$ BEGIN
    RAISE NOTICE 'Cleaning up duplicate FK constraints...';
END $$;

ALTER TABLE user_accounts 
DROP CONSTRAINT IF EXISTS user_accounts_lecturer_id_fkey;

ALTER TABLE user_accounts 
DROP CONSTRAINT IF EXISTS user_accounts_staff_id_fkey;

-- ====================================================================
-- Verify fixes
-- ====================================================================

DO $$ 
DECLARE
    updated_by_type TEXT;
    voting_pref_exists BOOLEAN;
    app_ts_type TEXT;
    user_ts_type TEXT;
BEGIN
    -- Check updated_by type
    SELECT data_type INTO updated_by_type
    FROM information_schema.columns
    WHERE table_name = 'app_settings' AND column_name = 'updated_by';
    RAISE NOTICE '  app_settings.updated_by type: %', updated_by_type;
    
    -- Check voting_method_preference exists
    SELECT EXISTS(
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'voters' AND column_name = 'voting_method_preference'
    ) INTO voting_pref_exists;
    RAISE NOTICE '  voters.voting_method_preference exists: %', voting_pref_exists;
    
    -- Check timestamp types
    SELECT data_type INTO app_ts_type
    FROM information_schema.columns
    WHERE table_name = 'app_settings' AND column_name = 'updated_at';
    RAISE NOTICE '  app_settings.updated_at type: %', app_ts_type;
    
    SELECT data_type INTO user_ts_type
    FROM information_schema.columns
    WHERE table_name = 'user_accounts' AND column_name = 'last_login_at';
    RAISE NOTICE '  user_accounts.last_login_at type: %', user_ts_type;
END $$;

COMMIT;

DO $$ BEGIN
    RAISE NOTICE 'All fixes applied successfully! ✅';
END $$;
