# Registration Endpoints Test Report

## Test Date
2025-11-26

## Summary
✅ All registration endpoints are working correctly for Students, Lecturers, and Staff.

---

## 1. Student Registration Endpoint

### Endpoint
```
POST /api/v1/auth/register/student
```

### Request Body
```json
{
  "nim": "202301001",
  "name": "Budi Santoso",
  "email": "budi@student.ac.id",
  "faculty_name": "Fakultas Teknik",
  "study_program_name": "Teknik Informatika",
  "semester": "5",
  "password": "password123"
}
```

### ✅ Success Response (201)
```json
{
  "message": "Registrasi mahasiswa berhasil.",
  "user": {
    "id": 51,
    "username": "202301001",
    "role": "STUDENT",
    "voter_id": 46,
    "profile": {
      "name": "Budi Santoso",
      "faculty_name": "Fakultas Teknik",
      "study_program_name": "Teknik Informatika",
      "semester": "5"
    }
  },
  "voting_mode": "ONLINE"
}
```

### Database Impact
- ✅ Record created in `voters` table
- ✅ Record created in `user_accounts` table with role "STUDENT"
- ✅ Voter linked to user account via `voter_id`

---

## 2. Lecturer Registration Endpoint

### Endpoint
```
POST /api/v1/auth/register/lecturer-staff
```

### Request Body
```json
{
  "type": "LECTURER",
  "nidn": "0123456789",
  "name": "Dr. Ahmad Fauzi, M.Kom",
  "email": "ahmad.fauzi@university.ac.id",
  "faculty_name": "Fakultas Teknik",
  "department_name": "Teknik Informatika",
  "position": "Lektor",
  "password": "password123"
}
```

### ✅ Success Response (201)
```json
{
  "message": "Registrasi berhasil.",
  "user": {
    "id": 52,
    "username": "0123456789",
    "role": "LECTURER",
    "voter_id": 47,
    "lecturer_id": 8,
    "profile": {
      "name": "Dr. Ahmad Fauzi, M.Kom",
      "faculty_name": "Fakultas Teknik",
      "department_name": "Teknik Informatika",
      "position": "Lektor"
    }
  },
  "voting_mode": "ONLINE"
}
```

### Database Impact
- ✅ Record created in `voters` table
- ✅ Record created in `lecturers` table
- ✅ Record created in `user_accounts` table with role "LECTURER"
- ✅ Lecturer linked to user account via `lecturer_id`
- ✅ Voter linked to user account via `voter_id`

**Note**: `unit_id` and `position_id` are NULL (still using text values)

---

## 3. Staff Registration Endpoint

### Endpoint
```
POST /api/v1/auth/register/lecturer-staff
```

### Request Body
```json
{
  "type": "STAFF",
  "nip": "202311001",
  "name": "Rina Kartika",
  "email": "rina.kartika@university.ac.id",
  "unit_name": "Biro Administrasi Umum",
  "position": "Staf",
  "password": "password123"
}
```

### ✅ Success Response (201)
```json
{
  "message": "Registrasi berhasil.",
  "user": {
    "id": 56,
    "username": "202311001",
    "role": "STAFF",
    "voter_id": 51,
    "staff_id": 12,
    "profile": {
      "name": "Rina Kartika",
      "position": "Staf",
      "unit_name": "Biro Administrasi Umum"
    }
  },
  "voting_mode": "ONLINE"
}
```

### Database Impact
- ✅ Record created in `voters` table
- ✅ Record created in `staff_members` table
- ✅ Record created in `user_accounts` table with role "STAFF"
- ✅ Staff linked to user account via `staff_id`
- ✅ Voter linked to user account via `voter_id`

**Note**: `unit_id` and `position_id` are NULL (still using text values)

---

## 4. Login Tests

All registered users can successfully login:

### Student Login ✅
```bash
POST /api/v1/auth/login
{
  "username": "202301001",
  "password": "password123"
}
```
**Result**: Returns access token and user profile

### Lecturer Login ✅
```bash
POST /api/v1/auth/login
{
  "username": "0123456789",
  "password": "password123"
}
```
**Result**: Returns access token and user profile

### Staff Login ✅
```bash
POST /api/v1/auth/login
{
  "username": "202311001",
  "password": "password123"
}
```
**Result**: Returns access token and user profile

---

## 5. Error Handling Tests

### ✅ Duplicate Registration
**Test**: Register with existing NIM/NIDN/NIP
```json
{
  "code": "USERNAME_EXISTS",
  "message": "Username sudah terdaftar."
}
```

### ✅ Missing Required Fields
**Test**: Register without `semester` field
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Data registrasi tidak lengkap atau tidak valid."
}
```

### ✅ Weak Password
**Test**: Register with password < 6 characters
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Data registrasi tidak lengkap atau tidak valid."
}
```

---

## 6. Database Verification

### Voters Table
```sql
SELECT nim, name, faculty_name, study_program_name FROM voters 
WHERE nim IN ('202301001', '0123456789', '202311001');
```

| NIM | Name | Faculty | Study Program |
|-----|------|---------|---------------|
| 202301001 | Budi Santoso | Fakultas Teknik | Teknik Informatika |
| 0123456789 | Dr. Ahmad Fauzi, M.Kom | Fakultas Teknik | Teknik Informatika |
| 202311001 | Rina Kartika | Biro Administrasi Umum | - |

### Lecturers Table
```sql
SELECT nidn, name, faculty_name, department_name, position, unit_id, position_id 
FROM lecturers WHERE nidn = '0123456789';
```

| NIDN | Name | Faculty | Department | Position | Unit ID | Position ID |
|------|------|---------|------------|----------|---------|-------------|
| 0123456789 | Dr. Ahmad Fauzi, M.Kom | Fakultas Teknik | Teknik Informatika | Lektor | NULL | NULL |

### Staff Members Table
```sql
SELECT nip, name, unit_name, position, unit_id, position_id 
FROM staff_members WHERE nip = '202311001';
```

| NIP | Name | Unit | Position | Unit ID | Position ID |
|-----|------|------|----------|---------|-------------|
| 202311001 | Rina Kartika | Biro Administrasi Umum | Staf | NULL | NULL |

---

## Issues & Recommendations

### ⚠️ Issue: Master Tables Not Used Yet
**Current Behavior**:
- Registration still saves text values for faculty, study_program, unit, position
- New columns `unit_id` and `position_id` in `lecturers` and `staff_members` are NULL
- No validation against master tables

**Recommendation**:
1. Update registration service to lookup IDs from master tables
2. Store foreign key IDs instead of text values
3. Add validation to ensure submitted values exist in master tables
4. Consider making text columns nullable and use only FK columns

### Example Fix Needed:
```go
// Instead of:
faculty_name: "Fakultas Teknik"

// Should lookup and use:
faculty_id: 5  // ID from faculties table
```

---

## Test Commands

```bash
# Start server
cd /home/noah/project/pemira-api
go run cmd/api/main.go

# Register Student
curl -X POST http://localhost:8080/api/v1/auth/register/student \
  -H "Content-Type: application/json" \
  -d '{"nim":"202301001","name":"Budi Santoso","email":"budi@student.ac.id","faculty_name":"Fakultas Teknik","study_program_name":"Teknik Informatika","semester":"5","password":"password123"}'

# Register Lecturer
curl -X POST http://localhost:8080/api/v1/auth/register/lecturer-staff \
  -H "Content-Type: application/json" \
  -d '{"type":"LECTURER","nidn":"0123456789","name":"Dr. Ahmad Fauzi, M.Kom","email":"ahmad.fauzi@university.ac.id","faculty_name":"Fakultas Teknik","department_name":"Teknik Informatika","position":"Lektor","password":"password123"}'

# Register Staff
curl -X POST http://localhost:8080/api/v1/auth/register/lecturer-staff \
  -H "Content-Type: application/json" \
  -d '{"type":"STAFF","nip":"202311001","name":"Rina Kartika","email":"rina.kartika@university.ac.id","unit_name":"Biro Administrasi Umum","position":"Staf","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"202301001","password":"password123"}'
```

---

## Conclusion

✅ **All registration endpoints work correctly**
✅ **Authentication and login functional**
✅ **Error handling working as expected**
✅ **Data properly stored in database**

⚠️ **Next Step**: Integrate master tables with registration logic to use foreign keys instead of text values.
