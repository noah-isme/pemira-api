# Integrasi Frontend Tanda Tangan Digital

Berikut adalah panduan integrasi untuk fitur tanda tangan digital pemilih.

## 1. Voter Side: Mengirim Tanda Tangan

Setelah pemilih sukses melakukan voting online (`POST /voting/online/cast`), frontend harus meminta pemilih untuk menggambar tanda tangan (misalnya menggunakan canvas).

### Endpoint
`POST /voting/online/signature`

### Headers
`Authorization: Bearer <token_pemilih>`

### Request Payload
```json
{
  "election_id": 105,
  "signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhORgAA..."
}
```
> **Note:** Format `signature` adalah string, disarankan menggunakan Base64 dari image/canvas.

### Response `200 OK`
```json
{
    "message": "Tanda tangan digital berhasil disimpan."
}
```

### Error Responses
- `400 Bad Request` ("VOTE_REQUIRED"): Pemilih belum melakukan voting.
- `409 Conflict` ("SIGNATURE_EXISTS"): Tanda tangan sudah pernah dikirim sebelumnya.

---

## 2. Admin Side: Menampilkan Tanda Tangan di DPT

Pada halaman list DPT Admin, data tanda tangan akan tersedia di dalam properti `status`.

### Endpoint
`GET /dpt/elections/{election_id}/voters`

### Response Payload (Contoh Item List)
```json
{
  "data": [
    {
      "voter_id": 1001,
      "nim": "12345678",
      "name": "Budi Santoso",
      "faculty_name": "Teknik",
      "study_program_name": "Informatika",
      "cohort_year": 2023,
      "email": "budi@student.univ.ac.id",
      "has_account": true,
      "voter_type": "STUDENT",
      "status": {
        "is_eligible": true,
        "has_voted": true,
        "last_vote_at": "2025-12-19T06:30:00Z",
        "voting_method": "ONLINE",
        "last_vote_channel": "ONLINE",
        "digital_signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhORgAA..." 
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 50,
    "total_pages": 5
  }
}
```

### Implementasi Admin
- Cek field `status.digital_signature`.
- Jika ada (tidak `null` atau `""`), tampilkan tombol "Lihat Tanda Tangan".
- Saat diklik, tampilkan gambar dari string Base64 tersebut.
