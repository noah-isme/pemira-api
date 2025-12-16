# Solusi: QR Code Tidak Muncul di Frontend

## 🚨 Masalah

Setelah membuat kandidat baru di admin panel, QR code tidak muncul di frontend.

**Contoh Kasus:**
- Kandidat ID: 77 (Ayu)
- Fakultas: Teknik Informatika
- Status: Terpublikasi
- **QR Code: Tidak muncul** ❌

---

## 🔍 Penyebab

QR code **TIDAK otomatis dibuat** saat kandidat dibuat atau diupdate. QR code hanya dibuat dalam kondisi berikut:

### Sebelumnya (❌):
- Create candidate → **Tidak ada QR**
- Update candidate → **Tidak ada QR**
- Publish candidate → **Tidak ada QR**

### Sekarang (✅):
- Create candidate → **Tidak ada QR** (normal)
- Publish candidate → **QR otomatis dibuat**
- Manual generate → **QR dibuat**

---

## ✅ Solusi

### Opsi 1: Auto-Generate saat Publish (Recommended)

QR code sekarang **otomatis dibuat** saat kandidat di-publish.

**Cara:**
1. Buka admin panel
2. Pilih kandidat
3. Klik tombol **"Publish"** atau **"Terbitkan"**
4. QR code otomatis dibuat ✅

**API:**
```bash
POST /api/v1/admin/elections/{electionID}/candidates/{candidateID}/publish
Authorization: Bearer {access_token}
```

### Opsi 2: Generate Manual (untuk kandidat yang sudah ada)

Untuk kandidat yang sudah terpublikasi tapi belum punya QR code:

**Endpoint Baru:**
```bash
POST /api/v1/admin/elections/{electionID}/candidates/{candidateID}/qr/generate
Authorization: Bearer {access_token}
```

**Contoh untuk kandidat Ayu (ID: 14, Election: 1):**
```bash
# Login dulu
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' | \
  jq -r '.access_token')

# Generate QR
curl -X POST "http://localhost:8080/api/v1/admin/elections/1/candidates/14/qr/generate" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

**Response:**
```json
{
  "id": 14,
  "name": "Ayu",
  "number": 77,
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

## 🚀 Bulk Generate untuk Semua Kandidat

Jika ada banyak kandidat yang belum punya QR code, gunakan script otomatis:

```bash
cd pemira-api
./scripts/generate-all-qr.sh 1
```

**Output:**
```
========================================
Generate QR Codes for All Candidates
========================================

Election ID: 1

✓ Login successful
✓ Found 9 candidate(s)

Generating QR Codes:
✓ Generated QR for #6 (Test Direct SQL): CAND06-u1lWt3-mUB7X
✓ Generated QR for #7 (Final Test): CAND07-CnPz-UTcNyMO
✓ Generated QR for #14 (Ayu): CAND14-NXLbfXqr7gtU
...

✓ All candidates now have QR codes!
```

---

## 🔍 Cara Cek QR Code

### 1. Via API
```bash
TOKEN="your-access-token"

# Get candidate detail
curl -s "http://localhost:8080/api/v1/admin/elections/1/candidates/14" \
  -H "Authorization: Bearer $TOKEN" | jq '.qr_code'
```

**Response jika ada QR:**
```json
{
  "id": 3,
  "token": "CAND14-NXLbfXqr7gtU",
  "payload": "PEMIRA-UNIWA|E:1|C:14|V:1",
  "version": 1,
  "is_active": true
}
```

**Response jika tidak ada QR:**
```json
null
```

### 2. Via Database
```bash
docker exec pemira-postgres psql -U pemira -d pemira -c \
  "SELECT c.id, c.name, c.number, qr.qr_token 
   FROM candidates c 
   LEFT JOIN candidate_qr_codes qr ON c.id = qr.candidate_id 
   WHERE c.id = 14;"
```

**Output:**
```
 id | name | number |      qr_token       
----+------+--------+---------------------
 14 | Ayu  |     77 | CAND14-NXLbfXqr7gtU
```

### 3. Lihat Semua Kandidat Tanpa QR
```sql
SELECT c.id, c.name, c.number, c.status
FROM candidates c
LEFT JOIN candidate_qr_codes qr ON c.id = qr.candidate_id AND qr.is_active = true
WHERE c.election_id = 1 AND qr.id IS NULL;
```

---

## 📱 Frontend Implementation

### Cek QR Code Exists

Frontend **HARUS** selalu cek apakah `qr_code` field ada:

```javascript
// ✅ BENAR - Cek dulu
if (candidate.qr_code && candidate.qr_code.payload) {
  // Render QR code
  <QRCode value={candidate.qr_code.payload} />
} else {
  // Show placeholder
  <p>QR Code tidak tersedia</p>
}
```

```jsx
// ✅ BENAR - React
{candidate.qr_code ? (
  <QRCode 
    value={candidate.qr_code.payload}
    size={256}
    level="H"
  />
) : (
  <div className="no-qr">
    QR Code tidak tersedia
  </div>
)}
```

```vue
<!-- ✅ BENAR - Vue.js -->
<div v-if="candidate.qr_code">
  <qrcode-vue :value="candidate.qr_code.payload" />
</div>
<div v-else>
  <p>QR Code tidak tersedia</p>
</div>
```

### ❌ SALAH - Langsung akses tanpa cek
```javascript
// ❌ SALAH - Bisa error jika qr_code null
<QRCode value={candidate.qr_code.payload} />
```

---

## 🔄 QR Code Rotation

Jika QR code perlu diganti (bocor/kompromi):

```bash
# Generate lagi - QR lama akan dinonaktifkan
POST /api/v1/admin/elections/1/candidates/14/qr/generate
```

**Proses:**
1. QR lama: `is_active = false`
2. Version bertambah: v1 → v2
3. Token baru dibuat
4. Response hanya return QR yang active

---

## 📊 Response Format

### Detail Kandidat dengan QR Code

**Endpoint:**
```
GET /api/v1/admin/elections/1/candidates/14
```

**Response:**
```json
{
  "id": 14,
  "election_id": 1,
  "number": 77,
  "name": "Ayu",
  "photo_url": "",
  "short_bio": "",
  "tagline": "",
  "faculty_name": "Teknik",
  "study_program_name": "Informatika",
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
```

**Key Points:**
- ✅ Response **langsung object**, TIDAK ada wrapper `data`
- ✅ Field `qr_code` di level root
- ✅ `qr_code` bisa `null` jika tidak ada
- ✅ Gunakan `qr_code.payload` untuk generate QR image

---

## 🛠️ Step-by-Step Fix

### Untuk Kandidat Ayu (ID: 14)

1. **Login ke server**
   ```bash
   docker exec pemira-postgres psql -U pemira -d pemira
   ```

2. **Cek apakah QR code ada**
   ```sql
   SELECT * FROM candidate_qr_codes WHERE candidate_id = 14;
   ```

3. **Jika tidak ada, generate via API**
   ```bash
   TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"password123"}' | \
     jq -r '.access_token')
   
   curl -X POST "http://localhost:8080/api/v1/admin/elections/1/candidates/14/qr/generate" \
     -H "Authorization: Bearer $TOKEN" | jq '.qr_code'
   ```

4. **Verify**
   ```bash
   curl -s "http://localhost:8080/api/v1/admin/elections/1/candidates/14" \
     -H "Authorization: Bearer $TOKEN" | jq '.qr_code'
   ```

5. **Refresh frontend** - QR code sekarang muncul ✅

---

## 📝 Checklist

- [ ] Server sudah running (`docker-compose up -d`)
- [ ] Login sebagai admin berhasil
- [ ] Generate QR code untuk kandidat yang belum punya
- [ ] Verify QR code muncul di API response
- [ ] Frontend cek `qr_code` exists sebelum render
- [ ] Test render QR code di frontend

---

## 🆘 Troubleshooting

### QR Code masih tidak muncul setelah generate

**Cek:**
1. Apakah generate berhasil?
   ```bash
   curl -X POST ".../qr/generate" -H "Authorization: Bearer $TOKEN"
   ```

2. Apakah ada di database?
   ```sql
   SELECT * FROM candidate_qr_codes WHERE candidate_id = 14;
   ```

3. Apakah API response ada `qr_code`?
   ```bash
   curl ".../candidates/14" -H "Authorization: Bearer $TOKEN" | jq '.qr_code'
   ```

4. Apakah frontend handle `null` dengan benar?
   ```javascript
   if (candidate.qr_code) { ... }
   ```

### Error saat generate

**Error: "CANDIDATE_NOT_FOUND"**
- Kandidat tidak ada atau sudah dihapus
- Cek: `SELECT * FROM candidates WHERE id = 14;`

**Error: "UNAUTHORIZED"**
- Token tidak valid atau expired
- Login ulang dan ambil token baru

**Error: "INVALID_REQUEST"**
- electionID atau candidateID salah
- Pastikan kandidat ada di election tersebut

---

## 📚 Dokumentasi Lengkap

- **API Response Format:** [docs/API_RESPONSE_EXAMPLES.md](docs/API_RESPONSE_EXAMPLES.md)
- **Frontend Integration:** [FRONTEND_QR_INTEGRATION.md](FRONTEND_QR_INTEGRATION.md)
- **QR Generation Guide:** [docs/QR_CODE_GENERATION_GUIDE.md](docs/QR_CODE_GENERATION_GUIDE.md)
- **Quick Start:** [QUICK_START.md](QUICK_START.md)

---

## 📞 Support

Jika masih ada masalah:
1. Cek logs: `docker-compose logs -f api`
2. Test endpoint: `./scripts/test-admin-qr.sh`
3. Generate QR: `./scripts/generate-all-qr.sh 1`
4. Baca dokumentasi di folder `/docs`

---

**Last Updated:** 16 Desember 2024  
**Status:** ✅ Fixed & Tested  
**Verified:** QR code generation working for all candidates