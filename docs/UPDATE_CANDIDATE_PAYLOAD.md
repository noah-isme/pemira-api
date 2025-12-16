# Update Candidate - Payload Examples

## Endpoint
```
PUT /admin/elections/{electionID}/candidates/{candidateID}
```

## Headers
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

---

## 1. UPDATE STATUS ONLY (Most Common)

### Change to PUBLISHED (make visible to public)
```json
{
  "status": "PUBLISHED"
}
```

### Change to DRAFT (hide from public)
```json
{
  "status": "DRAFT"
}
```

### Change to HIDDEN (temporarily hide)
```json
{
  "status": "HIDDEN"
}
```

### Change to ARCHIVED (archive candidate)
```json
{
  "status": "ARCHIVED"
}
```

### Change to APPROVED (legacy, same as PUBLISHED)
```json
{
  "status": "APPROVED"
}
```

---

## 2. UPDATE NAME AND STATUS
```json
{
  "name": "Updated Candidate Name",
  "status": "PUBLISHED"
}
```

---

## 3. UPDATE BASIC INFO
```json
{
  "number": 1,
  "name": "Pasangan Calon 1",
  "tagline": "Updated Tagline",
  "short_bio": "Updated short bio",
  "status": "PUBLISHED"
}
```

---

## 4. UPDATE COMPLETE PROFILE
```json
{
  "number": 1,
  "name": "Pasangan Calon A - Updated",
  "photo_url": "https://example.com/photo.jpg",
  "short_bio": "Bio singkat yang diupdate",
  "long_bio": "Bio lengkap yang diupdate dengan detail lebih banyak",
  "tagline": "Tagline Baru untuk Kampus",
  "faculty_name": "Fakultas Teknik",
  "study_program_name": "Teknik Informatika",
  "cohort_year": 2021,
  "vision": "Visi yang diperbarui untuk kampus yang lebih baik",
  "missions": [
    "Misi 1: Meningkatkan kualitas mahasiswa",
    "Misi 2: Membangun ekosistem kolaboratif",
    "Misi 3: Transparansi penuh dalam kegiatan"
  ],
  "main_programs": [
    {
      "title": "Program Aspirasi Digital",
      "description": "Membangun platform untuk menampung aspirasi mahasiswa",
      "category": "Teknologi"
    },
    {
      "title": "Beasiswa Prestasi",
      "description": "Memberikan beasiswa untuk mahasiswa berprestasi",
      "category": "Kesejahteraan"
    }
  ],
  "media": {
    "video_url": "https://youtube.com/watch?v=abc123",
    "gallery_photos": [
      "https://example.com/photo1.jpg",
      "https://example.com/photo2.jpg"
    ],
    "document_manifesto_url": "https://example.com/manifesto.pdf"
  },
  "social_links": [
    {
      "platform": "instagram",
      "url": "https://instagram.com/paslon_a"
    },
    {
      "platform": "tiktok",
      "url": "https://tiktok.com/@paslon_a"
    }
  ],
  "status": "PUBLISHED"
}
```

---

## 5. PARTIAL UPDATE EXAMPLES

### Update only vision and missions
```json
{
  "vision": "New vision for the university",
  "missions": [
    "Mission 1",
    "Mission 2",
    "Mission 3"
  ]
}
```

### Update only main programs
```json
{
  "main_programs": [
    {
      "title": "New Program 1",
      "description": "Description of program 1",
      "category": "Education"
    },
    {
      "title": "New Program 2",
      "description": "Description of program 2",
      "category": "Welfare"
    }
  ]
}
```

### Update only social media links
```json
{
  "social_links": [
    {
      "platform": "instagram",
      "url": "https://instagram.com/updated_account"
    },
    {
      "platform": "twitter",
      "url": "https://twitter.com/updated_account"
    }
  ]
}
```

### Update only media
```json
{
  "media": {
    "video_url": "https://youtube.com/watch?v=newvideo",
    "gallery_photos": [],
    "document_manifesto_url": ""
  }
}
```

---

## 6. CURL EXAMPLES

### Update status to PUBLISHED
```bash
curl -X PUT "https://pemira-api-noah-isme4297-1lkqsxtc.apn.leapcell.dev/api/v1/admin/elections/1/candidates/4" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PUBLISHED"
  }'
```

### Update name and status
```bash
curl -X PUT "https://pemira-api-noah-isme4297-1lkqsxtc.apn.leapcell.dev/api/v1/admin/elections/1/candidates/4" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Candidate Name",
    "status": "PUBLISHED"
  }'
```

### Update complete profile
```bash
curl -X PUT "https://pemira-api-noah-isme4297-1lkqsxtc.apn.leapcell.dev/api/v1/admin/elections/1/candidates/4" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "number": 1,
    "name": "Pasangan Calon A",
    "tagline": "Untuk Mahasiswa Lebih Baik",
    "vision": "Visi lengkap",
    "missions": ["Misi 1", "Misi 2"],
    "status": "PUBLISHED"
  }'
```

---

## VALIDATION RULES

- ✅ All fields are **OPTIONAL** (partial update supported)
- `number`: Must be positive integer, unique per election
- `name`: Minimum 3 characters if provided
- `status`: Must be valid enum value
  - **Valid values:** `DRAFT`, `PENDING`, `PUBLISHED`, `APPROVED`, `HIDDEN`, `REJECTED`, `WITHDRAWN`, `ARCHIVED`
- `missions`: Array of strings
- `main_programs`: Array of objects with `title`, `description`, `category`
- `social_links`: Array of objects with `platform`, `url`
- `media`: Object with optional `video_url`, `gallery_photos`, `document_manifesto_url`

---

## RESPONSE EXAMPLES

### Success (200 OK)
```json
{
  "id": 4,
  "election_id": 1,
  "number": 1,
  "name": "Updated Name",
  "photo_url": "https://...",
  "short_bio": "Bio",
  "long_bio": "Long bio",
  "tagline": "Tagline",
  "faculty_name": "Fakultas Teknik",
  "study_program_name": "Informatika",
  "cohort_year": 2021,
  "vision": "Vision text",
  "missions": ["Mission 1", "Mission 2"],
  "main_programs": [
    {
      "title": "Program 1",
      "description": "Desc",
      "category": "Cat"
    }
  ],
  "media": {
    "video_url": "https://...",
    "gallery_photos": [],
    "document_manifesto_url": ""
  },
  "social_links": [
    {
      "platform": "instagram",
      "url": "https://..."
    }
  ],
  "status": "PUBLISHED",
  "stats": {
    "total_votes": 0,
    "percentage": 0
  },
  "created_at": "2025-01-15T08:00:00Z",
  "updated_at": "2025-01-20T10:30:00Z"
}
```

### Error - Validation Failed (422 Unprocessable Entity)
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Nomor urut kandidat sudah digunakan"
}
```

### Error - Not Found (404 Not Found)
```json
{
  "code": "NOT_FOUND",
  "message": "Kandidat tidak ditemukan"
}
```

### Error - Unauthorized (401 Unauthorized)
```json
{
  "code": "UNAUTHORIZED",
  "message": "Token tidak valid atau expired"
}
```

---

## STATUS MEANINGS

| Status | Description | Public Visible |
|--------|-------------|----------------|
| **DRAFT** | Work in progress, not yet submitted | ❌ No |
| **PENDING** | Under review by admin | ❌ No |
| **PUBLISHED** | Published and visible to voters (preferred) | ✅ Yes |
| **APPROVED** | Approved and visible (legacy, same as PUBLISHED) | ✅ Yes |
| **HIDDEN** | Temporarily hidden from public view | ❌ No |
| **REJECTED** | Rejected by admin | ❌ No |
| **WITHDRAWN** | Withdrawn by candidate | ❌ No |
| **ARCHIVED** | Archived, no longer active | ❌ No |

---

## TYPICAL WORKFLOW

1. **Create candidate** → Status: `DRAFT` (hidden from public)
2. **Complete profile** → Update with full details
3. **Review & approve** → Change status to `PUBLISHED` (visible to public)
4. **Temporarily hide** → Change status to `HIDDEN` (if needed)
5. **Archive after election** → Change status to `ARCHIVED`

---

## FRONTEND INTEGRATION TIP

```typescript
// React/Next.js example
const updateCandidateStatus = async (
  candidateId: number, 
  newStatus: string
) => {
  const response = await fetch(
    `${API_BASE}/admin/elections/1/candidates/${candidateId}`,
    {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ status: newStatus })
    }
  );
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message);
  }
  
  return await response.json();
};

// Usage
await updateCandidateStatus(4, 'PUBLISHED');
```

---

For more details, see: `docs/CANDIDATE_ENDPOINTS_FRONTEND.md`
