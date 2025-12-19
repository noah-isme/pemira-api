# Panduan Pengujian Fitur Tanda Tangan Digital

Dokumen ini berisi langkah-langkah untuk Frontend Engineer (FE) dalam menguji fitur tanda tangan digital pemilih pada sistem PEMIRA.

## Prasyarat
- Server API berjalan di `http://localhost:8080` (atau sesuai konfigurasi).
- Akun Mahasiswa (untuk voting).
- Akun Admin (untuk verifikasi DPT).

---

## Skenario 1: Alur Pemilih (Student)

### 1. Login sebagai Mahasiswa
Gunakan endpoint login untuk mendapatkan `access_token`.
- **Endpoint**: `POST /auth/login`
- **Body**:
  ```json
  {
    "identity": "nim_mahasiswa",
    "password": "password",
    "role": "STUDENT"
  }
  ```

### 2. Cek Status Pemilih (Opsional tapi disarankan)
Pastikan mahasiswa terdaftar di pemilu.
- **Endpoint**: `GET /voting/config` (atau endpoint terkait status voting)

### 3. Test Case Negatif: Kirim Tanda Tangan Sebelum Memilih
Coba kirim tanda tangan sebelum melakukan voting. Harusnya **GAGAL**.
- **Endpoint**: `POST /voting/online/signature`
- **Header**: `Authorization: Bearer <token_mahasiswa>`
- **Body**:
  ```json
  {
    "election_id": 1,
    "signature": "data:image/png;base64,sample_signature_before_vote"
  }
  ```
- **Ekspektasi Response**: `400 Bad Request`
- **Message**: "Anda harus melakukan pemilihan terlebih dahulu."

### 4. Melakukan Voting Online
Lakukan voting terlebih dahulu.
- **Endpoint**: `POST /voting/online/cast`
- **Body**:
  ```json
  {
    "election_id": 1,
    "candidate_id": 1
  }
  ```
- **Ekspektasi Response**: `200 OK`

### 5. Kirim Tanda Tangan (Success Case)
Setelah voting berhasil, kirim tanda tangan.
- **Endpoint**: `POST /voting/online/signature`
- **Body**:
  ```json
  {
    "election_id": 1,
    "signature": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
  }
  ```
- **Ekspektasi Response**: `200 OK`
- **Message**: "Tanda tangan digital berhasil disimpan."

### 6. Test Case Negatif: Kirim Tanda Tangan Ganda
Coba kirim tanda tangan lagi dengan akun yang sama. Harusnya **GAGAL**.
- **Ekspektasi Response**: `409 Conflict`
- **Message**: "Tanda tangan digital sudah ada."

---

## Skenario 2: Verifikasi Admin

### 1. Login sebagai Admin
- **Endpoint**: `POST /auth/login`
- **Body**:
  ```json
  {
    "identity": "admin",
    "password": "password",
    "role": "ADMIN"
  }
  ```

### 2. Cek List DPT
Lihat data pemilih yang baru saja mengirim tanda tangan.
- **Endpoint**: `GET /dpt/elections/{election_id}/voters?search=nim_mahasiswa`
- **Header**: `Authorization: Bearer <token_admin>`
- **Ekspektasi Data**:
  Cek objek voter pada response. Field `status.digital_signature` harus berisi string base64 yang dikirim diatas.
  ```json
  {
    "data": [
      {
        "nim": "nim_mahasiswa",
        "status": {
          "has_voted": true,
          "voting_method": "ONLINE",
          "digital_signature": "data:image/png;base64,iVBORw0KGgoAAA..."
        }
      }
    ]
  }
  ```
