package hibe

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// HIBEKeyGenerator manages hierarchical key generation
type HIBEKeyGenerator struct {
	MasterKey *MasterKey
	PublicKey *PublicKey
	Cache     *KeyCache
	Metrics   *KeyGenMetrics
}

// NewHIBEKeyGenerator creates a new key generator
func NewHIBEKeyGenerator(params *SystemParams) (*HIBEKeyGenerator, error) {
	// Generate master key
	masterKey, err := generateMasterKey(params)
	if err != nil {
		return nil, err
	}
	
	// Generate public key
	publicKey, err := generatePublicKey(masterKey, params)
	if err != nil {
		return nil, err
	}
	
	return &HIBEKeyGenerator{
		MasterKey: masterKey,
		PublicKey: publicKey,
		Cache:     NewKeyCache(1000), // Cache up to 1000 keys
		Metrics:   &KeyGenMetrics{},
	}, nil
}

// GenerateWasteManagementKey generates a private key for waste-management access patterns
func (kg *HIBEKeyGenerator) GenerateWasteManagementKey(pattern *WasteManagementPattern) (*PrivateKey, time.Duration, error) {
	start := time.Now()
	
	// Check cache first
	cacheKey := kg.buildCacheKey(pattern)
	if cachedKey, found := kg.Cache.Get(cacheKey); found {
		duration := time.Since(start)
		kg.updateMetrics(duration, 0, true)
		return cachedKey, duration, nil
	}
	
	// Generate new key with wildcard optimization
	key, err := kg.generateOptimizedKey(pattern)
	if err != nil {
		return nil, 0, err
	}
	
	// Cache the result
	kg.Cache.Put(cacheKey, key)
	
	duration := time.Since(start)
	memoryUsed := kg.calculateMemoryUsage(key)
	kg.updateMetrics(duration, memoryUsed, false)
	
	return key, duration, nil
}

// generateOptimizedKey creates an optimized key based on wildcard patterns
func (kg *HIBEKeyGenerator) generateOptimizedKey(pattern *WasteManagementPattern) (*PrivateKey, error) {
	// Get key from pool
	key := kg.MasterKey.params.KeyPool.Get().(*PrivateKey)
	
	// Reset key properties
	key.Components = key.Components[:0]
	key.Identity = make([]string, 0, len(pattern.Components))
	key.IsWildcard = make([]bool, 0, len(pattern.WildcardMask))
	key.Timestamp = time.Now()
	
	// Algorithm 3: Optimized Key Generation for WasteManagement Patterns
	// Only process non-wildcard components
	activeComponents := 0
	for i, component := range pattern.Components {
		key.Identity = append(key.Identity, component)
		key.IsWildcard = append(key.IsWildcard, pattern.WildcardMask[i])
		
		if !pattern.WildcardMask[i] && component != "" {
			// Generate component key only for non-wildcard elements
			componentKey, err := kg.generateComponentKey(component, activeComponents)
			if err != nil {
				return nil, err
			}
			key.Components = append(key.Components, componentKey)
			activeComponents++
		} else {
			// For wildcards, use nil pointer (memory optimization)
			key.Components = append(key.Components, nil)
		}
	}
	
	key.Depth = activeComponents // Only count non-wildcard components
	
	return key, nil
}

// generateComponentKey generates a key for a single component
func (kg *HIBEKeyGenerator) generateComponentKey(component string, depth int) (*big.Int, error) {
	// Get big.Int from pool
	result := kg.MasterKey.params.BigIntPool.Get().(*big.Int)
	
	// Fast hash-based key generation optimized for waste-management components
	hash := kg.fastWasteManagementHash(component, depth)
	result.SetBytes(hash)
	
	// Apply HIBE key generation formula
	// K_component = g^(alpha + H(component) * beta) * u_depth^r
	kg.MasterKey.mu.RLock()
	alpha := kg.MasterKey.Alpha
	beta := kg.MasterKey.Beta
	kg.MasterKey.mu.RUnlock()
	
	// r is a random value
	r, err := rand.Int(rand.Reader, kg.MasterKey.params.Q)
	if err != nil {
		return nil, err
	}
	
	// Compute: alpha + H(component) * beta
	temp := kg.MasterKey.params.BigIntPool.Get().(*big.Int)
	temp.Mul(result, beta)
	temp.Add(alpha, temp)
	temp.Mod(temp, kg.MasterKey.params.Q)
	
	// Apply randomization with u_depth^r
	if depth < len(kg.PublicKey.U) {
		uDepth := kg.PublicKey.U[depth]
		rPow := kg.MasterKey.params.BigIntPool.Get().(*big.Int)
		rPow.Exp(uDepth, r, kg.MasterKey.params.P)
		result.Mul(temp, rPow)
		result.Mod(result, kg.MasterKey.params.P)
		
		// Return rPow to pool
		kg.MasterKey.params.BigIntPool.Put(rPow)
	} else {
		result.Set(temp)
	}
	
	// Return temp to pool
	kg.MasterKey.params.BigIntPool.Put(temp)
	
	return result, nil
}

// fastWasteManagementHash provides optimized hashing for waste-management components
func (kg *HIBEKeyGenerator) fastWasteManagementHash(component string, depth int) []byte {
	hasher := sha256.New()
	
	// Add depth for unique hashing per level
	hasher.Write([]byte{byte(depth)})
	
	// WasteManagement-specific optimizations
	switch {
	case component == "facility":
		hasher.Write([]byte{0x01}) // Fixed identifier for facility
	case component == "bin":
		hasher.Write([]byte{0x02}) // Fixed identifier for bin
	case isWasteManagementDepartment(component):
		hasher.Write([]byte{0x03}) // Department identifier
		hasher.Write([]byte(component))
	case isBinID(component):
		hasher.Write([]byte{0x04}) // Bin ID identifier
		hasher.Write([]byte(component))
	case isDataType(component):
		hasher.Write([]byte{0x05}) // Data type identifier
		hasher.Write([]byte(component))
	case isAccessLevel(component):
		hasher.Write([]byte{0x06}) // Access level identifier
		hasher.Write([]byte(component))
	default:
		hasher.Write([]byte(component))
	}
	
	return hasher.Sum(nil)
}

// Helper functions for waste-management component identification
func isWasteManagementDepartment(component string) bool {
	departments := map[string]bool{
		"cardiology": true,
		"neurology":  true,
		"oncology":   true,
		"emergency":  true,
		"general":    true,
	}
	return departments[component]
}

func isBinID(component string) bool {
	// Simple check for numeric bin IDs
	for _, r := range component {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(component) > 0 && len(component) <= 10
}

func isDataType(component string) bool {
	dataTypes := map[string]bool{
		"vitals":  true,
		"records": true,
		"imaging": true,
		"labs":    true,
	}
	return dataTypes[component]
}

func isAccessLevel(component string) bool {
	accessLevels := map[string]bool{
		"realtime":   true,
		"historical": true,
		"critical":   true,
		"routine":    true,
	}
	return accessLevels[component]
}

// buildCacheKey creates a cache key from the pattern
func (kg *HIBEKeyGenerator) buildCacheKey(pattern *WasteManagementPattern) string {
	var keyParts []string
	for i, component := range pattern.Components {
		if pattern.WildcardMask[i] {
			keyParts = append(keyParts, "*")
		} else {
			keyParts = append(keyParts, component)
		}
	}
	return strings.Join(keyParts, "/")
}

// calculateMemoryUsage estimates memory usage of a key
func (kg *HIBEKeyGenerator) calculateMemoryUsage(key *PrivateKey) int64 {
	size := int64(0)
	
	// Size of Components slice
	for _, comp := range key.Components {
		if comp != nil {
			size += int64(comp.BitLen()/8 + 1) // Approximate bytes
		}
	}
	
	// Size of Identity slice
	for _, id := range key.Identity {
		size += int64(len(id))
	}
	
	// Size of IsWildcard slice
	size += int64(len(key.IsWildcard))
	
	return size
}

// updateMetrics updates performance metrics
func (kg *HIBEKeyGenerator) updateMetrics(duration time.Duration, memoryUsed int64, cacheHit bool) {
	kg.Metrics.mu.Lock()
	defer kg.Metrics.mu.Unlock()
	
	kg.Metrics.TotalOperations++
	kg.Metrics.TotalDuration += duration
	kg.Metrics.MemoryAllocated += memoryUsed
	
	if kg.Metrics.MinDuration == 0 || duration < kg.Metrics.MinDuration {
		kg.Metrics.MinDuration = duration
	}
	if duration > kg.Metrics.MaxDuration {
		kg.Metrics.MaxDuration = duration
	}
	
	kg.Metrics.AverageDuration = kg.Metrics.TotalDuration / time.Duration(kg.Metrics.TotalOperations)
	
	// Update cache hit rate
	stats := kg.Cache.GetStats()
	total := stats.Hits + stats.Misses
	if total > 0 {
		kg.Metrics.CacheHitRate = float64(stats.Hits) / float64(total) * 100
	}
}

// GetMetrics returns current performance metrics
func (kg *HIBEKeyGenerator) GetMetrics() KeyGenMetrics {
	kg.Metrics.mu.RLock()
	defer kg.Metrics.mu.RUnlock()
	return *kg.Metrics
}

// generateMasterKey generates the master secret key
func generateMasterKey(params *SystemParams) (*MasterKey, error) {
	alpha, err := rand.Int(rand.Reader, params.Q)
	if err != nil {
		return nil, err
	}
	
	beta, err := rand.Int(rand.Reader, params.Q)
	if err != nil {
		return nil, err
	}
	
	gamma, err := rand.Int(rand.Reader, params.Q)
	if err != nil {
		return nil, err
	}
	
	return &MasterKey{
		Alpha:  alpha,
		Beta:   beta,
		Gamma:  gamma,
		params: params,
	}, nil
}

// generatePublicKey generates the public parameters
func generatePublicKey(masterKey *MasterKey, params *SystemParams) (*PublicKey, error) {
	// Generate U and H arrays for waste-management hierarchy
	U := make([]*big.Int, params.MaxDepth)
	H := make([]*big.Int, params.MaxDepth)
	
	for i := 0; i < params.MaxDepth; i++ {
		u, err := rand.Int(rand.Reader, params.P)
		if err != nil {
			return nil, err
		}
		U[i] = u
		
		h, err := rand.Int(rand.Reader, params.P)
		if err != nil {
			return nil, err
		}
		H[i] = h
	}
	
	return &PublicKey{
		G:      params.G1,
		G1:     params.G1,
		G2:     params.G2,
		U:      U,
		H:      H,
		Params: params,
	}, nil
}

// PrintMetrics displays performance metrics
func (kg *HIBEKeyGenerator) PrintMetrics() {
	metrics := kg.GetMetrics()
	fmt.Printf("\n=== HIBE Key Generation Metrics ===\n")
	fmt.Printf("Total Operations: %d\n", metrics.TotalOperations)
	fmt.Printf("Total Duration: %v\n", metrics.TotalDuration)
	fmt.Printf("Average Duration: %v\n", metrics.AverageDuration)
	fmt.Printf("Min Duration: %v\n", metrics.MinDuration)
	fmt.Printf("Max Duration: %v\n", metrics.MaxDuration)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemoryAllocated)
	fmt.Printf("Cache Hit Rate: %.2f%%\n", metrics.CacheHitRate)
}