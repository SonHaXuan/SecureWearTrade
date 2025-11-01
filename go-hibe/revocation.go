package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// RevocationEntry represents a revoked key entry
type RevocationEntry struct {
	KeyID          string    `json:"keyId"`          // Unique identifier for the delegated key
	URI            string    `json:"uri"`            // The URI pattern that was delegated
	Hierarchy      string    `json:"hierarchy"`      // The hierarchy of the key
	RevokedAt      time.Time `json:"revokedAt"`      // When the key was revoked
	RevokedBy      string    `json:"revokedBy"`      // Who revoked the key (optional)
	Reason         string    `json:"reason"`         // Reason for revocation
	EffectiveFrom  time.Time `json:"effectiveFrom"`  // When revocation takes effect
	EffectiveUntil time.Time `json:"effectiveUntil"` // When revocation expires (optional, 0 means permanent)
}

// RevocationList manages all revoked keys
type RevocationList struct {
	mu          sync.RWMutex
	revocations map[string]*RevocationEntry // KeyID -> RevocationEntry
	uriIndex    map[string][]string         // URI -> []KeyID for faster lookup
}

// Global revocation list instance
var globalRevocationList *RevocationList

// Initialize the revocation list
func init() {
	globalRevocationList = NewRevocationList()
}

// NewRevocationList creates a new revocation list
func NewRevocationList() *RevocationList {
	return &RevocationList{
		revocations: make(map[string]*RevocationEntry),
		uriIndex:    make(map[string][]string),
	}
}

// GenerateKeyID creates a unique identifier for a delegated key
// based on hierarchy, URI, and timestamp
func GenerateKeyID(hierarchy []byte, uri string, start time.Time, end time.Time) string {
	data := fmt.Sprintf("%s:%s:%d:%d", hierarchy, uri, start.Unix(), end.Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// RevokeKey adds a key to the revocation list
func (rl *RevocationList) RevokeKey(entry *RevocationEntry) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if already revoked
	if _, exists := rl.revocations[entry.KeyID]; exists {
		return fmt.Errorf("key %s is already revoked", entry.KeyID)
	}

	// Add to revocations map
	rl.revocations[entry.KeyID] = entry

	// Add to URI index for faster lookup
	rl.uriIndex[entry.URI] = append(rl.uriIndex[entry.URI], entry.KeyID)

	return nil
}

// IsKeyRevoked checks if a specific key is currently revoked
func (rl *RevocationList) IsKeyRevoked(keyID string) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	entry, exists := rl.revocations[keyID]
	if !exists {
		return false
	}

	now := time.Now()

	// Check if revocation is currently effective
	if now.Before(entry.EffectiveFrom) {
		return false // Revocation not yet effective
	}

	// Check if revocation has expired (0 means permanent)
	if !entry.EffectiveUntil.IsZero() && now.After(entry.EffectiveUntil) {
		return false // Revocation has expired
	}

	return true
}

// CheckRevocation checks if a key should be revoked based on URI and hierarchy
func (rl *RevocationList) CheckRevocation(hierarchy []byte, uri string, start time.Time, end time.Time) (bool, *RevocationEntry) {
	keyID := GenerateKeyID(hierarchy, uri, start, end)

	rl.mu.RLock()
	defer rl.mu.RUnlock()

	entry, exists := rl.revocations[keyID]
	if !exists {
		return false, nil
	}

	if rl.IsKeyRevoked(keyID) {
		return true, entry
	}

	return false, nil
}

// GetRevocationsByURI retrieves all revocations for a specific URI
func (rl *RevocationList) GetRevocationsByURI(uri string) []*RevocationEntry {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	keyIDs, exists := rl.uriIndex[uri]
	if !exists {
		return nil
	}

	var entries []*RevocationEntry
	for _, keyID := range keyIDs {
		if entry, exists := rl.revocations[keyID]; exists {
			entries = append(entries, entry)
		}
	}

	return entries
}

// RemoveExpiredRevocations removes revocations that are no longer effective
func (rl *RevocationList) RemoveExpiredRevocations() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	removed := 0

	for keyID, entry := range rl.revocations {
		// Remove if revocation has an expiry and it has passed
		if !entry.EffectiveUntil.IsZero() && now.After(entry.EffectiveUntil) {
			delete(rl.revocations, keyID)

			// Remove from URI index
			uriKeyIDs := rl.uriIndex[entry.URI]
			for i, id := range uriKeyIDs {
				if id == keyID {
					rl.uriIndex[entry.URI] = append(uriKeyIDs[:i], uriKeyIDs[i+1:]...)
					break
				}
			}

			removed++
		}
	}

	return removed
}

// GetAllRevocations returns all current revocations
func (rl *RevocationList) GetAllRevocations() []*RevocationEntry {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	entries := make([]*RevocationEntry, 0, len(rl.revocations))
	for _, entry := range rl.revocations {
		entries = append(entries, entry)
	}

	return entries
}

// GetActiveRevocations returns only currently active revocations
func (rl *RevocationList) GetActiveRevocations() []*RevocationEntry {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	var entries []*RevocationEntry

	for keyID, entry := range rl.revocations {
		// Check if currently effective
		if now.After(entry.EffectiveFrom) || now.Equal(entry.EffectiveFrom) {
			// Check if not expired
			if entry.EffectiveUntil.IsZero() || now.Before(entry.EffectiveUntil) {
				entries = append(entries, entry)
			}
		}
	}

	return entries
}

// RevokeByURI revokes all keys associated with a specific URI
func (rl *RevocationList) RevokeByURI(uri string, revokedBy string, reason string) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	keyIDs, exists := rl.uriIndex[uri]
	if !exists {
		return 0, fmt.Errorf("no keys found for URI: %s", uri)
	}

	revokedCount := 0
	now := time.Now()

	for _, keyID := range keyIDs {
		if entry, exists := rl.revocations[keyID]; exists {
			// Update existing revocation
			entry.RevokedAt = now
			entry.RevokedBy = revokedBy
			entry.Reason = reason
			entry.EffectiveFrom = now
		} else {
			// Create new revocation entry
			newEntry := &RevocationEntry{
				KeyID:         keyID,
				URI:           uri,
				RevokedAt:     now,
				RevokedBy:     revokedBy,
				Reason:        reason,
				EffectiveFrom: now,
			}
			rl.revocations[keyID] = newEntry
			revokedCount++
		}
	}

	return revokedCount, nil
}

// ClearRevocation removes a revocation entry (reinstate a key)
func (rl *RevocationList) ClearRevocation(keyID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.revocations[keyID]
	if !exists {
		return fmt.Errorf("revocation not found for key: %s", keyID)
	}

	// Remove from revocations map
	delete(rl.revocations, keyID)

	// Remove from URI index
	uriKeyIDs := rl.uriIndex[entry.URI]
	for i, id := range uriKeyIDs {
		if id == keyID {
			rl.uriIndex[entry.URI] = append(uriKeyIDs[:i], uriKeyIDs[i+1:]...)
			break
		}
	}

	return nil
}

// RevocationStats provides statistics about the revocation list
type RevocationStats struct {
	TotalRevocations   int       `json:"totalRevocations"`
	ActiveRevocations  int       `json:"activeRevocations"`
	ExpiredRevocations int       `json:"expiredRevocations"`
	PendingRevocations int       `json:"pendingRevocations"`
	UniqueURIs         int       `json:"uniqueUris"`
	LastUpdated        time.Time `json:"lastUpdated"`
}

// GetStats returns statistics about the revocation list
func (rl *RevocationList) GetStats() RevocationStats {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	stats := RevocationStats{
		TotalRevocations: len(rl.revocations),
		UniqueURIs:       len(rl.uriIndex),
		LastUpdated:      now,
	}

	for _, entry := range rl.revocations {
		if now.Before(entry.EffectiveFrom) {
			stats.PendingRevocations++
		} else if !entry.EffectiveUntil.IsZero() && now.After(entry.EffectiveUntil) {
			stats.ExpiredRevocations++
		} else {
			stats.ActiveRevocations++
		}
	}

	return stats
}

// ValidateRevocationRequest validates a revocation request
type RevocationRequest struct {
	KeyID         string    `json:"keyId,omitempty"`
	URI           string    `json:"uri,omitempty"`
	Hierarchy     string    `json:"hierarchy,omitempty"`
	RevokedBy     string    `json:"revokedBy"`
	Reason        string    `json:"reason"`
	EffectiveFrom time.Time `json:"effectiveFrom,omitempty"`
	EffectiveFor  int64     `json:"effectiveFor,omitempty"` // Duration in seconds, 0 means permanent
	StartTime     int64     `json:"startTime,omitempty"`    // Unix timestamp for key delegation start
	EndTime       int64     `json:"endTime,omitempty"`      // Unix timestamp for key delegation end
}

// ValidateRevocationRequest validates the revocation request
func ValidateRevocationRequest(req *RevocationRequest) error {
	if req.KeyID == "" && req.URI == "" {
		return fmt.Errorf("either keyId or uri must be provided")
	}

	if req.Reason == "" {
		return fmt.Errorf("reason is required")
	}

	if req.EffectiveFrom.IsZero() {
		req.EffectiveFrom = time.Now()
	}

	return nil
}

// CreateRevocationFromRequest creates a RevocationEntry from a request
func CreateRevocationFromRequest(req *RevocationRequest) (*RevocationEntry, error) {
	if err := ValidateRevocationRequest(req); err != nil {
		return nil, err
	}

	entry := &RevocationEntry{
		KeyID:         req.KeyID,
		URI:           req.URI,
		Hierarchy:     req.Hierarchy,
		RevokedAt:     time.Now(),
		RevokedBy:     req.RevokedBy,
		Reason:        req.Reason,
		EffectiveFrom: req.EffectiveFrom,
	}

	// If KeyID not provided, generate it from other parameters
	if entry.KeyID == "" && entry.URI != "" && req.StartTime != 0 && req.EndTime != 0 {
		start := time.Unix(req.StartTime, 0)
		end := time.Unix(req.EndTime, 0)
		entry.KeyID = GenerateKeyID([]byte(req.Hierarchy), req.URI, start, end)
	}

	// Set expiration if duration is specified
	if req.EffectiveFor > 0 {
		entry.EffectiveUntil = entry.EffectiveFrom.Add(time.Duration(req.EffectiveFor) * time.Second)
	}

	return entry, nil
}

// Context-aware revocation checking
func CheckRevocationWithContext(ctx context.Context, rl *RevocationList, keyID string) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if rl.IsKeyRevoked(keyID) {
		entry := rl.revocations[keyID]
		return fmt.Errorf("key revoked: %s (reason: %s, revoked by: %s)",
			entry.KeyID, entry.Reason, entry.RevokedBy)
	}

	return nil
}
