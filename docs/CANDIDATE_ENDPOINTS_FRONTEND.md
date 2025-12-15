# Candidate API Endpoints - Frontend Reference

## Base URL
```
Production: https://pemira-api-noah-isme4297-1lkqsxtc.apn.leapcell.dev/api/v1
Local: http://localhost:8080/api/v1
```

---

## 📱 Public Endpoints (For Students/Voters)

### 1. List Published Candidates
**Endpoint:** `GET /elections/{electionID}/candidates`

**Description:** Get list of published candidates for an election (only APPROVED or PUBLISHED status visible)

**Auth:** Optional (no auth required for public view)

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `search` | string | - | Search by name or tagline |
| `page` | integer | 1 | Page number |
| `limit` | integer | 10 | Items per page |

**Response:**
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "election_id": 1,
        "number": 1,
        "name": "Pasangan Calon A",
        "photo_url": "https://...",
        "photo_media_id": "uuid",
        "short_bio": "Bio singkat",
        "tagline": "Tagline kandidat",
        "faculty_name": "Fakultas Teknik",
        "study_program_name": "Informatika",
        "status": "PUBLISHED",
        "stats": {
          "total_votes": 0,
          "percentage": 0
        },
        "qr_code": {
          "id": 1,
          "token": "abc123",
          "url": "https://...",
          "payload": "...",
          "version": 1,
          "is_active": true
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total_items": 5,
      "total_pages": 1
    }
  }
}
```

---

### 2. Get Candidate Detail (Public)
**Endpoint:** `GET /elections/{electionID}/candidates/{candidateID}`

**Description:** Get detailed information about a specific candidate

**Auth:** Optional

**Response:**
```json
{
  "data": {
    "id": 1,
    "election_id": 1,
    "number": 1,
    "name": "Pasangan Calon A",
    "photo_url": "https://...",
    "short_bio": "Bio singkat...",
    "long_bio": "Bio lengkap...",
    "tagline": "Tagline...",
    "faculty_name": "Fakultas Teknik",
    "study_program_name": "Informatika",
    "cohort_year": 2021,
    "vision": "Visi lengkap...",
    "missions": [
      "Misi 1",
      "Misi 2",
      "Misi 3"
    ],
    "main_programs": [
      {
        "title": "Program Utama 1",
        "description": "Deskripsi...",
        "category": "Kategori"
      }
    ],
    "media": {
      "video_url": "https://youtube.com/...",
      "gallery_photos": [],
      "document_manifesto_url": ""
    },
    "social_links": [
      {
        "platform": "instagram",
        "url": "https://instagram.com/..."
      }
    ],
    "status": "PUBLISHED",
    "stats": {
      "total_votes": 0,
      "percentage": 0
    }
  }
}
```

---

### 3. Get Candidates with QR Codes
**Endpoint:** `GET /elections/{electionID}/qr-codes`

**Description:** Get all published candidates with their QR codes (for TPS voting)

**Auth:** Optional

**Response:** Same as List Published Candidates but guaranteed to include `qr_code` field

---

## 🔒 Admin Endpoints (For Election Administrators)

**Auth Required:** Bearer Token with ADMIN/PANITIA role

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

---

### 4. List All Candidates (Admin)
**Endpoint:** `GET /admin/elections/{electionID}/candidates`

**Description:** Get all candidates including draft, hidden, and deleted (soft deleted excluded by default)

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `search` | string | - | Search by name |
| `status` | string | - | Filter by status: DRAFT, PENDING, PUBLISHED, APPROVED, HIDDEN, REJECTED, WITHDRAWN, ARCHIVED |
| `page` | integer | 1 | Page number |
| `limit` | integer | 20 | Items per page |

**Response:**
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "election_id": 1,
        "number": 1,
        "name": "Pasangan Calon A",
        "status": "PUBLISHED",
        "created_at": "2025-01-15T08:00:00Z",
        "updated_at": "2025-01-20T10:30:00Z",
        ...
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total_items": 15,
      "total_pages": 1
    }
  }
}
```

---

### 5. Create Candidate
**Endpoint:** `POST /admin/elections/{electionID}/candidates`

**Request Body:**
```json
{
  "number": 1,
  "name": "Pasangan Calon A",
  "photo_url": "",
  "short_bio": "Bio singkat",
  "long_bio": "Bio lengkap",
  "tagline": "Tagline",
  "faculty_name": "Fakultas Teknik",
  "study_program_name": "Informatika",
  "cohort_year": 2021,
  "vision": "Visi lengkap",
  "missions": ["Misi 1", "Misi 2"],
  "main_programs": [
    {
      "title": "Program 1",
      "description": "Desc",
      "category": "Category"
    }
  ],
  "media": {
    "video_url": ""
  },
  "social_links": [
    {
      "platform": "instagram",
      "url": "https://instagram.com/..."
    }
  ],
  "status": "DRAFT"
}
```

**Response:** Returns created candidate object (201 Created)

**Validation:**
- `number`: Required, unique per election, positive integer
- `name`: Required, min 3 chars
- `status`: Optional, defaults to DRAFT
- Valid statuses: DRAFT, PENDING, PUBLISHED, APPROVED, HIDDEN, REJECTED, WITHDRAWN, ARCHIVED

---

### 6. Update Candidate
**Endpoint:** `PUT /admin/elections/{electionID}/candidates/{candidateID}`

**Request Body:** Same as Create, but all fields optional (partial update)

```json
{
  "status": "PUBLISHED",
  "name": "Updated Name"
}
```

**Response:** Returns updated candidate object (200 OK)

---

### 7. Delete Candidate (Soft Delete)
**Endpoint:** `DELETE /admin/elections/{electionID}/candidates/{candidateID}`

**Description:** Soft delete candidate (sets deleted_at timestamp, candidate won't appear in lists but data preserved)

**Response:**
```json
{
  "data": {
    "message": "Kandidat berhasil dihapus (soft delete)"
  }
}
```

**Note:** Soft deleted candidates are excluded from all listing endpoints (both public and admin)

---

### 8. Publish Candidate
**Endpoint:** `POST /admin/elections/{electionID}/candidates/{candidateID}/publish`

**Description:** Publish candidate (set status to APPROVED)

**Response:** Returns updated candidate object

---

### 9. Unpublish Candidate
**Endpoint:** `POST /admin/elections/{electionID}/candidates/{candidateID}/unpublish`

**Description:** Unpublish candidate (set status to PENDING)

**Response:** Returns updated candidate object

---

### 10. Get Candidate Detail (Admin)
**Endpoint:** `GET /admin/elections/{electionID}/candidates/{candidateID}`

**Description:** Get full candidate details including all statuses

**Response:** Returns complete candidate object

---

## 📊 Candidate Status Flow

```
DRAFT → PENDING → APPROVED/PUBLISHED → (Visible to Public)
                      ↓
                  HIDDEN (Temporarily hide from public)
                      ↓
                  ARCHIVED (Long-term archive)
```

**Status Meanings:**
- **DRAFT**: Not yet submitted, work in progress
- **PENDING**: Under review by admin
- **PUBLISHED**: Published and visible to public (preferred)
- **APPROVED**: Approved and visible to public (legacy, kept for compatibility)
- **HIDDEN**: Temporarily hidden from public view
- **REJECTED**: Rejected by admin
- **WITHDRAWN**: Withdrawn by candidate
- **ARCHIVED**: Archived, no longer active

**Public Visibility Rules:**
- ✅ Visible: PUBLISHED, APPROVED
- ❌ Hidden: DRAFT, PENDING, HIDDEN, REJECTED, WITHDRAWN, ARCHIVED, Soft Deleted

---

## 🔄 Frontend Implementation Example

### React/Next.js Example

```typescript
// api/candidates.ts
import axios from 'axios';

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

// Public: List candidates
export const getPublicCandidates = async (electionId: number, page = 1) => {
  const { data } = await axios.get(
    `${API_BASE}/elections/${electionId}/candidates`,
    { params: { page, limit: 10 } }
  );
  return data.data;
};

// Admin: List all candidates
export const getAdminCandidates = async (
  electionId: number, 
  token: string,
  filters?: { status?: string; search?: string; page?: number }
) => {
  const { data } = await axios.get(
    `${API_BASE}/admin/elections/${electionId}/candidates`,
    {
      headers: { Authorization: `Bearer ${token}` },
      params: filters
    }
  );
  return data.data;
};

// Admin: Create candidate
export const createCandidate = async (
  electionId: number,
  token: string,
  candidateData: any
) => {
  const { data } = await axios.post(
    `${API_BASE}/admin/elections/${electionId}/candidates`,
    candidateData,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return data.data;
};

// Admin: Update candidate status
export const updateCandidateStatus = async (
  electionId: number,
  candidateId: number,
  token: string,
  status: string
) => {
  const { data } = await axios.put(
    `${API_BASE}/admin/elections/${electionId}/candidates/${candidateId}`,
    { status },
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return data.data;
};

// Admin: Delete candidate (soft delete)
export const deleteCandidate = async (
  electionId: number,
  candidateId: number,
  token: string
) => {
  const { data } = await axios.delete(
    `${API_BASE}/admin/elections/${electionId}/candidates/${candidateId}`,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return data.data;
};
```

### Usage in Component

```typescript
// pages/admin/candidates.tsx
import { useState, useEffect } from 'react';
import { getAdminCandidates, deleteCandidate, updateCandidateStatus } from '@/api/candidates';

export default function AdminCandidatesPage() {
  const [candidates, setCandidates] = useState([]);
  const [filter, setFilter] = useState({ status: '', search: '' });
  const token = "..."; // Get from auth context

  const loadCandidates = async () => {
    const data = await getAdminCandidates(1, token, filter);
    setCandidates(data.items);
  };

  const handleDelete = async (candidateId: number) => {
    if (confirm('Hapus kandidat ini? (Soft delete)')) {
      await deleteCandidate(1, candidateId, token);
      loadCandidates();
    }
  };

  const handleStatusChange = async (candidateId: number, newStatus: string) => {
    await updateCandidateStatus(1, candidateId, token, newStatus);
    loadCandidates();
  };

  useEffect(() => {
    loadCandidates();
  }, [filter]);

  return (
    <div>
      {/* Filter by status */}
      <select onChange={e => setFilter({...filter, status: e.target.value})}>
        <option value="">All</option>
        <option value="DRAFT">Draft</option>
        <option value="PUBLISHED">Published</option>
        <option value="HIDDEN">Hidden</option>
        <option value="ARCHIVED">Archived</option>
      </select>

      {/* Candidate list */}
      {candidates.map(candidate => (
        <div key={candidate.id}>
          <h3>{candidate.number}. {candidate.name}</h3>
          <span>Status: {candidate.status}</span>
          
          {/* Status dropdown */}
          <select 
            value={candidate.status}
            onChange={e => handleStatusChange(candidate.id, e.target.value)}
          >
            <option value="DRAFT">Draft</option>
            <option value="PENDING">Pending</option>
            <option value="PUBLISHED">Published</option>
            <option value="HIDDEN">Hidden</option>
            <option value="ARCHIVED">Archived</option>
          </select>

          {/* Delete button */}
          <button onClick={() => handleDelete(candidate.id)}>
            Delete (Soft)
          </button>
        </div>
      ))}
    </div>
  );
}
```

---

## ⚠️ Important Notes

1. **Soft Delete**: Deleted candidates are not permanently removed, just marked as deleted with `deleted_at` timestamp
2. **Status Visibility**: Only PUBLISHED and APPROVED candidates are visible on public endpoints
3. **Authentication**: Admin endpoints require valid JWT token with ADMIN or PANITIA role
4. **Validation**: Number must be unique per election
5. **Default Status**: New candidates default to DRAFT status

---

## 🔧 Error Handling

**Common Error Codes:**
- `INVALID_REQUEST`: Invalid parameters or body
- `VALIDATION_ERROR`: Validation failed (e.g., duplicate number)
- `NOT_FOUND`: Candidate not found
- `UNAUTHORIZED`: Missing or invalid authentication
- `INTERNAL_ERROR`: Server error

**Example Error Response:**
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Nomor urut kandidat sudah digunakan"
}
```
