# API Payload: Kandidat dengan QR Code

## Overview
Dokumentasi ini menjelaskan struktur payload JSON yang dibutuhkan frontend untuk menampilkan daftar kandidat beserta QR code mereka dalam sistem pemilu TPS.

## Endpoint

### Get Candidates with QR Codes
**Method:** `GET`  
**Path:** `/api/v1/elections/{election_id}/qr-codes`  
**Authentication:** Optional (public endpoint)

> Catatan: endpoint `GET /api/v1/elections/{election_id}/candidates` juga dapat mengembalikan `qr_code` per kandidat,
> tetapi format responsnya berbeda (menggunakan wrapper `items`).

## Response Payload Structure

```json
{
  "election_id": 3,
  "election_name": "Simulasi Pemilu TPS",
  "candidates": [
    {
      "id": 4,
      "number": 1,
      "name": "Paslon Maju Bersama",
      "vision": "Membangun kampus yang lebih baik untuk semua mahasiswa",
      "missions": [
        "Meningkatkan fasilitas kampus",
        "Memperkuat organisasi kemahasiswaan",
        "Mensejahterakan mahasiswa"
      ],
      "photo_url": "https://storage.supabase.co/...",
      "short_bio": "Kami adalah pasangan calon yang berpengalaman...",
      "tagline": "Maju Bersama untuk Kampus Lebih Baik",
      "faculty_name": "Fakultas Teknik",
      "study_program_name": "S1 Teknik Informatika",
      "status": "APPROVED",
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
      "number": 2,
      "name": "Paslon Perubahan",
      "vision": "Kampus progresif dan inklusif",
      "missions": [
        "Digitalisasi layanan kampus",
        "Program beasiswa untuk semua",
        "Aksi peduli lingkungan"
      ],
      "photo_url": null,
      "short_bio": "",
      "tagline": "Perubahan Dimulai dari Kita",
      "faculty_name": "Fakultas Ekonomi",
      "study_program_name": "S1 Manajemen",
      "status": "APPROVED",
      "qr_code": {
        "id": 2,
        "token": "CAND02-DEF456UVW",
        "url": "https://pemira.local/ballot-qr/CAND02-DEF456UVW",
        "payload": "PEMIRA-UNIWA|E:3|C:5|V:1",
        "version": 1,
        "is_active": true
      }
    }
  ]
}
```

## Field Descriptions

### Root Level
| Field | Type | Description |
|-------|------|-------------|
| `election_id` | integer | ID pemilu |
| `election_name` | string | Nama pemilu |
| `candidates` | array | Array of candidate objects |

### Candidate Object
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | integer | Yes | Unique ID kandidat |
| `number` | integer | Yes | Nomor urut kandidat (1, 2, 3, ...) |
| `name` | string | Yes | Nama kandidat/paslon |
| `vision` | string | No | Visi kandidat |
| `missions` | array[string] | No | Array misi kandidat |
| `photo_url` | string | No | URL foto kandidat (nullable) |
| `short_bio` | string | No | Bio singkat |
| `tagline` | string | No | Tagline/slogan |
| `faculty_name` | string | No | Nama fakultas |
| `study_program_name` | string | No | Nama program studi |
| `status` | enum | Yes | Status kandidat: `APPROVED`, `PENDING`, `REJECTED`, `WITHDRAWN` |
| `qr_code` | object | Yes | QR code object |

### QR Code Object
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | integer | Yes | Unique ID QR code |
| `token` | string | Yes | QR code token (untuk scanning) |
| `url` | string | Yes | Full URL untuk generate QR image |
| `payload` | string | Yes | String yang di-encode ke QR (dipakai di `/voting/tps/ballots/*`) |
| `version` | integer | Yes | Versi QR code (untuk rotation) |
| `is_active` | boolean | Yes | Status aktif QR code |

## Frontend Implementation Guide

### 1. Fetch Kandidat dengan QR Code

```javascript
// React/Next.js Example
const fetchCandidates = async (electionId) => {
  const response = await fetch(`/api/v1/elections/${electionId}/candidates`);
  const data = await response.json();
  return data;
};
```

### 2. Generate QR Code Image

Gunakan library QR code generator seperti `qrcode.react` atau `react-qr-code`:

```jsx
import QRCode from 'react-qr-code';

function CandidateQRCard({ candidate }) {
  return (
    <div className="candidate-card">
      <h2>#{candidate.number} - {candidate.name}</h2>
      <p>{candidate.tagline}</p>
      
      {/* QR Code */}
      <div className="qr-code-container">
        <QRCode 
          value={candidate.qr_code.payload} 
          size={256}
          level="H"
        />
        <p className="qr-label">Scan untuk memilih</p>
      </div>
      
      <div className="candidate-info">
        <h3>Visi</h3>
        <p>{candidate.vision}</p>
        
        <h3>Misi</h3>
        <ul>
          {candidate.missions.map((mission, idx) => (
            <li key={idx}>{mission}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}
```

### 3. Display Kandidat List

```jsx
function CandidateList({ electionId }) {
  const [candidates, setCandidates] = useState([]);
  
  useEffect(() => {
    fetchCandidates(electionId).then(data => {
      setCandidates(data.candidates);
    });
  }, [electionId]);
  
  return (
    <div className="candidates-grid">
      {candidates.map(candidate => (
        <CandidateQRCard 
          key={candidate.id} 
          candidate={candidate} 
        />
      ))}
    </div>
  );
}
```

### 4. Print Mode untuk Ballot Paper

```jsx
function PrintableBallot({ candidates }) {
  return (
    <div className="printable-ballot">
      <h1>SURAT SUARA PEMILU</h1>
      <p>Pilih salah satu kandidat dengan cara scan QR Code</p>
      
      <div className="ballot-grid">
        {candidates.map(candidate => (
          <div key={candidate.id} className="ballot-item">
            <div className="number-box">
              {candidate.number}
            </div>
            <h3>{candidate.name}</h3>
            <QRCode 
              value={candidate.qr_code.payload}
              size={200}
              level="H"
            />
            <p className="qr-token">{candidate.qr_code.token}</p>
          </div>
        ))}
      </div>
      
      <style jsx>{`
        @media print {
          .ballot-item {
            page-break-inside: avoid;
            border: 2px solid #000;
            padding: 20px;
            margin: 10px;
          }
        }
      `}</style>
    </div>
  );
}
```

## QR Code Scanning Flow

### 1. Voter Side (Mobile App/Scanner)
```
1. Voter masuk ke TPS
2. Panitia scan QR voter untuk check-in
3. Voter menerima surat suara fisik
4. Voter scan QR kandidat pilihan
5. Sistem validasi dan record vote
```

### 2. QR Code Content
QR code berisi `token` string yang unik untuk setiap kandidat:
- Format: `CAND{number}-{random_hash}`
- Example: `CAND01-ABC123XYZ`

### 3. Validation
Saat QR di-scan, backend akan:
1. Validasi token exists dan active
2. Cek voter sudah check-in atau belum
3. Cek voter belum voting atau belum
4. Record vote jika semua validasi passed

## Security Considerations

### 1. QR Code Token
- Token harus **unique** dan **random**
- Token bisa di-rotate jika ada kebocoran
- `version` field untuk track rotation history

### 2. Active Status
- Hanya QR dengan `is_active: true` yang valid
- QR lama bisa di-deactivate saat rotation

### 3. One-Time Use
- Setiap QR code vote hanya bisa digunakan **satu kali** per voter
- Backend harus validasi `voter_status.has_voted`

## Example SQL Query

Untuk mendapatkan payload ini dari database:

```sql
SELECT json_build_object(
    'election_id', e.id,
    'election_name', e.name,
    'candidates', json_agg(
        json_build_object(
            'id', c.id,
            'number', c.number,
            'name', c.name,
            'vision', c.vision,
            'missions', c.missions,
            'photo_url', c.photo_url,
            'short_bio', c.short_bio,
            'tagline', c.tagline,
            'faculty_name', c.faculty_name,
            'study_program_name', c.study_program_name,
            'status', c.status,
            'qr_code', json_build_object(
                'id', qr.id,
                'token', qr.qr_token,
                'url', 'https://pemira.local/ballot-qr/' || qr.qr_token,
                'version', qr.version,
                'is_active', qr.is_active
            )
        ) ORDER BY c.number
    )
) as payload
FROM elections e
JOIN candidates c ON c.election_id = e.id
LEFT JOIN candidate_qr_codes qr ON qr.candidate_id = c.id AND qr.is_active = true
WHERE e.id = $1 AND c.status = 'APPROVED'
GROUP BY e.id, e.name;
```

## Testing

### Sample cURL Request
```bash
curl -X GET "http://localhost:8080/api/v1/elections/3/candidates" \
  -H "Content-Type: application/json" \
  | jq .
```

### Expected Response Status
- `200 OK` - Success with candidates data
- `404 Not Found` - Election not found
- `500 Internal Server Error` - Server error

## Related Documentation
- [Voting API Implementation](../VOTING_API_IMPLEMENTATION.md)
- [TPS Management](../ADMIN_TPS_API.md)
- [QR Code System Architecture](./QR_CODE_ARCHITECTURE.md)

---

**Last Updated:** 2025-12-11  
**Version:** 1.0  
**Maintainer:** PEMIRA API Team
