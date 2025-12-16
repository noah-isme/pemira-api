# QR Code Generation Guide - Kandidat

## 🎯 Masalah yang Diselesaikan

**Masalah:** Kandidat yang baru dibuat tidak otomatis memiliki QR code, sehingga tidak muncul di frontend.

**Solusi:** 
1. Auto-generate QR code saat kandidat di-publish
2. Endpoint manual untuk generate QR code

---

## ✅ Implementasi

### 1. Auto-Generate saat Publish

QR code akan **otomatis dibuat** saat kandidat di-publish (status menjadi APPROVED).

**Cara:**
```bash
POST /api/v1/admin/elections/{electionID}/candidates/{candidateID}/publish
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "id": 14,
  "name": "Ayu",
  "status": "APPROVED",
  "qr_code": {
    "id": 3,
    "token": "CAND14-NXLbfXqr7gtU",
    "url": "https://pemira.local/ballot-qr/CAND14-NXLbfXqr7gtU",
    "payload": "PEMIRA-UNIWA|E:1|C:14|V:1",
    "version": 1,
    "is_active": true
  }
}
```

### 2. Manual Generate (NEW)

Untuk kandidat yang sudah ada tapi belum punya QR code, gunakan endpoint baru:

**Endpoint:**
```bash
POST /api/v1/admin/elections/{electionID}/candidates/{candidateID}/qr/generate
Authorization: Bearer {access_token}
```

**Example:**
```bash
TOKEN="your-access-token"
curl -X POST "http://localhost:8080/api/v1/admin/elections/1/candidates/14/qr/generate" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

**Response:**
```json
{
  "id": 14,
  "election_id": 1,
  "number": 77,
  "name": "Ayu",
  "status": "PUBLISHED",
  "qr_code": {
    "id": 3,
    "token": "CAND14-NXLbfXqr7gtU",
    "url": "https://pemira.local/ballot-qr/CAND14-NXLbfXqr7gtU",
    "payload": "PEMIRA-UNIWA|E:1|C:14|V:1",
    "version": 1,
    "is_active": true
  }
}
```

---

## 🔧 Cara Generate QR untuk Kandidat yang Sudah Ada

### Opsi 1: Publish/Unpublish/Publish Lagi
```bash
# 1. Unpublish kandidat
POST /api/v1/admin/elections/1/candidates/14/unpublish

# 2. Publish lagi (QR akan auto-generate)
POST /api/v1/admin/elections/1/candidates/14/publish
```

### Opsi 2: Generate Manual (Recommended)
```bash
# Generate QR langsung tanpa ubah status
POST /api/v1/admin/elections/1/candidates/14/qr/generate
```

---

## 📝 Format QR Code

### Token Format
```
CAND{candidateID}-{random12chars}
```

**Contoh:**
- `CAND14-NXLbfXqr7gtU`
- `CAND01-ABC123XYZ456`
- `CAND05-DeF789GhI012`

### Payload Format
```
PEMIRA-UNIWA|E:{election_id}|C:{candidate_id}|V:{version}
```

**Contoh:**
- `PEMIRA-UNIWA|E:1|C:14|V:1`
- `PEMIRA-UNIWA|E:3|C:4|V:1`

---

## 🔄 QR Code Rotation

Jika perlu regenerate QR code (misal: bocor/kompromi):

**Cara:**
```bash
POST /api/v1/admin/elections/1/candidates/14/qr/generate
```

**Proses:**
1. QR lama di-set `is_active = false`
2. Version number bertambah (v1 → v2)
3. QR baru dibuat dengan token baru
4. Hanya QR dengan `is_active = true` yang muncul di response

**Contoh:**
```sql
-- Sebelum rotation
id | candidate_id | version | is_active
3  | 14          | 1       | true

-- Setelah rotation
id | candidate_id | version | is_active
3  | 14          | 1       | false    -- QR lama
4  | 14          | 2       | true     -- QR baru
```

---

## 🚀 Bulk Generate untuk Semua Kandidat

Jika ada banyak kandidat yang belum punya QR code:

### Script Bash
```bash
#!/bin/bash

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' | \
  jq -r '.access_token')

ELECTION_ID=1

# Get all candidates
CANDIDATES=$(curl -s "http://localhost:8080/api/v1/admin/elections/${ELECTION_ID}/candidates" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.items[].id')

# Generate QR for each candidate
for CANDIDATE_ID in $CANDIDATES; do
  echo "Generating QR for candidate ${CANDIDATE_ID}..."
  curl -s -X POST \
    "http://localhost:8080/api/v1/admin/elections/${ELECTION_ID}/candidates/${CANDIDATE_ID}/qr/generate" \
    -H "Authorization: Bearer $TOKEN" | jq -r '.qr_code.token'
done

echo "Done!"
```

### JavaScript/Node.js
```javascript
const axios = require('axios');

const API_URL = 'http://localhost:8080';
const ELECTION_ID = 1;

async function generateAllQRCodes() {
  // Login
  const loginRes = await axios.post(`${API_URL}/api/v1/auth/login`, {
    username: 'admin',
    password: 'password123'
  });
  
  const token = loginRes.data.access_token;
  const headers = { Authorization: `Bearer ${token}` };
  
  // Get all candidates
  const candidatesRes = await axios.get(
    `${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates`,
    { headers }
  );
  
  const candidates = candidatesRes.data.items;
  
  // Generate QR for each candidate without QR code
  for (const candidate of candidates) {
    if (!candidate.qr_code) {
      console.log(`Generating QR for ${candidate.name}...`);
      
      const res = await axios.post(
        `${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates/${candidate.id}/qr/generate`,
        {},
        { headers }
      );
      
      console.log(`✓ Generated: ${res.data.qr_code.token}`);
    } else {
      console.log(`✓ ${candidate.name} already has QR code`);
    }
  }
  
  console.log('Done!');
}

generateAllQRCodes();
```

---

## 🔍 Verifikasi QR Code

### Cek di Database
```sql
-- Lihat semua QR codes
SELECT 
  c.id,
  c.name,
  c.number,
  qr.id as qr_id,
  qr.qr_token,
  qr.version,
  qr.is_active
FROM candidates c
LEFT JOIN candidate_qr_codes qr ON c.id = qr.candidate_id
WHERE c.election_id = 1
ORDER BY c.id;

-- Kandidat tanpa QR code
SELECT c.id, c.name, c.number, c.status
FROM candidates c
LEFT JOIN candidate_qr_codes qr ON c.id = qr.candidate_id AND qr.is_active = true
WHERE c.election_id = 1 AND qr.id IS NULL;
```

### Cek via API
```bash
TOKEN="your-token"

# Get candidate detail
curl -s "http://localhost:8080/api/v1/admin/elections/1/candidates/14" \
  -H "Authorization: Bearer $TOKEN" | jq '.qr_code'

# List all candidates
curl -s "http://localhost:8080/api/v1/admin/elections/1/candidates" \
  -H "Authorization: Bearer $TOKEN" | jq '.items[] | {id, name, has_qr: (.qr_code != null)}'
```

---

## ⚠️ Important Notes

### 1. QR Code adalah Optional
Field `qr_code` bisa `null` jika kandidat belum punya QR code. Frontend harus handle ini:

```javascript
// ✅ BENAR
{candidate.qr_code && (
  <QRCode value={candidate.qr_code.payload} />
)}

// ❌ SALAH - bisa error
<QRCode value={candidate.qr_code.payload} />
```

### 2. Auto-Generate hanya pada Publish
QR code **TIDAK** otomatis dibuat saat:
- Create candidate
- Update candidate
- Unpublish candidate

QR code **OTOMATIS** dibuat saat:
- Publish candidate (status → APPROVED)

### 3. Regenerate = New Version
Setiap regenerate QR code akan:
- Increment version number
- Generate token baru
- Deactivate QR lama
- Hanya 1 QR active per kandidat

### 4. QR Code Tidak Dihapus
Saat kandidat dihapus (soft delete), QR code tetap ada di database untuk audit trail.

---

## 🐛 Troubleshooting

### QR Code tidak muncul di frontend
**Cek:**
1. Apakah kandidat punya QR code di database?
   ```sql
   SELECT * FROM candidate_qr_codes WHERE candidate_id = 14;
   ```

2. Apakah response API ada field `qr_code`?
   ```bash
   curl -s "http://localhost:8080/api/v1/admin/elections/1/candidates/14" \
     -H "Authorization: Bearer $TOKEN" | jq '.qr_code'
   ```

3. Apakah frontend cek `qr_code` exists?
   ```javascript
   if (candidate.qr_code) { ... }
   ```

**Solusi:**
```bash
# Generate QR code
curl -X POST "http://localhost:8080/api/v1/admin/elections/1/candidates/14/qr/generate" \
  -H "Authorization: Bearer $TOKEN"
```

### Generate QR Failed
**Possible errors:**
- `CANDIDATE_NOT_FOUND` - Kandidat tidak ada
- `INVALID_REQUEST` - electionID atau candidateID tidak valid
- `UNAUTHORIZED` - Token tidak valid

**Debug:**
```bash
# Full error response
curl -X POST "http://localhost:8080/api/v1/admin/elections/1/candidates/14/qr/generate" \
  -H "Authorization: Bearer $TOKEN" -v
```

---

## 📊 Database Schema

```sql
CREATE TABLE candidate_qr_codes (
    id BIGSERIAL PRIMARY KEY,
    election_id BIGINT NOT NULL REFERENCES elections(id),
    candidate_id BIGINT NOT NULL REFERENCES candidates(id),
    version INT NOT NULL DEFAULT 1,
    qr_token VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    rotated_at TIMESTAMP,
    
    CONSTRAINT uq_candidate_qr_active 
        UNIQUE (candidate_id, election_id, is_active) 
        DEFERRABLE INITIALLY IMMEDIATE
);

CREATE INDEX idx_candidate_qr_candidate 
    ON candidate_qr_codes(candidate_id);
    
CREATE INDEX idx_candidate_qr_active 
    ON candidate_qr_codes(candidate_id, is_active) 
    WHERE is_active = true;
```

---

## 📚 Related Documentation

- [API Response Examples](API_RESPONSE_EXAMPLES.md) - Format response lengkap
- [Frontend Integration](../FRONTEND_QR_INTEGRATION.md) - Cara integrase di frontend
- [Implementation Summary](../IMPLEMENTATION_SUMMARY.md) - Ringkasan teknis

---

**Created:** 16 Desember 2024  
**Version:** 1.0.0  
**Status:** ✅ Production Ready