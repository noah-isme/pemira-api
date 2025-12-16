# Running Pemira API with Docker Alpine

## Overview
Pemira API menggunakan multi-stage Docker build dengan Alpine Linux untuk image yang ringan dan aman.

## Quick Start

### 1. Start All Services
```bash
docker-compose up -d --build
```

Perintah ini akan:
- Build Docker image dari Alpine
- Start PostgreSQL (Alpine)
- Start Redis (Alpine)
- Start API server (Alpine)

### 2. Check Container Status
```bash
docker-compose ps
```

Expected output:
```
NAME              IMAGE                COMMAND                  SERVICE    CREATED         STATUS                   PORTS
pemira-api        pemira-api-api       "/api"                   api        5 seconds ago   Up 5 seconds            0.0.0.0:8080->8080/tcp
pemira-postgres   postgres:16-alpine   "docker-entrypoint.s…"   postgres   1 minute ago    Up 1 minute (healthy)   0.0.0.0:5432->5432/tcp
pemira-redis      redis:7-alpine       "docker-entrypoint.s…"   redis      1 minute ago    Up 1 minute (healthy)   0.0.0.0:6379->6379/tcp
```

### 3. View Logs
```bash
# All services
docker-compose logs -f

# API only
docker-compose logs -f api

# Last 50 lines
docker-compose logs api --tail=50
```

### 4. Stop Services
```bash
docker-compose down
```

### 5. Stop and Remove Volumes
```bash
docker-compose down -v
```

## Docker Image Details

### Base Image
- **Builder Stage**: `golang:alpine` - Compile Go application
- **Runtime Stage**: `alpine:latest` - Minimal runtime environment

### Image Size
- Builder stage: ~400MB (temporary)
- Final image: ~25MB (production)

### Security Features
- Runs as non-root user (`nonroot:nonroot`)
- User ID: 1000, Group ID: 1000
- No unnecessary packages
- Only required CA certificates and timezone data

## Environment Variables

### Default Configuration (docker-compose.yml)
```yaml
APP_ENV: development
HTTP_PORT: 8080
DATABASE_URL: postgres://pemira:pemira@postgres:5432/pemira?sslmode=disable
REDIS_URL: redis://redis:6379/0
JWT_SECRET: change-this-in-production
JWT_EXPIRATION: 24h
LOG_LEVEL: info
CORS_ALLOWED_ORIGINS: "*"
```

### Override Environment Variables
Create `.env.docker` file:
```bash
JWT_SECRET=your-super-secure-secret-key
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
LOG_LEVEL=debug
```

Then use:
```bash
docker-compose --env-file .env.docker up -d
```

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "data": {
    "status": "ok"
  }
}
```

### Public Endpoints

#### List Elections
```bash
curl http://localhost:8080/api/v1/elections | jq .
```

#### List Candidates
```bash
curl http://localhost:8080/api/v1/elections/1/candidates | jq .
```

#### Candidate Detail
```bash
curl http://localhost:8080/api/v1/elections/1/candidates/1 | jq .
```

### Admin Endpoints

#### Login as Admin
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }' | jq .
```

Response:
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "role": "SUPER_ADMIN"
    }
  }
}
```

#### List Candidates with QR Code (Admin)
```bash
TOKEN="your-admin-token"

curl http://localhost:8080/api/v1/admin/elections/1/candidates \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Response includes QR code:
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "election_id": 1,
        "number": 1,
        "name": "Kandidat A",
        "status": "APPROVED",
        "qr_code": {
          "id": 1,
          "token": "qr_e1c1_abc123",
          "url": "https://pemira.local/ballot-qr/qr_e1c1_abc123",
          "payload": "PEMIRA-UNIWA|E:1|C:1|V:1",
          "version": 1,
          "is_active": true
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total_items": 1,
      "total_pages": 1
    }
  }
}
```

#### Get Candidate Detail with QR Code (Admin)
```bash
TOKEN="your-admin-token"

curl http://localhost:8080/api/v1/admin/elections/1/candidates/1 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## Makefile Commands

### Start Docker Services
```bash
make docker-up
```

### Stop Docker Services
```bash
make docker-down
```

### Build Application
```bash
make build
```

### Run Tests
```bash
make test
```

## Troubleshooting

### Container Won't Start

#### Check logs
```bash
docker-compose logs api
```

#### Check if ports are available
```bash
# Check if port 8080 is in use
lsof -i :8080

# Check if PostgreSQL port is in use
lsof -i :5432

# Check if Redis port is in use
lsof -i :6379
```

#### Rebuild containers
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Database Connection Issues

#### Verify PostgreSQL is healthy
```bash
docker-compose ps postgres
```

#### Connect to PostgreSQL directly
```bash
docker exec -it pemira-postgres psql -U pemira -d pemira
```

#### Check database migrations
```bash
docker exec -it pemira-api /api --migrate
```

### Redis Connection Issues

#### Verify Redis is healthy
```bash
docker-compose ps redis
```

#### Connect to Redis directly
```bash
docker exec -it pemira-redis redis-cli
> PING
PONG
```

### API Not Responding

#### Check if container is running
```bash
docker ps | grep pemira-api
```

#### Check API logs
```bash
docker-compose logs api --tail=100
```

#### Restart API container
```bash
docker-compose restart api
```

#### Access container shell
```bash
docker exec -it pemira-api sh
```

## Development Workflow

### 1. Make Code Changes
Edit files in your local directory.

### 2. Rebuild and Restart
```bash
docker-compose up -d --build api
```

### 3. Watch Logs
```bash
docker-compose logs -f api
```

### 4. Test Changes
```bash
curl http://localhost:8080/health
```

## Production Considerations

### 1. Update Environment Variables
```yaml
environment:
  APP_ENV: production
  JWT_SECRET: use-strong-random-secret
  CORS_ALLOWED_ORIGINS: https://your-frontend.com
  LOG_LEVEL: warn
```

### 2. Use Docker Secrets
```yaml
secrets:
  jwt_secret:
    file: ./secrets/jwt_secret.txt
  db_password:
    file: ./secrets/db_password.txt
```

### 3. Resource Limits
```yaml
api:
  deploy:
    resources:
      limits:
        cpus: '1'
        memory: 512M
      reservations:
        cpus: '0.5'
        memory: 256M
```

### 4. Health Checks
```yaml
api:
  healthcheck:
    test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    start_period: 40s
```

## Monitoring

### Container Stats
```bash
docker stats pemira-api pemira-postgres pemira-redis
```

### Disk Usage
```bash
docker system df
```

### Clean Up
```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Remove everything unused
docker system prune -a
```

## Advanced Usage

### Custom Network
```bash
# Create custom network
docker network create pemira-network

# Update docker-compose.yml
networks:
  default:
    external:
      name: pemira-network
```

### Volume Backup

#### Backup PostgreSQL
```bash
docker exec pemira-postgres pg_dump -U pemira pemira > backup.sql
```

#### Restore PostgreSQL
```bash
cat backup.sql | docker exec -i pemira-postgres psql -U pemira pemira
```

### Scale Services
```bash
# Scale API to 3 instances
docker-compose up -d --scale api=3
```

## References

- [Dockerfile](../Dockerfile)
- [docker-compose.yml](../docker-compose.yml)
- [Makefile](../Makefile)
- [Admin QR Code Implementation](./changes/ADMIN_CANDIDATE_QR_CODE.md)
- [API Examples](./changes/ADMIN_QR_CODE_EXAMPLES.md)

## Support

If you encounter issues:
1. Check logs: `docker-compose logs api`
2. Verify health: `docker-compose ps`
3. Rebuild: `docker-compose up -d --build`
4. Check documentation in `/docs` folder