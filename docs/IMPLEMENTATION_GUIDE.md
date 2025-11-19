# PEMIRA API - Implementation Guide

## 🎯 Project Status

✅ **Architecture**: Secure, modular monolith designed for elections
✅ **Modules**: 11 core modules implemented
✅ **Build**: Passing
✅ **Documentation**: Complete

---

## 📁 Project Structure

```
pemira-api/
├── cmd/
│   ├── api/              # HTTP API server with WebSocket
│   └── worker/           # Background jobs (TODO)
├── internal/
│   ├── auth/             # ✓ JWT, login, session, RBAC
│   ├── election/         # ✓ Elections, phases, voting mode
│   ├── voter/            # ✓ DPT management, voter status
│   ├── candidate/        # ✓ Candidates, visi-misi, media
│   ├── tps/              # ✓ TPS, QR codes, check-ins
│   ├── voting/           # ✓ CRITICAL: Transaction-based voting
│   ├── monitoring/       # ✓ Live count, statistics, dashboard
│   ├── announcement/     # ✓ Pengumuman untuk mahasiswa
│   ├── audit/            # ✓ Audit logs untuk semua operasi sensitif
│   ├── ws/               # ✓ WebSocket hub untuk real-time
│   ├── http/
│   │   ├── middleware/   # ✓ Auth, RBAC, rate limit, logger
│   │   └── response/     # ✓ Standard JSON responses
│   ├── shared/           # ✓ Constants, errors, pagination, context keys
│   └── fileimport/       # TODO: CSV/XLSX parser
├── docs/
│   ├── ARCHITECTURE.md   # ✓ Arsitektur lengkap
│   ├── API.md            # ✓ API documentation
│   └── IMPLEMENTATION_GUIDE.md  # This file
├── migrations/           # TODO: SQL schema migrations
├── pkg/database/         # ✓ PostgreSQL pool utilities
└── docker-compose.yml    # ✓ Postgres + Redis
```

---

## 🔐 Security Features

### 1. Voting Security (MISSION CRITICAL)
```go
// internal/voting/transaction.go
CastVoteWithTransaction()
```

**Prevents double voting via:**
- ✅ Transaction isolation (READ COMMITTED)
- ✅ Row-level lock (`FOR UPDATE`)
- ✅ Atomic operations (no race conditions)
- ✅ Vote token generation (anonymous receipt)
- ✅ Audit logging

**Flow:**
1. Lock voter status row
2. Check: not voted, eligible, valid phase
3. Generate anonymous token
4. Insert vote + token
5. Update voter status
6. Update statistics
7. Audit log
8. Commit transaction

### 2. Authentication & Authorization
- JWT with short expiration
- Role-based access control (RBAC)
- Middleware: Auth, RequireRole, RequireAdmin
- Password hashing: bcrypt

### 3. Rate Limiting
```go
// internal/http/middleware/ratelimit.go
```
- Login: 5 req/min
- Voting: 3 req/min
- Others: 60 req/min

### 4. Audit Logging
All sensitive operations logged:
- Vote cast
- Election config changes
- TPS QR regeneration
- Candidate updates
- DPT imports

---

## 🚀 Next Steps (Implementation Checklist)

### Phase 1: Database Schema ⏳
```bash
# Create migrations
make migrate-create name=create_base_schema
```

**Tables needed:**
1. ✓ users (auth)
2. ✓ elections
3. ✓ election_phases
4. ✓ voters
5. ✓ voter_election_status
6. ✓ candidates
7. ✓ candidate_members
8. ✓ tps
9. ✓ tps_operators
10. ✓ tps_checkins
11. ✓ votes (CRITICAL: anonymous)
12. ✓ vote_tokens
13. ✓ vote_stats (materialized view)
14. ✓ announcements
15. ✓ audit_logs

**Constraints:**
- UNIQUE(voter_id, election_id) on voter_election_status
- UNIQUE(token_hash) on votes
- CHECK(has_voted IN (true, false))

### Phase 2: Repository Implementation ⏳
Implement all `Repository` interfaces using pgx:
- [ ] auth.Repository
- [ ] election.Repository
- [ ] voter.Repository
- [ ] candidate.Repository
- [ ] tps.Repository
- [ ] voting.Repository (use transaction.go)
- [ ] monitoring.Repository
- [ ] announcement.Repository
- [ ] audit.Repository

**Option: Use sqlc for type-safe queries**
```bash
make sqlc-generate
```

### Phase 3: Wire Dependencies in main.go ⏳
```go
// cmd/api/main.go
func main() {
    // Init DB
    pool := database.NewPostgresPool(...)
    
    // Init repos
    authRepo := auth.NewPostgresRepository(pool)
    voterRepo := voter.NewPostgresRepository(pool)
    // ... other repos
    
    // Init services
    authSvc := auth.NewService(authRepo, cfg.JWTSecret, 24*time.Hour)
    votingSvc := voting.NewService(votingRepo)
    // ... other services
    
    // Init handlers
    authHandler := auth.NewHandler(authSvc)
    votingHandler := voting.NewHandler(votingSvc)
    // ... other handlers
    
    // Register routes with middleware
    r := chi.NewRouter()
    r.Use(middleware.RequestLogger)
    r.Use(middleware.RealIP)
    r.Use(middleware.Recoverer)
    
    // Public routes
    authHandler.RegisterRoutes(r)
    electionHandler.RegisterRoutes(r)
    candidateHandler.RegisterRoutes(r)
    
    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(authMiddleware.Auth(authSvc))
        
        // Student routes
        voterHandler.RegisterRoutes(r)
        votingHandler.RegisterRoutes(r)
        
        // Admin routes
        r.Group(func(r chi.Router) {
            r.Use(middleware.RequireAdmin())
            // admin handlers...
        })
    })
    
    // WebSocket
    wsHandler.RegisterRoutes(r)
    
    // Start server
    srv.ListenAndServe()
}
```

### Phase 4: Testing 🧪
```bash
# Unit tests
go test ./internal/voting -v
go test ./internal/auth -v

# Integration tests (with test DB)
go test ./... -tags=integration
```

### Phase 5: Deployment 🚀
```bash
# Build
make build

# Docker
docker build -t pemira-api .
docker-compose up -d

# Run migrations
make migrate-up DATABASE_URL=$DATABASE_URL

# Start API
./bin/api
```

---

## 📊 API Endpoints Summary

| Category | Method | Endpoint | Auth | Role |
|----------|--------|----------|------|------|
| **Auth** | POST | /auth/login/student | - | - |
| | POST | /auth/login/admin | - | - |
| | GET | /auth/me | ✓ | All |
| **Election** | GET | /elections/current | - | - |
| | GET | /elections/{id}/candidates | - | - |
| **Voting** | POST | /voting/online/cast | ✓ | STUDENT |
| | POST | /voting/tps/cast | ✓ | STUDENT |
| **TPS** | POST | /tps/checkin/scan | ✓ | STUDENT |
| | GET | /tps/{id}/checkins | ✓ | TPS_OPERATOR |
| | POST | /tps/{id}/checkins/{id}/approve | ✓ | TPS_OPERATOR |
| **Monitoring** | GET | /admin/monitoring/summary | ✓ | ADMIN |
| | GET | /admin/monitoring/live-count/{id} | ✓ | ADMIN |
| **Admin** | GET | /admin/dpt | ✓ | ADMIN |
| | POST | /admin/dpt/import | ✓ | ADMIN |
| | GET | /admin/audit-logs | ✓ | SUPER_ADMIN |
| **WebSocket** | WS | /ws/tps/{id} | ✓ | TPS_OPERATOR |
| | WS | /ws/live-count/{id} | ✓ | ADMIN |

---

## 🔧 Configuration

### Environment Variables (.env)
```bash
# App
APP_ENV=development
HTTP_PORT=8080

# Database
DATABASE_URL=postgres://pemira:pemira@localhost:5432/pemira?sslmode=disable

# Auth
JWT_SECRET=your-secret-change-in-production
JWT_EXPIRATION=24h

# Redis (optional)
REDIS_URL=redis://localhost:6379/0

# Logging
LOG_LEVEL=info
```

### Feature Flags (TODO)
```go
type Config struct {
    // ... existing fields
    
    // Feature flags
    EnableLiveCount      bool
    EnableTPSMode        bool
    EnableOnlineMode     bool
    EnablePublicResults  bool
}
```

---

## 📈 Metrics & Observability

### Prometheus Metrics
```
# Votes
votes_total{candidate_id, voted_via}
votes_per_minute
voting_errors_total{error_type}

# TPS
tps_checkins_total{tps_id, status}
tps_queue_length{tps_id}

# HTTP
http_requests_total{method, path, status}
http_request_duration_seconds{method, path}
```

### Logging
```go
slog.Info("vote_cast", 
    "election_id", electionID,
    "voter_id", voterID,
    "voted_via", votedVia,
    "duration_ms", duration)
```

---

## 🧪 Testing Strategy

### Unit Tests
- Service logic
- Middleware
- Validators

### Integration Tests
- Repository with test DB
- HTTP endpoints
- WebSocket handlers

### Load Tests
- Voting endpoint: 100 req/s
- Live count broadcast: 1000 concurrent connections

---

## 📝 Development Workflow

```bash
# Start dependencies
make docker-up

# Run migrations
make migrate-up

# Run dev server (with hot reload)
make dev

# Run tests
make test

# Run linter
make lint

# Build
make build
```

---

## 🎓 Learning Resources

- **Go**: https://go.dev/tour/
- **chi router**: https://go-chi.io/
- **pgx**: https://pkg.go.dev/github.com/jackc/pgx/v5
- **sqlc**: https://sqlc.dev/
- **WebSocket**: https://nhooyr.io/websocket/

---

## 🤝 Contributing

1. Create feature branch
2. Implement with tests
3. Run linter and tests
4. Submit PR with description

---

## 📞 Support

- Documentation: `/docs`
- Issues: GitHub Issues
- Architecture questions: See ARCHITECTURE.md
- API questions: See API.md

---

**Status**: ✅ Ready for Phase 2 (Repository Implementation)

**Last Updated**: 2024-11-19
