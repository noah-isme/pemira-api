# Quick Reference - Identity Fields Update

## Field Mapping Per Voter Type

### 🎓 STUDENT (Mahasiswa)
| Profile Field | Identity Table Column | Editable | Example |
|---------------|----------------------|----------|---------|
| `faculty_code` | `students.faculty_code` | ✅ Yes | "FTI", "FEB" |
| `study_program_code` | `students.program_code` | ✅ Yes | "IF", "SI", "TI" |
| `cohort_year` | `students.cohort_year` | ✅ Yes | 2021, 2022 |
| `class_label` | `students.class_label` | ✅ Yes | "IF-A", "SI-B" |

**API Request Example:**
```json
{
  "faculty_code": "FTI",
  "study_program_code": "IF",
  "cohort_year": 2022,
  "class_label": "IF-B"
}
```

---

### 👨‍🏫 LECTURER (Dosen)
| Profile Field | Identity Table Column | Editable | Example |
|---------------|----------------------|----------|---------|
| `faculty_code` | `lecturers.faculty_code` | ✅ Yes | "FTI", "FEB" |
| `study_program_code` | `lecturers.department_code` | ✅ Yes | "Informatika" |
| `class_label` | `lecturers.position` | ✅ Yes | "Lektor", "Lektor Kepala" |
| `cohort_year` | - | ❌ No | Not applicable |

**API Request Example:**
```json
{
  "faculty_code": "FTI",
  "study_program_code": "Informatika",
  "class_label": "Lektor Kepala"
}
```

---

### 👔 STAFF (Staf)
| Profile Field | Identity Table Column | Editable | Example |
|---------------|----------------------|----------|---------|
| `faculty_code` | `staff_members.unit_code` | ✅ Yes | "BAU", "BAK" |
| `faculty_name` | `staff_members.unit_name` | ✅ Yes | "Bagian Administrasi Umum" |
| `class_label` | `staff_members.position` | ✅ Yes | "Koordinator", "Staf Senior" |
| `study_program_code` | - | ❌ No | Not applicable |
| `cohort_year` | - | ❌ No | Not applicable |

**API Request Example:**
```json
{
  "faculty_code": "BAU",
  "faculty_name": "Bagian Administrasi Umum",
  "class_label": "Koordinator"
}
```

---

## Complete Update Profile Request

### All Fields (Optional)
```json
{
  "email": "newemail@example.com",
  "phone": "081234567890",
  "photo_url": "https://storage.com/photo.jpg",
  "faculty_code": "FTI",
  "study_program_code": "IF",
  "cohort_year": 2022,
  "class_label": "IF-A"
}
```

### Only Identity Fields
```json
{
  "faculty_code": "FEB",
  "study_program_code": "Manajemen",
  "cohort_year": 2021,
  "class_label": "MJ-C"
}
```

### Only Contact Fields
```json
{
  "email": "newemail@example.com",
  "phone": "+6281234567890"
}
```

---

## Auto-Sync Behavior

When you update `voters` table fields, triggers automatically update identity tables:

```
PUT /voters/me/profile
    ↓
UPDATE voters SET faculty_code = 'FTI', cohort_year = 2022
    ↓
TRIGGER: sync_student_from_voter() fires
    ↓
UPDATE students SET faculty_code = 'FTI', cohort_year = 2022
    ↓
Response: { synced_to_identity: true }
```

---

## Backend Handler Example

```go
type UpdateProfileRequest struct {
    // Contact info
    Email    *string `json:"email,omitempty"`
    Phone    *string `json:"phone,omitempty"`
    PhotoURL *string `json:"photo_url,omitempty"`
    
    // Identity fields (auto-sync via trigger)
    FacultyCode      *string `json:"faculty_code,omitempty"`
    FacultyName      *string `json:"faculty_name,omitempty"`
    StudyProgramCode *string `json:"study_program_code,omitempty"`
    CohortYear       *int    `json:"cohort_year,omitempty"`
    ClassLabel       *string `json:"class_label,omitempty"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
    voterID := GetVoterIDFromToken(c)
    var req UpdateProfileRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // Update voters table (trigger will auto-sync to identity table)
    err := h.repo.UpdateVoter(voterID, req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "success": true,
        "data": gin.H{
            "message": "Profil berhasil diperbarui",
            "synced_to_identity": true,
        },
    })
}
```

---

## SQL Query Example

```sql
-- Update student via voters (trigger will sync)
UPDATE voters
SET 
    faculty_code = 'FEB',
    study_program_code = 'Akuntansi',
    cohort_year = 2022,
    class_label = 'AKT-B'
WHERE id = 1;

-- Check sync result
SELECT 
    v.id,
    v.faculty_code as v_faculty,
    v.cohort_year as v_cohort,
    s.faculty_code as s_faculty,
    s.cohort_year as s_cohort,
    s.program_code
FROM voters v
JOIN students s ON v.student_id = s.id
WHERE v.id = 1;
```

---

## Important Notes

1. **All fields are OPTIONAL** - Send only what needs to be updated
2. **Auto-sync is AUTOMATIC** - No manual update to identity tables needed
3. **Backward compatible** - Old API calls (without identity fields) still work
4. **Type-specific** - Only send fields applicable to voter type
5. **Triggers handle sync** - Backend just updates voters table

---

**Last Updated:** 2025-11-26
