-- Seed data for master tables
-- Based on actual hardcoded data from internal/auth/faculty_programs.go

-- =====================================================
-- 1. FACULTIES
-- =====================================================

INSERT INTO faculties (code, name) VALUES
    ('FAS', 'Fakultas Agama / Syariah'),
    ('FE', 'Fakultas Ekonomi'),
    ('FKIP', 'Fakultas Keguruan & Ilmu Pendidikan (FKIP)'),
    ('FKes', 'Fakultas Kesehatan'),
    ('FT', 'Fakultas Teknik'),
    ('FP', 'Fakultas Pertanian')
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 2. STUDY PROGRAMS
-- =====================================================

-- Fakultas Agama / Syariah
INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'HKI', 'Hukum Keluarga Islam', 'S1'
FROM faculties f WHERE f.code = 'FAS'
ON CONFLICT DO NOTHING;

-- Fakultas Ekonomi
INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'AK', 'Akuntansi', 'S1'
FROM faculties f WHERE f.code = 'FE'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'MJ', 'Manajemen', 'S1'
FROM faculties f WHERE f.code = 'FE'
ON CONFLICT DO NOTHING;

-- Fakultas Keguruan & Ilmu Pendidikan (FKIP)
INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'PKim', 'Pendidikan Kimia', 'S1'
FROM faculties f WHERE f.code = 'FKIP'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'PBI', 'Pendidikan Bahasa Inggris', 'S1'
FROM faculties f WHERE f.code = 'FKIP'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'PMat', 'Pendidikan Matematika', 'S1'
FROM faculties f WHERE f.code = 'FKIP'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'PGPAUD', 'PG PAUD', 'S1'
FROM faculties f WHERE f.code = 'FKIP'
ON CONFLICT DO NOTHING;

-- Fakultas Kesehatan
INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'Bid', 'Kebidanan', 'D3'
FROM faculties f WHERE f.code = 'FKes'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'Kep', 'Keperawatan', 'D3'
FROM faculties f WHERE f.code = 'FKes'
ON CONFLICT DO NOTHING;

-- Fakultas Teknik
INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'TI', 'Teknik Informatika', 'S1'
FROM faculties f WHERE f.code = 'FT'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'TIn', 'Teknik Industri', 'S1'
FROM faculties f WHERE f.code = 'FT'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'TM', 'Teknik Mesin', 'S1'
FROM faculties f WHERE f.code = 'FT'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'TS', 'Teknik Sipil', 'S1'
FROM faculties f WHERE f.code = 'FT'
ON CONFLICT DO NOTHING;

-- Fakultas Pertanian
INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'Agro', 'Agroteknologi', 'S1'
FROM faculties f WHERE f.code = 'FP'
ON CONFLICT DO NOTHING;

INSERT INTO study_programs (faculty_id, code, name, level) 
SELECT f.id, 'Agri', 'Agribisnis', 'S1'
FROM faculties f WHERE f.code = 'FP'
ON CONFLICT DO NOTHING;

-- =====================================================
-- 3. LECTURER UNITS
-- =====================================================

INSERT INTO lecturer_units (code, name) VALUES
    ('FAS', 'Fakultas Agama / Syariah'),
    ('FE', 'Fakultas Ekonomi'),
    ('FKIP', 'Fakultas Keguruan & Ilmu Pendidikan'),
    ('FKes', 'Fakultas Kesehatan'),
    ('FT', 'Fakultas Teknik'),
    ('FP', 'Fakultas Pertanian'),
    ('LPPM', 'Lembaga Penelitian dan Pengabdian Masyarakat'),
    ('LPM', 'Lembaga Penjaminan Mutu'),
    ('PPS', 'Program Pascasarjana')
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 4. LECTURER POSITIONS
-- =====================================================

-- Jabatan Fungsional
INSERT INTO lecturer_positions (category, code, name) VALUES
    ('FUNGSIONAL', 'AA', 'Asisten Ahli'),
    ('FUNGSIONAL', 'L', 'Lektor'),
    ('FUNGSIONAL', 'LK', 'Lektor Kepala'),
    ('FUNGSIONAL', 'GB', 'Guru Besar'),
    ('FUNGSIONAL', 'TA', 'Tenaga Pengajar')
ON CONFLICT DO NOTHING;

-- Jabatan Struktural
INSERT INTO lecturer_positions (category, code, name) VALUES
    ('STRUKTURAL', 'REKTOR', 'Rektor'),
    ('STRUKTURAL', 'WREK1', 'Wakil Rektor I'),
    ('STRUKTURAL', 'WREK2', 'Wakil Rektor II'),
    ('STRUKTURAL', 'WREK3', 'Wakil Rektor III'),
    ('STRUKTURAL', 'DEKAN', 'Dekan'),
    ('STRUKTURAL', 'WDEK1', 'Wakil Dekan I'),
    ('STRUKTURAL', 'WDEK2', 'Wakil Dekan II'),
    ('STRUKTURAL', 'WDEK3', 'Wakil Dekan III'),
    ('STRUKTURAL', 'KAPRODI', 'Kepala Program Studi'),
    ('STRUKTURAL', 'SEKPRODI', 'Sekretaris Program Studi'),
    ('STRUKTURAL', 'KELOMPOK', 'Kepala Kelompok Dosen')
ON CONFLICT DO NOTHING;

-- =====================================================
-- 5. STAFF UNITS
-- =====================================================

INSERT INTO staff_units (code, name) VALUES
    ('BAU', 'Biro Administrasi Umum'),
    ('BAAK', 'Biro Administrasi Akademik dan Kemahasiswaan'),
    ('BAK', 'Biro Administrasi Keuangan'),
    ('BAPSI', 'Biro Administrasi Perencanaan dan Sistem Informasi'),
    ('LPPM', 'Lembaga Penelitian dan Pengabdian Masyarakat'),
    ('LPM', 'Lembaga Penjaminan Mutu'),
    ('UPT-TIK', 'Unit Pelaksana Teknis - Teknologi Informasi dan Komunikasi'),
    ('UPT-Perpus', 'Unit Pelaksana Teknis - Perpustakaan'),
    ('UPT-Bahasa', 'Unit Pelaksana Teknis - Pusat Bahasa')
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 6. STAFF POSITIONS
-- =====================================================

INSERT INTO staff_positions (code, name) VALUES
    ('KEPALA_BIRO', 'Kepala Biro'),
    ('KABAG', 'Kepala Bagian'),
    ('KASUBAG', 'Kepala Sub Bagian'),
    ('KEPALA_UPT', 'Kepala Unit Pelaksana Teknis'),
    ('STAF_SENIOR', 'Staf Senior'),
    ('STAF', 'Staf'),
    ('OPERATOR', 'Operator'),
    ('TEKNISI', 'Teknisi')
ON CONFLICT (code) DO NOTHING;
