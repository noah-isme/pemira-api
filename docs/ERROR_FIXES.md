# Error Fixes - Frontend Integration Issues

## 🐛 Masalah yang Diperbaiki

### Error 1: 404 Not Found pada Detail Kandidat
```
GET /api/v1/elections/1/candidates/14
Response: 404 Not Found
```

### Error 2: 500 Internal Server Error pada Profile Media
```
GET /api/v1/elections/1/candidates/14/media/profile
Response: 500 Internal Server Error
```

---

## ✅ Solusi

### 1. Fix Query GetByID - Added deleted_at Filter

**Masalah:**
- Query `GetByID` tidak memfilter kandidat yang sudah dihapus (soft delete)
- Query `ListByElection` sudah ada filter `deleted_at IS NULL`
- Inkonsistensi ini menyebabkan kandidat muncul di list tapi tidak bisa diakses detail-nya

**File:** `internal/candidate/repository_pgx.go`

**Sebelum:**
```go
const qGetCandidateByID = `
SELECT ...
FROM candidates
WHERE election_id = $1 AND id = $2
`
```

**Sesudah:**
```go
const qGetCandidateByID = `
SELECT ...
FROM candidates
WHERE election_id = $1 AND id = $2 AND deleted_at IS NULL
`
```

**Queries yang diperbaiki:**
- ✅ `qGetCandidateByID`
- ✅ `qGetCandidateByIDNoPhotoMedia`
- ✅ `qGetCandidateByCandidateID`
- ✅ `qGetCandidateByCandidateIDNoPhotoMedia`

### 2. Fix Error Handling - Media Not Found

**Masalah:**
- Error `ErrCandidateMediaNotFound` tidak di-handle di public handler
- Semua error tidak dikenal dikembalikan sebagai 500 Internal Server Error
- Seharusnya 404 Not Found jika media tidak ada

**File:** `internal/candidate/http_handler.go`

**Sebelum:**
```go
func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
    switch {
    case errors.Is(err, ErrCandidateNotFound):
        response.NotFound(w, "NOT_FOUND", "Kandidat tidak ditemukan...")
    
    case errors.Is(err, ErrCandidateNotPublished):
        response.NotFound(w, "NOT_FOUND", "Kandidat tidak ditemukan...")
    
    default:
        // ❌ SEMUA error lain jadi 500
        response.InternalServerError(w, "INTERNAL_ERROR", "Terjadi kesalahan...")
    }
}
```

**Sesudah:**
```go
func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
    switch {
    case errors.Is(err, ErrCandidateNotFound):
        response.NotFound(w, "NOT_FOUND", "Kandidat tidak ditemukan...")
    
    case errors.Is(err, ErrCandidateNotPublished):
        response.NotFound(w, "NOT_FOUND", "Kandidat tidak ditemukan...")
    
    case errors.Is(err, ErrCandidateMediaNotFound):
        // ✅ Return 404 untuk media not found
        response.NotFound(w, "MEDIA_NOT_FOUND", "Media kandidat tidak ditemukan.")
    
    default:
        response.InternalServerError(w, "INTERNAL_ERROR", "Terjadi kesalahan...")
    }
}
```

---

## 🧪 Testing

### Test Detail Kandidat (Fixed ✅)

**Before:**
```bash
curl "http://localhost:8080/api/v1/elections/1/candidates/14"
# Response: {"code":"NOT_FOUND","message":"Kandidat tidak ditemukan..."}
# HTTP Status: 404 ❌
```

**After:**
```bash
curl "http://localhost:8080/api/v1/elections/1/candidates/14"
# Response:
{
  "data": {
    "id": 14,
    "name": "Ayu",
    "number": 77,
    "status": "PUBLISHED",
    "qr_code": {
      "id": 3,
      "token": "CAND14-NXLbfXqr7gtU",
      "payload": "PEMIRA-UNIWA|E:1|C:14|V:1",
      "version": 1,
      "is_active": true
    }
  }
}
# HTTP Status: 200 ✅
```

### Test Profile Media (Fixed ✅)

**Before:**
```bash
curl "http://localhost:8080/api/v1/elections/1/candidates/14/media/profile"
# Response: {"code":"INTERNAL_ERROR","message":"Terjadi kesalahan pada sistem."}
# HTTP Status: 500 ❌
```

**After:**
```bash
curl "http://localhost:8080/api/v1/elections/1/candidates/14/media/profile"
# Response: {"code":"MEDIA_NOT_FOUND","message":"Media kandidat tidak ditemukan."}
# HTTP Status: 404 ✅
```

---

## 📊 Response Format yang Benar

### Detail Kandidat dengan QR Code

**Endpoint:**
```
GET /api/v1/elections/{electionID}/candidates/{candidateID}
```

**Response Structure:**
```json
{
  "data": {
    "id": 14,
    "election_id": 1,
    "number": 77,
    "name": "Ayu",
    "photo_url": "",
    "short_bio": "test",
    "long_bio": "test",
    "tagline": "",
    "faculty_name": "Teknik",
    "study_program_name": "Informatika",
    "cohort_year": 2020,
    "vision": "test",
    "missions": ["test"],
    "main_programs": [
      {
        "title": "test",
        "description": "test",
        "category": ""
      }
    ],
    "media": {
      "video_url": ""
    },
    "social_links": [],
    "status": "PUBLISHED",
    "stats": {
      "total_votes": 0,
      "percentage": 0
    },
    "qr_code": {
      "id": 3,
      "token": "CAND14-NXLbfXqr7gtU",
      "url": "https://pemira.local/ballot-qr/CAND14-NXLbfXqr7gtU",
      "payload": "PEMIRA-UNIWA|E:1|C:14|V:1",
      "version": 1,
      "is_active": true
    }
  }
}
```

**Key Points:**
- ✅ Response dibungkus dalam `data` wrapper (public endpoint)
- ✅ Field `qr_code` ada di level `data`
- ✅ HTTP Status 200 jika berhasil
- ✅ HTTP Status 404 jika tidak ditemukan

---

## 🔍 Frontend Implementation

### React Example

```jsx
import { useEffect, useState } from 'react';
import QRCode from 'qrcode.react';

function CandidateDetail({ electionId, candidateId }) {
  const [candidate, setCandidate] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch(`/api/v1/elections/${electionId}/candidates/${candidateId}`)
      .then(res => {
        if (!res.ok) {
          throw new Error(`HTTP ${res.status}`);
        }
        return res.json();
      })
      .then(response => {
        // ⚠️ PENTING: Public endpoint punya wrapper 'data'
        setCandidate(response.data);
        setLoading(false);
      })
      .catch(err => {
        setError(err.message);
        setLoading(false);
      });
  }, [electionId, candidateId]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!candidate) return <div>Kandidat tidak ditemukan</div>;

  return (
    <div>
      <h2>{candidate.name}</h2>
      <p>Number: {candidate.number}</p>
      <p>Vision: {candidate.vision}</p>
      
      {/* QR Code - Always check if exists */}
      {candidate.qr_code && candidate.qr_code.payload ? (
        <div className="qr-section">
          <h3>QR Code untuk Voting</h3>
          <QRCode 
            value={candidate.qr_code.payload}
            size={256}
            level="H"
          />
          <p>Token: {candidate.qr_code.token}</p>
          <p>Version: {candidate.qr_code.version}</p>
        </div>
      ) : (
        <p>QR Code tidak tersedia</p>
      )}
    </div>
  );
}
```

### Vue.js Example

```vue
<template>
  <div>
    <div v-if="loading">Loading...</div>
    <div v-else-if="error">Error: {{ error }}</div>
    <div v-else-if="candidate">
      <h2>{{ candidate.name }}</h2>
      <p>Number: {{ candidate.number }}</p>
      
      <!-- QR Code -->
      <div v-if="candidate.qr_code?.payload" class="qr-section">
        <h3>QR Code untuk Voting</h3>
        <qrcode-vue 
          :value="candidate.qr_code.payload"
          :size="256"
          level="H"
        />
        <p>Token: {{ candidate.qr_code.token }}</p>
      </div>
      <p v-else>QR Code tidak tersedia</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import QrcodeVue from 'qrcode.vue';

const props = defineProps(['electionId', 'candidateId']);
const candidate = ref(null);
const loading = ref(true);
const error = ref(null);

onMounted(async () => {
  try {
    const response = await fetch(
      `/api/v1/elections/${props.electionId}/candidates/${props.candidateId}`
    );
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    const data = await response.json();
    // ⚠️ PENTING: Akses via data wrapper
    candidate.value = data.data;
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
});
</script>
```

---

## 🚨 Common Mistakes

### Mistake 1: Tidak Akses via `data` Wrapper

```javascript
// ❌ SALAH - Public endpoint punya wrapper
const qrCode = response.qr_code;

// ✅ BENAR - Akses via data
const qrCode = response.data.qr_code;
```

### Mistake 2: Tidak Cek QR Code Exists

```javascript
// ❌ SALAH - Bisa error jika qr_code null
<QRCode value={candidate.data.qr_code.payload} />

// ✅ BENAR - Cek dulu
{candidate.data.qr_code && (
  <QRCode value={candidate.data.qr_code.payload} />
)}
```

### Mistake 3: Salah Endpoint

```javascript
// ❌ SALAH - Admin endpoint (butuh auth)
GET /api/v1/admin/elections/1/candidates/14

// ✅ BENAR - Public endpoint (no auth)
GET /api/v1/elections/1/candidates/14
```

---

## 📝 Error Codes

### Candidate Endpoints

| HTTP Status | Error Code | Message | Meaning |
|-------------|------------|---------|---------|
| 200 | - | - | Success |
| 404 | `NOT_FOUND` | Kandidat tidak ditemukan | Kandidat tidak ada atau sudah dihapus |
| 404 | `MEDIA_NOT_FOUND` | Media kandidat tidak ditemukan | Kandidat tidak punya profile photo |
| 400 | `INVALID_REQUEST` | Invalid request | Parameter tidak valid |
| 500 | `INTERNAL_ERROR` | Terjadi kesalahan | Server error |

---

## 🔄 Comparison: Admin vs Public

### Admin Endpoint (Requires Auth)

**URL:** `/api/v1/admin/elections/{electionID}/candidates/{candidateID}`

**Response:** Langsung object (NO wrapper)
```json
{
  "id": 14,
  "name": "Ayu",
  "qr_code": { ... }
}
```

### Public Endpoint (No Auth)

**URL:** `/api/v1/elections/{electionID}/candidates/{candidateID}`

**Response:** Wrapped in `data`
```json
{
  "data": {
    "id": 14,
    "name": "Ayu",
    "qr_code": { ... }
  }
}
```

---

## 📚 Related Documentation

- [API Response Examples](API_RESPONSE_EXAMPLES.md)
- [Frontend QR Integration](../FRONTEND_QR_INTEGRATION.md)
- [QR Code Generation Guide](QR_CODE_GENERATION_GUIDE.md)
- [Solusi QR Tidak Muncul](../SOLUSI_QR_CODE_TIDAK_MUNCUL.md)

---

## ✅ Checklist Testing

- [x] Fix query GetByID dengan deleted_at filter
- [x] Fix error handling untuk media not found
- [x] Test detail kandidat (200 OK dengan QR code)
- [x] Test profile media tidak ada (404 Not Found)
- [x] Update dokumentasi
- [x] Deploy ke server

---

**Fixed:** 16 Desember 2024  
**Version:** 1.0.1  
**Status:** ✅ All Issues Resolved