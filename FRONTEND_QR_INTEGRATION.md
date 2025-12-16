# Frontend QR Code Integration - Quick Reference

## 🎯 Quick Summary

### Response Structure (PENTING!)

```javascript
// ❌ SALAH - Response TIDAK seperti ini:
{
  "data": {
    "qr_code": { ... }
  }
}

// ✅ BENAR - Response langsung seperti ini:
{
  "id": 4,
  "name": "Kandidat A",
  "qr_code": {
    "id": 1,
    "token": "CAND01-ABC123XYZ",
    "url": "https://pemira.local/ballot-qr/CAND01-ABC123XYZ",
    "payload": "PEMIRA-UNIWA|E:3|C:4|V:1",
    "version": 1,
    "is_active": true
  }
}
```

---

## 📡 API Endpoints

### 1. Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

**PENTING:** Token ada di `access_token`, bukan `token`!

### 2. Get Candidate Detail
```http
GET /api/v1/admin/elections/{electionID}/candidates/{candidateID}
Authorization: Bearer {access_token}
```

**Response:** Langsung object kandidat (tidak ada wrapper `data`)

### 3. List Candidates
```http
GET /api/v1/admin/elections/{electionID}/candidates
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "items": [ ... ],
  "pagination": { ... }
}
```

---

## ⚡ Quick Implementation

### React + TypeScript

```typescript
import { useState, useEffect } from 'react';
import QRCode from 'qrcode.react';

interface QRCodeData {
  id: number;
  token: string;
  payload: string;
  version: number;
  is_active: boolean;
}

interface Candidate {
  id: number;
  name: string;
  qr_code?: QRCodeData; // Optional!
}

function CandidateDetail({ electionId, candidateId }: Props) {
  const [candidate, setCandidate] = useState<Candidate | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    
    fetch(`/api/v1/admin/elections/${electionId}/candidates/${candidateId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
      .then(res => res.json())
      .then(data => setCandidate(data)); // Langsung data, tidak .data.data
  }, [electionId, candidateId]);

  if (!candidate) return <div>Loading...</div>;

  return (
    <div>
      <h2>{candidate.name}</h2>
      
      {/* Cek qr_code exists */}
      {candidate.qr_code && (
        <div className="qr-section">
          <h3>QR Code</h3>
          <QRCode 
            value={candidate.qr_code.payload}
            size={256}
            level="H"
          />
          <p>Version: {candidate.qr_code.version}</p>
        </div>
      )}
      
      {!candidate.qr_code && (
        <p>QR Code tidak tersedia</p>
      )}
    </div>
  );
}
```

### Vue.js 3

```vue
<template>
  <div>
    <h2>{{ candidate?.name }}</h2>
    
    <!-- QR Code Section -->
    <div v-if="candidate?.qr_code" class="qr-section">
      <h3>QR Code</h3>
      <qrcode-vue 
        :value="candidate.qr_code.payload"
        :size="256"
        level="H"
      />
      <p>Version: {{ candidate.qr_code.version }}</p>
    </div>
    
    <p v-else>QR Code tidak tersedia</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import QrcodeVue from 'qrcode.vue';

interface Candidate {
  id: number;
  name: string;
  qr_code?: {
    id: number;
    token: string;
    payload: string;
    version: number;
    is_active: boolean;
  };
}

const props = defineProps<{
  electionId: number;
  candidateId: number;
}>();

const candidate = ref<Candidate | null>(null);

onMounted(async () => {
  const token = localStorage.getItem('access_token');
  
  const response = await fetch(
    `/api/v1/admin/elections/${props.electionId}/candidates/${props.candidateId}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  candidate.value = await response.json(); // Langsung assign
});
</script>
```

### Vanilla JavaScript

```javascript
async function loadCandidate(electionId, candidateId) {
  const token = localStorage.getItem('access_token');
  
  const response = await fetch(
    `/api/v1/admin/elections/${electionId}/candidates/${candidateId}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  const candidate = await response.json(); // Langsung object
  
  // Check QR code
  if (candidate.qr_code && candidate.qr_code.payload) {
    // Generate QR code
    const canvas = document.getElementById('qr-canvas');
    QRCode.toCanvas(canvas, candidate.qr_code.payload, {
      width: 256,
      errorCorrectionLevel: 'H'
    });
  }
}
```

---

## 🚨 Common Mistakes & Fixes

### Mistake 1: Mencari data di wrapper yang tidak ada
```javascript
// ❌ SALAH
const qrCode = response.data.qr_code;

// ✅ BENAR
const qrCode = response.qr_code;
```

### Mistake 2: Tidak cek qr_code exists
```javascript
// ❌ SALAH - Bisa error jika qr_code null
<QRCode value={candidate.qr_code.payload} />

// ✅ BENAR - Cek dulu
{candidate.qr_code && (
  <QRCode value={candidate.qr_code.payload} />
)}
```

### Mistake 3: Salah ambil token dari login
```javascript
// ❌ SALAH
const token = loginResponse.token;

// ✅ BENAR
const token = loginResponse.access_token;
```

### Mistake 4: Format Authorization header salah
```javascript
// ❌ SALAH
headers: {
  'Authorization': token
}

// ✅ BENAR
headers: {
  'Authorization': `Bearer ${token}`
}
```

---

## 📦 Install QR Code Libraries

### React
```bash
npm install qrcode.react
# or
yarn add qrcode.react
```

### Vue.js
```bash
npm install qrcode.vue
# or
yarn add qrcode.vue
```

### Vanilla JS (CDN)
```html
<script src="https://cdn.jsdelivr.net/npm/qrcode@1.5.3/build/qrcode.min.js"></script>
```

---

## 🧪 Test dengan Browser Console

```javascript
// 1. Login
fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'admin',
    password: 'password123'
  })
})
  .then(r => r.json())
  .then(data => {
    console.log('Token:', data.access_token);
    localStorage.setItem('access_token', data.access_token);
  });

// 2. Get candidate
fetch('http://localhost:8080/api/v1/admin/elections/3/candidates/4', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('access_token')}`
  }
})
  .then(r => r.json())
  .then(candidate => {
    console.log('Candidate:', candidate);
    console.log('QR Code:', candidate.qr_code);
    console.log('QR Payload:', candidate.qr_code?.payload);
  });
```

---

## ✅ Checklist Integration

- [ ] Install QR code library
- [ ] Get `access_token` dari login (bukan `token`)
- [ ] Parse response langsung tanpa `.data` wrapper
- [ ] Cek `candidate.qr_code` exists sebelum render
- [ ] Gunakan `candidate.qr_code.payload` untuk generate QR
- [ ] Handle case ketika `qr_code` adalah `null`
- [ ] Set header `Authorization: Bearer {token}`
- [ ] Test di browser console dulu

---

## 📋 TypeScript Types

```typescript
interface LoginResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: {
    id: number;
    username: string;
    role: string;
  };
}

interface QRCode {
  id: number;
  token: string;
  url: string;
  payload: string;
  version: number;
  is_active: boolean;
}

interface Candidate {
  id: number;
  election_id: number;
  number: number;
  name: string;
  photo_url: string;
  short_bio: string;
  long_bio: string;
  tagline: string;
  faculty_name: string;
  study_program_name: string;
  vision: string;
  missions: string[];
  main_programs: any[];
  media: any;
  social_links: any[];
  status: string;
  stats: {
    total_votes: number;
    percentage: number;
  };
  qr_code?: QRCode; // Optional - bisa null!
}

interface CandidateListResponse {
  items: Candidate[];
  pagination: {
    page: number;
    limit: number;
    total_items: number;
    total_pages: number;
  };
}
```

---

## 🎨 Styling Example

```css
.qr-section {
  padding: 20px;
  border: 2px solid #ddd;
  border-radius: 8px;
  text-align: center;
  background: #fff;
}

.qr-section h3 {
  margin-bottom: 15px;
  color: #333;
}

.qr-section canvas {
  border: 1px solid #eee;
  padding: 10px;
  background: white;
}

.qr-info {
  margin-top: 15px;
  font-size: 14px;
  color: #666;
}

.no-qr {
  padding: 20px;
  text-align: center;
  color: #999;
  background: #f5f5f5;
  border-radius: 8px;
}
```

---

## 🔗 Resources

- **Full API Examples:** [docs/API_RESPONSE_EXAMPLES.md](docs/API_RESPONSE_EXAMPLES.md)
- **Implementation Details:** [docs/changes/ADMIN_CANDIDATE_QR_CODE.md](docs/changes/ADMIN_CANDIDATE_QR_CODE.md)
- **Docker Guide:** [docs/RUNNING_WITH_DOCKER.md](docs/RUNNING_WITH_DOCKER.md)

---

## 💡 Tips

1. **Always check if QR code exists** - Field `qr_code` is optional
2. **Use payload for QR generation** - `candidate.qr_code.payload`
3. **Test in browser console first** - Easier to debug
4. **Save token in localStorage** - For persistent sessions
5. **Handle errors gracefully** - Show message if QR not available

---

**Last Updated:** 16 Desember 2024  
**API Version:** 1.0.0  
**Status:** ✅ Production Ready