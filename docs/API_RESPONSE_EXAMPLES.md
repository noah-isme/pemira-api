# API Response Examples - Admin Candidate dengan QR Code

## Overview
Dokumentasi ini berisi contoh response **ACTUAL** dari server untuk membantu frontend developer mengintegrasikan fitur QR Code.

---

## 1. Admin Login

### Request
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

### Response (200 OK)
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU5NDY5MjgsImlhdCI6MTc2NTg2MDUyOCwicm9sZSI6IkFETUlOIiwic3ViIjoxLCJ0cHNfaWQiOjV9.-tfaaCWpPiGkh40lmJicxQxQG8BxHyiNEPwBfPFu9xQ",
  "refresh_token": "JxP3iF6TUyIApfRpgr1nnpNkf54vfh5MZSYyIsYVcmM=",
  "token_type": "Bearer",
  "expires_in": 86400,
  "user": {
    "id": 1,
    "username": "admin",
    "role": "ADMIN",
    "tps_id": 5,
    "profile": {}
  }
}
```

**Important:** Token ada di field `access_token`, bukan `token`!

---

## 2. Admin Get Candidate Detail dengan QR Code

### Request
```http
GET /api/v1/admin/elections/3/candidates/4
Authorization: Bearer {access_token}
```

### Response (200 OK)
```json
{
  "id": 4,
  "election_id": 3,
  "number": 1,
  "name": "Paslon Maju Bersama",
  "photo_url": "",
  "short_bio": "",
  "long_bio": "",
  "tagline": "",
  "faculty_name": "",
  "study_program_name": "",
  "vision": "Membangun kampus yang lebih baik untuk semua mahasiswa",
  "missions": [
    "Meningkatkan fasilitas kampus",
    "Memperkuat organisasi kemahasiswaan",
    "Mensejahterakan mahasiswa"
  ],
  "main_programs": [],
  "media": {},
  "social_links": [],
  "status": "APPROVED",
  "stats": {
    "total_votes": 1,
    "percentage": 100
  },
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

### Key Points untuk Frontend:
- ✅ **Response TIDAK dibungkus dalam `data`** - langsung object kandidat
- ✅ **Field `qr_code`** langsung di level root object
- ✅ **`qr_code.payload`** adalah string yang digunakan untuk generate QR code image
- ✅ **`qr_code`** bisa `null` jika kandidat tidak punya QR code

---

## 3. Admin List Candidates dengan QR Code

### Request
```http
GET /api/v1/admin/elections/3/candidates?page=1&limit=20
Authorization: Bearer {access_token}
```

### Response (200 OK)
```json
{
  "items": [
    {
      "id": 4,
      "election_id": 3,
      "number": 1,
      "name": "Paslon Maju Bersama",
      "photo_url": "",
      "short_bio": "",
      "long_bio": "",
      "tagline": "",
      "faculty_name": "",
      "study_program_name": "",
      "vision": "Membangun kampus yang lebih baik untuk semua mahasiswa",
      "missions": [
        "Meningkatkan fasilitas kampus",
        "Memperkuat organisasi kemahasiswaan",
        "Mensejahterakan mahasiswa"
      ],
      "main_programs": [],
      "media": {},
      "social_links": [],
      "status": "APPROVED",
      "stats": {
        "total_votes": 1,
        "percentage": 100
      },
      "qr_code": {
        "id": 1,
        "token": "CAND01-ABC123XYZ",
        "url": "https://pemira.local/ballot-qr/CAND01-ABC123XYZ",
        "payload": "PEMIRA-UNIWA|E:3|C:4|V:1",
        "version": 1,
        "is_active": true
      }
    },
    {
      "id": 5,
      "election_id": 3,
      "number": 2,
      "name": "Paslon Perubahan",
      "photo_url": "",
      "short_bio": "",
      "long_bio": "",
      "tagline": "",
      "faculty_name": "",
      "study_program_name": "",
      "vision": "Kampus progresif dan inklusif",
      "missions": [
        "Digitalisasi layanan kampus",
        "Program beasiswa untuk semua",
        "Aksi peduli lingkungan"
      ],
      "main_programs": [],
      "media": {},
      "social_links": [],
      "status": "APPROVED",
      "stats": {
        "total_votes": 0,
        "percentage": 0
      },
      "qr_code": {
        "id": 2,
        "token": "CAND02-DEF456UVW",
        "url": "https://pemira.local/ballot-qr/CAND02-DEF456UVW",
        "payload": "PEMIRA-UNIWA|E:3|C:5|V:1",
        "version": 1,
        "is_active": true
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_items": 2,
    "total_pages": 1
  }
}
```

### Key Points untuk Frontend:
- ✅ **Response ada `items` dan `pagination`** di root level
- ✅ **Setiap item dalam `items[]`** memiliki field `qr_code`
- ✅ **Tidak ada wrapper `data`** di response ini

---

## Frontend Integration Examples

### React/TypeScript Example

```typescript
// Types
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
  qr_code?: QRCode; // Optional!
}

// API Service
async function getCandidateDetail(
  electionId: number, 
  candidateId: number,
  token: string
): Promise<Candidate> {
  const response = await fetch(
    `http://localhost:8080/api/v1/admin/elections/${electionId}/candidates/${candidateId}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  if (!response.ok) {
    throw new Error('Failed to fetch candidate');
  }
  
  // Response langsung return object, TIDAK dibungkus data
  return await response.json();
}

// Component
import QRCode from 'qrcode.react';

function CandidateDetail({ electionId, candidateId, token }) {
  const [candidate, setCandidate] = useState<Candidate | null>(null);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    getCandidateDetail(electionId, candidateId, token)
      .then(data => {
        setCandidate(data);
        setLoading(false);
      })
      .catch(error => {
        console.error(error);
        setLoading(false);
      });
  }, [electionId, candidateId, token]);
  
  if (loading) return <div>Loading...</div>;
  if (!candidate) return <div>Candidate not found</div>;
  
  return (
    <div className="candidate-detail">
      <h2>{candidate.name}</h2>
      <p>Number: {candidate.number}</p>
      <p>Vision: {candidate.vision}</p>
      
      {/* QR Code Section - PENTING: cek qr_code exists */}
      {candidate.qr_code && (
        <div className="qr-code-section">
          <h3>QR Code for Voting</h3>
          <QRCode 
            value={candidate.qr_code.payload}
            size={256}
            level="H"
            includeMargin={true}
          />
          <div className="qr-info">
            <p>Token: {candidate.qr_code.token}</p>
            <p>Version: {candidate.qr_code.version}</p>
            <p>Status: {candidate.qr_code.is_active ? 'Active' : 'Inactive'}</p>
          </div>
        </div>
      )}
      
      {!candidate.qr_code && (
        <div className="no-qr">
          <p>QR Code belum tersedia untuk kandidat ini</p>
        </div>
      )}
    </div>
  );
}
```

### Vue.js Example

```vue
<template>
  <div class="candidate-detail">
    <div v-if="loading">Loading...</div>
    
    <div v-else-if="candidate">
      <h2>{{ candidate.name }}</h2>
      <p>Number: {{ candidate.number }}</p>
      <p>Vision: {{ candidate.vision }}</p>
      
      <!-- QR Code Section -->
      <div v-if="candidate.qr_code" class="qr-section">
        <h3>QR Code for Voting</h3>
        <qrcode-vue 
          :value="candidate.qr_code.payload"
          :size="256"
          level="H"
        />
        <div class="qr-info">
          <p>Token: {{ candidate.qr_code.token }}</p>
          <p>Version: {{ candidate.qr_code.version }}</p>
          <p>Active: {{ candidate.qr_code.is_active ? 'Yes' : 'No' }}</p>
        </div>
      </div>
      
      <div v-else class="no-qr">
        <p>QR Code belum tersedia</p>
      </div>
    </div>
  </div>
</template>

<script>
import QrcodeVue from 'qrcode.vue';

export default {
  components: {
    QrcodeVue
  },
  data() {
    return {
      candidate: null,
      loading: true
    };
  },
  mounted() {
    this.fetchCandidate();
  },
  methods: {
    async fetchCandidate() {
      try {
        const token = localStorage.getItem('access_token');
        const response = await fetch(
          `http://localhost:8080/api/v1/admin/elections/${this.electionId}/candidates/${this.candidateId}`,
          {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          }
        );
        
        // Response langsung object, tidak ada wrapper
        this.candidate = await response.json();
      } catch (error) {
        console.error('Failed to fetch candidate:', error);
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>
```

### Vanilla JavaScript Example

```javascript
// Fetch candidate detail
async function loadCandidateDetail(electionId, candidateId, accessToken) {
  try {
    const response = await fetch(
      `http://localhost:8080/api/v1/admin/elections/${electionId}/candidates/${candidateId}`,
      {
        headers: {
          'Authorization': `Bearer ${accessToken}`
        }
      }
    );
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    // Parse JSON - langsung return object
    const candidate = await response.json();
    
    // Check if QR code exists
    if (candidate.qr_code) {
      console.log('QR Code payload:', candidate.qr_code.payload);
      
      // Generate QR code using qrcode.js library
      const qrContainer = document.getElementById('qr-code-container');
      QRCode.toCanvas(qrContainer, candidate.qr_code.payload, {
        width: 256,
        errorCorrectionLevel: 'H'
      });
    } else {
      console.log('No QR code available for this candidate');
    }
    
    return candidate;
  } catch (error) {
    console.error('Error fetching candidate:', error);
    throw error;
  }
}

// Usage
const accessToken = localStorage.getItem('access_token');
loadCandidateDetail(3, 4, accessToken)
  .then(candidate => {
    console.log('Candidate loaded:', candidate);
  });
```

---

## Common Issues & Solutions

### Issue 1: "Cannot read property 'qr_code' of undefined"

**Penyebab:** Frontend mencoba akses `data.qr_code` padahal response langsung object.

**Salah:**
```javascript
const qrCode = response.data.qr_code; // ❌ SALAH
```

**Benar:**
```javascript
const qrCode = response.qr_code; // ✅ BENAR
```

### Issue 2: "QR Code tidak muncul padahal data ada"

**Penyebab:** Tidak cek apakah `qr_code` field exists.

**Benar:**
```javascript
if (candidate.qr_code && candidate.qr_code.payload) {
  // Generate QR code
  QRCode.toCanvas(canvas, candidate.qr_code.payload);
}
```

### Issue 3: "Token invalid"

**Penyebab:** Menggunakan field `token` padahal yang benar `access_token`.

**Login response:**
```javascript
const loginResponse = await response.json();
const token = loginResponse.access_token; // ✅ BENAR, bukan .token
```

### Issue 4: "CORS error"

**Solusi:** Pastikan backend sudah set CORS allowed origins yang benar.

---

## Testing dengan cURL

### Get Access Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' \
  | jq -r '.access_token'
```

### Get Candidate Detail
```bash
TOKEN="your-access-token-here"

curl -s "http://localhost:8080/api/v1/admin/elections/3/candidates/4" \
  -H "Authorization: Bearer $TOKEN" \
  | jq .
```

### Check if QR Code exists
```bash
TOKEN="your-access-token-here"

curl -s "http://localhost:8080/api/v1/admin/elections/3/candidates/4" \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.qr_code'
```

---

## QR Code Libraries Recommendations

### React
```bash
npm install qrcode.react
# or
npm install react-qr-code
```

### Vue.js
```bash
npm install qrcode.vue
```

### Vanilla JS
```bash
npm install qrcode
# or use CDN
<script src="https://cdn.jsdelivr.net/npm/qrcode@1.5.3/build/qrcode.min.js"></script>
```

---

## Summary

### ✅ Key Points untuk Frontend Developer:

1. **Response Detail Candidate:** Langsung object, TIDAK ada wrapper `data`
2. **Response List Candidates:** Ada `items[]` dan `pagination`, TIDAK ada wrapper `data`
3. **QR Code Field:** Optional (`qr_code?`), bisa `null` atau `undefined`
4. **QR Code Payload:** Gunakan `candidate.qr_code.payload` untuk generate QR image
5. **Token:** Ada di `access_token` bukan `token` pada login response
6. **Authorization Header:** `Bearer {access_token}`

### 📝 Checklist Integration:

- [ ] Parse response tanpa `data` wrapper
- [ ] Cek `qr_code` exists sebelum render
- [ ] Gunakan `qr_code.payload` untuk QR library
- [ ] Handle case ketika `qr_code` null
- [ ] Simpan `access_token` dari login
- [ ] Set header `Authorization: Bearer {token}`

---

**Last Updated:** 16 Desember 2024  
**Tested On:** Docker Alpine Linux  
**API Version:** 1.0.0