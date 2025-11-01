#!/bin/bash

# HIBE Key Revocation System - Test Script
# This script demonstrates the complete workflow of the revocation system

BASE_URL="http://localhost:8080"
HIERARCHY="testHierarchy"
URI="facility/bin123/record"
START_TIME=1565119330
END_TIME=1565219330

echo "======================================"
echo "HIBE Key Revocation System Test Suite"
echo "======================================"
echo ""

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test results
print_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_section() {
    echo ""
    echo "======================================"
    echo "$1"
    echo "======================================"
    echo ""
}

# Check if server is running
print_section "1. Server Health Check"
print_test "Checking if HIBE API server is running..."
HEALTH=$(curl -s "$BASE_URL/health" 2>/dev/null)
if [ $? -eq 0 ]; then
    print_success "Server is running"
    echo "$HEALTH" | jq '.' 2>/dev/null || echo "$HEALTH"
else
    print_error "Server is not running. Please start the server first."
    exit 1
fi

# Test 1: Generate Key ID
print_section "2. Generate Key ID"
print_test "Generating key ID for delegation parameters..."
KEYGEN_RESPONSE=$(curl -s -X POST "$BASE_URL/revoke/generate-key-id" \
  -H "Content-Type: application/json" \
  -d "{
    \"hierarchy\": \"$HIERARCHY\",
    \"uri\": \"$URI\",
    \"startTime\": $START_TIME,
    \"endTime\": $END_TIME
  }")

KEY_ID=$(echo "$KEYGEN_RESPONSE" | jq -r '.keyId')

if [ ! -z "$KEY_ID" ] && [ "$KEY_ID" != "null" ]; then
    print_success "Key ID generated: $KEY_ID"
    echo "$KEYGEN_RESPONSE" | jq '.'
else
    print_error "Failed to generate key ID"
    exit 1
fi

# Test 2: Create a Delegation
print_section "3. Create Delegation"
print_test "Creating new delegation..."
DELEGATE_RESPONSE=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
  -H "Content-Type: application/json" \
  -d "{
    \"uri\": \"$URI\",
    \"hierarchy\": \"$HIERARCHY\",
    \"startTime\": $START_TIME,
    \"endTime\": $END_TIME
  }")

SUCCESS=$(echo "$DELEGATE_RESPONSE" | jq -r '.success')
if [ "$SUCCESS" == "true" ]; then
    print_success "Delegation created successfully"
    echo "$DELEGATE_RESPONSE" | jq '.'
else
    print_error "Failed to create delegation"
    echo "$DELEGATE_RESPONSE" | jq '.'
fi

# Test 3: Check Revocation Status (should not be revoked)
print_section "4. Check Initial Revocation Status"
print_test "Checking if key is revoked (should be false)..."
CHECK_RESPONSE=$(curl -s "$BASE_URL/revoke/check/$KEY_ID")
IS_REVOKED=$(echo "$CHECK_RESPONSE" | jq -r '.isRevoked')

if [ "$IS_REVOKED" == "false" ]; then
    print_success "Key is not revoked (as expected)"
    echo "$CHECK_RESPONSE" | jq '.'
else
    print_error "Key is unexpectedly revoked"
    echo "$CHECK_RESPONSE" | jq '.'
fi

# Test 4: Get Initial Statistics
print_section "5. Initial Revocation Statistics"
print_test "Getting revocation statistics..."
STATS_RESPONSE=$(curl -s "$BASE_URL/revocations/stats")
print_success "Statistics retrieved"
echo "$STATS_RESPONSE" | jq '.'

# Test 5: Revoke the Key
print_section "6. Revoke the Key"
print_test "Revoking the delegated key..."
REVOKE_RESPONSE=$(curl -s -X POST "$BASE_URL/revoke" \
  -H "Content-Type: application/json" \
  -d "{
    \"uri\": \"$URI\",
    \"hierarchy\": \"$HIERARCHY\",
    \"startTime\": $START_TIME,
    \"endTime\": $END_TIME,
    \"revokedBy\": \"test-script\",
    \"reason\": \"Testing revocation functionality\",
    \"effectiveFor\": 0
  }")

REVOKE_SUCCESS=$(echo "$REVOKE_RESPONSE" | jq -r '.success')
if [ "$REVOKE_SUCCESS" == "true" ]; then
    print_success "Key revoked successfully"
    echo "$REVOKE_RESPONSE" | jq '.'
else
    print_error "Failed to revoke key"
    echo "$REVOKE_RESPONSE" | jq '.'
fi

# Test 6: Check Revocation Status Again (should be revoked now)
print_section "7. Verify Revocation"
print_test "Checking if key is now revoked (should be true)..."
CHECK2_RESPONSE=$(curl -s "$BASE_URL/revoke/check/$KEY_ID")
IS_REVOKED2=$(echo "$CHECK2_RESPONSE" | jq -r '.isRevoked')

if [ "$IS_REVOKED2" == "true" ]; then
    print_success "Key is revoked (as expected)"
    echo "$CHECK2_RESPONSE" | jq '.'
else
    print_error "Key is not revoked (unexpected)"
    echo "$CHECK2_RESPONSE" | jq '.'
fi

# Test 7: Try to Delegate Again (should fail)
print_section "8. Attempt to Delegate Revoked Key"
print_test "Trying to delegate with revoked parameters (should fail)..."
DELEGATE2_RESPONSE=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
  -H "Content-Type: application/json" \
  -d "{
    \"uri\": \"$URI\",
    \"hierarchy\": \"$HIERARCHY\",
    \"startTime\": $START_TIME,
    \"endTime\": $END_TIME
  }")

SUCCESS2=$(echo "$DELEGATE2_RESPONSE" | jq -r '.success')
if [ "$SUCCESS2" == "false" ]; then
    print_success "Delegation correctly denied for revoked key"
    echo "$DELEGATE2_RESPONSE" | jq '.'
else
    print_error "Delegation incorrectly allowed for revoked key"
    echo "$DELEGATE2_RESPONSE" | jq '.'
fi

# Test 8: List All Revocations
print_section "9. List All Revocations"
print_test "Retrieving all revocations..."
LIST_RESPONSE=$(curl -s "$BASE_URL/revocations")
COUNT=$(echo "$LIST_RESPONSE" | jq -r '.count')
print_success "Found $COUNT revocation(s)"
echo "$LIST_RESPONSE" | jq '.'

# Test 9: Get Revocations by URI
print_section "10. Get Revocations by URI"
print_test "Getting revocations for URI: $URI..."
URI_ENCODED=$(echo "$URI" | sed 's/\//%2F/g')
URI_RESPONSE=$(curl -s "$BASE_URL/revocations/uri/$URI_ENCODED")
URI_COUNT=$(echo "$URI_RESPONSE" | jq -r '.count')
print_success "Found $URI_COUNT revocation(s) for URI"
echo "$URI_RESPONSE" | jq '.'

# Test 10: Get Updated Statistics
print_section "11. Updated Revocation Statistics"
print_test "Getting updated statistics..."
STATS2_RESPONSE=$(curl -s "$BASE_URL/revocations/stats")
ACTIVE=$(echo "$STATS2_RESPONSE" | jq -r '.activeRevocations')
print_success "Statistics updated - Active revocations: $ACTIVE"
echo "$STATS2_RESPONSE" | jq '.'

# Test 11: Clear Revocation
print_section "12. Clear Revocation"
print_test "Clearing the revocation (reinstating the key)..."
CLEAR_RESPONSE=$(curl -s -X DELETE "$BASE_URL/revoke/$KEY_ID")
CLEAR_SUCCESS=$(echo "$CLEAR_RESPONSE" | jq -r '.success')

if [ "$CLEAR_SUCCESS" == "true" ]; then
    print_success "Revocation cleared successfully"
    echo "$CLEAR_RESPONSE" | jq '.'
else
    print_error "Failed to clear revocation"
    echo "$CLEAR_RESPONSE" | jq '.'
fi

# Test 12: Verify Key is No Longer Revoked
print_section "13. Verify Key Reinstatement"
print_test "Checking if key is no longer revoked..."
CHECK3_RESPONSE=$(curl -s "$BASE_URL/revoke/check/$KEY_ID")
IS_REVOKED3=$(echo "$CHECK3_RESPONSE" | jq -r '.isRevoked')

if [ "$IS_REVOKED3" == "false" ]; then
    print_success "Key is no longer revoked (as expected)"
    echo "$CHECK3_RESPONSE" | jq '.'
else
    print_error "Key is still revoked (unexpected)"
    echo "$CHECK3_RESPONSE" | jq '.'
fi

# Test 13: Test URI-based Revocation
print_section "14. Test URI-based Revocation"
print_test "Revoking all keys for URI: $URI..."
URI_REVOKE_RESPONSE=$(curl -s -X POST "$BASE_URL/revoke-by-uri" \
  -H "Content-Type: application/json" \
  -d "{
    \"uri\": \"$URI\",
    \"revokedBy\": \"test-script\",
    \"reason\": \"Testing URI-based revocation\"
  }")

URI_REVOKE_SUCCESS=$(echo "$URI_REVOKE_RESPONSE" | jq -r '.success')
if [ "$URI_REVOKE_SUCCESS" == "true" ]; then
    REVOKED_COUNT=$(echo "$URI_REVOKE_RESPONSE" | jq -r '.revokedCount')
    print_success "Revoked $REVOKED_COUNT key(s) for URI"
    echo "$URI_REVOKE_RESPONSE" | jq '.'
else
    print_error "Failed to revoke by URI"
    echo "$URI_REVOKE_RESPONSE" | jq '.'
fi

# Test 14: Temporary Revocation
print_section "15. Test Temporary Revocation (10 seconds)"
NEW_URI="test/temporary"
NEW_START=$START_TIME
NEW_END=$END_TIME

print_test "Creating temporary revocation (expires in 10 seconds)..."
TEMP_REVOKE_RESPONSE=$(curl -s -X POST "$BASE_URL/revoke" \
  -H "Content-Type: application/json" \
  -d "{
    \"uri\": \"$NEW_URI\",
    \"hierarchy\": \"$HIERARCHY\",
    \"startTime\": $NEW_START,
    \"endTime\": $NEW_END,
    \"revokedBy\": \"test-script\",
    \"reason\": \"Testing temporary revocation\",
    \"effectiveFor\": 10
  }")

TEMP_KEY_ID=$(echo "$TEMP_REVOKE_RESPONSE" | jq -r '.keyId')
if [ ! -z "$TEMP_KEY_ID" ] && [ "$TEMP_KEY_ID" != "null" ]; then
    print_success "Temporary revocation created: $TEMP_KEY_ID"
    echo "$TEMP_REVOKE_RESPONSE" | jq '.'

    print_test "Waiting 12 seconds for revocation to expire..."
    sleep 12

    print_test "Running cleanup to remove expired revocations..."
    CLEANUP_RESPONSE=$(curl -s -X POST "$BASE_URL/revocations/cleanup")
    REMOVED=$(echo "$CLEANUP_RESPONSE" | jq -r '.removedCount')
    print_success "Cleanup removed $REMOVED expired revocation(s)"
    echo "$CLEANUP_RESPONSE" | jq '.'
fi

# Test 15: Final Statistics
print_section "16. Final Statistics"
print_test "Getting final statistics..."
FINAL_STATS=$(curl -s "$BASE_URL/revocations/stats")
print_success "Final statistics retrieved"
echo "$FINAL_STATS" | jq '.'

# Test 16: List Active Revocations
print_section "17. List Active Revocations"
print_test "Getting only active revocations..."
ACTIVE_RESPONSE=$(curl -s "$BASE_URL/revocations?status=active")
ACTIVE_COUNT=$(echo "$ACTIVE_RESPONSE" | jq -r '.count')
print_success "Found $ACTIVE_COUNT active revocation(s)"
echo "$ACTIVE_RESPONSE" | jq '.'

# Summary
print_section "Test Summary"
echo "All tests completed!"
echo ""
echo "Key operations tested:"
echo "  ✓ Key ID generation"
echo "  ✓ Delegation creation"
echo "  ✓ Revocation status checking"
echo "  ✓ Key revocation"
echo "  ✓ Delegation denial for revoked keys"
echo "  ✓ Revocation listing"
echo "  ✓ URI-based revocation"
echo "  ✓ Revocation clearing"
echo "  ✓ Statistics retrieval"
echo "  ✓ Temporary revocations"
echo "  ✓ Cleanup of expired revocations"
echo ""
echo "For detailed API documentation, see REVOCATION_GUIDE.md"
echo "For integration instructions, see INTEGRATION_INSTRUCTIONS.md"
echo ""
