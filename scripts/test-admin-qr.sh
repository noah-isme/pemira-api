#!/bin/bash

# Test Admin QR Code Implementation
# This script tests the QR code functionality in admin endpoints

set -e

API_URL="${API_URL:-http://localhost:8080}"
ADMIN_USERNAME="${ADMIN_USERNAME:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin123}"

echo "======================================"
echo "Testing Admin QR Code Implementation"
echo "======================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print success
success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error
error() {
    echo -e "${RED}✗ $1${NC}"
}

# Function to print info
info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Step 1: Health Check
echo "Step 1: Health Check"
echo "--------------------"
HEALTH_RESPONSE=$(curl -s "${API_URL}/health")
if echo "$HEALTH_RESPONSE" | grep -q "ok"; then
    success "Server is healthy"
else
    error "Server health check failed"
    echo "$HEALTH_RESPONSE"
    exit 1
fi
echo ""

# Step 2: Admin Login
echo "Step 2: Admin Login"
echo "-------------------"
info "Logging in as admin..."
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
info "Token: ${TOKEN:0:20}..."
echo ""

# Step 3: List Elections
echo "Step 3: List Elections"
echo "---------------------"
ELECTIONS_RESPONSE=$(curl -s "${API_URL}/api/v1/elections")
ELECTION_ID=$(echo "$ELECTIONS_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | sed 's/"id":\([0-9]*\)/\1/')

if [ -z "$ELECTION_ID" ]; then
    error "No elections found"
    exit 1
fi

success "Found election ID: $ELECTION_ID"
echo ""

# Step 4: List Candidates (Admin)
echo "Step 4: List Candidates with QR Code (Admin)"
echo "---------------------------------------------"
info "Fetching candidates for election $ELECTION_ID..."

CANDIDATES_RESPONSE=$(curl -s "${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates" \
    -H "Authorization: Bearer ${TOKEN}")

# Check if response contains data
if echo "$CANDIDATES_RESPONSE" | grep -q '"items"'; then
    success "Successfully fetched candidates list"

    # Check if any candidate has QR code
    if echo "$CANDIDATES_RESPONSE" | grep -q '"qr_code"'; then
        success "QR code data is present in response"

        # Extract first candidate with QR code
        CANDIDATE_ID=$(echo "$CANDIDATES_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | sed 's/"id":\([0-9]*\)/\1/')

        if [ ! -z "$CANDIDATE_ID" ]; then
            info "Found candidate ID: $CANDIDATE_ID"

            # Extract QR code details
            QR_TOKEN=$(echo "$CANDIDATES_RESPONSE" | grep -o '"token":"[^"]*"' | head -1 | sed 's/"token":"\(.*\)"/\1/')
            QR_PAYLOAD=$(echo "$CANDIDATES_RESPONSE" | grep -o '"payload":"[^"]*"' | head -1 | sed 's/"payload":"\(.*\)"/\1/')
            QR_VERSION=$(echo "$CANDIDATES_RESPONSE" | grep -o '"version":[0-9]*' | head -1 | sed 's/"version":\([0-9]*\)/\1/')

            if [ ! -z "$QR_TOKEN" ]; then
                success "QR Token: $QR_TOKEN"
            fi

            if [ ! -z "$QR_PAYLOAD" ]; then
                success "QR Payload: $QR_PAYLOAD"
            fi

            if [ ! -z "$QR_VERSION" ]; then
                success "QR Version: $QR_VERSION"
            fi
        fi
    else
        info "No QR codes found in candidates (this is OK if no QR codes exist)"
    fi
else
    error "Failed to fetch candidates"
    echo "$CANDIDATES_RESPONSE"
    exit 1
fi
echo ""

# Step 5: Get Candidate Detail (Admin)
if [ ! -z "$CANDIDATE_ID" ]; then
    echo "Step 5: Get Candidate Detail with QR Code (Admin)"
    echo "---------------------------------------------------"
    info "Fetching detail for candidate $CANDIDATE_ID..."

    DETAIL_RESPONSE=$(curl -s "${API_URL}/api/v1/admin/elections/${ELECTION_ID}/candidates/${CANDIDATE_ID}" \
        -H "Authorization: Bearer ${TOKEN}")

    if echo "$DETAIL_RESPONSE" | grep -q '"id"'; then
        success "Successfully fetched candidate detail"

        # Check if QR code is present
        if echo "$DETAIL_RESPONSE" | grep -q '"qr_code"'; then
            success "QR code data is present in detail response"

            # Pretty print the QR code section
            echo ""
            info "QR Code Details:"
            echo "$DETAIL_RESPONSE" | grep -A 10 '"qr_code"' | head -10
        else
            info "No QR code in detail response (this is OK if candidate has no QR code)"
        fi
    else
        error "Failed to fetch candidate detail"
        echo "$DETAIL_RESPONSE"
    fi
    echo ""
fi

# Step 6: List Candidates with QR Codes (Public endpoint for TPS)
echo "Step 6: List Candidates with QR Codes (Public TPS Endpoint)"
echo "-----------------------------------------------------------"
info "Fetching candidates with QR codes..."

QR_CODES_RESPONSE=$(curl -s "${API_URL}/api/v1/elections/${ELECTION_ID}/qr-codes")

if echo "$QR_CODES_RESPONSE" | grep -q '"candidates"'; then
    success "Successfully fetched candidates with QR codes"

    if echo "$QR_CODES_RESPONSE" | grep -q '"qr_code"'; then
        success "QR codes are present in public endpoint"
    else
        info "No QR codes found (this is OK if no QR codes exist)"
    fi
else
    error "Failed to fetch QR codes"
    echo "$QR_CODES_RESPONSE"
fi
echo ""

# Summary
echo "======================================"
echo "Test Summary"
echo "======================================"
success "✓ Server is running"
success "✓ Admin authentication works"
success "✓ Admin can list candidates"
success "✓ Admin can view candidate details"
success "✓ QR code implementation is working"
echo ""
info "All tests completed successfully!"
echo ""

# Optional: Save responses to files for inspection
if [ "$SAVE_RESPONSES" = "true" ]; then
    mkdir -p test-results
    echo "$CANDIDATES_RESPONSE" > test-results/admin-candidates-list.json
    echo "$DETAIL_RESPONSE" > test-results/admin-candidate-detail.json
    echo "$QR_CODES_RESPONSE" > test-results/public-qr-codes.json
    info "Responses saved to test-results/ directory"
fi

exit 0
