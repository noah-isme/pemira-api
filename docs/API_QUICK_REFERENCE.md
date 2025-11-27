# API Quick Reference - Frontend Integration

**Version:** 2.0  
**Date:** 2025-11-27  
**Base URL:** `/api/v1`

---

## Registration Endpoints

### 1. Register Student
```
POST /auth/register/student
```
```json
{
  "nim": "2024001002",
  "name": "Budi Santoso",
  "email": "",
  "faculty_name": "Fakultas Ekonomi",
  "study_program_name": "Akuntansi",
  "semester": "5",
  "password": "password123",
  "voting_mode": "ONLINE"
}
```

### 2. Register Lecturer
```
POST /auth/register/lecturer-staff
```
```json
{
  "type": "LECTURER",
  "nidn": "0020129001",
  "name": "Prof. Dr. Hartono, M.Sc",
  "email": "",
  "faculty_name": "Fakultas Kesehatan",
  "department_name": "Keperawatan",
  "position": "Guru Besar",
  "password": "password123",
  "voting_mode": "ONLINE"
}
```

### 3. Register Staff
```
POST /auth/register/lecturer-staff
```
```json
{
  "type": "STAFF",
  "nip": "199203151234568",
  "name": "Indah Permata Sari",
  "email": "",
  "unit_name": "Biro Administrasi Keuangan",
  "position": "Staf",
  "password": "password123",
  "voting_mode": "TPS"
}
```

---

## Master Data Endpoints

### Get Faculties & Programs (Legacy)
```
GET /meta/faculties-programs
```
Response: List of faculties with programs array

### Get Faculties (with IDs)
```
GET /master/faculties
```

### Get Study Programs
```
GET /master/study-programs?faculty_id={id}
```

### Get Lecturer Units
```
GET /master/lecturer-units
```

### Get Lecturer Positions
```
GET /master/lecturer-positions?category={FUNGSIONAL|STRUKTURAL}
```

### Get Staff Units
```
GET /master/staff-units
```

### Get Staff Positions
```
GET /master/staff-positions
```

---

## Authentication

### Login
```
POST /auth/login
```
```json
{
  "username": "2024001002",
  "password": "password123"
}
```

Response includes `access_token` and user profile.

---

## Frontend Implementation Guide

### Student Registration Form

```typescript
// 1. Fetch faculties
const faculties = await fetch('/api/v1/master/faculties')
  .then(r => r.json())

// 2. On faculty selected, fetch programs
const programs = await fetch(`/api/v1/master/study-programs?faculty_id=${facultyId}`)
  .then(r => r.json())

// 3. Submit registration
const response = await fetch('/api/v1/auth/register/student', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    nim: formData.nim,
    name: formData.name,
    email: formData.email || '', // Empty = auto-generate
    faculty_name: selectedFaculty.name,
    study_program_name: selectedProgram.name,
    semester: formData.semester,
    password: formData.password,
    voting_mode: formData.votingMode // "ONLINE" or "TPS"
  })
})
```

### Lecturer Registration Form

```typescript
// 1. Fetch units
const units = await fetch('/api/v1/master/lecturer-units')
  .then(r => r.json())

// 2. Fetch positions (optional: filter by category)
const positions = await fetch('/api/v1/master/lecturer-positions?category=FUNGSIONAL')
  .then(r => r.json())

// 3. Submit registration
const response = await fetch('/api/v1/auth/register/lecturer-staff', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    type: 'LECTURER',
    nidn: formData.nidn,
    name: formData.name,
    email: formData.email || '',
    faculty_name: selectedUnit.name,
    department_name: formData.department,
    position: selectedPosition.name,
    password: formData.password,
    voting_mode: formData.votingMode
  })
})
```

### Staff Registration Form

```typescript
// 1. Fetch units
const units = await fetch('/api/v1/master/staff-units')
  .then(r => r.json())

// 2. Fetch positions
const positions = await fetch('/api/v1/master/staff-positions')
  .then(r => r.json())

// 3. Submit registration
const response = await fetch('/api/v1/auth/register/lecturer-staff', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    type: 'STAFF',
    nip: formData.nip,
    name: formData.name,
    email: formData.email || '',
    unit_name: selectedUnit.name,
    position: selectedPosition.name,
    password: formData.password,
    voting_mode: formData.votingMode
  })
})
```

---

## Important Notes

### Email Auto-generation
- Leave `email` field empty or send empty string
- System will generate: `{username}@pemira.ac.id`
- Example: `2024001002@pemira.ac.id`

### Master Data Matching
- Use exact name from master data endpoints
- Foreign key IDs populated automatically if name matches
- If no match, stored as text only (backward compatible)

### Validation
- **Password:** Minimum 6 characters
- **Semester:** 1-14 (student only)
- **Username:** Must be unique (NIM/NIDN/NIP)
- **Voting Mode:** "ONLINE" or "TPS", default "ONLINE"

### Error Handling
```typescript
try {
  const response = await fetch('/api/v1/auth/register/student', options)
  const data = await response.json()
  
  if (!response.ok) {
    if (data.code === 'USERNAME_EXISTS') {
      // Handle duplicate username
    } else if (data.code === 'VALIDATION_ERROR') {
      // Handle validation errors
    }
  }
} catch (error) {
  // Handle network errors
}
```

---

## Testing

### Test Credentials
```
Student:
- Username: 2024001002
- Password: password123

Lecturer:
- Username: 0020129001
- Password: password123

Staff:
- Username: 199203151234568
- Password: password123
```

---

For detailed API documentation, see:
- `API_CONTRACT_VOTER_REGISTRATION.md` - Full registration endpoints
- `API_CONTRACT_VOTER_PROFILE.md` - Profile management
- `API_CONTRACT_DPT.md` - DPT management
