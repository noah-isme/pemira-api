-- Rollback: Remove master tables and references

-- Remove columns from staff_members
ALTER TABLE staff_members
    DROP COLUMN IF EXISTS position_id,
    DROP COLUMN IF EXISTS unit_id;

-- Remove columns from lecturers
ALTER TABLE lecturers
    DROP COLUMN IF EXISTS position_id,
    DROP COLUMN IF EXISTS unit_id;

-- Drop staff tables
DROP TABLE IF EXISTS staff_positions;
DROP TABLE IF EXISTS staff_units;

-- Drop lecturer tables
DROP TABLE IF EXISTS lecturer_positions;
DROP TABLE IF EXISTS lecturer_units;

-- Drop faculty/prodi tables
DROP TABLE IF EXISTS study_programs;
DROP TABLE IF EXISTS faculties;
