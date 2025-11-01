# Integration Instructions for Key Revocation System

## Overview

This document provides step-by-step instructions for integrating the key revocation system into your existing HIBE API (`main.go`).

## Files Added

The revocation system consists of three new files:

1. **`revocation.go`** - Core revocation list and data structures
2. **`revocation_endpoints.go`** - API endpoints for revocation management
3. **`delegation_with_revocation.go`** - Enhanced delegation and decryption with revocation checking

## Integration Steps

### Step 1: Add Revocation Endpoints to main.go

Add the following code in your `main()` function, after you create the Gin router (`r := gin.Default()`):

```go
func main() {
	ctx := context.Background()
	_, store := NewTestKeyStore()
	encoder := hibe.NewDefaultPatternEncoder(TestPatternSize - hibe.MaxTimeLength)

	state := NewTestState()
	now := time.Now()

	r := gin.Default()

	// Add security headers middleware
	r.Use(securityHeaders())

	// ===== ADD THESE LINES =====
	// Register revocation management endpoints
	RegisterRevocationEndpoints(r)

	// Register enhanced delegation endpoints with revocation support
	RegisterDelegationWithRevocationEndpoint(r, ctx, store, encoder)

	// Register enhanced decrypt endpoint with revocation checking
	RegisterEnhancedDecryptEndpoint(r, ctx, state, now)

	// Register delegation tracking endpoints
	RegisterDelegationManagementEndpoints(r)
	// ===== END NEW LINES =====

	// ... rest of your existing endpoints ...

	r.Run() // listen and serve on 0.0.0.0:8080
}
```

### Step 2: Update Existing Delegation Endpoint (Optional)

If you want to add revocation checking to the existing `/hibe-private-key` endpoint, modify it as follows:

**Before:**
```go
r.GET("/hibe-private-key", func(c *gin.Context) {
	uri := "a/b/c"
	start := time.Unix(1565119330, 0)
	end := time.Unix(1565219330, 0)

	delegation, err := hibe.Delegate(ctx, store, encoder, TestHierarchy, uri, start, end, hibe.DecryptPermission|hibe.SignPermission)
	// ... rest of code
})
```

**After:**
```go
r.GET("/hibe-private-key", func(c *gin.Context) {
	uri := "a/b/c"
	start := time.Unix(1565119330, 0)
	end := time.Unix(1565219330, 0)

	// ADD: Check revocation before delegating
	keyID, err := checkAndRecordDelegation(TestHierarchy, uri, start, end)
	if err != nil {
		c.JSON(403, gin.H{
			"error": err.Error(),
			"keyId": keyID,
		})
		return
	}

	delegation, err := hibe.Delegate(ctx, store, encoder, TestHierarchy, uri, start, end, hibe.DecryptPermission|hibe.SignPermission)
	// ... rest of code
})
```

### Step 3: Update Existing Decrypt Endpoint (Optional)

To add revocation checking to the existing `/decrypt` endpoint:

**Before:**
```go
r.POST("/decrypt", func(c *gin.Context) {
	// ... existing code ...

	var decrypted []byte
	if decrypted, err = state.Decrypt(ctx, TestHierarchy, uri, now, encrypted); err != nil {
		// handle error
	}
	// ... rest of code
})
```

**After:**
```go
r.POST("/decrypt", func(c *gin.Context) {
	var decryptRequest DecryptRequest
	c.BindJSON(&decryptRequest)

	encrypted, err := base64.StdEncoding.DecodeString(decryptRequest.ENCRYPTEDMESSAGE)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uri := decryptRequest.URI

	// ADD: Get key info (you may need to modify DecryptRequest to include these)
	// For now, we'll assume you have these values
	start := time.Unix(1565119330, 0) // Get from request
	end := time.Unix(1565219330, 0)   // Get from request

	// ADD: Check revocation before decrypting
	keyID := GenerateKeyID(TestHierarchy, uri, start, end)
	if err := CheckRevocationWithContext(ctx, globalRevocationList, keyID); err != nil {
		c.JSON(403, gin.H{
			"error": fmt.Sprintf("Decryption denied: %v", err),
		})
		return
	}

	var decrypted []byte
	if decrypted, err = state.Decrypt(ctx, TestHierarchy, uri, now, encrypted); err != nil {
		// handle error
	}
	// ... rest of code
})
```

### Step 4: Update DecryptRequest Structure

To support revocation checking in decrypt operations, update the `DecryptRequest` structure:

```go
type DecryptRequest struct {
	URI              string `json:"uri" binding:"required"`
	ENCRYPTEDMESSAGE string `json:"encryptedMessage" binding:"required"`
	KEY              string `json:"key" binding:"required"`
	// ADD these fields for revocation checking
	StartTime        int64  `json:"startTime,omitempty"`
	EndTime          int64  `json:"endTime,omitempty"`
	Hierarchy        string `json:"hierarchy,omitempty"`
}
```

## Available Endpoints After Integration

### Revocation Management
- `POST /revoke` - Revoke a specific key
- `POST /revoke-by-uri` - Revoke all keys for a URI
- `GET /revoke/check/:keyId` - Check if a key is revoked
- `POST /revoke/check` - Check revocation by parameters
- `GET /revocations` - List all revocations
- `GET /revocations/uri/:uri` - Get revocations for a URI
- `DELETE /revoke/:keyId` - Clear a revocation
- `GET /revocations/stats` - Get statistics
- `POST /revocations/cleanup` - Clean up expired revocations
- `POST /revoke/generate-key-id` - Generate key ID utility

### Enhanced Delegation
- `POST /hibe-delegate` - Delegate with revocation check
- `GET /hibe-delegate-info/:keyId` - Get delegation info
- `GET /delegations` - List all delegations
- `GET /delegations/:keyId` - Get specific delegation

### Enhanced Decryption
- `POST /decrypt-with-revocation` - Decrypt with revocation check

## Testing the Integration

### 1. Test Basic Revocation

```bash
# 1. Create a delegation
curl -X POST http://localhost:8080/hibe-delegate \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "test/data",
    "hierarchy": "testHierarchy",
    "startTime": 1565119330,
    "endTime": 1565219330
  }'

# Save the keyId from response

# 2. Check it's not revoked
curl http://localhost:8080/revoke/check/YOUR_KEY_ID

# 3. Revoke it
curl -X POST http://localhost:8080/revoke \
  -H "Content-Type: application/json" \
  -d '{
    "keyId": "YOUR_KEY_ID",
    "revokedBy": "test-admin",
    "reason": "Testing revocation system"
  }'

# 4. Verify it's revoked
curl http://localhost:8080/revoke/check/YOUR_KEY_ID

# 5. Try to delegate again (should fail)
curl -X POST http://localhost:8080/hibe-delegate \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "test/data",
    "hierarchy": "testHierarchy",
    "startTime": 1565119330,
    "endTime": 1565219330
  }'
```

### 2. Test URI-based Revocation

```bash
# Revoke all keys for a URI
curl -X POST http://localhost:8080/revoke-by-uri \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "test/data",
    "revokedBy": "admin",
    "reason": "URI access terminated"
  }'

# Check revocations for that URI
curl http://localhost:8080/revocations/uri/test/data
```

### 3. Test Statistics

```bash
# Get revocation stats
curl http://localhost:8080/revocations/stats

# List all revocations
curl http://localhost:8080/revocations

# List only active revocations
curl http://localhost:8080/revocations?status=active
```

### 4. Test Cleanup

```bash
# Clean up expired revocations
curl -X POST http://localhost:8080/revocations/cleanup
```

## Troubleshooting

### Import Issues

If you get import errors, ensure your `go.mod` includes:
```
module hibe-api

go 1.20

replace hibe => ./packages/hibe

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/ucbrise/hibe-pairing v0.0.0-20220312033002-c4bf151b8d2b
	hibe v0.0.0-00010101000000-000000000000
)
```

### Build Issues

If you encounter build errors:
```bash
# Clean and rebuild
go clean
go mod tidy
go build
```

### Runtime Issues

If endpoints are not accessible:
1. Check that all registration functions are called in `main()`
2. Verify Gin router is set up correctly
3. Check for port conflicts
4. Review logs for error messages

## Minimal Integration (Quick Start)

If you want just basic revocation without all features:

```go
func main() {
	// ... existing setup ...

	r := gin.Default()

	// Just add basic revocation endpoints
	RegisterRevocationEndpoints(r)

	// ... rest of your code ...
}
```

Then manually check revocations in your existing endpoints as needed.

## Complete Example main.go Snippet

Here's a complete example showing where to add the revocation system:

```go
func main() {
	ctx := context.Background()
	_, store := NewTestKeyStore()
	encoder := hibe.NewDefaultPatternEncoder(TestPatternSize - hibe.MaxTimeLength)
	state := NewTestState()
	now := time.Now()

	r := gin.Default()
	r.Use(securityHeaders())

	// ========== REVOCATION SYSTEM INTEGRATION ==========
	RegisterRevocationEndpoints(r)
	RegisterDelegationWithRevocationEndpoint(r, ctx, store, encoder)
	RegisterEnhancedDecryptEndpoint(r, ctx, state, now)
	RegisterDelegationManagementEndpoints(r)
	// ====================================================

	// Your existing endpoints
	r.GET("/hibe-private-key", func(c *gin.Context) {
		// ... existing code ...
	})

	r.POST("/encrypt", func(c *gin.Context) {
		// ... existing code ...
	})

	r.POST("/decrypt", func(c *gin.Context) {
		// ... existing code ...
	})

	// ... more endpoints ...

	r.Run() // listen and serve on 0.0.0.0:8080
	fmt.Println("DONE!")
}
```

## Next Steps

1. **Test thoroughly** - Run all test cases
2. **Add persistence** - Implement database backend for production
3. **Add authentication** - Secure revocation endpoints with proper auth
4. **Monitor** - Set up monitoring for revocation statistics
5. **Document** - Update your API documentation to include revocation endpoints

## References

- See `REVOCATION_GUIDE.md` for detailed API documentation
- See `revocation.go` for implementation details
- See `revocation_endpoints.go` for endpoint implementations
- See `delegation_with_revocation.go` for enhanced delegation features

---

**Need Help?**

If you encounter issues:
1. Check the error messages in console
2. Verify all files are in the correct location
3. Run `go mod tidy` to ensure dependencies are correct
4. Check that port 8080 is available
