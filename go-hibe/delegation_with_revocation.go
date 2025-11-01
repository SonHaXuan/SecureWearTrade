package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"hibe-api"
)

// DelegationRequest represents a request to delegate a key
type DelegationRequest struct {
	URI       string `json:"uri" binding:"required"`
	Hierarchy string `json:"hierarchy"`
	StartTime int64  `json:"startTime" binding:"required"` // Unix timestamp
	EndTime   int64  `json:"endTime" binding:"required"`   // Unix timestamp
	Parent    string `json:"parent,omitempty"`             // For hierarchical delegation
}

// DelegationResponse represents the response from a delegation request
type DelegationResponse struct {
	Success       bool      `json:"success"`
	KeyID         string    `json:"keyId"`
	Data          []byte    `json:"data"`
	URI           string    `json:"uri"`
	Hierarchy     string    `json:"hierarchy"`
	StartTime     int64     `json:"startTime"`
	EndTime       int64     `json:"endTime"`
	ExecutionTime int64     `json:"executionTime"` // in microseconds
	Message       string    `json:"message,omitempty"`
	Error         string    `json:"error,omitempty"`
}

// RegisterDelegationWithRevocationEndpoint adds the enhanced delegation endpoint
func RegisterDelegationWithRevocationEndpoint(r *gin.Engine, ctx context.Context, store *TestKeyStore, encoder hibe.PatternEncoder) {

	// POST /hibe-delegate - New endpoint for key delegation with revocation support
	r.POST("/hibe-delegate", func(c *gin.Context) {
		var req DelegationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, DelegationResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid request: %v", err),
			})
			return
		}

		// Use default hierarchy if not provided
		hierarchy := TestHierarchy
		if req.Hierarchy != "" {
			hierarchy = []byte(req.Hierarchy)
		}

		// Convert timestamps to time.Time
		start := time.Unix(req.StartTime, 0)
		end := time.Unix(req.EndTime, 0)

		// Validate time range
		if end.Before(start) {
			c.JSON(400, DelegationResponse{
				Success: false,
				Error:   "endTime must be after startTime",
			})
			return
		}

		// Check if this delegation would create a revoked key
		keyID, err := checkAndRecordDelegation(hierarchy, req.URI, start, end)
		if err != nil {
			c.JSON(403, DelegationResponse{
				Success: false,
				KeyID:   keyID,
				Error:   err.Error(),
			})
			return
		}

		// Perform the actual delegation
		startExecution := time.Now()
		delegation, err := hibe.Delegate(ctx, store, encoder, hierarchy, req.URI, start, end, hibe.DecryptPermission|hibe.SignPermission)
		if err != nil {
			c.JSON(500, DelegationResponse{
				Success: false,
				KeyID:   keyID,
				Error:   fmt.Sprintf("Delegation failed: %v", err),
			})
			return
		}
		executionTime := time.Since(startExecution).Microseconds()

		// Marshal the delegation
		marshalled := delegation.Marshal()

		// Return successful response
		c.JSON(200, DelegationResponse{
			Success:       true,
			KeyID:         keyID,
			Data:          marshalled,
			URI:           req.URI,
			Hierarchy:     string(hierarchy),
			StartTime:     req.StartTime,
			EndTime:       req.EndTime,
			ExecutionTime: executionTime,
			Message:       "Key delegated successfully",
		})
	})

	// GET /hibe-delegate-info/:keyId - Get delegation information for a key ID
	r.GET("/hibe-delegate-info/:keyId", func(c *gin.Context) {
		keyID := c.Param("keyId")

		// Check if key is revoked
		isRevoked := globalRevocationList.IsKeyRevoked(keyID)

		response := gin.H{
			"keyId":     keyID,
			"isRevoked": isRevoked,
		}

		if isRevoked {
			globalRevocationList.mu.RLock()
			entry := globalRevocationList.revocations[keyID]
			globalRevocationList.mu.RUnlock()

			if entry != nil {
				response["revocationDetails"] = entry
				response["status"] = "revoked"
			}
		} else {
			response["status"] = "active"
		}

		c.JSON(200, response)
	})
}

// DecryptWithRevocationCheck performs decryption with revocation checking
func DecryptWithRevocationCheck(
	ctx context.Context,
	state *hibe.ClientState,
	hierarchy []byte,
	uri string,
	timestamp time.Time,
	encrypted []byte,
	startTime time.Time,
	endTime time.Time,
) ([]byte, error) {

	// Generate key ID for revocation check
	keyID := GenerateKeyID(hierarchy, uri, startTime, endTime)

	// Check if key is revoked
	if err := CheckRevocationWithContext(ctx, globalRevocationList, keyID); err != nil {
		return nil, fmt.Errorf("decryption denied: %v", err)
	}

	// Perform actual decryption
	decrypted, err := state.Decrypt(ctx, hierarchy, uri, timestamp, encrypted)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	return decrypted, nil
}

// Enhanced decrypt endpoint with revocation checking
func RegisterEnhancedDecryptEndpoint(r *gin.Engine, ctx context.Context, state *hibe.ClientState, now time.Time) {

	r.POST("/decrypt-with-revocation", func(c *gin.Context) {
		var req struct {
			URI              string `json:"uri" binding:"required"`
			EncryptedMessage string `json:"encryptedMessage" binding:"required"`
			Hierarchy        string `json:"hierarchy"`
			StartTime        int64  `json:"startTime" binding:"required"`
			EndTime          int64  `json:"endTime" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid request: %v", err),
			})
			return
		}

		// Decode encrypted message
		encrypted, err := decodeBase64(req.EncryptedMessage)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid encrypted message: %v", err),
			})
			return
		}

		// Use default hierarchy if not provided
		hierarchy := TestHierarchy
		if req.Hierarchy != "" {
			hierarchy = []byte(req.Hierarchy)
		}

		start := time.Unix(req.StartTime, 0)
		end := time.Unix(req.EndTime, 0)

		// Perform decryption with revocation check
		startExecution := time.Now()
		decrypted, err := DecryptWithRevocationCheck(
			ctx,
			state,
			hierarchy,
			req.URI,
			now,
			encrypted,
			start,
			end,
		)
		executionTime := time.Since(startExecution).Microseconds()

		if err != nil {
			c.JSON(403, gin.H{
				"success":       false,
				"error":         err.Error(),
				"executionTime": executionTime,
			})
			return
		}

		// Measure resource usage
		measureUsage := measureMemoryUsage()
		measureUsage.executionTime = executionTime

		// Calculate energy consumption
		executionTimeSeconds := float64(executionTime) / 1000000.0
		measureUsage.energyConsumption = calculateEnergyConsumption(measureUsage.powerUsage, executionTimeSeconds)

		// Record power usage
		recordPowerUsage("decrypt-with-revocation", measureUsage, len(encrypted))

		c.JSON(200, gin.H{
			"success":                true,
			"data":                   string(decrypted),
			"executionTime":          executionTime,
			"memoryUsage":            measureUsage.memory,
			"cpuPercentage":          measureUsage.cpuPercentage,
			"powerUsageWatts":        measureUsage.powerUsage,
			"energyConsumptionJoules": measureUsage.energyConsumption,
		})
	})
}

// Helper function to decode base64
func decodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// DelegationInfo stores information about active delegations
type DelegationInfo struct {
	KeyID         string    `json:"keyId"`
	URI           string    `json:"uri"`
	Hierarchy     string    `json:"hierarchy"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUsed      time.Time `json:"lastUsed,omitempty"`
	UsageCount    int       `json:"usageCount"`
	IsRevoked     bool      `json:"isRevoked"`
}

// DelegationRegistry keeps track of active delegations
type DelegationRegistry struct {
	mu          sync.RWMutex
	delegations map[string]*DelegationInfo
}

var globalDelegationRegistry = &DelegationRegistry{
	delegations: make(map[string]*DelegationInfo),
}

// RecordDelegation records a new delegation
func (dr *DelegationRegistry) RecordDelegation(info *DelegationInfo) {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	dr.delegations[info.KeyID] = info
}

// UpdateUsage updates the usage statistics for a delegation
func (dr *DelegationRegistry) UpdateUsage(keyID string) {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	if info, exists := dr.delegations[keyID]; exists {
		info.LastUsed = time.Now()
		info.UsageCount++
	}
}

// GetDelegation retrieves delegation information
func (dr *DelegationRegistry) GetDelegation(keyID string) (*DelegationInfo, bool) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	info, exists := dr.delegations[keyID]
	return info, exists
}

// GetAllDelegations returns all registered delegations
func (dr *DelegationRegistry) GetAllDelegations() []*DelegationInfo {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	delegations := make([]*DelegationInfo, 0, len(dr.delegations))
	for _, info := range dr.delegations {
		delegations = append(delegations, info)
	}

	return delegations
}

// RegisterDelegationManagementEndpoints adds delegation management endpoints
func RegisterDelegationManagementEndpoints(r *gin.Engine) {

	// GET /delegations - List all delegations
	r.GET("/delegations", func(c *gin.Context) {
		delegations := globalDelegationRegistry.GetAllDelegations()

		// Update revocation status
		for _, delegation := range delegations {
			delegation.IsRevoked = globalRevocationList.IsKeyRevoked(delegation.KeyID)
		}

		c.JSON(200, gin.H{
			"count":       len(delegations),
			"delegations": delegations,
		})
	})

	// GET /delegations/:keyId - Get specific delegation info
	r.GET("/delegations/:keyId", func(c *gin.Context) {
		keyID := c.Param("keyId")

		info, exists := globalDelegationRegistry.GetDelegation(keyID)
		if !exists {
			c.JSON(404, gin.H{
				"error": "Delegation not found",
			})
			return
		}

		// Update revocation status
		info.IsRevoked = globalRevocationList.IsKeyRevoked(keyID)

		c.JSON(200, info)
	})
}
