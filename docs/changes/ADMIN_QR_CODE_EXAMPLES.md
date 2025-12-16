# Admin Candidate QR Code - Example Responses

## Overview
Dokumentasi ini berisi contoh response dari endpoint admin kandidat setelah implementasi QR code.

## 1. Get Candidate Detail (Admin)

### Endpoint
```
GET /api/v1/admin/elections/{electionID}/candidates/{candidateID}
```

### Headers
```
Authorization: Bearer <admin_jwt_token>
```

### Response dengan QR Code (200 OK)
```json
{
  "id": 1,
  "election_id": 1,
  "number": 1,
  "name": "Ahmad Pratama & Siti Nurhaliza",
  "photo_url": "https://storage.example.com/candidates/1/profile.jpg",
  "photo_media_id": "profile_001",
  "short_bio": "Mahasiswa aktif yang peduli dengan kesejahteraan kampus",
  "long_bio": "Ahmad Pratama adalah mahasiswa semester 6...",
  "tagline": "Bersama Membangun Kampus yang Lebih Baik",
  "faculty_name": "Fakultas Teknik",
  "study_program_name": "Teknik Informatika",
  "cohort_year": 2021,
  "vision": "Mewujudkan kampus yang inklusif dan inovatif",
  "missions": [
    "Meningkatkan fasilitas kampus",
    "Memberdayakan organisasi mahasiswa",
    "Transparansi anggaran kampus"
  ],
  "main_programs": [
    {
      "title": "Program Beasiswa Digital",
      "description": "Memberikan beasiswa untuk mahasiswa berprestasi",
      "icon": "scholarship"
    },
    {
      "title": "Renovasi Perpustakaan",
      "description": "Modernisasi perpustakaan dengan teknologi terkini",
      "icon": "library"
    }
  ],
  "media": {
    "profile_photo": "https://storage.example.com/candidates/1/profile.jpg",
    "banner": "https://storage.example.com/candidates/1/banner.jpg",
    "video_url": "https://youtube.com/watch?v=example"
  },
  "media_files": [
    {
      "id": "media_001",
      "type": "image",
      "url": "https://storage.example.com/candidates/1/gallery1.jpg",
      "size": 2048576,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "social_links": [
    {
      "platform": "instagram",
      "url": "https://instagram.com/ahmad_pratama"
    },
    {
      "platform": "twitter",
      "url": "https://twitter.com/ahmad_pratama"
    }
  ],
  "status": "APPROVED",
  "stats": {
    "total_votes": 150,
    "percentage": 35.5,
    "rank": 1
  },
  "qr_code": {
    "id": 1,
    "token": "qr_e1c1_abc123def456ghi789jkl012mno345",
    "url": "https://pemira.local/ballot-qr/qr_e1c1_abc123def456ghi789jkl012mno345",
    "payload": "PEMIRA-UNIWA|E:1|C:1|V:1",
    "version": 1,
    "is_active": true
  }
}
```

### Response tanpa QR Code (200 OK)
Jika kandidat belum memiliki QR code yang aktif:
```json
{
  "id": 2,
  "election_id": 1,
  "number": 2,
  "name": "Budi Santoso & Ani Wijaya",
  "photo_url": "https://storage.example.com/candidates/2/profile.jpg",
  "short_bio": "Aktivis mahasiswa yang berpengalaman",
  "long_bio": "Budi Santoso telah aktif di berbagai organisasi...",
  "tagline": "Perubahan Dimulai dari Kita",
  "faculty_name": "Fakultas Ekonomi",
  "study_program_name": "Manajemen",
  "cohort_year": 2020,
  "vision": "Kampus yang demokratis dan transparan",
  "missions": [
    "Meningkatkan partisipasi mahasiswa",
    "Digitalisasi layanan kampus"
  ],
  "main_programs": [],
  "media": {},
  "media_files": [],
  "social_links": [],
  "status": "PENDING",
  "stats": {
    "total_votes": 0,
    "percentage": 0,
    "rank": 0
  }
}
```

## 2. List Candidates (Admin)

### Endpoint
```
GET /api/v1/admin/elections/{electionID}/candidates?page=1&limit=10&status=APPROVED
```

### Headers
```
Authorization: Bearer <admin_jwt_token>
```

### Query Parameters
- `page` (optional): Halaman (default: 1)
- `limit` (optional): Jumlah per halaman (default: 20)
- `status` (optional): Filter status (`PENDING`, `APPROVED`, `REJECTED`, `HIDDEN`)
- `search` (optional): Pencarian berdasarkan nama

### Response (200 OK)
```json
{
  "items": [
    {
      "id": 1,
      "election_id": 1,
      "number": 1,
      "name": "Ahmad Pratama & Siti Nurhaliza",
      "photo_url": "https://storage.example.com/candidates/1/profile.jpg",
      "photo_media_id": "profile_001",
      "short_bio": "Mahasiswa aktif yang peduli dengan kesejahteraan kampus",
      "long_bio": "Ahmad Pratama adalah mahasiswa semester 6...",
      "tagline": "Bersama Membangun Kampus yang Lebih Baik",
      "faculty_name": "Fakultas Teknik",
      "study_program_name": "Teknik Informatika",
      "cohort_year": 2021,
      "vision": "Mewujudkan kampus yang inklusif dan inovatif",
      "missions": [
        "Meningkatkan fasilitas kampus",
        "Memberdayakan organisasi mahasiswa"
      ],
      "main_programs": [
        {
          "title": "Program Beasiswa Digital",
          "description": "Memberikan beasiswa untuk mahasiswa berprestasi",
          "icon": "scholarship"
        }
      ],
      "media": {
        "profile_photo": "https://storage.example.com/candidates/1/profile.jpg"
      },
      "media_files": [],
      "social_links": [
        {
          "platform": "instagram",
          "url": "https://instagram.com/ahmad_pratama"
        }
      ],
      "status": "APPROVED",
      "stats": {
        "total_votes": 150,
        "percentage": 35.5,
        "rank": 1
      },
      "qr_code": {
        "id": 1,
        "token": "qr_e1c1_abc123def456ghi789jkl012mno345",
        "url": "https://pemira.local/ballot-qr/qr_e1c1_abc123def456ghi789jkl012mno345",
        "payload": "PEMIRA-UNIWA|E:1|C:1|V:1",
        "version": 1,
        "is_active": true
      }
    },
    {
      "id": 3,
      "election_id": 1,
      "number": 3,
      "name": "Citra Dewi & Dedi Kurniawan",
      "photo_url": "https://storage.example.com/candidates/3/profile.jpg",
      "short_bio": "Pemimpin masa depan untuk kampus yang lebih baik",
      "long_bio": "Citra Dewi memiliki visi untuk...",
      "tagline": "Aksi Nyata untuk Perubahan",
      "faculty_name": "Fakultas Hukum",
      "study_program_name": "Ilmu Hukum",
      "cohort_year": 2021,
      "vision": "Kampus yang adil dan merata",
      "missions": [
        "Perbaikan sistem informasi akademik",
        "Peningkatan soft skill mahasiswa"
      ],
      "main_programs": [],
      "media": {},
      "media_files": [],
      "social_links": [],
      "status": "APPROVED",
      "stats": {
        "total_votes": 120,
        "percentage": 28.4,
        "rank": 2
      },
      "qr_code": {
        "id": 3,
        "token": "qr_e1c3_xyz789uvw456rst123opq890lmn567",
        "url": "https://pemira.local/ballot-qr/qr_e1c3_xyz789uvw456rst123opq890lmn567",
        "payload": "PEMIRA-UNIWA|E:1|C:3|V:1",
        "version": 1,
        "is_active": true
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 2,
    "total_pages": 1
  }
}
```

## QR Code Payload Format

### Format
```
PEMIRA-UNIWA|E:{election_id}|C:{candidate_id}|V:{version}
```

### Penjelasan
- `PEMIRA-UNIWA`: Prefix identifikasi sistem
- `E:{election_id}`: ID pemilihan
- `C:{candidate_id}`: ID kandidat
- `V:{version}`: Versi QR code (untuk rotasi/invalidasi)

### Contoh
```
PEMIRA-UNIWA|E:1|C:1|V:1
PEMIRA-UNIWA|E:2|C:5|V:2
PEMIRA-UNIWA|E:1|C:10|V:1
```

## Frontend Implementation

### React dengan qrcode.react

```jsx
import React from 'react';
import QRCode from 'qrcode.react';

function CandidateQRCodeDisplay({ candidate }) {
  if (!candidate.qr_code) {
    return (
      <div className="qr-code-empty">
        <p>QR Code belum tersedia untuk kandidat ini</p>
      </div>
    );
  }

  return (
    <div className="qr-code-container">
      <h3>QR Code untuk Voting TPS</h3>
      
      {/* Render QR Code */}
      <QRCode 
        value={candidate.qr_code.payload}
        size={256}
        level="H" // High error correction
        includeMargin={true}
      />
      
      {/* Metadata */}
      <div className="qr-metadata">
        <p><strong>Token:</strong> {candidate.qr_code.token}</p>
        <p><strong>Version:</strong> {candidate.qr_code.version}</p>
        <p><strong>Status:</strong> {candidate.qr_code.is_active ? 'Active' : 'Inactive'}</p>
      </div>
      
      {/* Download button */}
      <button onClick={() => downloadQRCode(candidate)}>
        Download QR Code
      </button>
    </div>
  );
}

function downloadQRCode(candidate) {
  const canvas = document.querySelector('canvas');
  const pngUrl = canvas.toDataURL('image/png');
  const downloadLink = document.createElement('a');
  downloadLink.href = pngUrl;
  downloadLink.download = `qr-candidate-${candidate.number}.png`;
  document.body.appendChild(downloadLink);
  downloadLink.click();
  document.body.removeChild(downloadLink);
}

export default CandidateQRCodeDisplay;
```

### Vue.js dengan qrcode.vue

```vue
<template>
  <div class="qr-code-section">
    <div v-if="candidate.qr_code" class="qr-code-container">
      <h3>QR Code untuk Voting TPS</h3>
      
      <qrcode-vue 
        :value="candidate.qr_code.payload"
        :size="256"
        level="H"
        render-as="canvas"
      />
      
      <div class="qr-info">
        <p><strong>Version:</strong> {{ candidate.qr_code.version }}</p>
        <p><strong>Status:</strong> 
          <span :class="candidate.qr_code.is_active ? 'active' : 'inactive'">
            {{ candidate.qr_code.is_active ? 'Aktif' : 'Tidak Aktif' }}
          </span>
        </p>
      </div>
      
      <button @click="downloadQR" class="btn-download">
        Download QR Code
      </button>
    </div>
    
    <div v-else class="qr-code-empty">
      <p>QR Code belum tersedia</p>
    </div>
  </div>
</template>

<script>
import QrcodeVue from 'qrcode.vue';

export default {
  name: 'CandidateQRCode',
  components: {
    QrcodeVue
  },
  props: {
    candidate: {
      type: Object,
      required: true
    }
  },
  methods: {
    downloadQR() {
      const canvas = this.$el.querySelector('canvas');
      const url = canvas.toDataURL('image/png');
      const link = document.createElement('a');
      link.download = `candidate-${this.candidate.number}-qr.png`;
      link.href = url;
      link.click();
    }
  }
};
</script>

<style scoped>
.qr-code-container {
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  text-align: center;
}

.qr-info {
  margin-top: 15px;
}

.active {
  color: green;
  font-weight: bold;
}

.inactive {
  color: red;
  font-weight: bold;
}

.btn-download {
  margin-top: 15px;
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-download:hover {
  background-color: #0056b3;
}

.qr-code-empty {
  padding: 20px;
  text-align: center;
  color: #666;
}
</style>
```

## Use Cases

### 1. Admin Panel - Candidate Management
Admin dapat melihat dan mendownload QR code untuk setiap kandidat untuk:
- Pencetakan ballot QR untuk voting TPS
- Verifikasi QR code kandidat
- Monitoring status QR code (active/inactive)

### 2. TPS Voting Setup
QR code ini digunakan untuk:
- Pemilih scan QR code kandidat di TPS
- Sistem validasi kandidat berdasarkan payload QR
- Pencatatan vote dengan referensi QR code

### 3. QR Code Rotation
Jika diperlukan rotasi QR code untuk keamanan:
- Version number akan bertambah
- QR lama di-set `is_active = false`
- QR baru dibuat dengan version baru

## Error Responses

### Candidate Not Found (404)
```json
{
  "error": {
    "code": "CANDIDATE_NOT_FOUND",
    "message": "Kandidat tidak ditemukan"
  }
}
```

### Unauthorized (401)
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token tidak valid atau sudah kadaluarsa"
  }
}
```

### Forbidden (403)
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke resource ini"
  }
}
```

## Notes

1. **QR Code adalah Optional**: Field `qr_code` hanya muncul jika kandidat memiliki QR code aktif
2. **Performance**: List endpoint menggunakan bulk query untuk efisiensi
3. **Security**: QR token bersifat unik dan tidak dapat diprediksi
4. **Version Control**: Setiap QR code memiliki version untuk tracking dan rotasi