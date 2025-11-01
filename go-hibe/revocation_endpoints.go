package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRevocationEndpoints registers all revocation-related endpoints
func RegisterRevocationEndpoints(r *gin.Engine) {

	// POST /revoke - Revoke a delegated key
	r.POST("/revoke", func(c *gin.Context) {
		var req RevocationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid request: %v", err),
			})
			return
		}

		// Create revocation entry from request
		entry, err := CreateRevocationFromRequest(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// Add to revocation list
		if err := globalRevocationList.RevokeKey(entry); err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Key revoked successfully",
			"keyId":    entry.KeyID,
			"revokedAt": entry.RevokedAt,
			"effectiveFrom": entry.EffectiveFrom,
			"effectiveUntil": entry.EffectiveUntil,
		})
	})

	// POST /revoke-by-uri - Revoke all keys for a specific URI
	r.POST("/revoke-by-uri", func(c *gin.Context) {
		var req struct {
			URI       string `json:"uri" binding:"required"`
			RevokedBy string `json:"revokedBy"`
			Reason    string `json:"reason" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid request: %v", err),
			})
			return
		}

		count, err := globalRevocationList.RevokeByURI(req.URI, req.RevokedBy, req.Reason)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"message":      fmt.Sprintf("Revoked %d key(s) for URI: %s", count, req.URI),
			"revokedCount": count,
			"uri":          req.URI,
		})
	})

	// GET /revoke/check/:keyId - Check if a specific key is revoked
	r.GET("/revoke/check/:keyId", func(c *gin.Context) {
		keyID := c.Param("keyId")

		isRevoked := globalRevocationList.IsKeyRevoked(keyID)

		response := gin.H{
			"keyId":     keyID,
			"isRevoked": isRevoked,
		}

		if isRevoked {
			// Get revocation details
			globalRevocationList.mu.RLock()
			entry := globalRevocationList.revocations[keyID]
			globalRevocationList.mu.RUnlock()

			if entry != nil {
				response["revocationDetails"] = entry
			}
		}

		c.JSON(http.StatusOK, response)
	})

	// POST /revoke/check - Check revocation for key parameters
	r.POST("/revoke/check", func(c *gin.Context) {
		var req struct {
			Hierarchy string `json:"hierarchy" binding:"required"`
			URI       string `json:"uri" binding:"required"`
			StartTime int64  `json:"startTime" binding:"required"`
			EndTime   int64  `json:"endTime" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid request: %v", err),
			})
			return
		}

		start := time.Unix(req.StartTime, 0)
		end := time.Unix(req.EndTime, 0)

		isRevoked, entry := globalRevocationList.CheckRevocation(
			[]byte(req.Hierarchy),
			req.URI,
			start,
			end,
		)

		response := gin.H{
			"isRevoked": isRevoked,
			"hierarchy": req.Hierarchy,
			"uri":       req.URI,
			"startTime": req.StartTime,
			"endTime":   req.EndTime,
		}

		if isRevoked && entry != nil {
			response["revocationDetails"] = entry
		}

		c.JSON(http.StatusOK, response)
	})

	// GET /revocations - List all revocations
	r.GET("/revocations", func(c *gin.Context) {
		status := c.DefaultQuery("status", "all") // all, active, expired, pending

		var entries []*RevocationEntry

		switch status {
		case "active":
			entries = globalRevocationList.GetActiveRevocations()
		case "all":
			entries = globalRevocationList.GetAllRevocations()
		default:
			entries = globalRevocationList.GetAllRevocations()
		}

		c.JSON(http.StatusOK, gin.H{
			"count":       len(entries),
			"status":      status,
			"revocations": entries,
		})
	})

	// GET /revocations/uri/:uri - Get revocations for specific URI
	r.GET("/revocations/uri/:uri", func(c *gin.Context) {
		uri := c.Param("uri")

		entries := globalRevocationList.GetRevocationsByURI(uri)

		c.JSON(http.StatusOK, gin.H{
			"uri":         uri,
			"count":       len(entries),
			"revocations": entries,
		})
	})

	// DELETE /revoke/:keyId - Clear/reinstate a revocation
	r.DELETE("/revoke/:keyId", func(c *gin.Context) {
		keyID := c.Param("keyId")

		if err := globalRevocationList.ClearRevocation(keyID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Revocation cleared for key: %s", keyID),
			"keyId":   keyID,
		})
	})

	// GET /revocations/stats - Get revocation statistics
	r.GET("/revocations/stats", func(c *gin.Context) {
		stats := globalRevocationList.GetStats()

		c.JSON(http.StatusOK, stats)
	})

	// POST /revocations/cleanup - Remove expired revocations
	r.POST("/revocations/cleanup", func(c *gin.Context) {
		removed := globalRevocationList.RemoveExpiredRevocations()

		c.JSON(http.StatusOK, gin.H{
			"success":        true,
			"message":        "Cleanup completed",
			"removedCount":   removed,
			"remainingCount": len(globalRevocationList.GetAllRevocations()),
		})
	})

	// GET /revoke/generate-key-id - Generate a key ID from parameters (utility endpoint)
	r.POST("/revoke/generate-key-id", func(c *gin.Context) {
		var req struct {
			Hierarchy string `json:"hierarchy" binding:"required"`
			URI       string `json:"uri" binding:"required"`
			StartTime int64  `json:"startTime" binding:"required"`
			EndTime   int64  `json:"endTime" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid request: %v", err),
			})
			return
		}

		start := time.Unix(req.StartTime, 0)
		end := time.Unix(req.EndTime, 0)

		keyID := GenerateKeyID([]byte(req.Hierarchy), req.URI, start, end)

		c.JSON(http.StatusOK, gin.H{
			"keyId":     keyID,
			"hierarchy": req.Hierarchy,
			"uri":       req.URI,
			"startTime": req.StartTime,
			"endTime":   req.EndTime,
		})
	})
}

// Helper function to integrate revocation checking into delegation endpoint
func checkAndRecordDelegation(hierarchy []byte, uri string, start, end time.Time) (string, error) {
	// Generate key ID for this delegation
	keyID := GenerateKeyID(hierarchy, uri, start, end)

	// Check if this key is revoked
	if globalRevocationList.IsKeyRevoked(keyID) {
		globalRevocationList.mu.RLock()
		entry := globalRevocationList.revocations[keyID]
		globalRevocationList.mu.RUnlock()

		return keyID, fmt.Errorf("cannot delegate: key is revoked (reason: %s)", entry.Reason)
	}

	return keyID, nil
}

// Middleware to check revocation for decrypt operations
func revocationCheckMiddleware(c *gin.Context) {
	// Only apply to decrypt endpoint
	if c.Request.URL.Path != "/decrypt" {
		c.Next()
		return
	}

	// This is a placeholder - actual implementation would need to extract
	// key information from the decrypt request and check revocation
	// For now, we'll just pass through
	c.Next()
}
