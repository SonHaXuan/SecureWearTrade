# HIBE Key Revocation System Guide

## Overview

The HIBE Key Revocation System provides a comprehensive mechanism to revoke and manage delegated cryptographic keys in the HIBE encryption infrastructure. This system allows administrators to revoke keys that have been compromised, are no longer needed, or should be temporarily suspended.

## Table of Contents

- [Architecture](#architecture)
- [Key Concepts](#key-concepts)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Integration Guide](#integration-guide)
- [Best Practices](#best-practices)
- [Security Considerations](#security-considerations)

## Architecture

### Components

1. **RevocationList**: Core data structure that maintains all revoked keys
   - Thread-safe with RWMutex protection
   - Indexed by both KeyID and URI for fast lookup
   - Supports time-based revocations (immediate, scheduled, temporary)

2. **RevocationEntry**: Individual revocation record containing:
   - KeyID: Unique identifier for the revoked key
   - URI: The resource pattern associated with the key
   - Hierarchy: The hierarchy context
   - Timestamps: Creation, effective start, and expiration
   - Metadata: Reason, revoked by, etc.

3. **Key Generation**: Deterministic key ID generation based on:
   - Hierarchy
   - URI pattern
   - Start time
   - End time

### Revocation Flow

```
┌─────────────────┐
│  Delegation     │
│   Request       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│  Generate       │────▶│  Check          │
│  KeyID          │     │  Revocation     │
└─────────────────┘     └────────┬────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                    ▼                         ▼
            ┌───────────────┐         ┌──────────────┐
            │   Revoked?    │         │  Not Revoked │
            │   Deny        │         │   Allow      │
            └───────────────┘         └──────────────┘
```

## Key Concepts

### Key ID Generation

Each delegated key is assigned a unique identifier based on its parameters:

```go
keyID = SHA256(hierarchy + ":" + uri + ":" + startTime + ":" + endTime)
```

This ensures:
- Deterministic key identification
- No need to store delegation metadata
- Easy verification without database lookups

### Revocation Types

1. **Immediate Revocation**: Takes effect immediately
2. **Scheduled Revocation**: Takes effect at a future time
3. **Temporary Revocation**: Automatically expires after a duration
4. **Permanent Revocation**: Never expires (default)

### Revocation States

- **Pending**: Revocation scheduled but not yet effective
- **Active**: Currently in effect
- **Expired**: Was active but has now expired

## API Endpoints

### 1. Revoke a Key

**Endpoint**: `POST /revoke`

**Description**: Revoke a specific delegated key.

**Request Body**:
```json
{
  "keyId": "abc123...",           // Optional if providing other params
  "uri": "facility/bin123",   // Required if keyId not provided
  "hierarchy": "testHierarchy",   // Required if keyId not provided
  "startTime": 1565119330,        // Unix timestamp
  "endTime": 1565219330,          // Unix timestamp
  "revokedBy": "admin@example.com",
  "reason": "Security breach detected",
  "effectiveFrom": "2025-10-31T10:00:00Z",  // Optional, default: now
  "effectiveFor": 86400           // Optional, duration in seconds, 0 = permanent
}
```

**Response**:
```json
{
  "success": true,
  "message": "Key revoked successfully",
  "keyId": "abc123...",
  "revokedAt": "2025-10-31T09:30:00Z",
  "effectiveFrom": "2025-10-31T10:00:00Z",
  "effectiveUntil": "2025-11-01T10:00:00Z"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/revoke \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "facility/bin123/record",
    "hierarchy": "testHierarchy",
    "startTime": 1565119330,
    "endTime": 1565219330,
    "revokedBy": "security-admin",
    "reason": "Suspected key compromise",
    "effectiveFor": 0
  }'
```

### 2. Revoke by URI

**Endpoint**: `POST /revoke-by-uri`

**Description**: Revoke all keys associated with a specific URI pattern.

**Request Body**:
```json
{
  "uri": "facility/bin123",
  "revokedBy": "admin@example.com",
  "reason": "Bin data access terminated"
}
```

**Response**:
```json
{
  "success": true,
  "message": "Revoked 3 key(s) for URI: facility/bin123",
  "revokedCount": 3,
  "uri": "facility/bin123"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/revoke-by-uri \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "facility/bin123",
    "revokedBy": "admin",
    "reason": "Access rights expired"
  }'
```

### 3. Check Revocation Status

**Endpoint**: `GET /revoke/check/:keyId`

**Description**: Check if a specific key is revoked.

**Response**:
```json
{
  "keyId": "abc123...",
  "isRevoked": true,
  "revocationDetails": {
    "keyId": "abc123...",
    "uri": "facility/bin123",
    "hierarchy": "testHierarchy",
    "revokedAt": "2025-10-31T09:30:00Z",
    "revokedBy": "admin@example.com",
    "reason": "Security breach detected",
    "effectiveFrom": "2025-10-31T10:00:00Z",
    "effectiveUntil": "0001-01-01T00:00:00Z"
  }
}
```

**Example**:
```bash
curl http://localhost:8080/revoke/check/abc123...
```

### 4. Check Revocation by Parameters

**Endpoint**: `POST /revoke/check`

**Description**: Check revocation status using key parameters instead of key ID.

**Request Body**:
```json
{
  "hierarchy": "testHierarchy",
  "uri": "facility/bin123",
  "startTime": 1565119330,
  "endTime": 1565219330
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/revoke/check \
  -H "Content-Type: application/json" \
  -d '{
    "hierarchy": "testHierarchy",
    "uri": "facility/bin123/record",
    "startTime": 1565119330,
    "endTime": 1565219330
  }'
```

### 5. List Revocations

**Endpoint**: `GET /revocations?status={all|active|expired|pending}`

**Description**: List all revocations with optional filtering.

**Query Parameters**:
- `status`: Filter by status (default: "all")

**Response**:
```json
{
  "count": 5,
  "status": "active",
  "revocations": [
    {
      "keyId": "abc123...",
      "uri": "facility/bin123",
      "hierarchy": "testHierarchy",
      "revokedAt": "2025-10-31T09:30:00Z",
      "revokedBy": "admin",
      "reason": "Security breach",
      "effectiveFrom": "2025-10-31T10:00:00Z",
      "effectiveUntil": "0001-01-01T00:00:00Z"
    }
  ]
}
```

**Example**:
```bash
# Get all active revocations
curl http://localhost:8080/revocations?status=active

# Get all revocations
curl http://localhost:8080/revocations
```

### 6. Get Revocations by URI

**Endpoint**: `GET /revocations/uri/:uri`

**Description**: Get all revocations for a specific URI.

**Example**:
```bash
curl http://localhost:8080/revocations/uri/facility/bin123
```

### 7. Clear Revocation

**Endpoint**: `DELETE /revoke/:keyId`

**Description**: Remove a revocation entry (reinstate a key).

**Response**:
```json
{
  "success": true,
  "message": "Revocation cleared for key: abc123...",
  "keyId": "abc123..."
}
```

**Example**:
```bash
curl -X DELETE http://localhost:8080/revoke/abc123...
```

### 8. Get Revocation Statistics

**Endpoint**: `GET /revocations/stats`

**Description**: Get statistics about the revocation list.

**Response**:
```json
{
  "totalRevocations": 10,
  "activeRevocations": 7,
  "expiredRevocations": 2,
  "pendingRevocations": 1,
  "uniqueUris": 5,
  "lastUpdated": "2025-10-31T10:00:00Z"
}
```

**Example**:
```bash
curl http://localhost:8080/revocations/stats
```

### 9. Cleanup Expired Revocations

**Endpoint**: `POST /revocations/cleanup`

**Description**: Remove all expired revocations from the list.

**Response**:
```json
{
  "success": true,
  "message": "Cleanup completed",
  "removedCount": 3,
  "remainingCount": 7
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/revocations/cleanup
```

### 10. Generate Key ID

**Endpoint**: `POST /revoke/generate-key-id`

**Description**: Utility endpoint to generate a key ID from parameters.

**Request Body**:
```json
{
  "hierarchy": "testHierarchy",
  "uri": "facility/bin123",
  "startTime": 1565119330,
  "endTime": 1565219330
}
```

**Response**:
```json
{
  "keyId": "abc123...",
  "hierarchy": "testHierarchy",
  "uri": "facility/bin123",
  "startTime": 1565119330,
  "endTime": 1565219330
}
```

## Enhanced Delegation Endpoints

### 1. Delegate with Revocation Check

**Endpoint**: `POST /hibe-delegate`

**Description**: Create a new delegation with automatic revocation checking.

**Request Body**:
```json
{
  "uri": "facility/bin123",
  "hierarchy": "testHierarchy",
  "startTime": 1565119330,
  "endTime": 1565219330
}
```

**Response**:
```json
{
  "success": true,
  "keyId": "abc123...",
  "data": "base64_encoded_delegation_key",
  "uri": "facility/bin123",
  "hierarchy": "testHierarchy",
  "startTime": 1565119330,
  "endTime": 1565219330,
  "executionTime": 1234,
  "message": "Key delegated successfully"
}
```

**Error Response** (if revoked):
```json
{
  "success": false,
  "keyId": "abc123...",
  "error": "cannot delegate: key is revoked (reason: Security breach)"
}
```

### 2. Decrypt with Revocation Check

**Endpoint**: `POST /decrypt-with-revocation`

**Description**: Decrypt data with automatic revocation checking.

**Request Body**:
```json
{
  "uri": "facility/bin123",
  "encryptedMessage": "base64_encoded_ciphertext",
  "hierarchy": "testHierarchy",
  "startTime": 1565119330,
  "endTime": 1565219330
}
```

**Response**:
```json
{
  "success": true,
  "data": "decrypted plaintext",
  "executionTime": 456,
  "memoryUsage": 52428800,
  "cpuPercentage": 15.5,
  "powerUsageWatts": 2.5,
  "energyConsumptionJoules": 0.00114
}
```

## Integration with main.go

To integrate the revocation system into your main.go file, add the following after setting up the Gin router:

```go
// In main.go, after r := gin.Default()

// Register revocation endpoints
RegisterRevocationEndpoints(r)

// Register enhanced delegation endpoints
RegisterDelegationWithRevocationEndpoint(r, ctx, store, encoder)

// Register enhanced decrypt endpoint
RegisterEnhancedDecryptEndpoint(r, ctx, state, now)

// Register delegation management endpoints
RegisterDelegationManagementEndpoints(r)
```

## Best Practices

### 1. Revocation Reasons

Always provide clear, detailed reasons for revocations:
- ✅ "Security breach detected on 2025-10-31 - unauthorized access attempt"
- ✅ "Employee termination - access rights revoked per policy"
- ❌ "revoked"
- ❌ "no longer needed"

### 2. Time-based Revocations

Use appropriate revocation timing:
- **Immediate**: For security incidents
- **Scheduled**: For planned access termination
- **Temporary**: For temporary suspensions (audits, investigations)

### 3. Regular Cleanup

Schedule regular cleanup of expired revocations:
```bash
# Run daily cleanup via cron
0 2 * * * curl -X POST http://localhost:8080/revocations/cleanup
```

### 4. Monitoring

Monitor revocation statistics:
```bash
# Check stats every hour
*/60 * * * * curl http://localhost:8080/revocations/stats
```

### 5. Audit Logging

Always include audit information:
- Who revoked the key
- When it was revoked
- Why it was revoked
- When it takes effect

## Security Considerations

### 1. Access Control

Implement authentication and authorization for revocation endpoints:
- Only authorized administrators should revoke keys
- Implement role-based access control (RBAC)
- Log all revocation operations

### 2. Rate Limiting

Prevent abuse by implementing rate limiting:
```go
// Example rate limiting middleware
r.Use(rateLimit())
```

### 3. Validation

Always validate input parameters:
- Verify URI patterns
- Validate time ranges
- Check hierarchy permissions

### 4. Persistence

The current implementation uses in-memory storage. For production:
- Implement database persistence
- Add replication for high availability
- Implement backup and recovery

### 5. Distributed Systems

For distributed deployments:
- Implement revocation list synchronization
- Use distributed caching (Redis)
- Implement eventual consistency handling

## Example Workflows

### Workflow 1: Emergency Key Revocation

```bash
# 1. Identify compromised key
KEY_ID="abc123..."

# 2. Immediately revoke
curl -X POST http://localhost:8080/revoke \
  -H "Content-Type: application/json" \
  -d "{
    \"keyId\": \"$KEY_ID\",
    \"revokedBy\": \"security-team\",
    \"reason\": \"SECURITY INCIDENT: Unauthorized access detected\",
    \"effectiveFrom\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  }"

# 3. Verify revocation
curl http://localhost:8080/revoke/check/$KEY_ID

# 4. Notify stakeholders
# (external notification system)
```

### Workflow 2: Scheduled Access Termination

```bash
# 1. Schedule revocation for employee departure
curl -X POST http://localhost:8080/revoke-by-uri \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "company/employee/john-doe",
    "revokedBy": "hr-system",
    "reason": "Employee departure scheduled for 2025-11-30"
  }'

# 2. Monitor pending revocations
curl http://localhost:8080/revocations?status=pending
```

### Workflow 3: Temporary Suspension

```bash
# 1. Temporarily revoke for 24 hours during audit
curl -X POST http://localhost:8080/revoke \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "finance/records",
    "hierarchy": "company",
    "startTime": 1565119330,
    "endTime": 1565219330,
    "revokedBy": "audit-team",
    "reason": "Temporary suspension during quarterly audit",
    "effectiveFor": 86400
  }'

# 2. After 24 hours, automatically expires
```

## Troubleshooting

### Issue: Revocation not taking effect

**Check**:
```bash
# Verify effective time
curl http://localhost:8080/revoke/check/KEY_ID

# Check system time
date -u
```

### Issue: Unable to clear revocation

**Solution**:
```bash
# Verify key ID is correct
curl -X DELETE http://localhost:8080/revoke/CORRECT_KEY_ID
```

### Issue: Performance degradation

**Solutions**:
- Run cleanup regularly
- Implement database indexing
- Use caching for frequent lookups

## Performance Metrics

Expected performance (based on in-memory implementation):
- Revocation check: < 1ms
- Add revocation: < 2ms
- List revocations: < 5ms (for 10,000 entries)
- Cleanup: < 10ms (for 1,000 expired entries)

## Future Enhancements

1. **Database Persistence**: PostgreSQL/MySQL backend
2. **Distributed Cache**: Redis integration
3. **Webhook Notifications**: Real-time notifications on revocations
4. **Certificate Revocation Lists (CRL)**: Standard CRL format export
5. **OCSP Support**: Online Certificate Status Protocol
6. **Blockchain Integration**: Immutable revocation audit trail

## Support

For issues or questions:
- Check existing revocations: `GET /revocations`
- Review statistics: `GET /revocations/stats`
- Verify implementation: Review `revocation.go` and `revocation_endpoints.go`

---

**Version**: 1.0.0
**Last Updated**: 2025-10-31
**Author**: HIBE Development Team
