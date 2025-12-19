# PERBANDINGAN DATABASE PRODUCTION vs LOCAL
Tanggal: 2025-12-19

## RINGKASAN PERBEDAAN

### 1. SKEMA DATABASE
- **PRODUCTION**: Menggunakan 2 schema
  - `myschema`: 25 tabel (tabel lama/utama)
  - `public`: 6 tabel (tabel baru)
- **LOCAL**: Semua 32 tabel di schema `public`

### 2. TABEL YANG HILANG DI PRODUCTION
❌ **vote_stats** - Tabel ini ADA di local tapi TIDAK ADA di production

Definisi tabel vote_stats (local):
```sql
CREATE TABLE vote_stats (
    election_id BIGINT NOT NULL,
    candidate_id BIGINT NOT NULL,
    total_votes BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    PRIMARY KEY (election_id, candidate_id)
);
```

### 3. DISTRIBUSI TABEL

#### PRODUCTION (myschema) - 25 tabel:
- app_settings
- branding_files, branding_settings
- candidate_media, candidate_qr_codes, candidates
- election_voters, elections
- faculties
- lecturer_positions, lecturer_units, lecturers
- migration_history, registration_tokens, schema_migrations
- staff_members, staff_positions, staff_units
- students, study_programs
- tps, tps_checkins
- user_accounts, voters, votes

#### PRODUCTION (public) - 6 tabel:
- tps_ballot_scans
- tps_qr
- user_sessions
- vote_tokens
- voter_status
- voter_tps_qr

#### LOCAL (public) - 32 tabel:
Semua tabel di atas + **vote_stats**

### 4. MASALAH KRITIS

1. **Tabel vote_stats tidak ada di production**
   - Tabel ini mungkin diperlukan untuk statistik voting
   - Perlu migrasi untuk membuat tabel ini

2. **Perbedaan schema**
   - Production: split antara myschema dan public
   - Local: semua di public
   - Ini bisa menyebabkan masalah query jika tidak handle dengan benar

3. **Migrasi belum sinkron**
   - Local punya tabel baru yang belum ada di production
   - Perlu apply migrasi yang missing

## REKOMENDASI ACTION

### URGENT:
1. ✅ Cek apakah aplikasi menggunakan tabel vote_stats
2. ✅ Jika ya, buat migrasi untuk production
3. ✅ Apply semua migrasi yang missing ke production

### YANG PERLU DILAKUKAN:
```bash
# 1. Cek file migrasi yang membuat vote_stats
ls -la migrations/ | grep vote_stats

# 2. Buat SQL untuk production (dengan schema myschema atau public sesuai kebutuhan)
# 3. Apply migrasi ke production dengan hati-hati
```

### VERIFIKASI:
- Cek apakah ada code yang menggunakan vote_stats
- Pastikan semua referensi ke tabel menggunakan schema yang benar
- Test di staging dulu sebelum production

## NEXT STEPS

Apakah Anda ingin:
1. Mencari file migrasi yang membuat vote_stats?
2. Membuat SQL script untuk menambah tabel ke production?
3. Memeriksa code yang menggunakan vote_stats?
