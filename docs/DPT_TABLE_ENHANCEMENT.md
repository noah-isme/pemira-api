# DPT Table Enhancement - Complete Data for Frontend

**Date:** 2025-11-26  
**Status:** ✅ COMPLETED  

---

## Frontend Table Requirements

Frontend DPT table displays these columns:

| No | Column | Source | Status |
|----|--------|--------|--------|
| 1 | No | Row number (UI) | ✅ N/A |
| 2 | NIM/NIDN/NIP | `ev.nim` | ✅ Already included |
| 3 | Nama | `v.name` | ✅ Already included |
| 4 | Fakultas | `v.faculty_name` | ⚠️ **ADDED** |
| 5 | Prodi | `v.study_program_name` | ⚠️ **ADDED** |
| 6 | Semester | Calculated from `v.cohort_year` | ℹ️ Frontend calc |
| 7 | Tipe Pemilih | `v.voter_type` | ✅ Already included |
| 8 | Akademik | `v.academic_status` | ⚠️ **ADDED** |
| 9 | Status Verifikasi | `ev.status` | ✅ Already included |
| 10 | Status Suara | `vs.has_voted` | ⚠️ **ADDED** |
| 11 | Metode | `ev.voting_method` | ✅ Already included |
| 12 | Aksi | UI actions | ✅ N/A |
| 13 | Terakhir Vote | `ev.voted_at` | ✅ Already included |

---

## Changes Made

### 1. Updated Model

**File:** `internal/electionvoter/models.go`

**Added fields:**
```go
type ElectionVoter struct {
    // ... existing fields ...
    FacultyName      *string   `json:"faculty_name,omitempty"`        // NEW
    StudyProgramName *string   `json:"study_program_name,omitempty"`  // NEW
    AcademicStatus   *string   `json:"academic_status,omitempty"`     // NEW
    HasVoted         *bool     `json:"has_voted,omitempty"`           // NEW
}
```

### 2. Enhanced SQL Query

**File:** `internal/electionvoter/repository_pgx.go`

**Before:**
```sql
SELECT
    ev.id, ev.election_id, ev.voter_id, ev.nim,
    ev.status, ev.voting_method, ev.tps_id,
    ev.checked_in_at, ev.voted_at, ev.updated_at,
    v.voter_type, v.name, v.email, v.faculty_code, 
    v.study_program_code, v.cohort_year
FROM election_voters ev
JOIN voters v ON v.id = ev.voter_id
```

**After:**
```sql
SELECT
    ev.id, ev.election_id, ev.voter_id, ev.nim,
    ev.status, ev.voting_method, ev.tps_id,
    ev.checked_in_at, ev.voted_at, ev.updated_at,
    v.voter_type, v.name, v.email, 
    v.faculty_code, v.faculty_name,                    -- ADDED
    v.study_program_code, v.study_program_name,        -- ADDED
    v.cohort_year, v.academic_status,                  -- ADDED
    vs.has_voted                                        -- ADDED
FROM election_voters ev
JOIN voters v ON v.id = ev.voter_id
LEFT JOIN voter_status vs ON vs.election_id = ev.election_id 
                          AND vs.voter_id = ev.voter_id  -- ADDED JOIN
```

### 3. Updated Scan Logic

Added scanning for new fields:
```go
var facultyName sql.NullString
var studyProgramName sql.NullString
var academicStatus sql.NullString
var hasVoted sql.NullBool

// ... scan all fields including new ones ...

item.FacultyName = nullableStringPtr(facultyName)
item.StudyProgramName = nullableStringPtr(studyProgramName)
item.AcademicStatus = nullableStringPtr(academicStatus)
item.HasVoted = nullableBoolPtr(hasVoted)
```

### 4. Added Helper Function

```go
func nullableBoolPtr(nb sql.NullBool) *bool {
    if nb.Valid {
        val := nb.Bool
        return &val
    }
    return nil
}
```

---

## API Response Example

**Before:**
```json
{
  "items": [
    {
      "election_voter_id": 6,
      "nim": "2021101",
      "name": "Agus Santoso",
      "voter_type": "STUDENT",
      "faculty_code": "FT",
      "study_program_code": "TI",
      "cohort_year": 2021,
      "status": "VERIFIED",
      "voting_method": "ONLINE"
    }
  ]
}
```

**After:**
```json
{
  "items": [
    {
      "election_voter_id": 6,
      "nim": "2021101",
      "name": "Agus Santoso",
      "voter_type": "STUDENT",
      "faculty_code": "FT",
      "faculty_name": "Fakultas Teknik",              // NEW
      "study_program_code": "TI",
      "study_program_name": "Teknik Informatika",     // NEW
      "cohort_year": 2021,
      "academic_status": "ACTIVE",                    // NEW
      "status": "VERIFIED",
      "voting_method": "ONLINE",
      "has_voted": false                              // NEW
    }
  ]
}
```

---

## Frontend Integration

### Semester Calculation

Frontend should calculate semester from `cohort_year`:

```typescript
function calculateSemester(cohortYear: number | null): string {
  if (!cohortYear) return '-';
  
  const currentYear = new Date().getFullYear();
  const currentMonth = new Date().getMonth() + 1;
  
  const yearsEnrolled = currentYear - cohortYear;
  let semester = yearsEnrolled * 2;
  
  // If after August, add 1 semester (ganjil)
  if (currentMonth >= 8) {
    semester += 1;
  }
  
  return semester.toString();
}
```

### Table Mapping

```typescript
interface DPTRow {
  no: number;                           // index + 1
  nim: string;                          // nim
  nama: string;                         // name
  fakultas: string;                     // faculty_name
  prodi: string;                        // study_program_name
  semester: string;                     // calculated from cohort_year
  tipePemilih: string;                  // voter_type
  akademik: string;                     // academic_status
  statusVerifikasi: string;             // status
  statusSuara: boolean;                 // has_voted
  metode: string;                       // voting_method
  terakhirVote: string | null;          // voted_at
}
```

---

## Verification

### Test Query

```bash
curl 'http://localhost:8080/api/v1/admin/elections/1/voters?page=1&limit=5' \
  -H "Authorization: Bearer YOUR_TOKEN" | jq '.'
```

### Expected Fields in Response

- ✅ `nim`
- ✅ `name`
- ✅ `faculty_code`
- ✅ `faculty_name` ← **NEW**
- ✅ `study_program_code`
- ✅ `study_program_name` ← **NEW**
- ✅ `cohort_year`
- ✅ `voter_type`
- ✅ `academic_status` ← **NEW**
- ✅ `status`
- ✅ `voting_method`
- ✅ `has_voted` ← **NEW**
- ✅ `voted_at`

---

## Database Schema Notes

### Tables Involved

1. **election_voters** - Main table for voter enrollment
   - Columns: id, election_id, voter_id, nim, status, voting_method, tps_id, voted_at, etc.

2. **voters** - Voter personal information
   - Columns: id, nim, name, faculty_code, faculty_name, study_program_code, study_program_name, cohort_year, voter_type, academic_status, etc.

3. **voter_status** - Voting status tracking
   - Columns: id, election_id, voter_id, has_voted, voted_at, etc.

### Data Consistency

**Important:** System now uses TWO tables for tracking:
- `election_voters`: Registration & enrollment status
- `voter_status`: Actual voting status (has_voted)

Both must be kept in sync when voter votes!

---

## Deployment Checklist

### Development
- [x] Model updated
- [x] Repository query enhanced
- [x] Scan logic updated
- [x] Helper function added
- [x] Code compiled successfully
- [ ] Application restarted
- [ ] Frontend tested with new fields

### Production
- [ ] Code deployed
- [ ] API tested
- [ ] Frontend updated to use new fields
- [ ] Semester calculation implemented
- [ ] Table display verified

---

## Troubleshooting

### If fields still missing:

1. **Restart application:**
   ```bash
   pkill -f "go run cmd/api/main.go"
   go run cmd/api/main.go
   ```

2. **Check database has data:**
   ```sql
   SELECT v.faculty_name, v.study_program_name, v.academic_status
   FROM voters v
   LIMIT 5;
   ```

3. **Verify JOIN works:**
   ```sql
   SELECT ev.nim, v.name, vs.has_voted
   FROM election_voters ev
   JOIN voters v ON v.id = ev.voter_id
   LEFT JOIN voter_status vs ON vs.election_id = ev.election_id 
                             AND vs.voter_id = ev.voter_id
   WHERE ev.election_id = 1
   LIMIT 5;
   ```

---

## Related Files

- Model: `internal/electionvoter/models.go`
- Repository: `internal/electionvoter/repository_pgx.go`
- API Handler: `internal/electionvoter/http_handler.go`
- Endpoint: `GET /api/v1/admin/elections/{electionID}/voters`

---

**Enhanced By:** AI Assistant  
**Impact:** Complete data for frontend DPT table  
**Breaking Changes:** None (added fields only)  
**Status:** ✅ Ready to Deploy
