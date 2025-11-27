-- +goose Up
-- =========================================================
-- ADD IDENTITY FIELDS EDITABLE FROM PROFILE
-- Allow voters to update their identity information
-- =========================================================

-- 1. Add function to sync student data when voter updates profile
CREATE OR REPLACE FUNCTION sync_student_from_voter()
RETURNS TRIGGER AS $$
BEGIN
    -- If voter has student_id, update students table
    IF NEW.student_id IS NOT NULL THEN
        UPDATE students
        SET 
            faculty_code = NEW.faculty_code,
            program_code = NEW.study_program_code,
            cohort_year = NEW.cohort_year,
            class_label = NEW.class_label,
            updated_at = NOW()
        WHERE id = NEW.student_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Add function to sync lecturer data when voter updates profile
CREATE OR REPLACE FUNCTION sync_lecturer_from_voter()
RETURNS TRIGGER AS $$
BEGIN
    -- If voter has lecturer_id, update lecturers table
    IF NEW.lecturer_id IS NOT NULL THEN
        UPDATE lecturers
        SET 
            faculty_code = NEW.faculty_code,
            department_code = COALESCE(NEW.study_program_code, department_code),
            position = COALESCE(NEW.class_label, position),
            updated_at = NOW()
        WHERE id = NEW.lecturer_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 3. Add function to sync staff data when voter updates profile
CREATE OR REPLACE FUNCTION sync_staff_from_voter()
RETURNS TRIGGER AS $$
BEGIN
    -- If voter has staff_id, update staff_members table
    IF NEW.staff_id IS NOT NULL THEN
        UPDATE staff_members
        SET 
            unit_code = NEW.faculty_code,
            unit_name = NEW.faculty_name,
            position = COALESCE(NEW.class_label, position),
            updated_at = NOW()
        WHERE id = NEW.staff_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. Add triggers to voters table for auto-sync
DROP TRIGGER IF EXISTS sync_student_on_voter_update ON voters;
CREATE TRIGGER sync_student_on_voter_update
    AFTER UPDATE OF faculty_code, study_program_code, cohort_year, class_label ON voters
    FOR EACH ROW
    WHEN (NEW.voter_type = 'STUDENT')
    EXECUTE FUNCTION sync_student_from_voter();

DROP TRIGGER IF EXISTS sync_lecturer_on_voter_update ON voters;
CREATE TRIGGER sync_lecturer_on_voter_update
    AFTER UPDATE OF faculty_code, study_program_code, class_label ON voters
    FOR EACH ROW
    WHEN (NEW.voter_type = 'LECTURER')
    EXECUTE FUNCTION sync_lecturer_from_voter();

DROP TRIGGER IF EXISTS sync_staff_on_voter_update ON voters;
CREATE TRIGGER sync_staff_on_voter_update
    AFTER UPDATE OF faculty_code, faculty_name, class_label ON voters
    FOR EACH ROW
    WHEN (NEW.voter_type = 'STAFF')
    EXECUTE FUNCTION sync_staff_from_voter();

-- 5. Sync current data from identity tables to voters (for display)
-- For students
UPDATE voters v
SET 
    faculty_code = s.faculty_code,
    study_program_code = s.program_code,
    cohort_year = s.cohort_year,
    class_label = s.class_label
FROM students s
WHERE v.student_id = s.id 
  AND v.voter_type = 'STUDENT'
  AND (
      v.faculty_code IS DISTINCT FROM s.faculty_code OR
      v.study_program_code IS DISTINCT FROM s.program_code OR
      v.cohort_year IS DISTINCT FROM s.cohort_year OR
      v.class_label IS DISTINCT FROM s.class_label
  );

-- For lecturers (map department_code to study_program_code, position to class_label)
UPDATE voters v
SET 
    faculty_code = l.faculty_code,
    study_program_code = l.department_code,
    class_label = l.position
FROM lecturers l
WHERE v.lecturer_id = l.id 
  AND v.voter_type = 'LECTURER'
  AND (
      v.faculty_code IS DISTINCT FROM l.faculty_code OR
      v.study_program_code IS DISTINCT FROM l.department_code OR
      v.class_label IS DISTINCT FROM l.position
  );

-- For staff (map unit_code to faculty_code, unit_name to faculty_name, position to class_label)
UPDATE voters v
SET 
    faculty_code = sm.unit_code,
    faculty_name = sm.unit_name,
    class_label = sm.position
FROM staff_members sm
WHERE v.staff_id = sm.id 
  AND v.voter_type = 'STAFF'
  AND (
      v.faculty_code IS DISTINCT FROM sm.unit_code OR
      v.faculty_name IS DISTINCT FROM sm.unit_name OR
      v.class_label IS DISTINCT FROM sm.position
  );

-- 6. Add comment for documentation
COMMENT ON COLUMN voters.faculty_code IS 'Editable by voter. For STUDENT: faculty code, LECTURER: faculty code, STAFF: unit code';
COMMENT ON COLUMN voters.study_program_code IS 'Editable by voter. For STUDENT: program code, LECTURER: department, STAFF: not used';
COMMENT ON COLUMN voters.cohort_year IS 'Editable by voter. For STUDENT only: enrollment year';
COMMENT ON COLUMN voters.class_label IS 'Editable by voter. For STUDENT: class, LECTURER: position/rank, STAFF: job position';
