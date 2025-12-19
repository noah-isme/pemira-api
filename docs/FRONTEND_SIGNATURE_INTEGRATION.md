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
> **Note:** Format `signature` adalah string Base64 dari canvas/image. Backend akan mengupload ke Supabase Storage.

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

Pada halaman list DPT Admin, data tanda tangan tersedia sebagai **URL** (bukan base64).

### Endpoint
`GET /admin/elections/{election_id}/voters`

### Response Payload (Contoh Item)
```json
{
  "data": {
    "items": [
      {
        "voter_id": 1001,
        "nim": "12345678",
        "name": "Budi Santoso",
        "status": {
          "is_eligible": true,
          "has_voted": true,
          "voting_method": "ONLINE",
          "digital_signature_url": "https://xxx.supabase.co/storage/v1/object/public/pemira/signatures/1/1001.png"
        }
      }
    ]
  }
}
```

> **PENTING:** Field berubah dari `digital_signature` (base64) menjadi `digital_signature_url` (URL Supabase Storage).

### Implementasi Admin
- Cek field `status.digital_signature_url`
- Jika ada (tidak `null`), tampilkan tombol "Lihat Tanda Tangan"
- Render langsung sebagai `<img src={digital_signature_url} />` (bukan base64 decode)

