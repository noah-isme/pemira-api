# Production Environment Variables

## 📋 Complete List

Total: **11 environment variables** (10 required + 1 optional)

---

## 🔴 REQUIRED VARIABLES

### 1. APP_ENV
```
APP_ENV=production
```
- **Description**: Application environment mode
- **Required**: Yes
- **Default**: None

### 2. HTTP_PORT
```
HTTP_PORT=8080
```
- **Description**: Port untuk aplikasi
- **Required**: Yes
- **Default**: 8080

### 3. DATABASE_URL
```
DATABASE_URL=postgresql://kvokrhrtikcshvtbrwiu:egqsndudnxifwwiaoiucitxgnbvjyo@9qasp5v56q8ckkf5dc.apn.leapcellpool.com:6438/wgompdtswdesvswbydon?sslmode=require
```
- **Description**: PostgreSQL connection string
- **Required**: Yes
- **Status**: ✅ Already migrated and ready
- **Note**: Include `?sslmode=require` for production

### 4. JWT_SECRET
```
JWT_SECRET=<GENERATE-NEW-SECRET>
```
- **Description**: Secret key untuk JWT token authentication
- **Required**: Yes
- **Security**: ⚠️ MUST be unique and random (min 32 characters)
- **Generate with**:
  ```bash
  openssl rand -base64 32
  ```
- **Example output**:
  ```
  4jeNG90OoKN7CC1UfM2f4aWHLWxjL3/9bMcolyqMXWM=
  ```

### 5. JWT_EXPIRATION
```
JWT_EXPIRATION=24h
```
- **Description**: Token validity duration
- **Required**: Yes
- **Format**: `24h` (24 hours), `1h` (1 hour), `30m` (30 minutes)

### 6. SUPABASE_URL
```
SUPABASE_URL=https://xqzfrodnznhjstfstvyz.supabase.co
```
- **Description**: Supabase project URL
- **Required**: Yes (for file uploads)
- **Location**: Supabase Dashboard > Settings > API > Project URL

### 7. SUPABASE_SECRET_KEY
```
SUPABASE_SECRET_KEY=<YOUR-SERVICE-ROLE-KEY>
```
- **Description**: Supabase service role secret key
- **Required**: Yes (for file uploads)
- **Security**: ⚠️ Use `service_role` key, NOT `anon` key!
- **Location**: Supabase Dashboard > Settings > API > Project API keys > `service_role` (click reveal)
- **Format**: Starts with `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`

### 8. SUPABASE_MEDIA_BUCKET
```
SUPABASE_MEDIA_BUCKET=pemira
```
- **Description**: Bucket name for candidate photos/media
- **Required**: Yes
- **Default**: pemira

### 9. SUPABASE_BRANDING_BUCKET
```
SUPABASE_BRANDING_BUCKET=pemira
```
- **Description**: Bucket name for branding logos
- **Required**: Yes
- **Default**: pemira
- **Note**: Can use same bucket as media

### 10. CORS_ALLOWED_ORIGINS
```
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com
```
- **Description**: Allowed frontend domains for CORS
- **Required**: Yes
- **Format**: 
  - Single domain: `https://app.yourdomain.com`
  - Multiple domains: `https://app.com,https://www.app.com`
- **⚠️ Important**: 
  - Must include protocol (`https://`)
  - No trailing slash
  - Replace with actual production domain

### 11. LOG_LEVEL
```
LOG_LEVEL=info
```
- **Description**: Logging level
- **Required**: Yes
- **Options**: `debug`, `info`, `warn`, `error`
- **Recommended**: `info` for production

---

## 🟡 OPTIONAL VARIABLES

### 12. REDIS_URL
```
REDIS_URL=redis://default:password@host:6379
```
- **Description**: Redis connection for caching/sessions
- **Required**: No
- **Note**: Application works without Redis

---

## 📝 Copy-Paste Template for Leapcell

```bash
APP_ENV=production
HTTP_PORT=8080
DATABASE_URL=postgresql://kvokrhrtikcshvtbrwiu:egqsndudnxifwwiaoiucitxgnbvjyo@9qasp5v56q8ckkf5dc.apn.leapcellpool.com:6438/wgompdtswdesvswbydon?sslmode=require
JWT_SECRET=<PASTE-GENERATED-SECRET-HERE>
JWT_EXPIRATION=24h
SUPABASE_URL=https://xqzfrodnznhjstfstvyz.supabase.co
SUPABASE_SECRET_KEY=<PASTE-SERVICE-ROLE-KEY-HERE>
SUPABASE_MEDIA_BUCKET=pemira
SUPABASE_BRANDING_BUCKET=pemira
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com
LOG_LEVEL=info
```

---

## ⚠️ BEFORE DEPLOYMENT - Must Replace

### 1. Generate JWT_SECRET
```bash
# Run this command locally:
openssl rand -base64 32

# Example output (DO NOT USE THIS):
4jeNG90OoKN7CC1UfM2f4aWHLWxjL3/9bMcolyqMXWM=
```

**Copy the output and paste as `JWT_SECRET` value**

### 2. Get Supabase Service Role Key

1. Go to: https://supabase.com/dashboard
2. Select your project
3. Navigate to: **Settings** > **API**
4. Under "Project API keys", find `service_role`
5. Click **Reveal** and copy the key
6. Paste as `SUPABASE_SECRET_KEY` value

**⚠️ Important**: 
- Use `service_role` key (has full access)
- DO NOT use `anon` key (limited access)
- Keep this key secret!

### 3. Set Frontend Domain

Replace `https://your-frontend-domain.com` with your actual frontend URL.

**Examples**:
- `https://pemira.youruniversity.edu`
- `https://vote.example.com`
- Multiple: `https://app.com,https://www.app.com`

**⚠️ Common mistakes**:
- ❌ `http://localhost:3000` (don't use localhost)
- ❌ `your-domain.com` (missing protocol)
- ❌ `https://your-domain.com/` (trailing slash)
- ✅ `https://your-domain.com` (correct)

---

## 🔍 Supabase Bucket Setup

Before deployment, create the storage bucket:

### Steps:

1. **Login to Supabase Dashboard**
   - URL: https://supabase.com/dashboard
   - Select project: `xqzfrodnznhjstfstvyz`

2. **Go to Storage**
   - Click **Storage** in sidebar

3. **Create Bucket**
   - Click **New bucket**
   - Name: `pemira`
   - Settings:
     - ✅ Public bucket: **Yes** (important!)
     - File size limit: **10 MB**
     - Allowed MIME types: `image/jpeg, image/png`

4. **Verify**
   - Bucket should appear in list
   - Status: Public
   - Access: Anyone can read

---

## ✅ Pre-Deployment Checklist

Before deploying to Leapcell, verify:

- [ ] `DATABASE_URL` is set (already migrated ✓)
- [ ] `JWT_SECRET` generated with `openssl rand -base64 32`
- [ ] `SUPABASE_SECRET_KEY` is **service_role** key (not anon!)
- [ ] Supabase bucket `pemira` created and set to **public**
- [ ] `CORS_ALLOWED_ORIGINS` updated with actual frontend domain
- [ ] All 11 environment variables pasted to Leapcell
- [ ] No placeholder values (`<...>`) remaining

---

## 🚀 Deployment Steps

### 1. Set Environment Variables in Leapcell

1. Login to Leapcell Dashboard
2. Go to your project
3. Navigate to **Settings** > **Environment Variables**
4. Paste all variables from template above
5. Replace placeholder values:
   - `JWT_SECRET`
   - `SUPABASE_SECRET_KEY`
   - `CORS_ALLOWED_ORIGINS`

### 2. Deploy

1. Click **Deploy** button
2. Wait for build to complete
3. Check deployment logs for errors

### 3. Verify Deployment

```bash
# Replace with your Leapcell URL
API_URL="https://your-app.leapcell.io"

# Health check
curl $API_URL/health

# Login test
curl -X POST $API_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

Expected responses:
- Health: `200 OK`
- Login: `200 OK` with `access_token`

---

## 🔐 Security Notes

### Critical Security Rules:

1. **JWT_SECRET**
   - ⚠️ Must be random and unique
   - ⚠️ Never commit to git
   - ⚠️ Never share publicly
   - ✅ Generate new for each environment

2. **SUPABASE_SECRET_KEY**
   - ⚠️ Use `service_role` key only
   - ⚠️ Never commit to git
   - ⚠️ Never expose to frontend
   - ✅ Only use in backend

3. **DATABASE_URL**
   - ⚠️ Contains password - keep secret
   - ⚠️ Never commit to git
   - ✅ Use `?sslmode=require` in production

4. **CORS_ALLOWED_ORIGINS**
   - ⚠️ Only list trusted domains
   - ⚠️ Don't use wildcard (`*`) in production
   - ✅ Explicitly list each domain

---

## 🆘 Troubleshooting

### Database Connection Failed
- Check `DATABASE_URL` format
- Ensure `sslmode=require` is included
- Verify database is accessible from Leapcell

### CORS Errors
- Verify `CORS_ALLOWED_ORIGINS` has correct domain
- Include protocol: `https://`
- No trailing slash
- Check browser console for exact origin

### Upload Fails (500 Error)
- Verify `SUPABASE_SECRET_KEY` is **service_role** key
- Check bucket `pemira` exists and is **public**
- Verify `SUPABASE_URL` is correct
- Check Supabase dashboard for storage errors

### JWT Token Invalid
- Verify `JWT_SECRET` is set and not empty
- Check `JWT_EXPIRATION` format is valid (e.g., `24h`)
- Token expires after set duration - try fresh login

---

## 📚 Related Documentation

- **Quick Start**: `QUICK_DEPLOY_GUIDE.md`
- **Full Guide**: `DEPLOYMENT.md`
- **Checklist**: `DEPLOYMENT_CHECKLIST.md`
- **Technical Details**: `BRANDING_LOGO_IMPLEMENTATION.md`

---

**Last Updated**: 2025-11-26  
**Status**: Production Ready ✅
