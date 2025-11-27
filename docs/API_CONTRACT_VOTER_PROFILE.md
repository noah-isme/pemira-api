# API Contract - Voter Profile Management

**Date:** 2025-11-26  
**Version:** 3.1  
**Base URL:** `/api/v1`

**Changelog v3.1:**
- ✨ Added editable identity fields (faculty, program, cohort, class/position)
- 🔄 Auto-sync to identity tables via database triggers
- 📝 Updated field specifications and examples

---

## Table of Contents

1. [Overview](#overview)
2. [Profile Endpoints](#profile-endpoints)
3. [Data Models](#data-models)
4. [Field Specifications](#field-specifications)
5. [Error Codes](#error-codes)

---

## Overview

Voter Profile API allows authenticated voters to:
- View complete profile information
- Update personal information (email, phone, photo)
- Change voting method preference
- Change password
- View participation statistics
- Delete profile photo

### Authentication

All endpoints require JWT authentication:
```
Authorization: Bearer {voter_token}
```

### Base Path

All profile endpoints are under:
```
/api/v1/voters/me/
```

---

## Profile Endpoints

### 1. Get Complete Profile

**Endpoint:** `GET /voters/me/complete-profile`

**Description:** Get comprehensive voter profile including personal info, voting preferences, participation stats, and account info.

**Authentication:** Required (Voter role)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "personal_info": {
      "voter_id": 1,
      "name": "Ahmad Zulfikar",
      "username": "2021001",
      "email": "ahmad@example.com",
      "phone": "081234567890",
      "faculty_name": "Fakultas Teknologi Informasi",
      "study_program_name": "Teknik Informatika",
      "cohort_year": 2021,
      "semester": "7",
      "photo_url": "https://storage.example.com/photos/ahmad.jpg",
      "voter_type": "STUDENT"
    },
    "voting_info": {
      "preferred_method": "ONLINE",
      "has_voted": false,
      "voted_at": null,
      "tps_name": null,
      "tps_location": null
    },
    "participation": {
      "total_elections": 5,
      "participated_elections": 3,
      "participation_rate": 60.0,
      "last_participation": "2024-10-15T14:30:00Z"
    },
    "account_info": {
      "created_at": "2024-01-01T00:00:00Z",
      "last_login": "2024-11-26T08:00:00Z",
      "login_count": 25,
      "account_status": "active"
    }
  }
}
```

**Response 401:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token tidak valid atau tidak ditemukan."
  }
}
```

**Response 403:**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Hanya voter yang dapat mengakses profil."
  }
}
```

---

### 2. Update Profile

**Endpoint:** `PUT /voters/me/profile`

**Description:** Update voter's personal information (email, phone, photo).

**Authentication:** Required (Voter role)

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "phone": "081234567891",
  "photo_url": "https://storage.example.com/photos/new-photo.jpg",
  "faculty_code": "FTI",
  "study_program_code": "IF",
  "cohort_year": 2021,
  "class_label": "IF-A"
}
```

**Field Rules:**
- All fields are **optional** (partial update supported)
- `email`: Valid email format (e.g., user@domain.com)
- `phone`: Format `08xxx` or `+62xxx` (10-15 digits)
- `photo_url`: Valid URL (https://...)
- `faculty_code`: Faculty/unit code (editable based on voter type)
- `study_program_code`: Program/department code (for STUDENT/LECTURER)
- `cohort_year`: Enrollment year (for STUDENT only)
- `class_label`: Class/position label (editable based on voter type)

**Editable Fields:**
- ✅ `email` - Email address
- ✅ `phone` - Phone number
- ✅ `photo_url` - Profile photo URL
- ✅ `faculty_code` - Faculty/unit code (auto-syncs to identity table)
- ✅ `study_program_code` - Program/department code (STUDENT/LECTURER only)
- ✅ `cohort_year` - Enrollment year (STUDENT only)
- ✅ `class_label` - Class/position label (auto-syncs to identity table)

**Non-Editable Fields (Read-Only):**
- ❌ `nim` - NIM/NIDN/NIP (system assigned)
- ❌ `name` - Full name (from registration)
- ❌ `faculty_name` - Faculty name (lookup from code)
- ❌ `study_program_name` - Program name (lookup from code)
- ❌ `voter_type` - Voter type (STUDENT/LECTURER/STAFF)
- ❌ `academic_status` - Academic status (from system)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Profil berhasil diperbarui",
    "updated_fields": ["email", "phone", "faculty_code", "cohort_year"],
    "synced_to_identity": true
  }
}
```

**Response 400 - Invalid Email:**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_EMAIL",
    "message": "Format email tidak valid."
  }
}
```

**Response 400 - Invalid Phone:**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PHONE",
    "message": "Format nomor telepon tidak valid. Gunakan format 08xxx atau +62xxx."
  }
}
```

**Example - Update Email Only:**
```json
{
  "email": "newemail@example.com"
}
```

**Example - Update Phone Only:**
```json
{
  "phone": "+6281234567890"
}
```

**Example - Update Photo Only:**
```json
{
  "photo_url": "https://storage.example.com/photos/profile-123.jpg"
}
```

**Example - Update Student Identity:**
```json
{
  "faculty_code": "FTI",
  "study_program_code": "SI",
  "cohort_year": 2022,
  "class_label": "SI-B"
}
```

**Example - Update Lecturer Identity:**
```json
{
  "faculty_code": "FTI",
  "study_program_code": "Informatika",
  "class_label": "Lektor Kepala"
}
```

**Example - Update Staff Identity:**
```json
{
  "faculty_code": "BAU",
  "class_label": "Koordinator"
}
```

---

### 3. Update Voting Method

**Endpoint:** `PUT /voters/me/voting-method`

**Description:** Change preferred voting method for a specific election.

**Authentication:** Required (Voter role)

**Request Body:**
```json
{
  "election_id": 1,
  "preferred_method": "ONLINE"
}
```

**Field Rules:**
- `election_id`: Required, integer (must be valid election ID)
- `preferred_method`: Required, enum: `"ONLINE"` | `"TPS"`

**Response 200:**
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Metode voting berhasil diubah ke ONLINE",
    "new_method": "ONLINE",
    "warning": "Jika sudah check-in TPS, perubahan tidak berlaku untuk election ini"
  }
}
```

**Response 400 - Already Voted:**
```json
{
  "success": false,
  "error": {
    "code": "ALREADY_VOTED",
    "message": "Tidak dapat mengubah metode voting karena sudah voting."
  }
}
```

**Response 400 - Already Checked In:**
```json
{
  "success": false,
  "error": {
    "code": "ALREADY_CHECKED_IN",
    "message": "Tidak dapat mengubah ke ONLINE karena sudah check-in di TPS."
  }
}
```

**Response 400 - Invalid Method:**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_METHOD",
    "message": "Metode voting tidak valid. Gunakan ONLINE atau TPS."
  }
}
```

**Business Rules:**
1. Can change before voting
2. Cannot change after voting
3. Cannot change from TPS to ONLINE after TPS check-in
4. Can change from ONLINE to TPS anytime before voting

---

### 4. Change Password

**Endpoint:** `POST /voters/me/change-password`

**Description:** Change voter's account password.

**Authentication:** Required (Voter role)

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Field Rules:**
- `current_password`: Required, string
- `new_password`: Required, min 8 characters
- `confirm_password`: Required, must match `new_password`

**Password Requirements:**
- Minimum 8 characters
- Cannot be the same as current password
- Must match confirmation

**Response 200:**
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Password berhasil diubah"
  }
}
```

**Response 400 - Password Mismatch:**
```json
{
  "success": false,
  "error": {
    "code": "PASSWORD_MISMATCH",
    "message": "Konfirmasi password tidak cocok."
  }
}
```

**Response 400 - Password Too Short:**
```json
{
  "success": false,
  "error": {
    "code": "PASSWORD_TOO_SHORT",
    "message": "Password minimal 8 karakter."
  }
}
```

**Response 400 - Same Password:**
```json
{
  "success": false,
  "error": {
    "code": "PASSWORD_SAME",
    "message": "Password baru tidak boleh sama dengan password lama."
  }
}
```

**Response 401 - Wrong Current Password:**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PASSWORD",
    "message": "Password saat ini salah."
  }
}
```

---

### 5. Get Participation Stats

**Endpoint:** `GET /voters/me/participation-stats`

**Description:** Get voter's election participation history and statistics.

**Authentication:** Required (Voter role)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "summary": {
      "total_elections": 5,
      "participated": 3,
      "not_participated": 2,
      "participation_rate": 60.0
    },
    "elections": [
      {
        "election_id": 1,
        "election_name": "Pemilihan Ketua BEM 2024",
        "year": 2024,
        "voted": true,
        "voted_at": "2024-10-15T14:30:00Z",
        "method": "ONLINE"
      },
      {
        "election_id": 2,
        "election_name": "Pemilihan Ketua Ormawa 2024",
        "year": 2024,
        "voted": false,
        "voted_at": null,
        "method": "NONE"
      },
      {
        "election_id": 3,
        "election_name": "Pemilihan Rektor 2024",
        "year": 2024,
        "voted": true,
        "voted_at": "2024-09-20T10:15:00Z",
        "method": "TPS"
      }
    ]
  }
}
```

**Statistics Included:**
- Total elections eligible to vote
- Number of elections participated
- Number of elections skipped
- Overall participation rate (percentage)
- Detailed history per election

---

### 6. Delete Profile Photo

**Endpoint:** `DELETE /voters/me/photo`

**Description:** Remove profile photo (set to null).

**Authentication:** Required (Voter role)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Foto profil berhasil dihapus"
  }
}
```

**Response 404:**
```json
{
  "success": false,
  "error": {
    "code": "VOTER_NOT_FOUND",
    "message": "Voter tidak ditemukan."
  }
}
```

---

## Data Models

### CompleteProfileResponse

```typescript
interface CompleteProfileResponse {
  personal_info: PersonalInfo;
  voting_info: VotingInfo;
  participation: ParticipationInfo;
  account_info: AccountInfo;
}
```

### PersonalInfo

```typescript
interface PersonalInfo {
  voter_id: number;
  name: string;                     // Full name (read-only)
  username: string;                 // NIM/NIDN/NIP (read-only, from identity table)
  email: string | null;             // Editable
  phone: string | null;             // Editable
  faculty_name: string | null;      // Read-only (from identity table)
  study_program_name: string | null; // Read-only (from identity table, STUDENT only)
  cohort_year: number | null;       // Read-only (from identity table, STUDENT only)
  semester: string;                 // Calculated (read-only, STUDENT only)
  photo_url: string | null;         // Editable
  voter_type: "STUDENT" | "LECTURER" | "STAFF"; // Read-only
  
  // Additional fields for LECTURER
  department: string | null;        // Read-only (LECTURER only)
  position: string | null;          // Read-only (LECTURER only)
  
  // Additional fields for STAFF
  unit: string | null;              // Read-only (STAFF only)
  job_title: string | null;         // Read-only (STAFF only)
}
```

### VotingInfo

```typescript
interface VotingInfo {
  preferred_method: "ONLINE" | "TPS" | null;
  has_voted: boolean;
  voted_at: string | null;           // ISO 8601
  tps_name: string | null;
  tps_location: string | null;
}
```

### ParticipationInfo

```typescript
interface ParticipationInfo {
  total_elections: number;
  participated_elections: number;
  participation_rate: number;        // 0-100
  last_participation: string | null; // ISO 8601
}
```

### AccountInfo

```typescript
interface AccountInfo {
  created_at: string;                // ISO 8601
  last_login: string | null;         // ISO 8601
  login_count: number;
  account_status: "active" | "inactive";
}
```

### UpdateProfileRequest

```typescript
interface UpdateProfileRequest {
  email?: string;                    // Optional
  phone?: string;                    // Optional
  photo_url?: string;                // Optional
  faculty_code?: string;             // Optional, syncs to identity
  study_program_code?: string;       // Optional, STUDENT/LECTURER only
  cohort_year?: number;              // Optional, STUDENT only
  class_label?: string;              // Optional, syncs to identity
}
```

### UpdateVotingMethodRequest

```typescript
interface UpdateVotingMethodRequest {
  election_id: number;               // Required
  preferred_method: "ONLINE" | "TPS"; // Required
}
```

### ChangePasswordRequest

```typescript
interface ChangePasswordRequest {
  current_password: string;          // Required
  new_password: string;              // Required, min 8 chars
  confirm_password: string;          // Required, must match
}
```

### ParticipationStatsResponse

```typescript
interface ParticipationStatsResponse {
  summary: {
    total_elections: number;
    participated: number;
    not_participated: number;
    participation_rate: number;
  };
  elections: Array<{
    election_id: number;
    election_name: string;
    year: number;
    voted: boolean;
    voted_at: string | null;
    method: "ONLINE" | "TPS" | "NONE";
  }>;
}
```

---

## Field Specifications

### Database Schema (voters table)

```sql
-- Main voters table (clean, only voter-specific data)
CREATE TABLE voters (
  id                 BIGSERIAL PRIMARY KEY,
  name               TEXT NOT NULL,
  email              TEXT,                           -- Editable
  phone              VARCHAR(20),                    -- Editable
  photo_url          TEXT,                           -- Editable
  voter_type         TEXT NOT NULL CHECK (voter_type IN ('STUDENT','LECTURER','STAFF')),
  
  -- Identity linkage (only one will be filled based on voter_type)
  student_id         BIGINT REFERENCES students(id) ON DELETE SET NULL,
  lecturer_id        BIGINT REFERENCES lecturers(id) ON DELETE SET NULL,
  staff_id           BIGINT REFERENCES staff_members(id) ON DELETE SET NULL,
  
  created_at         TIMESTAMPTZ DEFAULT NOW(),
  updated_at         TIMESTAMPTZ DEFAULT NOW()
);

-- Identity tables (source of truth for academic data)
CREATE TABLE students (
  id                 BIGSERIAL PRIMARY KEY,
  nim                TEXT UNIQUE NOT NULL,
  name               TEXT NOT NULL,
  faculty_code       TEXT,
  program_code       TEXT,
  cohort_year        INT,
  class_label        TEXT,
  created_at         TIMESTAMPTZ DEFAULT NOW(),
  updated_at         TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE lecturers (
  id                 BIGSERIAL PRIMARY KEY,
  nidn               TEXT UNIQUE,
  name               TEXT NOT NULL,
  faculty_code       TEXT,
  department         TEXT,
  position           TEXT,
  created_at         TIMESTAMPTZ DEFAULT NOW(),
  updated_at         TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE staff_members (
  id                 BIGSERIAL PRIMARY KEY,
  nip                TEXT UNIQUE,
  name               TEXT NOT NULL,
  unit               TEXT,
  job_title          TEXT,
  created_at         TIMESTAMPTZ DEFAULT NOW(),
  updated_at         TIMESTAMPTZ DEFAULT NOW()
);
```

### Editable vs Read-Only Fields

| Field | Editable | Via Endpoint | Notes |
|-------|----------|--------------|-------|
| `email` | ✅ Yes | PUT /voters/me/profile | Valid email format |
| `phone` | ✅ Yes | PUT /voters/me/profile | Format: 08xxx or +62xxx |
| `photo_url` | ✅ Yes | PUT /voters/me/profile | Valid URL |
| `faculty_code` | ✅ Yes | PUT /voters/me/profile | Syncs to identity table |
| `study_program_code` | ✅ Yes | PUT /voters/me/profile | For STUDENT/LECTURER, syncs to identity |
| `cohort_year` | ✅ Yes | PUT /voters/me/profile | For STUDENT only, syncs to identity |
| `class_label` | ✅ Yes | PUT /voters/me/profile | Maps to position/job, syncs to identity |
| `voting_method` | ✅ Yes | PUT /voters/me/voting-method | Per election |
| `password` | ✅ Yes | POST /voters/me/change-password | Min 8 chars |
| `nim/nidn/nip` | ❌ No | - | From identity table |
| `name` | ❌ No | - | From identity table |
| `faculty_name` | ❌ No | - | Lookup from faculty_code |
| `study_program_name` | ❌ No | - | Lookup from study_program_code |
| `voter_type` | ❌ No | - | System assigned |

---

## Validation Rules

### Email Validation

**Format:** Standard email format
```
user@domain.com
user.name@subdomain.domain.com
```

**Regex:** 
```regex
^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
```

**Examples:**
- ✅ Valid: `ahmad@example.com`, `test.user@mail.com`
- ❌ Invalid: `invalid`, `@example.com`, `user@`

---

### Phone Validation

**Format:** Indonesian phone numbers
```
08xxxxxxxxxx       (11-13 digits)
+628xxxxxxxxxx     (12-15 digits)
```

**Regex:**
```regex
^(08\d{8,11}|\+628\d{8,12})$
```

**Examples:**
- ✅ Valid: `081234567890`, `+6281234567890`
- ❌ Invalid: `08123`, `1234567890`, `08123456789012345`

---

### Photo URL Validation

**Format:** Valid HTTPS URL
```
https://domain.com/path/to/image.jpg
```

**Requirements:**
- Must start with `https://`
- Valid URL format
- Recommended: Image file extension (jpg, jpeg, png, webp)

**Examples:**
- ✅ Valid: `https://storage.example.com/photos/user-123.jpg`
- ❌ Invalid: `http://insecure.com/photo.jpg`, `not-a-url`

---

### Password Validation

**Requirements:**
- Minimum 8 characters
- No maximum length limit
- Can contain any characters
- Must not be same as current password
- Must match confirmation

**Examples:**
- ✅ Valid: `password123`, `MyP@ssw0rd!`, `longpassword`
- ❌ Invalid: `pass` (too short), `1234567` (7 chars)

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid token |
| `FORBIDDEN` | 403 | Not a voter account |
| `VOTER_NOT_FOUND` | 404 | Voter not found |
| `INVALID_EMAIL` | 400 | Invalid email format |
| `INVALID_PHONE` | 400 | Invalid phone format |
| `INVALID_METHOD` | 400 | Invalid voting method |
| `PASSWORD_MISMATCH` | 400 | Password confirmation doesn't match |
| `PASSWORD_TOO_SHORT` | 400 | Password less than 8 characters |
| `PASSWORD_SAME` | 400 | New password same as current |
| `INVALID_PASSWORD` | 401 | Current password incorrect |
| `ALREADY_VOTED` | 400 | Cannot change after voting |
| `ALREADY_CHECKED_IN` | 400 | Cannot change after TPS check-in |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Usage Examples

### Get Complete Profile

```bash
curl 'http://localhost:8080/api/v1/voters/me/complete-profile' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN"
```

### Update Email

```bash
curl -X PUT 'http://localhost:8080/api/v1/voters/me/profile' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }'
```

### Update Phone

```bash
curl -X PUT 'http://localhost:8080/api/v1/voters/me/profile' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+6281234567890"
  }'
```

### Update Multiple Fields

```bash
curl -X PUT 'http://localhost:8080/api/v1/voters/me/profile' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com",
    "phone": "081234567891",
    "photo_url": "https://storage.example.com/photos/new.jpg"
  }'
```

### Change Voting Method

```bash
curl -X PUT 'http://localhost:8080/api/v1/voters/me/voting-method' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "election_id": 1,
    "preferred_method": "TPS"
  }'
```

### Change Password

```bash
curl -X POST 'http://localhost:8080/api/v1/voters/me/change-password' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "oldpassword123",
    "new_password": "newpassword456",
    "confirm_password": "newpassword456"
  }'
```

### Get Participation Stats

```bash
curl 'http://localhost:8080/api/v1/voters/me/participation-stats' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN"
```

### Delete Photo

```bash
curl -X DELETE 'http://localhost:8080/api/v1/voters/me/photo' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN"
```

---

## Frontend Integration Guide

### Profile Form Example

```typescript
// Component: ProfileEditForm.tsx

interface ProfileFormData {
  email: string;
  phone: string;
  photo_url: string;
}

const ProfileEditForm: React.FC = () => {
  const [formData, setFormData] = useState<ProfileFormData>({
    email: '',
    phone: '',
    photo_url: ''
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const response = await fetch('/api/v1/voters/me/profile', {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });
      
      const result = await response.json();
      
      if (result.success) {
        alert('Profile updated successfully!');
      } else {
        alert(`Error: ${result.error.message}`);
      }
    } catch (error) {
      console.error('Update failed:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={formData.email}
        onChange={(e) => setFormData({...formData, email: e.target.value})}
        placeholder="Email"
      />
      <input
        type="tel"
        value={formData.phone}
        onChange={(e) => setFormData({...formData, phone: e.target.value})}
        placeholder="Phone"
        pattern="^(08\d{8,11}|\+628\d{8,12})$"
      />
      <input
        type="url"
        value={formData.photo_url}
        onChange={(e) => setFormData({...formData, photo_url: e.target.value})}
        placeholder="Photo URL"
      />
      <button type="submit">Update Profile</button>
    </form>
  );
};
```

### Password Change Form Example

```typescript
// Component: ChangePasswordForm.tsx

const ChangePasswordForm: React.FC = () => {
  const [passwords, setPasswords] = useState({
    current_password: '',
    new_password: '',
    confirm_password: ''
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (passwords.new_password !== passwords.confirm_password) {
      alert('New password and confirmation do not match!');
      return;
    }
    
    if (passwords.new_password.length < 8) {
      alert('Password must be at least 8 characters!');
      return;
    }
    
    try {
      const response = await fetch('/api/v1/voters/me/change-password', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(passwords)
      });
      
      const result = await response.json();
      
      if (result.success) {
        alert('Password changed successfully!');
        setPasswords({ current_password: '', new_password: '', confirm_password: '' });
      } else {
        alert(`Error: ${result.error.message}`);
      }
    } catch (error) {
      console.error('Password change failed:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="password"
        value={passwords.current_password}
        onChange={(e) => setPasswords({...passwords, current_password: e.target.value})}
        placeholder="Current Password"
        required
      />
      <input
        type="password"
        value={passwords.new_password}
        onChange={(e) => setPasswords({...passwords, new_password: e.target.value})}
        placeholder="New Password"
        minLength={8}
        required
      />
      <input
        type="password"
        value={passwords.confirm_password}
        onChange={(e) => setPasswords({...passwords, confirm_password: e.target.value})}
        placeholder="Confirm New Password"
        required
      />
      <button type="submit">Change Password</button>
    </form>
  );
};
```

---

## Important Notes

### 1. Partial Updates Supported

You can update individual fields without sending all fields:

```json
// Update only email
{ "email": "new@example.com" }

// Update only phone
{ "phone": "081234567890" }

// Update email and phone
{ "email": "new@example.com", "phone": "081234567890" }
```

### 2. Editable Identity Fields

The following identity fields **can be updated** by voters:
- `faculty_code` - Faculty/unit code
- `study_program_code` - Program/department code (STUDENT/LECTURER)
- `cohort_year` - Enrollment year (STUDENT only)
- `class_label` - Class/position/job title

**Auto-Sync:** Changes are automatically synchronized to identity tables via database triggers.

**Mapping:**
- STUDENT: Updates `students` table (faculty_code, program_code, cohort_year, class_label)
- LECTURER: Updates `lecturers` table (faculty_code, department_code, position)
- STAFF: Updates `staff_members` table (unit_code, unit_name, position)

### 3. Read-Only Fields

The following fields **cannot** be updated via profile API:
- NIM/NIDN/NIP (unique identifier)
- Full name (requires admin verification)
- Voter type (immutable)

### 4. Voting Method Per Election

Voting method preference is **per election**, not global:
- Each election can have different voting method
- Change affects only specified election
- Cannot change after voting in that election

### 5. Password Security

- Current password required for verification
- New password must be different from current
- Minimum 8 characters enforced
- No password strength requirements (only length)

### 6. Photo Management

- Photo URL stored as string (not file upload)
- Use separate file upload service for actual photos
- Can be deleted (set to null)
- No file size/format validation in API

---

## Testing Checklist

### Profile Update
- [ ] Update email with valid format
- [ ] Update email with invalid format (should fail)
- [ ] Update phone with valid format
- [ ] Update phone with invalid format (should fail)
- [ ] Update photo URL
- [ ] Update multiple fields at once
- [ ] Update with empty values (should clear)

### Voting Method
- [ ] Change from ONLINE to TPS
- [ ] Change from TPS to ONLINE
- [ ] Change before voting (should succeed)
- [ ] Change after voting (should fail)
- [ ] Change after TPS check-in (should fail)

### Password Change
- [ ] Change with correct current password
- [ ] Change with wrong current password (should fail)
- [ ] Use password shorter than 8 chars (should fail)
- [ ] Use same password as current (should fail)
- [ ] Mismatch confirmation (should fail)

### Data Retrieval
- [ ] Get complete profile
- [ ] Get participation stats
- [ ] Check all fields returned correctly
- [ ] Verify semester calculation

### Photo Management
- [ ] Delete photo
- [ ] Photo returns null after deletion

---

**Last Updated:** 2025-11-26  
**API Version:** 2.0  
**Maintained By:** Backend Team
