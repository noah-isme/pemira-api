# Master Tables Migration

## Overview
Migrasi ini menggantikan hardcoded values untuk fakultas, prodi, unit, dan jabatan dengan tabel master yang relasional.

## Struktur Tabel

### 1. Fakultas & Program Studi (Mahasiswa & Dosen)

#### `faculties`
- `id`: Primary key
- `code`: Kode fakultas (FTI, FEB, FKIP, dst) - UNIQUE
- `name`: Nama lengkap fakultas

#### `study_programs`
- `id`: Primary key
- `faculty_id`: Foreign key ke faculties
- `code`: Kode prodi (TI, SI, MI, dst)
- `name`: Nama lengkap program studi
- `level`: Jenjang (S1, D3, S2, S3)
- **UNIQUE INDEX**: `(faculty_id, code)`

### 2. Unit & Jabatan Dosen

#### `lecturer_units`
- `id`: Primary key
- `code`: Kode unit (FTI, LPPM, PPS, dst) - UNIQUE
- `name`: Nama lengkap unit

#### `lecturer_positions`
- `id`: Primary key
- `category`: FUNGSIONAL / STRUKTURAL (CHECK constraint)
- `code`: Kode jabatan (AA, LEKTOR, DEKAN, dst)
- `name`: Nama lengkap jabatan
- **UNIQUE INDEX**: `(category, code)`

#### Perubahan pada `lecturers` table
```sql
ALTER TABLE lecturers
  ADD COLUMN unit_id BIGINT REFERENCES lecturer_units(id) ON DELETE SET NULL,
  ADD COLUMN position_id BIGINT REFERENCES lecturer_positions(id) ON DELETE SET NULL;
```

### 3. Unit & Jabatan Staf

#### `staff_units`
- `id`: Primary key
- `code`: Kode unit (BAU, BAAK, LPPM, UPT-TIK, dst) - UNIQUE
- `name`: Nama lengkap unit

#### `staff_positions`
- `id`: Primary key
- `code`: Kode jabatan (KEPALA_BIRO, KABAG, STAF, dst) - UNIQUE
- `name`: Nama lengkap jabatan

#### Perubahan pada `staff_members` table
```sql
ALTER TABLE staff_members
  ADD COLUMN unit_id BIGINT REFERENCES staff_units(id) ON DELETE SET NULL,
  ADD COLUMN position_id BIGINT REFERENCES staff_positions(id) ON DELETE SET NULL;
```

## Cara Menjalankan

### 1. Jalankan Migrasi
```bash
# Menggunakan golang-migrate atau tool migrasi lainnya
migrate -path migrations -database "postgresql://..." up
```

### 2. Seed Data Master
```bash
# Jalankan seed file untuk data awal
psql -U user -d database -f seeds/028_seed_master_tables.sql
```

## Integrasi dengan Frontend

### Form Mahasiswa
1. **Dropdown Fakultas**: Query dari `faculties` table
2. **Dropdown Prodi**: Query dari `study_programs` WHERE `faculty_id = selected_faculty`

### Form Dosen
1. **Dropdown Unit**: Query dari `lecturer_units` table
2. **Dropdown Jabatan**: Query dari `lecturer_positions` table
   - Optional: Filter by `category` (FUNGSIONAL/STRUKTURAL)

### Form Staf
1. **Dropdown Unit**: Query dari `staff_units` table
2. **Dropdown Jabatan**: Query dari `staff_positions` table

## API Endpoints yang Perlu Ditambahkan

```
GET /api/master/faculties
GET /api/master/study-programs?faculty_id={id}
GET /api/master/lecturer-units
GET /api/master/lecturer-positions?category={FUNGSIONAL|STRUKTURAL}
GET /api/master/staff-units
GET /api/master/staff-positions
```

## Rollback

Jika perlu rollback:
```bash
migrate -path migrations -database "postgresql://..." down 1
```

Migration down akan:
1. Remove kolom baru dari `staff_members` dan `lecturers`
2. Drop semua master tables

## Data Migration (Existing Records)

Untuk data yang sudah ada (jika ada), perlu dibuat script terpisah untuk:
1. Map hardcoded values ke master table IDs
2. Update existing records dengan foreign keys yang sesuai
3. Remove old hardcoded columns (jika ada)

## Notes

- Semua tabel master memiliki `created_at` dan `updated_at` timestamps
- Foreign keys menggunakan `ON DELETE SET NULL` untuk mencegah data loss
- `ON DELETE CASCADE` digunakan pada `study_programs` → `faculties` karena prodi tidak bisa exist tanpa fakultas
- Unique constraints mencegah duplikasi data
