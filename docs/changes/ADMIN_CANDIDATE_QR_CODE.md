# Admin Candidate QR Code Implementation

## Tanggal: 2024
## Status: ✅ Selesai

## Deskripsi
Implementasi kode QR untuk detail kandidat di admin panel. Sebelumnya, endpoint admin untuk list dan detail kandidat tidak menampilkan data QR code meskipun field `qr_code` sudah tersedia di DTO.

## Perubahan

### File: `internal/candidate/service.go`

#### 1. Fungsi `AdminGetCandidate`
**Sebelum:**
```go
return &CandidateDetailDTO{
    ID:               c.ID,
    ElectionID:       c.ElectionID,
    // ... fields lainnya
    Status:           string(c.Status),
    Stats:            stats,
}, nil
```

**Sesudah:**
```go
dto := &CandidateDetailDTO{
    ID:               c.ID,
    ElectionID:       c.ElectionID,
    // ... fields lainnya
    Status:           string(c.Status),
    Stats:            stats,
}

// Get QR code if available
if qrCode, err := s.repo.GetActiveQRCode(ctx, c.ID); err == nil && qrCode != nil {
    dto.QRCode = &QRCodeDTO{
        ID:       qrCode.ID,
        Token:    qrCode.QRToken,
        URL:      fmt.Sprintf("https://pemira.local/ballot-qr/%s", qrCode.QRToken),
        Payload:  buildBallotQRPayload(qrCode.ElectionID, qrCode.CandidateID, qrCode.Version),
        Version:  qrCode.Version,
        IsActive: qrCode.IsActive,
    }
}

return dto, nil
```

#### 2. Fungsi `AdminListCandidates`
**Sebelum:**
```go
dtos := make([]CandidateDetailDTO, 0, len(candidates))
for _, c := range candidates {
    stats := statsMap[c.ID]
    dtos = append(dtos, CandidateDetailDTO{
        // ... fields
        Stats:            stats,
    })
}
```

**Sesudah:**
```go
// Get QR codes for all candidates in this election
qrCodesMap, err := s.repo.GetQRCodesByElection(ctx, electionID)
if err != nil {
    qrCodesMap = make(map[int64]*CandidateQRCode)
}

dtos := make([]CandidateDetailDTO, 0, len(candidates))
for _, c := range candidates {
    stats := statsMap[c.ID]
    dto := CandidateDetailDTO{
        // ... fields
        Stats:            stats,
    }

    // Add QR code if available
    if qrCode, exists := qrCodesMap[c.ID]; exists {
        dto.QRCode = &QRCodeDTO{
            ID:       qrCode.ID,
            Token:    qrCode.QRToken,
            URL:      fmt.Sprintf("https://pemira.local/ballot-qr/%s", qrCode.QRToken),
            Payload:  buildBallotQRPayload(qrCode.ElectionID, qrCode.CandidateID, qrCode.Version),
            Version:  qrCode.Version,
            IsActive: qrCode.IsActive,
        }
    }

    dtos = append(dtos, dto)
}
```

## Endpoint yang Terpengaruh

### 1. Get Candidate Detail (Admin)
```
GET /admin/elections/{electionID}/candidates/{candidateID}
```

**Response baru:**
```json
{
  "id": 1,
  "election_id": 1,
  "number": 1,
  "name": "Kandidat A",
  // ... fields lainnya
  "qr_code": {
    "id": 1,
    "token": "abc123token",
    "url": "https://pemira.local/ballot-qr/abc123token",
    "payload": "PEMIRA-UNIWA|E:1|C:1|V:1",
    "version": 1,
    "is_active": true
  }
}
```

### 2. List Candidates (Admin)
```
GET /admin/elections/{electionID}/candidates
```

**Response baru:** Setiap item dalam array `items` sekarang juga menyertakan field `qr_code`.

## Format QR Code

### Payload Format
```
PEMIRA-UNIWA|E:{election_id}|C:{candidate_id}|V:{version}
```

Contoh:
```
PEMIRA-UNIWA|E:1|C:5|V:1
```

### URL Format
```
https://pemira.local/ballot-qr/{token}
```

## Catatan Penting

1. **Nullable Field**: Field `qr_code` bersifat optional (`omitempty`). Jika kandidat belum memiliki QR code, field ini tidak akan muncul dalam response.

2. **Active QR Only**: Hanya QR code yang aktif (`is_active = true`) yang akan ditampilkan.

3. **Performance**: 
   - Untuk detail kandidat: menggunakan `GetActiveQRCode` (single query)
   - Untuk list kandidat: menggunakan `GetQRCodesByElection` (bulk query untuk efisiensi)

4. **Repository Methods**:
   - `GetActiveQRCode(ctx, candidateID)`: Mengambil QR code aktif untuk satu kandidat
   - `GetQRCodesByElection(ctx, electionID)`: Mengambil semua QR code aktif untuk election (return map)

## Testing

### Manual Testing
```bash
# Test detail kandidat
curl -X GET "http://localhost:8080/api/v1/admin/elections/1/candidates/1" \
  -H "Authorization: Bearer {admin_token}"

# Test list kandidat
curl -X GET "http://localhost:8080/api/v1/admin/elections/1/candidates" \
  -H "Authorization: Bearer {admin_token}"
```

## Dependencies
- Tabel: `candidate_qr_codes`
- Repository: `internal/candidate/repository.go` & `repository_pgx.go`
- Helper: `buildBallotQRPayload()` function

## Frontend Integration

Frontend dapat menggunakan library seperti `qrcode.react` atau `qrcode` untuk render QR code:

```jsx
import QRCode from 'qrcode.react';

function CandidateDetail({ candidate }) {
  return (
    <div>
      <h2>{candidate.name}</h2>
      {candidate.qr_code && (
        <div>
          <h3>QR Code untuk Voting</h3>
          <QRCode 
            value={candidate.qr_code.payload} 
            size={256}
            level="H"
          />
          <p>Version: {candidate.qr_code.version}</p>
        </div>
      )}
    </div>
  );
}
```

## Referensi
- `docs/CANDIDATE_ENDPOINTS_FRONTEND.md`
- `docs/api/CANDIDATE_QR_PAYLOAD.md`
- `migrations/011_ballot_qr_schema.up.sql`
