-- Migration: Create master tables for faculties, study programs, lecturer and staff units/positions
-- Replaces hardcoded values with relational master data

-- =====================================================
-- 1. MASTER FAKULTAS & PRODI (for mahasiswa & dosen)
-- =====================================================

CREATE TABLE faculties (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,       -- FTI, FEB, FKIP, dst
    name TEXT NOT NULL,              -- Fakultas Teknik, Fakultas Ekonomi dan Bisnis, ...
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE study_programs (
    id BIGSERIAL PRIMARY KEY,
    faculty_id BIGINT NOT NULL REFERENCES faculties(id) ON DELETE CASCADE,
    code TEXT NOT NULL,              -- TI, SI, MI, etc
    name TEXT NOT NULL,              -- Teknik Informatika, Sistem Informasi, ...
    level TEXT NOT NULL,             -- S1, D3, S2, S3
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX ux_study_programs_faculty_code
    ON study_programs (faculty_id, code);

-- =====================================================
-- 2. MASTER UNIT & JABATAN DOSEN
-- =====================================================

CREATE TABLE lecturer_units (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,       -- "FTI", "LP2M", "PPS", dll
    name TEXT NOT NULL,              -- "Fakultas Teknik & Informatika", "LPPM", dll
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lecturer_positions (
    id BIGSERIAL PRIMARY KEY,
    category TEXT NOT NULL           -- 'FUNGSIONAL' / 'STRUKTURAL'
        CHECK (category IN ('FUNGSIONAL','STRUKTURAL')),
    code TEXT NOT NULL,              -- "AA", "LEKTOR", "DEKAN", dst
    name TEXT NOT NULL,              -- "Asisten Ahli", "Lektor", "Dekan", ...
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX ux_lecturer_positions_code
    ON lecturer_positions (category, code);

-- Add references to lecturers table
ALTER TABLE lecturers
    ADD COLUMN unit_id BIGINT REFERENCES lecturer_units(id) ON DELETE SET NULL,
    ADD COLUMN position_id BIGINT REFERENCES lecturer_positions(id) ON DELETE SET NULL;

-- =====================================================
-- 3. MASTER UNIT & JABATAN STAF
-- =====================================================

CREATE TABLE staff_units (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,       -- "BAU", "BAAK", "LPPM", "LPM", "UPT-TIK", ...
    name TEXT NOT NULL,              -- "Biro Administrasi Umum", ...
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE staff_positions (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,       -- "KEPALA_BIRO", "KABAG", "KASUBAG", "STAF", ...
    name TEXT NOT NULL,              -- "Kepala Biro", "Kepala Bagian", ...
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add references to staff_members table
ALTER TABLE staff_members
    ADD COLUMN unit_id BIGINT REFERENCES staff_units(id) ON DELETE SET NULL,
    ADD COLUMN position_id BIGINT REFERENCES staff_positions(id) ON DELETE SET NULL;
