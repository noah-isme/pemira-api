-- +goose Down
-- =========================================================
-- ROLLBACK VOTERS SCHEMA REFACTOR
-- =========================================================

-- 1. Drop constraint
ALTER TABLE voters DROP CONSTRAINT IF EXISTS chk_voters_single_identity;

-- 2. Drop indexes
DROP INDEX IF EXISTS idx_voters_student_id;
DROP INDEX IF EXISTS idx_voters_lecturer_id;
DROP INDEX IF EXISTS idx_voters_staff_id;

-- 3. Drop foreign key columns
ALTER TABLE voters DROP COLUMN IF EXISTS student_id;
ALTER TABLE voters DROP COLUMN IF EXISTS lecturer_id;
ALTER TABLE voters DROP COLUMN IF EXISTS staff_id;

-- 4. Recreate old indexes
CREATE UNIQUE INDEX IF NOT EXISTS ux_voters_student_nim
    ON voters (nim)
    WHERE voter_type = 'STUDENT' AND nim IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_voters_faculty
    ON voters (faculty_code, study_program_code);

CREATE INDEX IF NOT EXISTS idx_voters_cohort
    ON voters (cohort_year);

-- 5. Drop identity tables (careful: this will lose data!)
DROP TRIGGER IF EXISTS update_staff_members_updated_at ON staff_members;
DROP TRIGGER IF EXISTS update_lecturers_updated_at ON lecturers;
DROP TRIGGER IF EXISTS update_students_updated_at ON students;

DROP TABLE IF EXISTS staff_members CASCADE;
DROP TABLE IF EXISTS lecturers CASCADE;
DROP TABLE IF EXISTS students CASCADE;
