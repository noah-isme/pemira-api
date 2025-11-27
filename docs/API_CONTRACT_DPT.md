# API Contract - DPT (Daftar Pemilih Tetap) Management

**Date:** 2025-11-26  
**Version:** 2.0  
**Base URL:** `/api/v1`

---

## Table of Contents

1. [Overview](#overview)
2. [Admin Endpoints](#admin-endpoints)
3. [Voter Self-Service Endpoints](#voter-self-service-endpoints)
4. [Data Models](#data-models)
5. [Error Codes](#error-codes)

---

## Overview

DPT Management API provides endpoints for:
- **Admin:** CRUD operations for voter enrollment in elections
- **Voters:** Self-registration and status checking

### Authentication

All endpoints require JWT authentication via `Authorization: Bearer {token}` header.

### Base Path

All DPT endpoints are under:
- Admin: `/api/v1/admin/elections/{electionID}/voters`
- Voter: `/api/v1/voters/me/elections/{electionID}`

---

## Admin Endpoints

### 1. Get DPT List (Paginated)

**Endpoint:** `GET /admin/elections/{electionID}/voters`

**Description:** Get paginated list of voters enrolled in specific election with filtering and search capabilities.

**Authentication:** Admin role required

**Path Parameters:**
```typescript
{
  electionID: number; // Election ID
}
```

**Query Parameters:**
```typescript
{
  page?: number;                // Page number, default: 1
  limit?: number;               // Items per page, default: 50, max: 100
  search?: string;              // Search by NIM or name
  voter_type?: string;          // Filter: STUDENT | LECTURER | STAFF
  status?: string;              // Filter: PENDING | VERIFIED | REJECTED | VOTED | BLOCKED
  voting_method?: string;       // Filter: ONLINE | TPS
  faculty_code?: string;        // Filter by faculty code
  study_program_code?: string;  // Filter by study program code
  cohort_year?: number;         // Filter by cohort year
  tps_id?: number;              // Filter by TPS ID
}
```

**Response 200:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "election_voter_id": 6,
        "election_id": 1,
        "voter_id": 1,
        "nim": "2021101",
        "name": "Agus Santoso",
        "email": "agus@example.com",
        "voter_type": "STUDENT",
        "faculty_code": "FT",
        "faculty_name": "Fakultas Teknik",
        "study_program_code": "TI",
        "study_program_name": "Teknik Informatika",
        "cohort_year": 2021,
        "academic_status": "ACTIVE",
        "status": "VERIFIED",
        "voting_method": "ONLINE",
        "tps_id": null,
        "checked_in_at": null,
        "voted_at": null,
        "has_voted": false,
        "updated_at": "2024-11-26T10:00:00Z"
      }
    ],
    "page": 1,
    "limit": 50,
    "total_items": 41,
    "total_pages": 1
  }
}
```

**Response 401:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing authorization header."
  }
}
```

**Response 403:**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Insufficient permissions."
  }
}
```

---

### 2. Lookup Voter by NIM

**Endpoint:** `GET /admin/elections/{electionID}/voters/lookup`

**Description:** Search for a voter by NIM and check their enrollment status in the election.

**Authentication:** Admin role required

**Path Parameters:**
```typescript
{
  electionID: number;
}
```

**Query Parameters:**
```typescript
{
  nim: string; // Required - NIM/NIDN/NIP to search
}
```

**Response 200:**
```json
{
  "success": true,
  "data": {
    "voter": {
      "id": 1,
      "nim": "2021101",
      "name": "Agus Santoso",
      "voter_type": "STUDENT",
      "email": "agus@example.com",
      "faculty_code": "FT",
      "study_program_code": "TI",
      "cohort_year": 2021,
      "academic_status": "ACTIVE",
      "has_account": true,
      "lecturer_id": null,
      "staff_id": null,
      "voting_method": "ONLINE"
    },
    "election_voter": {
      "election_voter_id": 6,
      "election_id": 1,
      "voter_id": 1,
      "nim": "2021101",
      "status": "VERIFIED",
      "voting_method": "ONLINE",
      "tps_id": null,
      "checked_in_at": null,
      "voted_at": null,
      "updated_at": "2024-11-26T10:00:00Z",
      "voter_type": "STUDENT",
      "name": "Agus Santoso",
      "email": "agus@example.com",
      "faculty_code": "FT",
      "study_program_code": "TI",
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

**Response 400:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "parameter nim wajib diisi"
  }
}
```

---

### 3. Add/Update Voter (Upsert)

**Endpoint:** `POST /admin/elections/{electionID}/voters`

**Description:** Add a new voter or update existing voter and enroll them in the election.

**Authentication:** Admin role required

**Path Parameters:**
```typescript
{
  electionID: number;
}
```

**Request Body:**
```json
{
  "voter_type": "STUDENT",
  "nim": "2021101",
  "name": "Agus Santoso",
  "email": "agus@example.com",
  "phone": "081234567890",
  "faculty_code": "FT",
  "faculty_name": "Fakultas Teknik",
  "study_program_code": "TI",
  "study_program_name": "Teknik Informatika",
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
- `voter_type`: Required, enum: `"STUDENT"` | `"LECTURER"` | `"STAFF"`
- `nim`: Required, unique string
- `name`: Required, string
- `email`: Optional, valid email format
- `phone`: Optional, phone number format
- `voting_method`: Required, enum: `"ONLINE"` | `"TPS"`
- `status`: Required, enum: `"PENDING"` | `"VERIFIED"` | `"REJECTED"` | `"BLOCKED"`
- `tps_id`: Optional, number (required if voting_method is TPS)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "voter_id": 1,
    "election_voter_id": 6,
    "status": "PENDING",
    "voting_method": "ONLINE",
    "tps_id": null,
    "created_voter": false,
    "created_election_voter": true,
    "duplicate_in_election": false
  }
}
```

**Response 409:**
```json
{
  "success": false,
  "error": {
    "code": "DUPLICATE",
    "message": "NIM sudah terdaftar di pemilu ini"
  }
}
```

**Response 400:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Data wajib diisi atau tidak valid"
  }
}
```

---

### 4. Update Voter Status

**Endpoint:** `PATCH /admin/elections/{electionID}/voters/{voterID}`

**Description:** Update voter enrollment status, voting method, or TPS assignment.

**Authentication:** Admin role required

**Path Parameters:**
```typescript
{
  electionID: number;  // Election ID
  voterID: number;     // Election Voter ID (not voter_id!)
}
```

**Request Body:**
```json
{
  "status": "VERIFIED",
  "voting_method": "TPS",
  "tps_id": 5
}
```

**Field Rules:**
- All fields are optional
- `status`: enum: `"PENDING"` | `"VERIFIED"` | `"REJECTED"` | `"VOTED"` | `"BLOCKED"`
- `voting_method`: enum: `"ONLINE"` | `"TPS"`
- `tps_id`: number or null

**Response 200:**
```json
{
  "success": true,
  "data": {
    "election_voter_id": 6,
    "election_id": 1,
    "voter_id": 1,
    "nim": "2021101",
    "status": "VERIFIED",
    "voting_method": "TPS",
    "tps_id": 5,
    "checked_in_at": null,
    "voted_at": null,
    "updated_at": "2024-11-26T10:05:00Z"
  }
}
```

**Response 404:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Data pemilih tidak ditemukan"
  }
}
```

**Response 400:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "status / voting_method tidak valid"
  }
}
```

---

### 5. Import DPT from CSV

**Endpoint:** `POST /admin/elections/{electionID}/voters/import`

**Description:** Bulk import voters from CSV file.

**Authentication:** Admin role required

**Path Parameters:**
```typescript
{
  electionID: number;
}
```

**Request:** `multipart/form-data`

**Form Fields:**
```typescript
{
  file: File; // CSV file
}
```

**CSV Format:**
```csv
nim,name,faculty,study_program,cohort_year
2021101,Agus Santoso,Fakultas Teknik,Teknik Informatika,2021
2021102,Budi Pratama,Fakultas Teknik,Teknik Elektro,2021
```

**Required Columns:**
- `nim` - NIM/NIDN/NIP
- `name` - Full name
- `faculty` - Faculty name
- `study_program` - Study program name
- `cohort_year` - Year enrolled (number)

**Response 200:**
```json
{
  "success": true,
  "data": {
    "success": 45,
    "failed": 2,
    "total": 47,
    "errors": [
      {
        "row": 3,
        "nim": "2021103",
        "error": "NIM sudah terdaftar"
      },
      {
        "row": 15,
        "nim": "2021115",
        "error": "cohort_year tidak valid"
      }
    ]
  }
}
```

**Response 400:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Field file wajib diisi."
  }
}
```

**Response 422:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Kolom 'nim' wajib ada di CSV."
  }
}
```

---

### 6. Export DPT to CSV

**Endpoint:** `GET /admin/elections/{electionID}/voters/export`

**Description:** Export all voters in election to CSV file.

**Authentication:** Admin role required

**Path Parameters:**
```typescript
{
  electionID: number;
}
```

**Query Parameters:**
```typescript
{
  // Same filters as list endpoint
  voter_type?: string;
  status?: string;
  voting_method?: string;
  faculty_code?: string;
  study_program_code?: string;
  cohort_year?: number;
}
```

**Response 200:**
```
Content-Type: text/csv
Content-Disposition: attachment; filename="dpt-election-1-20241126.csv"

nim,name,email,voter_type,faculty_name,study_program_name,cohort_year,academic_status,status,voting_method,has_voted,voted_at
2021101,Agus Santoso,agus@example.com,STUDENT,Fakultas Teknik,Teknik Informatika,2021,ACTIVE,VERIFIED,ONLINE,false,
2021102,Budi Pratama,budi@example.com,STUDENT,Fakultas Teknik,Teknik Elektro,2021,ACTIVE,VERIFIED,TPS,true,2024-11-26T14:30:00Z
```

---

## Voter Self-Service Endpoints

### 7. Self Register to Election

**Endpoint:** `POST /voters/me/elections/{electionID}/register`

**Description:** Voter self-registers for an election.

**Authentication:** Voter role required

**Path Parameters:**
```typescript
{
  electionID: number;
}
```

**Request Body:**
```json
{
  "voting_method": "ONLINE",
  "tps_id": null
}
```

**Field Rules:**
- `voting_method`: Required, enum: `"ONLINE"` | `"TPS"`
- `tps_id`: Required if voting_method is `"TPS"`, otherwise null

**Response 200:**
```json
{
  "success": true,
  "data": {
    "election_voter_id": 10,
    "election_id": 1,
    "voter_id": 5,
    "nim": "2021105",
    "status": "PENDING",
    "voting_method": "ONLINE",
    "tps_id": null,
    "checked_in_at": null,
    "voted_at": null,
    "updated_at": "2024-11-26T10:00:00Z"
  }
}
```

**Response 403:**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Akses tidak diizinkan"
  }
}
```

**Response 409:**
```json
{
  "success": false,
  "error": {
    "code": "DUPLICATE",
    "message": "Sudah terdaftar di pemilu ini"
  }
}
```

---

### 8. Get My Election Status

**Endpoint:** `GET /voters/me/elections/{electionID}/status`

**Description:** Get voter's enrollment status in specific election.

**Authentication:** Voter role required

**Path Parameters:**
```typescript
{
  electionID: number;
}
```

**Response 200:**
```json
{
  "success": true,
  "data": {
    "election_voter_id": 10,
    "election_id": 1,
    "voter_id": 5,
    "nim": "2021105",
    "status": "VERIFIED",
    "voting_method": "ONLINE",
    "tps_id": null,
    "checked_in_at": null,
    "voted_at": null,
    "updated_at": "2024-11-26T10:00:00Z",
    "voter_type": "STUDENT",
    "name": "Eka Putri",
    "email": "eka@example.com",
    "faculty_code": "FEB",
    "study_program_code": "AK",
    "cohort_year": 2021
  }
}
```

**Response 404:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Belum terdaftar di pemilu ini"
  }
}
```

---

## Data Models

### ElectionVoter (Full)

```typescript
interface ElectionVoter {
  election_voter_id: number;
  election_id: number;
  voter_id: number;
  nim: string;
  name: string;
  email: string | null;
  voter_type: "STUDENT" | "LECTURER" | "STAFF";
  faculty_code: string | null;
  faculty_name: string | null;
  study_program_code: string | null;
  study_program_name: string | null;
  cohort_year: number | null;
  academic_status: "ACTIVE" | "GRADUATED" | "ON_LEAVE" | "DROPPED" | "INACTIVE" | null;
  status: "PENDING" | "VERIFIED" | "REJECTED" | "VOTED" | "BLOCKED";
  voting_method: "ONLINE" | "TPS";
  tps_id: number | null;
  checked_in_at: string | null;  // ISO 8601
  voted_at: string | null;        // ISO 8601
  has_voted: boolean | null;
  updated_at: string;              // ISO 8601
}
```

### VoterSummary

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
  voting_method?: "ONLINE" | "TPS";
}
```

### LookupResult

```typescript
interface LookupResult {
  voter: VoterSummary;
  election_voter?: ElectionVoter;
}
```

### UpsertResult

```typescript
interface UpsertResult {
  voter_id: number;
  election_voter_id: number;
  status: string;
  voting_method: string;
  tps_id: number | null;
  created_voter: boolean;
  created_election_voter: boolean;
  duplicate_in_election: boolean;
}
```

### ImportResult

```typescript
interface ImportResult {
  success: number;
  failed: number;
  total: number;
  errors: Array<{
    row: number;
    nim: string;
    error: string;
  }>;
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid input data |
| `DUPLICATE` | 409 | NIM already registered in election |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Status Enum Values

### Enrollment Status

- `PENDING`: Newly registered, awaiting verification
- `VERIFIED`: Admin approved, can vote
- `REJECTED`: Admin rejected registration
- `VOTED`: Already voted
- `BLOCKED`: Temporarily blocked from voting

### Voting Method

- `ONLINE`: Vote via online system
- `TPS`: Vote at physical polling station (TPS)

### Voter Type

- `STUDENT`: Student voter (has NIM, cohort_year)
- `LECTURER`: Lecturer voter (has NIDN)
- `STAFF`: Staff voter (has NIP)

### Academic Status

- `ACTIVE`: Currently active
- `GRADUATED`: Already graduated
- `ON_LEAVE`: On leave (cuti)
- `DROPPED`: Dropped out (DO)
- `INACTIVE`: Inactive status

---

## Usage Examples

### Admin: Get DPT List with Filters

```bash
curl 'http://localhost:8080/api/v1/admin/elections/1/voters?page=1&limit=20&voter_type=STUDENT&status=VERIFIED' \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Admin: Add New Voter

```bash
curl -X POST 'http://localhost:8080/api/v1/admin/elections/1/voters' \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "voter_type": "STUDENT",
    "nim": "2021999",
    "name": "Test User",
    "email": "test@example.com",
    "faculty_code": "FT",
    "faculty_name": "Fakultas Teknik",
    "study_program_code": "TI",
    "study_program_name": "Teknik Informatika",
    "cohort_year": 2021,
    "voting_method": "ONLINE",
    "status": "PENDING"
  }'
```

### Admin: Update Voter Status

```bash
curl -X PATCH 'http://localhost:8080/api/v1/admin/elections/1/voters/10' \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "VERIFIED"
  }'
```

### Voter: Self Register

```bash
curl -X POST 'http://localhost:8080/api/v1/voters/me/elections/1/register' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "voting_method": "ONLINE",
    "tps_id": null
  }'
```

### Voter: Check Status

```bash
curl 'http://localhost:8080/api/v1/voters/me/elections/1/status' \
  -H "Authorization: Bearer YOUR_VOTER_TOKEN"
```

---

## Notes

### Important Considerations

1. **Enrollment vs Voting Status**
   - `election_voters.status`: Registration/enrollment status (PENDING, VERIFIED, etc.)
   - `has_voted`: Actual voting status from `voter_status` table

2. **Data Consistency**
   - System maintains two tables: `election_voters` and `voter_status`
   - Keep both in sync when voter votes

3. **NIM Uniqueness**
   - NIM must be unique per election (not globally)
   - Same NIM can be in multiple elections

4. **TPS Assignment**
   - If `voting_method` is `"TPS"`, `tps_id` must be provided
   - If `voting_method` is `"ONLINE"`, `tps_id` should be null

5. **Semester Calculation**
   - Frontend should calculate semester from `cohort_year`
   - Formula: `(current_year - cohort_year) * 2 + (month >= 8 ? 1 : 0)`

---

**Last Updated:** 2025-11-26  
**API Version:** 2.0  
**Maintained By:** Backend Team
