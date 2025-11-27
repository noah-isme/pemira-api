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

-- 2. Create lecturers table (for LECTURER identity)
CREATE TABLE IF NOT EXISTS lecturers (
    id BIGSERIAL PRIMARY KEY,
    nidn TEXT UNIQUE,
    name TEXT NOT NULL,
    faculty_code TEXT,
    department TEXT,
    position TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lecturers_nidn ON lecturers(nidn);
CREATE INDEX IF NOT EXISTS idx_lecturers_faculty ON lecturers(faculty_code);

-- 3. Create staff_members table (for STAFF identity)
CREATE TABLE IF NOT EXISTS staff_members (
    id BIGSERIAL PRIMARY KEY,
    nip TEXT UNIQUE,
    name TEXT NOT NULL,
    unit TEXT,
    job_title TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_staff_members_nip ON staff_members(nip);
CREATE INDEX IF NOT EXISTS idx_staff_members_unit ON staff_members(unit);

-- 4. Add triggers for updated_at
CREATE TRIGGER update_students_updated_at
    BEFORE UPDATE ON students
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lecturers_updated_at
    BEFORE UPDATE ON lecturers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_staff_members_updated_at
    BEFORE UPDATE ON staff_members
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
CREATE INDEX IF NOT EXISTS idx_voters_lecturer_id ON voters(lecturer_id);
CREATE INDEX IF NOT EXISTS idx_voters_staff_id ON voters(staff_id);

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
