-- Migration: Create vote_stats table for real-time vote counting
-- Date: 2025-12-19
-- Description: This table tracks vote statistics per candidate per election
--              Used by internal/voting/repository_stats.go

CREATE TABLE IF NOT EXISTS vote_stats (
    election_id BIGINT NOT NULL,
    candidate_id BIGINT NOT NULL,
    total_votes BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT vote_stats_pkey PRIMARY KEY (election_id, candidate_id)
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_vote_stats_election ON vote_stats(election_id);
CREATE INDEX IF NOT EXISTS idx_vote_stats_candidate ON vote_stats(candidate_id);

-- Add foreign key constraints (optional, comment out if causes issues)
-- ALTER TABLE vote_stats 
--     ADD CONSTRAINT fk_vote_stats_election 
--     FOREIGN KEY (election_id) REFERENCES elections(id) ON DELETE CASCADE;

-- ALTER TABLE vote_stats 
--     ADD CONSTRAINT fk_vote_stats_candidate 
--     FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE;

-- Add comment
COMMENT ON TABLE vote_stats IS 'Real-time vote counting statistics per candidate per election';
