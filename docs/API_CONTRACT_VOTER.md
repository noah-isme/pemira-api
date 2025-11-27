# API Contract - Voter & Election Voter Endpoints

**Tanggal Update:** 2025-11-26  
**Versi:** 2.0 (Post voting_method_preference removal)  
**Breaking Changes:** ⚠️ Ya - Field name changed

---

## 🔴 Breaking Changes

### Field Name Changed

**Old (v1.x):**
```json
{
  "voting_method_preference": "ONLINE"
}
```

**New (v2.0):**
```json
{
  "voting_method": "ONLINE"
}
```

### Affected Endpoints
- All voter endpoints returning voter objects
- All election voter endpoints
- Profile endpoints

---

## 📋 Table of Contents

1. [Voter Endpoints](#voter-endpoints)
2. [Voter Profile Endpoints](#voter-profile-endpoints)
3. [Election Voter Endpoints](#election-voter-endpoints)
4. [Data Models](#data-models)

---

## Voter Endpoints

### 1. Get Voter List

**Endpoint:** `GET /api/voters`

**Authentication:** Required

**Query Parameters:**
```typescript
{
  page?: number;        // default: 1
  per_page?: number;    // default: 20, max: 100
}
```

**Response 200:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "nim": "2021001",
        "name": "Ahmad Zulfikar",
        "email": "ahmad@example.com",
        "phone": "081234567890",
        "faculty_code": "FTI",
        "faculty_name": "Fakultas Teknologi Informasi",
        "study_program_code": "IF",
        "study_program_name": "Informatika",
        "cohort_year": 2021,
        "class_label": "IF-A",
        "photo_url": "https://storage.example.com/photos/ahmad.jpg",
        "bio": "Mahasiswa aktif",
        "voting_method": "ONLINE",
        "academic_status": "ACTIVE",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-11-26T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

---

### 2. Get Voter by NIM

**Endpoint:** `GET /api/voters/nim/{nim}`

**Authentication:** Required

**Response 200:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "nim": "2021001",
    "name": "Ahmad Zulfikar",
    "email": "ahmad@example.com",
    "phone": "081234567890",
    "faculty_code": "FTI",
    "faculty_name": "Fakultas Teknologi Informasi",
    "study_program_code": "IF",
    "study_program_name": "Informatika",
    "cohort_year": 2021,
    "class_label": "IF-A",
    "photo_url": "https://storage.example.com/photos/ahmad.jpg",
    "bio": "Mahasiswa aktif",
    "voting_method": "ONLINE",
    "academic_status": "ACTIVE",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-11-26T10:00:00Z"
  }
}
```

**Response 404:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Voter not found"
  }
}
```

---

## Voter Profile Endpoints

### 3. Get Complete Profile

**Endpoint:** `GET /api/voters/me/complete-profile`

**Authentication:** Required (Bearer Token)

**Description:** Get complete voter profile including personal info, voting info, participation stats, and account info.

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
      "study_program_name": "Informatika",
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

### 4. Update Profile

**Endpoint:** `PUT /api/voters/me/profile`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "phone": "081234567891",
  "photo_url": "https://storage.example.com/photos/new-photo.jpg"
}
```

**Field Rules:**
- All fields are optional
- `email`: Valid email format
- `phone`: Format 08xxx or +62xxx
- `photo_url`: Valid URL

**Response 200:**
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Profil berhasil diperbarui",
    "updated_fields": ["email", "phone", "photo_url"]
  }
}
```

**Response 400:**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_EMAIL",
    "message": "Format email tidak valid."
  }
}
```

---

### 5. Update Voting Method

**Endpoint:** `PUT /api/voters/me/voting-method`

**Authentication:** Required (Bearer Token)

**Description:** Update preferred voting method for a specific election.

**Request Body:**
```json
{
  "election_id": 1,
  "preferred_method": "ONLINE"
}
```

**Field Rules:**
- `election_id`: Required, integer
- `preferred_method`: Required, enum: "ONLINE" | "TPS"

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

---

### 6. Get Participation Stats

**Endpoint:** `GET /api/voters/me/participation-stats`

**Authentication:** Required (Bearer Token)

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
      }
    ]
  }
}
```

---

### 7. Change Password

**Endpoint:** `POST /api/voters/me/change-password`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Field Rules:**
- `current_password`: Required
- `new_password`: Required, min 8 characters
- `confirm_password`: Required, must match new_password

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

**Response 401 - Invalid Current Password:**
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

### 8. Delete Photo

**Endpoint:** `DELETE /api/voters/me/photo`

**Authentication:** Required (Bearer Token)

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

---

## Election Voter Endpoints

### 9. Admin Lookup Voter

**Endpoint:** `GET /admin/elections/{electionID}/voters/lookup`

**Authentication:** Required (Admin)

**Query Parameters:**
```typescript
{
  nim: string; // Required
}
```

**Response 200:**
```json
{
  "success": true,
  "data": {
    "voter": {
      "id": 1,
      "nim": "2021001",
      "name": "Ahmad Zulfikar",
      "voter_type": "STUDENT",
      "email": "ahmad@example.com",
      "faculty_code": "FTI",
      "study_program_code": "IF",
      "cohort_year": 2021,
      "academic_status": "ACTIVE",
      "has_account": true,
      "lecturer_id": null,
      "staff_id": null,
      "voting_method": "ONLINE"
    },
    "election_voter": {
      "election_voter_id": 10,
      "election_id": 1,
      "voter_id": 1,
      "nim": "2021001",
      "status": "PENDING",
      "voting_method": "ONLINE",
      "tps_id": null,
      "checked_in_at": null,
      "voted_at": null,
      "updated_at": "2024-11-26T10:00:00Z",
      "voter_type": "STUDENT",
      "name": "Ahmad Zulfikar",
      "email": "ahmad@example.com",
      "faculty_code": "FTI",
      "study_program_code": "IF",
      "cohort_year": 2021
    }
  }
}
```

**Response 404:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Pemilih dengan NIM tersebut tidak ditemukan"
  }
}
```

---

### 10. Admin Upsert and Enroll Voter

**Endpoint:** `POST /admin/elections/{electionID}/voters`

**Authentication:** Required (Admin)

**Description:** Create or update voter and enroll to election.

**Request Body:**
```json
{
  "voter_type": "STUDENT",
  "nim": "2021001",
  "name": "Ahmad Zulfikar",
  "email": "ahmad@example.com",
  "phone": "081234567890",
  "faculty_code": "FTI",
  "faculty_name": "Fakultas Teknologi Informasi",
  "study_program_code": "IF",
  "study_program_name": "Informatika",
  "cohort_year": 2021,
  "academic_status": "ACTIVE",
  "lecturer_id": null,
  "staff_id": null,
  "voting_method": "ONLINE",
  "status": "PENDING",
  "tps_id": null
}
```

**Field Rules:**
- `voter_type`: Required, enum: "STUDENT" | "LECTURER" | "STAFF"
- `nim`: Required, string
- `name`: Required, string
- `voting_method`: Required, enum: "ONLINE" | "TPS"
- `status`: Required, enum: "PENDING" | "APPROVED" | "REJECTED"

**Response 200:**
```json
{
  "success": true,
  "data": {
    "voter_id": 1,
    "election_voter_id": 10,
    "status": "PENDING",
    "voting_method": "ONLINE",
    "tps_id": null,
    "created_voter": false,
    "created_election_voter": true,
    "duplicate_in_election": false
  }
}
```

**Response 409 - Duplicate:**
```json
{
  "success": false,
  "error": {
    "code": "DUPLICATE",
    "message": "NIM sudah terdaftar di pemilu ini"
  }
}
```

---

### 11. Admin List Election Voters

**Endpoint:** `GET /admin/elections/{electionID}/voters`

**Authentication:** Required (Admin)

**Query Parameters:**
```typescript
{
  search?: string;              // Search by NIM or name
  voter_type?: string;          // STUDENT | LECTURER | STAFF
  status?: string;              // PENDING | APPROVED | REJECTED
  voting_method?: string;       // ONLINE | TPS
  faculty_code?: string;
  study_program_code?: string;
  cohort_year?: number;
  tps_id?: number;
  page?: number;                // default: 1
  limit?: number;               // default: 50, max: 100
}
```

**Response 200:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "election_voter_id": 10,
        "election_id": 1,
        "voter_id": 1,
        "nim": "2021001",
        "status": "PENDING",
        "voting_method": "ONLINE",
        "tps_id": null,
        "checked_in_at": null,
        "voted_at": null,
        "updated_at": "2024-11-26T10:00:00Z",
        "voter_type": "STUDENT",
        "name": "Ahmad Zulfikar",
        "email": "ahmad@example.com",
        "faculty_code": "FTI",
        "study_program_code": "IF",
        "cohort_year": 2021
      }
    ],
    "page": 1,
    "limit": 50,
    "total_items": 100,
    "total_pages": 2
  }
}
```

---

### 12. Admin Update Election Voter

**Endpoint:** `PATCH /admin/elections/{electionID}/voters/{voterID}`

**Authentication:** Required (Admin)

**Request Body:**
```json
{
  "status": "APPROVED",
  "voting_method": "TPS",
  "tps_id": 5
}
```

**Field Rules:**
- All fields are optional
- `status`: enum: "PENDING" | "APPROVED" | "REJECTED"
- `voting_method`: enum: "ONLINE" | "TPS"
- `tps_id`: integer or null

**Response 200:**
```json
{
  "success": true,
  "data": {
    "election_voter_id": 10,
    "election_id": 1,
    "voter_id": 1,
    "nim": "2021001",
    "status": "APPROVED",
    "voting_method": "TPS",
    "tps_id": 5,
    "checked_in_at": null,
    "voted_at": null,
    "updated_at": "2024-11-26T10:05:00Z"
  }
}
```

---

### 13. Voter Self Register

**Endpoint:** `POST /elections/{electionID}/voters/register`

**Authentication:** Required (Voter Bearer Token)

**Request Body:**
```json
{
  "voting_method": "ONLINE",
  "tps_id": null
}
```

**Field Rules:**
- `voting_method`: Required, enum: "ONLINE" | "TPS"
- `tps_id`: Required if voting_method is "TPS", otherwise null

**Response 200:**
```json
{
  "success": true,
  "data": {
    "election_voter_id": 10,
    "election_id": 1,
    "voter_id": 1,
    "nim": "2021001",
    "status": "PENDING",
    "voting_method": "ONLINE",
    "tps_id": null,
    "checked_in_at": null,
    "voted_at": null,
    "updated_at": "2024-11-26T10:00:00Z"
  }
}
```

---

### 14. Get Voter Election Status

**Endpoint:** `GET /elections/{electionID}/voters/me/status`

**Authentication:** Required (Voter Bearer Token)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "election_voter_id": 10,
    "election_id": 1,
    "voter_id": 1,
    "nim": "2021001",
    "status": "APPROVED",
    "voting_method": "ONLINE",
    "tps_id": null,
    "checked_in_at": null,
    "voted_at": null,
    "updated_at": "2024-11-26T10:00:00Z",
    "voter_type": "STUDENT",
    "name": "Ahmad Zulfikar",
    "email": "ahmad@example.com",
    "faculty_code": "FTI",
    "study_program_code": "IF",
    "cohort_year": 2021
  }
}
```

---

## Data Models

### Voter Model

```typescript
interface Voter {
  id: number;
  nim: string;
  name: string;
  email: string | null;
  phone: string | null;
  faculty_code: string | null;
  faculty_name: string | null;
  study_program_code: string | null;
  study_program_name: string | null;
  cohort_year: number | null;
  class_label: string | null;
  photo_url: string | null;
  bio: string | null;
  voting_method: "ONLINE" | "TPS" | null;  // ⚠️ NEW: Changed from voting_method_preference
  academic_status: "ACTIVE" | "GRADUATED" | "ON_LEAVE" | "DROPPED" | "INACTIVE";
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
}
```

### VoterSummary Model

```typescript
interface VoterSummary {
  id: number;
  nim: string;
  name: string;
  voter_type: "STUDENT" | "LECTURER" | "STAFF";
  email?: string;
  faculty_code?: string;
  study_program_code?: string;
  cohort_year?: number;
  academic_status?: string;
  has_account: boolean;
  lecturer_id?: number;
  staff_id?: number;
  voting_method?: "ONLINE" | "TPS";  // ⚠️ NEW: Single field
}
```

### ElectionVoter Model

```typescript
interface ElectionVoter {
  election_voter_id: number;
  election_id: number;
  voter_id: number;
  nim: string;
  status: "PENDING" | "APPROVED" | "REJECTED";
  voting_method: "ONLINE" | "TPS";
  tps_id: number | null;
  checked_in_at: string | null; // ISO 8601
  voted_at: string | null; // ISO 8601
  updated_at: string; // ISO 8601
  voter_type?: "STUDENT" | "LECTURER" | "STAFF";
  name?: string;
  email?: string;
  faculty_code?: string;
  study_program_code?: string;
  cohort_year?: number;
}
```

---

## Migration Guide for Frontend

### Step 1: Find and Replace

Search for `voting_method_preference` in your codebase and replace with `voting_method`.

**Example:**

```typescript
// Before (v1.x)
interface OldVoter {
  voting_method_preference: "ONLINE" | "TPS";
}

// After (v2.0)
interface NewVoter {
  voting_method: "ONLINE" | "TPS";
}
```

### Step 2: Update API Calls

```typescript
// Before (v1.x)
const updateVotingMethod = async (electionId: number, preference: string) => {
  return api.put(`/voters/me/voting-method`, {
    election_id: electionId,
    preferred_method: preference,  // This field name is still correct
  });
};

// After (v2.0) - No change needed in request body
// But response will have "voting_method" in voter objects
```

### Step 3: Update State Management

```typescript
// Before (v1.x)
const [voter, setVoter] = useState<Voter>({
  id: 0,
  nim: "",
  name: "",
  voting_method_preference: "ONLINE",
  // ...
});

// After (v2.0)
const [voter, setVoter] = useState<Voter>({
  id: 0,
  nim: "",
  name: "",
  voting_method: "ONLINE",  // ⚠️ Changed
  // ...
});
```

### Step 4: Update UI Display

```typescript
// Before (v1.x)
<p>Metode Voting: {voter.voting_method_preference}</p>

// After (v2.0)
<p>Metode Voting: {voter.voting_method}</p>
```

---

## Error Codes Reference

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Token tidak valid atau tidak ada |
| `FORBIDDEN` | 403 | Akses ditolak untuk role ini |
| `NOT_FOUND` | 404 | Resource tidak ditemukan |
| `VALIDATION_ERROR` | 400 | Input tidak valid |
| `INVALID_EMAIL` | 400 | Format email tidak valid |
| `INVALID_PHONE` | 400 | Format telepon tidak valid |
| `INVALID_METHOD` | 400 | Metode voting tidak valid |
| `PASSWORD_MISMATCH` | 400 | Password tidak cocok |
| `PASSWORD_TOO_SHORT` | 400 | Password terlalu pendek |
| `PASSWORD_SAME` | 400 | Password baru sama dengan lama |
| `INVALID_PASSWORD` | 401 | Password saat ini salah |
| `ALREADY_VOTED` | 400 | Sudah melakukan voting |
| `ALREADY_CHECKED_IN` | 400 | Sudah check-in di TPS |
| `DUPLICATE` | 409 | Data sudah ada |
| `INTERNAL_ERROR` | 500 | Kesalahan server |

---

## Testing Checklist

### ✅ Critical Tests

- [ ] GET /api/voters - Returns voters with `voting_method` field
- [ ] GET /api/voters/me/complete-profile - Returns profile with correct structure
- [ ] PUT /api/voters/me/voting-method - Updates voting method successfully
- [ ] GET /admin/elections/{id}/voters/lookup - Returns voter with `voting_method`
- [ ] POST /admin/elections/{id}/voters - Creates voter with `voting_method`
- [ ] GET /admin/elections/{id}/voters - Lists voters with `voting_method`

### ✅ Edge Cases

- [ ] Voter without voting_method (null value handled)
- [ ] Update voting method after already voted (should fail)
- [ ] Update to ONLINE after TPS check-in (should fail)
- [ ] Profile update with invalid email/phone format

---

## Support

**Questions?** Contact backend team or check:
- Full report: `DATABASE_FIX_REPORT.md`
- SQL fixes: `fix_db_issues.sql`

**Version History:**
- v2.0 (2025-11-26): Removed `voting_method_preference`, use `voting_method`
- v1.x: Legacy version with `voting_method_preference`

---

**Last Updated:** 2025-11-26  
**Maintained By:** Backend Team
