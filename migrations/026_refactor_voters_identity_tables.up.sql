-- +goose Up
-- =========================================================
-- REFACTOR VOTERS SCHEMA
-- Create identity tables and refactor voters table
-- =========================================================

-- 1. Create students table (for STUDENT identity)
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    nim TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    faculty_code TEXT,
    program_code TEXT,
    cohort_year INT,
    class_label TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_students_nim ON students(nim);
CREATE INDEX IF NOT EXISTS idx_students_faculty ON students(faculty_code, program_code);
CREATE INDEX IF NOT EXISTS idx_students_cohort ON students(cohort_year);

-- 2-4. lecturers/staff_members sudah dibuat di migrasi 008; tidak diubah di sini
--     hanya tambahkan trigger untuk students
CREATE TRIGGER update_students_updated_at
    BEFORE UPDATE ON students
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 5. Migrate existing voters data to students table
-- Only migrate STUDENT type voters
INSERT INTO students (nim, name, faculty_code, program_code, cohort_year, class_label, created_at, updated_at)
SELECT 
    nim,
    name,
    faculty_code,
    study_program_code,
    cohort_year,
    class_label,
    created_at,
    updated_at
FROM voters
WHERE voter_type = 'STUDENT'
ON CONFLICT (nim) DO NOTHING;

-- 6. Add foreign key columns to voters table
ALTER TABLE voters ADD COLUMN IF NOT EXISTS student_id BIGINT REFERENCES students(id) ON DELETE SET NULL;
ALTER TABLE voters ADD COLUMN IF NOT EXISTS lecturer_id BIGINT REFERENCES lecturers(id) ON DELETE SET NULL;
ALTER TABLE voters ADD COLUMN IF NOT EXISTS staff_id BIGINT REFERENCES staff_members(id) ON DELETE SET NULL;

-- 7. Link existing voters to students
UPDATE voters v
SET student_id = s.id
FROM students s
WHERE v.nim = s.nim 
  AND v.voter_type = 'STUDENT'
  AND v.student_id IS NULL;

-- 8. Add indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_voters_student_id ON voters(student_id);
-- idx_voters_lecturer_id dan idx_voters_staff_id sudah ada dari migrasi 008

-- 9. Add constraint to ensure one identity type per voter
ALTER TABLE voters ADD CONSTRAINT chk_voters_single_identity
    CHECK (
        (CASE WHEN student_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN lecturer_id IS NOT NULL THEN 1 ELSE 0 END +
         CASE WHEN staff_id IS NOT NULL THEN 1 ELSE 0 END) <= 1
    );

-- 10. Drop old indexes that are no longer needed
DROP INDEX IF EXISTS ux_voters_nim;
DROP INDEX IF EXISTS ux_voters_student_nim;
DROP INDEX IF EXISTS idx_voters_faculty;
DROP INDEX IF EXISTS idx_voters_cohort;

-- Note: We keep the old columns (nim, faculty_code, etc.) for backward compatibility
-- They can be removed in a future migration after full transition
