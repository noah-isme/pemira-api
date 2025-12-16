# Quick Start Guide - Pemira API

## 🚀 Langkah Cepat (5 Menit)

### 1. Clone & Navigate
```bash
cd pemira-api
```

### 2. Start dengan Docker
```bash
docker-compose up -d --build
```

### 3. Verifikasi Server Berjalan
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"data":{"status":"ok"}}
```

### 4. Test API
```bash
# List elections
curl http://localhost:8080/api/v1/elections | jq .

# Login sebagai admin
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' | jq .
```

## ✅ Fitur Baru: QR Code di Admin Panel

### Test QR Code Implementation
```bash
# Automated test
ADMIN_PASSWORD="password123" ./scripts/test-admin-qr.sh
```

### Manual Test
```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' | \
  grep -o '"access_token":"[^"]*"' | sed 's/"access_token":"\(.*\)"/\1/')

# 2. List kandidat dengan QR code
curl -s http://localhost:8080/api/v1/admin/elections/1/candidates \
  -H "Authorization: Bearer $TOKEN" | jq '.data.items[0].qr_code'

# 3. Detail kandidat dengan QR code
curl -s http://localhost:8080/api/v1/admin/elections/1/candidates/1 \
  -H "Authorization: Bearer $TOKEN" | jq '.qr_code'
```

## 📦 Yang Sudah Berjalan

Setelah `docker-compose up -d`, Anda akan memiliki:

| Service | Port | Status | Image |
|---------|------|--------|-------|
| API Server | 8080 | ✅ Running | Alpine Linux |
| PostgreSQL | 5432 | ✅ Running | postgres:16-alpine |
| Redis | 6379 | ✅ Running | redis:7-alpine |

## 🔑 Default Credentials

### Admin
- **Username:** `admin`
- **Password:** `password123`

### Panitia
- **Username:** `panitia`
- **Password:** `password123`

### Operator TPS
- **Username:** `operator`
- **Password:** `password123`

## 📡 Endpoint Utama

### Public Endpoints
```bash
GET  /health                                          # Health check
GET  /api/v1/elections                                # List elections
GET  /api/v1/elections/{id}/candidates                # List candidates
GET  /api/v1/elections/{id}/candidates/{candidateId}  # Candidate detail
GET  /api/v1/elections/{id}/qr-codes                  # Candidates with QR
```

### Admin Endpoints (Requires Auth)
```bash
POST /api/v1/auth/login                                      # Login
GET  /api/v1/admin/elections/{id}/candidates                 # List (with QR)
GET  /api/v1/admin/elections/{id}/candidates/{candidateId}   # Detail (with QR)
POST /api/v1/admin/elections/{id}/candidates                 # Create
PUT  /api/v1/admin/elections/{id}/candidates/{candidateId}   # Update
DELETE /api/v1/admin/elections/{id}/candidates/{candidateId} # Delete
```

## 🛠️ Useful Commands

### Docker
```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# Restart API only
docker-compose restart api

# View logs
docker-compose logs -f api

# Check status
docker-compose ps
```

### Makefile
```bash
make docker-up      # Start docker services
make docker-down    # Stop docker services
make build          # Build Go binary
make test           # Run tests
```

### Database
```bash
# Connect to PostgreSQL
docker exec -it pemira-postgres psql -U pemira -d pemira

# Connect to Redis
docker exec -it pemira-redis redis-cli
```

## 📊 Response Format QR Code

```json
{
  "id": 1,
  "election_id": 1,
  "number": 1,
  "name": "Kandidat A",
  "status": "APPROVED",
  "qr_code": {
    "id": 1,
    "token": "CAND01-ABC123XYZ",
    "url": "https://pemira.local/ballot-qr/CAND01-ABC123XYZ",
    "payload": "PEMIRA-UNIWA|E:1|C:1|V:1",
    "version": 1,
    "is_active": true
  }
}
```

## 🔍 Troubleshooting

### Port Already in Use
```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port in docker-compose.yml
```

### Container Won't Start
```bash
# Check logs
docker-compose logs api

# Rebuild from scratch
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

### Database Connection Failed
```bash
# Check PostgreSQL is healthy
docker-compose ps postgres

# Check database logs
docker-compose logs postgres
```

## 📚 Documentation

- **Implementation Guide:** [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)
- **Docker Guide:** [docs/RUNNING_WITH_DOCKER.md](docs/RUNNING_WITH_DOCKER.md)
- **QR Code Implementation:** [docs/changes/ADMIN_CANDIDATE_QR_CODE.md](docs/changes/ADMIN_CANDIDATE_QR_CODE.md)
- **API Examples:** [docs/changes/ADMIN_QR_CODE_EXAMPLES.md](docs/changes/ADMIN_QR_CODE_EXAMPLES.md)

## 🎯 What's New

### QR Code Implementation (Latest)
✅ Admin panel sekarang menampilkan QR code untuk setiap kandidat  
✅ Endpoint list dan detail sudah terintegrasi  
✅ Format QR: `PEMIRA-UNIWA|E:{election_id}|C:{candidate_id}|V:{version}`  
✅ Support untuk QR code rotation dengan version control  
✅ Automated testing script tersedia  

## 🚦 Health Check

Server sehat jika:
```bash
curl http://localhost:8080/health
# Returns: {"data":{"status":"ok"}}
```

Database sehat jika:
```bash
docker exec pemira-postgres pg_isready -U pemira
# Returns: accepting connections
```

Redis sehat jika:
```bash
docker exec pemira-redis redis-cli ping
# Returns: PONG
```

## 💡 Tips

1. **Gunakan jq untuk format JSON:**
   ```bash
   curl http://localhost:8080/api/v1/elections | jq .
   ```

2. **Save token ke variable:**
   ```bash
   export TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"password123"}' | \
     grep -o '"access_token":"[^"]*"' | sed 's/"access_token":"\(.*\)"/\1/')
   ```

3. **Test semua endpoint sekaligus:**
   ```bash
   ./scripts/test-admin-qr.sh
   ```

## 📞 Support

Jika ada masalah:
1. Cek logs: `docker-compose logs -f api`
2. Cek health: `curl http://localhost:8080/health`
3. Rebuild: `docker-compose up -d --build`
4. Baca dokumentasi di folder `/docs`

---

**Last Updated:** 16 Desember 2024  
**Version:** 1.0.0  
**Status:** ✅ Production Ready