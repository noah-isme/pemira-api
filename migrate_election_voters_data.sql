-- ====================================================================
-- MIGRATE DATA: voter_status → election_voters
-- Date: 2025-11-26
-- Critical: Populate election_voters table from existing voter_status
-- ====================================================================

BEGIN;

-- Insert voters from voter_status into election_voters
INSERT INTO election_voters (
    election_id,
    voter_id,
    nim,
    status,
    voting_method,
    tps_id,
    checked_in_at,
    voted_at,
    created_at,
    updated_at
)
SELECT 
    vs.election_id,
    vs.voter_id,
    v.nim,
    CASE 
        WHEN vs.has_voted THEN 'VOTED'::election_voter_status
        WHEN vs.is_eligible = false THEN 'REJECTED'::election_voter_status
        ELSE 'VERIFIED'::election_voter_status
    END as status,
    COALESCE(vs.voting_method, v.voting_method, 'ONLINE')::voting_method as voting_method,
    vs.tps_id,
    NULL as checked_in_at,  -- No check-in data in voter_status
    vs.voted_at,
    vs.created_at,
    vs.updated_at
FROM voter_status vs
INNER JOIN voters v ON vs.voter_id = v.id
WHERE NOT EXISTS (
    SELECT 1 FROM election_voters ev 
    WHERE ev.election_id = vs.election_id 
    AND ev.voter_id = vs.voter_id
)
AND v.nim IS NOT NULL;

-- Show migration results
DO $$
DECLARE
    migrated_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO migrated_count FROM election_voters;
    RAISE NOTICE 'Migration completed: % voters migrated to election_voters', migrated_count;
END $$;

COMMIT;

-- Verification
SELECT 
    'election_voters' as table_name,
    election_id,
    COUNT(*) as voter_count
FROM election_voters
GROUP BY election_id
ORDER BY election_id;

SELECT 
    'voter_status' as table_name,
    election_id,
    COUNT(*) as voter_count
FROM voter_status
GROUP BY election_id
ORDER BY election_id;
