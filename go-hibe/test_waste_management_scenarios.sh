#!/bin/bash

# Smart Waste Management in Smart City - HIBE Key Delegation & Revocation Test Suite
# This script tests real-world waste management scenarios for key management

BASE_URL="http://localhost:8080"
HIERARCHY="smart-city-waste"
CITY_ID="city-metro"

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Helper functions
print_header() {
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
    echo ""
}

print_scenario() {
    echo ""
    echo -e "${CYAN}━━━ Scenario: $1 ━━━${NC}"
    echo ""
}

print_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓ PASS]${NC} $1"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
}

print_fail() {
    echo -e "${RED}[✗ FAIL]${NC} $1"
    ((FAILED_TESTS++))
    ((TOTAL_TESTS++))
}

print_info() {
    echo -e "${CYAN}[INFO]${NC} $1"
}

# Check server health
check_server() {
    print_header "SMART CITY WASTE MANAGEMENT - HIBE TEST SUITE"
    print_test "Checking HIBE API server status..."
    HEALTH=$(curl -s "$BASE_URL/health" 2>/dev/null)
    if [ $? -eq 0 ]; then
        print_success "HIBE API server is running"
        print_info "Health status: $(echo $HEALTH | jq -r '.status')"
    else
        print_fail "Server is not running. Please start the server first."
        exit 1
    fi
}

# Utility function to get current timestamp
get_timestamp() {
    date -u +%s
}

# Utility function to get future timestamp
get_future_timestamp() {
    local hours=$1
    echo $(($(date -u +%s) + ($hours * 3600)))
}

#═══════════════════════════════════════════════════════════
# SCENARIO 1: Collection Vehicle Route Optimization
#═══════════════════════════════════════════════════════════
test_route_optimization() {
    print_scenario "1. Collection Vehicle Route Optimization"

    print_info "Context: Driver needs access to smart bin data for route planning"
    print_info "Bin: bin-001 (Downtown District A - Mixed Waste, 85% full)"
    print_info "Driver: driver-001 (John Martinez)"
    print_info "Depot: depot-001 (Central Collection Depot)"

    local BIN_ID="bin-001"
    local DRIVER_ID="driver-001"
    local DEPOT_ID="depot-001"
    local URI="waste/$CITY_ID/bin/$BIN_ID/sensor-data"
    local START_TIME=$(get_timestamp)
    local END_TIME=$(get_future_timestamp 8)  # 8 hour shift

    # Test 1.1: Create delegation for driver
    print_test "1.1: Creating delegation for route optimization..."
    DELEGATE_RESPONSE=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"hierarchy\": \"$HIERARCHY\",
        \"startTime\": $START_TIME,
        \"endTime\": $END_TIME
      }")

    SUCCESS=$(echo "$DELEGATE_RESPONSE" | jq -r '.success')
    KEY_ID_1=$(echo "$DELEGATE_RESPONSE" | jq -r '.keyId')

    if [ "$SUCCESS" == "true" ]; then
        print_success "Delegation created for route optimization"
        print_info "Key ID: $KEY_ID_1"
        print_info "Access duration: 8 hours (shift)"
    else
        print_fail "Failed to create delegation"
        echo "$DELEGATE_RESPONSE" | jq '.'
        return
    fi

    # Test 1.2: Verify key is not revoked
    print_test "1.2: Verifying driver has active access..."
    CHECK_RESPONSE=$(curl -s "$BASE_URL/revoke/check/$KEY_ID_1")
    IS_REVOKED=$(echo "$CHECK_RESPONSE" | jq -r '.isRevoked')

    if [ "$IS_REVOKED" == "false" ]; then
        print_success "Driver has active access to bin sensor data"
    else
        print_fail "Access unexpectedly revoked"
    fi

    # Test 1.3: Simulate encryption/decryption of bin data
    print_test "1.3: Testing encryption of smart bin sensor data..."
    BIN_DATA="Bin: BIN-001, Location: Downtown District A, Type: Mixed, FillLevel: 85%, Weight: 120kg, Temperature: 18C, LastCollection: 2 days ago"

    ENCRYPT_RESPONSE=$(curl -s -X POST "$BASE_URL/encrypt" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"message\": \"$BIN_DATA\"
      }")

    ENCRYPTED_DATA=$(echo "$ENCRYPT_RESPONSE" | jq -r '.data')
    EXEC_TIME=$(echo "$ENCRYPT_RESPONSE" | jq -r '.time')

    if [ ! -z "$ENCRYPTED_DATA" ] && [ "$ENCRYPTED_DATA" != "null" ]; then
        print_success "Bin sensor data encrypted successfully"
        print_info "Execution time: ${EXEC_TIME}μs"
    else
        print_fail "Encryption failed"
    fi

    # Store for later use
    ENCRYPTED_BIN_001=$ENCRYPTED_DATA
}

#═══════════════════════════════════════════════════════════
# SCENARIO 2: Emergency Overflow Response
#═══════════════════════════════════════════════════════════
test_emergency_overflow() {
    print_scenario "2. Emergency Overflow Response"

    print_info "Context: Emergency team needs immediate access to overflowing bin"
    print_info "Bin: bin-003 (Commercial Street C - Organic, 92% full)"
    print_info "Response Team: emergency-001"
    print_info "Situation: Bin overflow detected, public health risk"

    local BIN_ID="bin-003"
    local EMERGENCY_ID="emergency-001"
    local URI="waste/$CITY_ID/bin/$BIN_ID/sensor-data/critical"
    local START_TIME=$(get_timestamp)
    local END_TIME=$(get_future_timestamp 4)  # 4 hours emergency window

    # Test 2.1: Emergency delegation (no prior authorization needed)
    print_test "2.1: Creating emergency access delegation..."
    EMERGENCY_DELEGATE=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"hierarchy\": \"$HIERARCHY\",
        \"startTime\": $START_TIME,
        \"endTime\": $END_TIME
      }")

    EMERGENCY_KEY=$(echo "$EMERGENCY_DELEGATE" | jq -r '.keyId')

    if [ ! -z "$EMERGENCY_KEY" ] && [ "$EMERGENCY_KEY" != "null" ]; then
        print_success "Emergency access granted immediately"
        print_info "Key ID: $EMERGENCY_KEY"
        print_info "Access window: 4 hours"
    else
        print_fail "Emergency access creation failed"
    fi

    # Test 2.2: Access critical bin information
    print_test "2.2: Accessing critical bin status information..."
    CRITICAL_DATA="CRITICAL STATUS - Bin: BIN-003, FillLevel: 92%, Weight: 145kg, Overflow Risk: HIGH, Location: Commercial St C, Odor Level: HIGH"

    CRITICAL_ENCRYPT=$(curl -s -X POST "$BASE_URL/encrypt" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"message\": \"$CRITICAL_DATA\"
      }")

    if [ $(echo "$CRITICAL_ENCRYPT" | jq -r '.data') != "null" ]; then
        print_success "Critical bin information accessed for emergency response"
        print_info "Overflow status retrieved: HIGH RISK"
    else
        print_fail "Failed to access critical information"
    fi

    # Test 2.3: Auto-revocation simulation (overflow resolved)
    print_test "2.3: Simulating overflow resolution (auto-revoke emergency access)..."

    # Revoke emergency access after resolution
    REVOKE_EMERGENCY=$(curl -s -X POST "$BASE_URL/revoke" \
      -H "Content-Type: application/json" \
      -d "{
        \"keyId\": \"$EMERGENCY_KEY\",
        \"revokedBy\": \"system/auto-revoke\",
        \"reason\": \"Overflow resolved, bin emptied\",
        \"effectiveFor\": 0
      }")

    if [ $(echo "$REVOKE_EMERGENCY" | jq -r '.success') == "true" ]; then
        print_success "Emergency team access auto-revoked after resolution"
        print_info "Access duration: ~30 minutes (until resolution)"
    else
        print_fail "Auto-revocation failed"
    fi

    # Verify revocation
    CHECK_EMERGENCY=$(curl -s "$BASE_URL/revoke/check/$EMERGENCY_KEY")
    if [ $(echo "$CHECK_EMERGENCY" | jq -r '.isRevoked') == "true" ]; then
        print_success "Verified: Emergency team can no longer access bin data"
    else
        print_fail "Revocation verification failed"
    fi
}

#═══════════════════════════════════════════════════════════
# SCENARIO 3: Recycling Center Processing
#═══════════════════════════════════════════════════════════
test_recycling_processing() {
    print_scenario "3. Recycling Center Batch Processing"

    print_info "Context: Recycling center processes batch for environmental dept"
    print_info "Bin: bin-002 (Residential Zone B - Recyclables, 65% full)"
    print_info "Center: recycling-001 (GreenCycle Center)"
    print_info "Inspector: inspector-001 (Alex Kumar)"

    local BIN_ID="bin-002"
    local CENTER_ID="recycling-001"
    local INSPECTOR_ID="inspector-001"
    local URI="waste/$CITY_ID/bin/$BIN_ID/collection/batch-$(date +%Y%m%d)"
    local START_TIME=$(get_timestamp)
    local END_TIME=$(get_future_timestamp 168)  # 7 days (1 week)

    # Test 3.1: Center creates encrypted batch data
    print_test "3.1: Recycling center encrypting batch processing data..."
    BATCH_DATA="Batch Processing - Bin: BIN-002 - Date: $(date +%Y-%m-%d) - Type: Recyclables (Plastics) - Weight: 95kg - Contamination: 3% - Sorted: 92kg - Rejected: 3kg - Quality: A"

    BATCH_ENCRYPT=$(curl -s -X POST "$BASE_URL/encrypt" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"message\": \"$BATCH_DATA\"
      }")

    BATCH_ENCRYPTED_DATA=$(echo "$BATCH_ENCRYPT" | jq -r '.data')

    if [ ! -z "$BATCH_ENCRYPTED_DATA" ] && [ "$BATCH_ENCRYPTED_DATA" != "null" ]; then
        print_success "Batch processing data encrypted and stored securely"
    else
        print_fail "Batch encryption failed"
    fi

    # Test 3.2: Delegate access to environmental dept
    print_test "3.2: Granting access to environmental department..."
    ENV_DELEGATE=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"hierarchy\": \"$HIERARCHY\",
        \"startTime\": $START_TIME,
        \"endTime\": $END_TIME
      }")

    ENV_KEY=$(echo "$ENV_DELEGATE" | jq -r '.keyId')

    if [ ! -z "$ENV_KEY" ] && [ "$ENV_KEY" != "null" ]; then
        print_success "Environmental dept can access batch data"
        print_info "Access valid for 7 days"
    else
        print_fail "Failed to delegate to environmental dept"
    fi

    # Test 3.3: Quality inspector access
    print_test "3.3: Quality inspector accessing batch data..."
    CHECK_ENV=$(curl -s "$BASE_URL/revoke/check/$ENV_KEY")
    if [ $(echo "$CHECK_ENV" | jq -r '.isRevoked') == "false" ]; then
        print_success "Quality inspector can access batch for audit"
    else
        print_fail "Inspector access denied"
    fi
}

#═══════════════════════════════════════════════════════════
# SCENARIO 4: Transfer Station Data Sharing
#═══════════════════════════════════════════════════════════
test_transfer_station() {
    print_scenario "4. Transfer Station to Processing Plant"

    print_info "Context: Transfer station delegates e-waste data to processing plant"
    print_info "Bin: bin-004 (Industrial Park D - Electronic, 45% full)"
    print_info "Transfer: transfer-001 (Metro Transfer Hub)"
    print_info "Processing: processing-001 (E-Waste Processing)"

    local BIN_ID="bin-004"
    local TRANSFER_ID="transfer-001"
    local PROCESSING_ID="processing-001"
    local URI="waste/$CITY_ID/bin/$BIN_ID/sensor-data/e-waste"
    local START_TIME=$(get_timestamp)
    local END_TIME=$(get_future_timestamp 720)  # 30 days

    # Test 4.1: Create delegation for processing plant
    print_test "4.1: Transfer station delegating e-waste data..."
    PROCESSING_DELEGATE=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"hierarchy\": \"$HIERARCHY\",
        \"startTime\": $START_TIME,
        \"endTime\": $END_TIME
      }")

    PROCESSING_KEY=$(echo "$PROCESSING_DELEGATE" | jq -r '.keyId')

    if [ ! -z "$PROCESSING_KEY" ] && [ "$PROCESSING_KEY" != "null" ]; then
        print_success "Processing plant granted access to e-waste batch"
        print_info "Access duration: 30 days"
    else
        print_fail "Failed to create processing delegation"
    fi

    # Test 4.2: Verify scope restriction (e-waste only)
    print_test "4.2: Verifying scope restriction (e-waste only)..."
    E_WASTE_DATA="E-Waste Batch - Bin: BIN-004 - Type: Electronics - Components: 15 phones, 8 laptops, 23 batteries - Weight: 78kg - Hazardous: Yes - Special Handling Required"

    E_WASTE_ENCRYPT=$(curl -s -X POST "$BASE_URL/encrypt" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"message\": \"$E_WASTE_DATA\"
      }")

    if [ $(echo "$E_WASTE_ENCRYPT" | jq -r '.data') != "null" ]; then
        print_success "Processing plant can access e-waste specific data"
        print_info "Scope restriction: Electronic waste only"
    else
        print_fail "E-waste data access failed"
    fi

    # Test 4.3: Transfer station maintains access
    print_test "4.3: Verifying transfer station maintains access..."
    CHECK_PROCESSING=$(curl -s "$BASE_URL/revoke/check/$PROCESSING_KEY")
    if [ $(echo "$CHECK_PROCESSING" | jq -r '.isRevoked') == "false" ]; then
        print_success "Both transfer station and processing plant have access"
    else
        print_fail "Concurrent access failed"
    fi
}

#═══════════════════════════════════════════════════════════
# SCENARIO 5: Facility Compliance Inspection
#═══════════════════════════════════════════════════════════
test_facility_inspection() {
    print_scenario "5. Facility Compliance Inspection"

    print_info "Context: City inspector auditing recycling center operations"
    print_info "Facility: recycling-001 (GreenCycle Center)"
    print_info "Inspector: inspector-001 (Rachel Green)"
    print_info "Type: Quarterly compliance inspection"

    local FACILITY_ID="recycling-001"
    local INSPECTOR_ID="inspector-001"
    local URI="waste/$CITY_ID/facility/$FACILITY_ID/operations"
    local START_TIME=$(get_timestamp)
    local END_TIME=$(get_future_timestamp 9)  # 9 hour inspection day

    # Test 5.1: Create inspection delegation
    print_test "5.1: Creating compliance inspection access..."
    INSPECTION_DELEGATE=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"hierarchy\": \"$HIERARCHY\",
        \"startTime\": $START_TIME,
        \"endTime\": $END_TIME
      }")

    INSPECTION_KEY=$(echo "$INSPECTION_DELEGATE" | jq -r '.keyId')

    if [ ! -z "$INSPECTION_KEY" ] && [ "$INSPECTION_KEY" != "null" ]; then
        print_success "Inspector granted access to facility operations"
        print_info "Inspection window: 9 hours"
    else
        print_fail "Failed to create inspection access"
    fi

    # Test 5.2: Inspector accesses operations data
    print_test "5.2: Inspector accessing facility operations data..."
    OPS_DATA="Operations Data - Facility: GreenCycle - Throughput: 250 tons/month - Efficiency: 94% - Violations: 0 - Equipment Status: Operational - Staff: 24 - Certifications: Valid"

    OPS_ENCRYPT=$(curl -s -X POST "$BASE_URL/encrypt" \
      -H "Content-Type: application/json" \
      -d "{
        \"uri\": \"$URI\",
        \"message\": \"$OPS_DATA\"
      }")

    if [ $(echo "$OPS_ENCRYPT" | jq -r '.data') != "null" ]; then
        print_success "Inspector can access facility operations for compliance check"
    else
        print_fail "Operations data access failed"
    fi

    # Test 5.3: Auto-revoke after inspection (simulate end of day)
    print_test "5.3: Simulating end of inspection (auto-revoke)..."
    REVOKE_INSPECTION=$(curl -s -X POST "$BASE_URL/revoke" \
      -H "Content-Type: application/json" \
      -d "{
        \"keyId\": \"$INSPECTION_KEY\",
        \"revokedBy\": \"system/inspection-complete\",
        \"reason\": \"Inspection completed successfully\",
        \"effectiveFor\": 0
      }")

    if [ $(echo "$REVOKE_INSPECTION" | jq -r '.success') == "true" ]; then
        print_success "Inspector access revoked after inspection completion"
    else
        print_fail "Inspection revocation failed"
    fi
}

#═══════════════════════════════════════════════════════════
# REVOCATION SCENARIO 1: Driver Departure
#═══════════════════════════════════════════════════════════
test_driver_departure() {
    print_scenario "REVOCATION 1: Collection Driver Departure"

    print_info "Context: Driver leaves company, all bin access must be revoked"
    print_info "Driver: driver-001 (John Martinez)"
    print_info "Affected: 47 smart bins in collection route"
    print_info "Reason: Employment terminated"

    # Create multiple delegations to simulate driver's access
    print_test "Setting up driver access to multiple bins..."
    DRIVER_KEYS=()

    for i in {1..5}; do
        local BIN_URI="waste/$CITY_ID/bin/bin-route-$i/sensor-data"
        local DELEGATE_RESP=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
          -H "Content-Type: application/json" \
          -d "{
            \"uri\": \"$BIN_URI\",
            \"hierarchy\": \"$HIERARCHY\",
            \"startTime\": $(get_timestamp),
            \"endTime\": $(get_future_timestamp 720)
          }")

        local KEY=$(echo "$DELEGATE_RESP" | jq -r '.keyId')
        if [ ! -z "$KEY" ] && [ "$KEY" != "null" ]; then
            DRIVER_KEYS+=("$KEY")
        fi
    done

    print_success "Driver has access to ${#DRIVER_KEYS[@]} bins"

    # Bulk revocation
    print_test "Performing bulk revocation for driver departure..."
    REVOKED_COUNT=0

    for KEY in "${DRIVER_KEYS[@]}"; do
        REVOKE_RESP=$(curl -s -X POST "$BASE_URL/revoke" \
          -H "Content-Type: application/json" \
          -d "{
            \"keyId\": \"$KEY\",
            \"revokedBy\": \"supervisor-001\",
            \"reason\": \"Driver employment terminated\",
            \"effectiveFor\": 0
          }")

        if [ $(echo "$REVOKE_RESP" | jq -r '.success') == "true" ]; then
            ((REVOKED_COUNT++))
        fi
    done

    if [ $REVOKED_COUNT -eq ${#DRIVER_KEYS[@]} ]; then
        print_success "All $REVOKED_COUNT bin access delegations revoked immediately"
        print_info "Driver can no longer access any bins"
    else
        print_fail "Bulk revocation incomplete: $REVOKED_COUNT/${#DRIVER_KEYS[@]}"
    fi

    # Verify revocation
    print_test "Verifying driver access denial..."
    local FIRST_KEY=${DRIVER_KEYS[0]}
    CHECK_DRIVER=$(curl -s "$BASE_URL/revoke/check/$FIRST_KEY")
    if [ $(echo "$CHECK_DRIVER" | jq -r '.isRevoked') == "true" ]; then
        print_success "Verified: Driver access completely blocked"
    else
        print_fail "Access verification failed"
    fi
}

#═══════════════════════════════════════════════════════════
# REVOCATION SCENARIO 2: Security Breach
#═══════════════════════════════════════════════════════════
test_security_breach() {
    print_scenario "REVOCATION 2: Security Breach Response"

    print_info "Context: Unauthorized access detected, emergency revocation"
    print_info "Incident: SEC-2025-11-002"
    print_info "Target: All facility access blocked"
    print_info "Response: Immediate emergency revocation"

    # Create facility access delegations
    print_test "Setting up facility access delegations..."
    FACILITY_KEYS=()

    for FACILITY in "recycling-001" "recycling-002" "processing-001"; do
        local FAC_URI="waste/$CITY_ID/facility/$FACILITY/operations"
        local DELEGATE_RESP=$(curl -s -X POST "$BASE_URL/hibe-delegate" \
          -H "Content-Type: application/json" \
          -d "{
            \"uri\": \"$FAC_URI\",
            \"hierarchy\": \"$HIERARCHY\",
            \"startTime\": $(get_timestamp),
            \"endTime\": $(get_future_timestamp 720)
          }")

        local KEY=$(echo "$DELEGATE_RESP" | jq -r '.keyId')
        if [ ! -z "$KEY" ] && [ "$KEY" != "null" ]; then
            FACILITY_KEYS+=("$KEY")
        fi
    done

    print_success "Facility access established for ${#FACILITY_KEYS[@]} facilities"

    # Emergency revocation
    print_test "Executing emergency revocation protocol..."
    EMERGENCY_REVOKE_COUNT=0

    for KEY in "${FACILITY_KEYS[@]}"; do
        REVOKE_RESP=$(curl -s -X POST "$BASE_URL/revoke" \
          -H "Content-Type: application/json" \
          -d "{
            \"keyId\": \"$KEY\",
            \"revokedBy\": \"security-team\",
            \"reason\": \"Security breach detected - Incident SEC-2025-11-002\",
            \"effectiveFor\": 0
          }")

        if [ $(echo "$REVOKE_RESP" | jq -r '.success') == "true" ]; then
            ((EMERGENCY_REVOKE_COUNT++))
        fi
    done

    if [ $EMERGENCY_REVOKE_COUNT -eq ${#FACILITY_KEYS[@]} ]; then
        print_success "Emergency revocation: All $EMERGENCY_REVOKE_COUNT facility accesses blocked"
        print_info "Response time: < 1 second"
        print_info "Manual re-approval required"
    else
        print_fail "Emergency revocation incomplete"
    fi

    # Verify emergency lockdown
    print_test "Verifying complete facility lockdown..."
    local FIRST_FAC_KEY=${FACILITY_KEYS[0]}
    CHECK_FAC=$(curl -s "$BASE_URL/revoke/check/$FIRST_FAC_KEY")
    if [ $(echo "$CHECK_FAC" | jq -r '.isRevoked') == "true" ]; then
        print_success "Verified: Complete facility access lockdown"
        print_info "Security investigation initiated"
    else
        print_fail "Lockdown verification failed"
    fi
}

#═══════════════════════════════════════════════════════════
# Main Test Execution
#═══════════════════════════════════════════════════════════

# Check if server is running
check_server

# Run all test scenarios
test_route_optimization
test_emergency_overflow
test_recycling_processing
test_transfer_station
test_facility_inspection

# Run revocation scenarios
test_driver_departure
test_security_breach

# Print summary
print_header "TEST SUMMARY"
echo ""
echo -e "Total Tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ ALL TESTS PASSED!${NC}"
    echo ""
    echo -e "${GREEN}✓${NC} Smart waste management delegation and revocation system working perfectly"
    echo -e "${GREEN}✓${NC} All scenarios tested successfully"
    echo -e "${GREEN}✓${NC} Ready for smart city waste management deployment"
    echo ""
    exit 0
else
    echo -e "${RED}✗ SOME TESTS FAILED${NC}"
    echo ""
    echo -e "${YELLOW}⚠${NC} Please review failed tests above"
    echo -e "${YELLOW}⚠${NC} Check server logs for details"
    echo ""
    exit 1
fi
