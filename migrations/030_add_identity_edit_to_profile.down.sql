-- +goose Down
-- =========================================================
-- ROLLBACK IDENTITY FIELDS EDITABLE FROM PROFILE
-- =========================================================

-- Drop triggers
DROP TRIGGER IF EXISTS sync_staff_on_voter_update ON voters;
DROP TRIGGER IF EXISTS sync_lecturer_on_voter_update ON voters;
DROP TRIGGER IF EXISTS sync_student_on_voter_update ON voters;

-- Drop functions
DROP FUNCTION IF EXISTS sync_staff_from_voter();
DROP FUNCTION IF EXISTS sync_lecturer_from_voter();
DROP FUNCTION IF EXISTS sync_student_from_voter();

-- Remove comments
COMMENT ON COLUMN voters.faculty_code IS NULL;
COMMENT ON COLUMN voters.study_program_code IS NULL;
COMMENT ON COLUMN voters.cohort_year IS NULL;
COMMENT ON COLUMN voters.class_label IS NULL;
