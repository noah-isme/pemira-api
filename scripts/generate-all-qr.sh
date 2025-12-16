#!/bin/bash

# Script to generate QR codes for all candidates
# Usage: ./scripts/generate-all-qr.sh [election_id]

set -e

API_URL="${API_URL:-http://localhost:8080}"
ADMIN_USERNAME="${ADMIN_USERNAME:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-password123}"
ELECTION_ID="${1:-1}"

echo "========================================"
echo "Generate QR Codes for All Candidates"
echo "========================================"
echo ""
echo "API URL: $API_URL"
echo "Election ID: $ELECTION_ID"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

success() {
    echo -e "${GREEN}✓ $1${NC}"
}

error() {
    echo -e "${RED}✗ $1${NC}"
}

info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Step 1: Login
echo "Step 1: Login"
echo "-------------"
LOGIN_RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"${ADMIN_USERNAME}\",\"password\":\"${ADMIN_PASSWORD}\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | sed 's/"access_token":"\(.*\)"/\1/')

if [ -z "$TOKEN" ]; then
    error "Failed to login"
    echo "$LOGIN_RESPONSE"
    exit 1
fi

success "Login successful"
echo ""

# Step 2: Get all candidates
echo "Step 2: Get Candidates"
echo "----------------------"
CANDIDATES_RESPONSE=$(curl -s "${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates?limit=100" \
    -H "Authorization: Bearer ${TOKEN}")

# Extract candidate IDs and check for QR codes
CANDIDATE_IDS=$(echo "$CANDIDATES_RESPONSE" | grep -o '"id":[0-9]*' | sed 's/"id":\([0-9]*\)/\1/' | head -100)

if [ -z "$CANDIDATE_IDS" ]; then
    error "No candidates found"
    exit 1
fi

TOTAL_CANDIDATES=$(echo "$CANDIDATE_IDS" | wc -l)
success "Found $TOTAL_CANDIDATES candidate(s)"
echo ""

# Step 3: Generate QR for each candidate
echo "Step 3: Generate QR Codes"
echo "-------------------------"

GENERATED=0
SKIPPED=0
FAILED=0

for CANDIDATE_ID in $CANDIDATE_IDS; do
    # Get candidate detail to check if QR exists
    DETAIL=$(curl -s "${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates/${CANDIDATE_ID}" \
        -H "Authorization: Bearer ${TOKEN}")

    CANDIDATE_NAME=$(echo "$DETAIL" | grep -o '"name":"[^"]*"' | head -1 | sed 's/"name":"\(.*\)"/\1/')
    HAS_QR=$(echo "$DETAIL" | grep -c '"qr_code"' || true)

    if [ "$HAS_QR" -gt 0 ]; then
        QR_TOKEN=$(echo "$DETAIL" | grep -o '"token":"[^"]*"' | head -1 | sed 's/"token":"\(.*\)"/\1/')
        info "Candidate #${CANDIDATE_ID} (${CANDIDATE_NAME}) - Already has QR: ${QR_TOKEN}"
        SKIPPED=$((SKIPPED + 1))
        continue
    fi

    # Generate QR code
    info "Generating QR for candidate #${CANDIDATE_ID} (${CANDIDATE_NAME})..."

    GENERATE_RESPONSE=$(curl -s -X POST \
        "${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates/${CANDIDATE_ID}/qr/generate" \
        -H "Authorization: Bearer ${TOKEN}")

    # Check if generation was successful
    NEW_QR_TOKEN=$(echo "$GENERATE_RESPONSE" | grep -o '"token":"[^"]*"' | head -1 | sed 's/"token":"\(.*\)"/\1/')

    if [ -n "$NEW_QR_TOKEN" ]; then
        success "Generated QR for #${CANDIDATE_ID} (${CANDIDATE_NAME}): ${NEW_QR_TOKEN}"
        GENERATED=$((GENERATED + 1))
    else
        error "Failed to generate QR for #${CANDIDATE_ID} (${CANDIDATE_NAME})"
        echo "$GENERATE_RESPONSE"
        FAILED=$((FAILED + 1))
    fi

    # Small delay to avoid overwhelming the server
    sleep 0.2
done

echo ""
echo "========================================"
echo "Summary"
echo "========================================"
echo "Total candidates: $TOTAL_CANDIDATES"
success "Generated: $GENERATED"
info "Skipped (already has QR): $SKIPPED"
if [ $FAILED -gt 0 ]; then
    error "Failed: $FAILED"
fi
echo ""

# Step 4: Verify
echo "Step 4: Verification"
echo "--------------------"
info "Checking candidates without QR code..."

CANDIDATES_RESPONSE=$(curl -s "${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates?limit=100" \
    -H "Authorization: Bearer ${TOKEN}")

WITHOUT_QR=$(echo "$CANDIDATES_RESPONSE" | grep -o '"qr_code":null' | wc -l)

if [ "$WITHOUT_QR" -eq 0 ]; then
    success "All candidates now have QR codes!"
else
    error "Still $WITHOUT_QR candidate(s) without QR code"
fi

echo ""
info "Done!"
exit 0
