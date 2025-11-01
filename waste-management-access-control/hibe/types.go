package hibe

import (
	"math/big"
	"sync"
	"time"
)

// SystemParams contains the system-wide parameters for HIBE
type SystemParams struct {
	// Pairing parameters
	P    *big.Int // Large prime
	Q    *big.Int // Order of subgroup  
	G1   *big.Int // Generator of G1
	G2   *big.Int // Generator of G2
	GT   *big.Int // Target group
	
	// System parameters
	MaxDepth    int
	SecurityLevel int
	
	// Memory pools for optimization
	KeyPool     *sync.Pool
	BigIntPool  *sync.Pool
}

// MasterKey represents the master secret key
type MasterKey struct {
	Alpha *big.Int
	Beta  *big.Int
	Gamma *big.Int
	params *SystemParams
	mu    sync.RWMutex
}

// PublicKey represents the public parameters
type PublicKey struct {
	G       *big.Int
	G1      *big.Int  
	G2      *big.Int
	U       []*big.Int // U_1, U_2, ..., U_maxdepth
	H       []*big.Int // H_1, H_2, ..., H_maxdepth
	Params  *SystemParams
}

// PrivateKey represents a hierarchical private key
type PrivateKey struct {
	Components []*big.Int
	Depth      int
	Identity   []string
	IsWildcard []bool
	Timestamp  time.Time
	mu         sync.RWMutex
}

// WasteManagementPattern represents a waste-management access pattern
type WasteManagementPattern struct {
	Components    []string
	WildcardMask  []bool
	Depth         int
	PatternType   string
	Facility      string
	Department    string
	BinID     string
	DataType      string
	AccessLevel   string
}

// KeyCache for caching frequently used keys
type KeyCache struct {
	cache    map[string]*PrivateKey
	maxSize  int
	mu       sync.RWMutex
	stats    *CacheStats
}

type CacheStats struct {
	Hits      int64
	Misses    int64
	Evictions int64
	mu        sync.RWMutex
}

// Performance metrics
type KeyGenMetrics struct {
	TotalOperations    int64
	TotalDuration     time.Duration
	AverageDuration   time.Duration
	MinDuration       time.Duration
	MaxDuration       time.Duration
	MemoryAllocated   int64
	CacheHitRate      float64
	mu                sync.RWMutex
}

// NewSystemParams initializes system parameters
func NewSystemParams(maxDepth, securityLevel int) *SystemParams {
	return &SystemParams{
		MaxDepth:      maxDepth,
		SecurityLevel: securityLevel,
		KeyPool: &sync.Pool{
			New: func() interface{} {
				return &PrivateKey{
					Components: make([]*big.Int, maxDepth),
				}
			},
		},
		BigIntPool: &sync.Pool{
			New: func() interface{} {
				return new(big.Int)
			},
		},
	}
}

// NewKeyCache creates a new key cache
func NewKeyCache(maxSize int) *KeyCache {
	return &KeyCache{
		cache:   make(map[string]*PrivateKey),
		maxSize: maxSize,
		stats:   &CacheStats{},
	}
}

// Get retrieves a key from cache
func (kc *KeyCache) Get(key string) (*PrivateKey, bool) {
	kc.mu.RLock()
	defer kc.mu.RUnlock()
	
	if privKey, exists := kc.cache[key]; exists {
		kc.stats.mu.Lock()
		kc.stats.Hits++
		kc.stats.mu.Unlock()
		return privKey, true
	}
	
	kc.stats.mu.Lock()
	kc.stats.Misses++
	kc.stats.mu.Unlock()
	return nil, false
}

// Put stores a key in cache
func (kc *KeyCache) Put(key string, privKey *PrivateKey) {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	
	if len(kc.cache) >= kc.maxSize {
		// Simple LRU eviction - remove oldest
		for k := range kc.cache {
			delete(kc.cache, k)
			kc.stats.mu.Lock()
			kc.stats.Evictions++
			kc.stats.mu.Unlock()
			break
		}
	}
	
	kc.cache[key] = privKey
}

// GetStats returns cache statistics
func (kc *KeyCache) GetStats() CacheStats {
	kc.stats.mu.RLock()
	defer kc.stats.mu.RUnlock()
	return *kc.stats
}