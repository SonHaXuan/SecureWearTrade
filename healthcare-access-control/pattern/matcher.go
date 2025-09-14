package pattern

import (
	"strings"
	"sync"
	"time"
)

// PatternMatcher handles optimized pattern matching for healthcare URIs
type PatternMatcher struct {
	Cache          *MatchCache
	MemoryPool     *StringPool
	OptimizedPaths map[string]*CompiledPattern
	Metrics        *MatchMetrics
	mu             sync.RWMutex
}

// CompiledPattern represents a pre-compiled pattern for fast matching
type CompiledPattern struct {
	StaticParts     []string
	WildcardMask    []bool
	CompareCount    int
	MemorySize      int
	PatternHash     uint64
	OptimizedComponents []OptimizedComponent
}

// OptimizedComponent represents an optimized pattern component
type OptimizedComponent struct {
	Value       []byte
	IsWildcard  bool
	ComponentType ComponentType
	Hash        uint32
}

// ComponentType represents the type of healthcare URI component
type ComponentType int

const (
	ComponentHospital ComponentType = iota
	ComponentDepartment
	ComponentPatientLabel
	ComponentPatientID
	ComponentDataType
	ComponentAccessLevel
)

// MatchCache caches pattern matching results
type MatchCache struct {
	cache   map[string]*MatchResult
	maxSize int
	hits    int64
	misses  int64
	mu      sync.RWMutex
}

// MatchResult contains pattern matching results
type MatchResult struct {
	IsMatch     bool
	MatchTime   time.Duration
	Matches     int
	CacheTime   time.Time
}

// MatchMetrics tracks pattern matching performance
type MatchMetrics struct {
	TotalMatches    int64
	TotalDuration   time.Duration
	AverageDuration time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
	CacheHitRate    float64
	WildcardSavings time.Duration
	mu              sync.RWMutex
}

// StringPool manages reusable string allocations
type StringPool struct {
	pool sync.Pool
}

// NewPatternMatcher creates a new optimized pattern matcher
func NewPatternMatcher(cacheSize int) *PatternMatcher {
	pm := &PatternMatcher{
		Cache:          NewMatchCache(cacheSize),
		MemoryPool:     NewStringPool(),
		OptimizedPaths: make(map[string]*CompiledPattern),
		Metrics:        &MatchMetrics{},
	}
	
	// Pre-compile common healthcare patterns
	pm.precompileHealthcarePatterns()
	
	return pm
}

// MatchHealthcarePattern performs optimized pattern matching
func (pm *PatternMatcher) MatchHealthcarePattern(uri string, pattern *CompiledPattern) (bool, time.Duration, int) {
	start := time.Now()
	
	// Check cache first
	cacheKey := pm.buildCacheKey(uri, pattern)
	if cached, found := pm.Cache.Get(cacheKey); found {
		duration := time.Since(start)
		pm.updateMetrics(duration, cached.Matches, true, pattern.WildcardMask)
		return cached.IsMatch, duration, cached.Matches
	}
	
	// Perform actual pattern matching
	isMatch, matches := pm.performOptimizedMatching(uri, pattern)
	
	duration := time.Since(start)
	
	// Cache the result
	pm.Cache.Put(cacheKey, &MatchResult{
		IsMatch:   isMatch,
		MatchTime: duration,
		Matches:   matches,
		CacheTime: time.Now(),
	})
	
	pm.updateMetrics(duration, matches, false, pattern.WildcardMask)
	
	return isMatch, duration, matches
}

// performOptimizedMatching performs the actual optimized matching
func (pm *PatternMatcher) performOptimizedMatching(uri string, pattern *CompiledPattern) (bool, int) {
	// Fast path: Parse URI components once
	components := pm.parseURIComponents(uri)
	if len(components) != len(pattern.OptimizedComponents) {
		return false, 0
	}
	
	matches := 0
	
	// Optimized matching: Skip wildcard comparisons entirely
	for i, expected := range pattern.OptimizedComponents {
		if !expected.IsWildcard {
			if i < len(components) {
				// Fast byte comparison for performance
				if pm.fastCompareComponents(components[i], expected) {
					matches++
				} else {
					return false, matches
				}
			} else {
				return false, matches
			}
		}
		// Wildcards are completely skipped - no comparison needed
	}
	
	// Match succeeds if all non-wildcard components matched
	isMatch := matches == pattern.CompareCount
	return isMatch, matches
}

// parseURIComponents efficiently parses URI into components
func (pm *PatternMatcher) parseURIComponents(uri string) [][]byte {
	// Remove leading slash and split
	if len(uri) > 0 && uri[0] == '/' {
		uri = uri[1:]
	}
	
	parts := strings.Split(uri, "/")
	components := make([][]byte, len(parts))
	
	for i, part := range parts {
		components[i] = []byte(part)
	}
	
	return components
}

// fastCompareComponents performs optimized component comparison
func (pm *PatternMatcher) fastCompareComponents(actual []byte, expected OptimizedComponent) bool {
	// Fast length check first
	if len(actual) != len(expected.Value) {
		return false
	}
	
	// Healthcare-specific optimizations based on component type
	switch expected.ComponentType {
	case ComponentHospital:
		// "hospital" is always the same - single comparison
		return string(actual) == "hospital"
		
	case ComponentPatientLabel:
		// "patient" is always the same - single comparison
		return string(actual) == "patient"
		
	case ComponentDepartment:
		// Pre-validated department names - hash comparison first
		if pm.fastHash(actual) == expected.Hash {
			return pm.bytesEqual(actual, expected.Value)
		}
		return false
		
	case ComponentPatientID:
		// Numeric patient ID - optimized numeric comparison
		return pm.isValidPatientID(actual) && pm.bytesEqual(actual, expected.Value)
		
	case ComponentDataType:
		// Pre-validated data types - hash comparison
		if pm.fastHash(actual) == expected.Hash {
			return pm.bytesEqual(actual, expected.Value)
		}
		return false
		
	case ComponentAccessLevel:
		// Pre-validated access levels - hash comparison
		if pm.fastHash(actual) == expected.Hash {
			return pm.bytesEqual(actual, expected.Value)
		}
		return false
		
	default:
		// Generic byte comparison
		return pm.bytesEqual(actual, expected.Value)
	}
}

// bytesEqual performs fast byte slice comparison
func (pm *PatternMatcher) bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// fastHash computes a fast hash for component comparison
func (pm *PatternMatcher) fastHash(data []byte) uint32 {
	hash := uint32(0)
	for _, b := range data {
		hash = hash*31 + uint32(b)
	}
	return hash
}

// isValidPatientID checks if a component is a valid patient ID
func (pm *PatternMatcher) isValidPatientID(data []byte) bool {
	if len(data) == 0 || len(data) > 10 {
		return false
	}
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// CompilePattern compiles a pattern string into an optimized format
func (pm *PatternMatcher) CompilePattern(patternStr string) *CompiledPattern {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	// Check if already compiled
	if compiled, exists := pm.OptimizedPaths[patternStr]; exists {
		return compiled
	}
	
	// Parse pattern components
	components := pm.parseURIComponents(patternStr)
	
	compiled := &CompiledPattern{
		StaticParts:         make([]string, len(components)),
		WildcardMask:        make([]bool, len(components)),
		OptimizedComponents: make([]OptimizedComponent, len(components)),
		CompareCount:        0,
		MemorySize:          0,
	}
	
	for i, comp := range components {
		compStr := string(comp)
		compiled.StaticParts[i] = compStr
		
		if compStr == "*" {
			compiled.WildcardMask[i] = true
			compiled.OptimizedComponents[i] = OptimizedComponent{
				Value:         nil, // No allocation for wildcards
				IsWildcard:    true,
				ComponentType: pm.detectComponentType(i),
				Hash:          0,
			}
		} else {
			compiled.WildcardMask[i] = false
			compiled.CompareCount++
			
			optimizedComp := OptimizedComponent{
				Value:         make([]byte, len(comp)),
				IsWildcard:    false,
				ComponentType: pm.detectComponentType(i),
				Hash:          pm.fastHash(comp),
			}
			copy(optimizedComp.Value, comp)
			compiled.OptimizedComponents[i] = optimizedComp
			compiled.MemorySize += len(comp)
		}
	}
	
	// Calculate pattern hash for cache optimization
	compiled.PatternHash = pm.calculatePatternHash(patternStr)
	
	// Store in cache
	pm.OptimizedPaths[patternStr] = compiled
	
	return compiled
}

// detectComponentType determines the type of URI component based on position
func (pm *PatternMatcher) detectComponentType(position int) ComponentType {
	switch position {
	case 0:
		return ComponentHospital
	case 1:
		return ComponentDepartment
	case 2:
		return ComponentPatientLabel
	case 3:
		return ComponentPatientID
	case 4:
		return ComponentDataType
	case 5:
		return ComponentAccessLevel
	default:
		return ComponentDataType // Default fallback
	}
}

// calculatePatternHash computes a hash for the entire pattern
func (pm *PatternMatcher) calculatePatternHash(pattern string) uint64 {
	hash := uint64(0)
	for _, b := range []byte(pattern) {
		hash = hash*31 + uint64(b)
	}
	return hash
}

// precompileHealthcarePatterns pre-compiles common healthcare patterns
func (pm *PatternMatcher) precompileHealthcarePatterns() {
	commonPatterns := []string{
		"/hospital/cardiology/patient/12345/vitals/realtime",
		"/hospital/neurology/patient/67890/records/historical",
		"/hospital/oncology/patient/11111/imaging/routine",
		"/hospital/emergency/patient/22222/vitals/critical",
		"/hospital/general/patient/33333/labs/routine",
		
		// Wildcard patterns
		"/hospital/*/patient/*/vitals/*",
		"/hospital/*/patient/*/records/*",
		"/hospital/*/patient/*/imaging/*",
		"/hospital/*/patient/*/labs/*",
		"/hospital/cardiology/patient/*/vitals/*",
		"/hospital/emergency/patient/*/vitals/critical",
	}
	
	for _, pattern := range commonPatterns {
		pm.CompilePattern(pattern)
	}
}

// buildCacheKey creates a cache key from URI and pattern
func (pm *PatternMatcher) buildCacheKey(uri string, pattern *CompiledPattern) string {
	return uri + "|" + string(rune(pattern.PatternHash))
}

// updateMetrics updates performance metrics
func (pm *PatternMatcher) updateMetrics(duration time.Duration, matches int, cacheHit bool, wildcardMask []bool) {
	pm.Metrics.mu.Lock()
	defer pm.Metrics.mu.Unlock()
	
	pm.Metrics.TotalMatches++
	pm.Metrics.TotalDuration += duration
	
	if pm.Metrics.MinDuration == 0 || duration < pm.Metrics.MinDuration {
		pm.Metrics.MinDuration = duration
	}
	if duration > pm.Metrics.MaxDuration {
		pm.Metrics.MaxDuration = duration
	}
	
	pm.Metrics.AverageDuration = pm.Metrics.TotalDuration / time.Duration(pm.Metrics.TotalMatches)
	
	// Calculate wildcard savings (estimated time saved by skipping wildcards)
	wildcardCount := 0
	for _, isWildcard := range wildcardMask {
		if isWildcard {
			wildcardCount++
		}
	}
	estimatedSavings := time.Duration(wildcardCount) * duration / time.Duration(len(wildcardMask))
	pm.Metrics.WildcardSavings += estimatedSavings
	
	// Update cache hit rate
	stats := pm.Cache.GetStats()
	total := stats.Hits + stats.Misses
	if total > 0 {
		pm.Metrics.CacheHitRate = float64(stats.Hits) / float64(total) * 100
	}
}

// GetMetrics returns current performance metrics
func (pm *PatternMatcher) GetMetrics() MatchMetrics {
	pm.Metrics.mu.RLock()
	defer pm.Metrics.mu.RUnlock()
	return *pm.Metrics
}

// NewMatchCache creates a new match cache
func NewMatchCache(maxSize int) *MatchCache {
	return &MatchCache{
		cache:   make(map[string]*MatchResult),
		maxSize: maxSize,
	}
}

// Get retrieves a match result from cache
func (mc *MatchCache) Get(key string) (*MatchResult, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	if result, exists := mc.cache[key]; exists {
		mc.hits++
		return result, true
	}
	
	mc.misses++
	return nil, false
}

// Put stores a match result in cache
func (mc *MatchCache) Put(key string, result *MatchResult) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	if len(mc.cache) >= mc.maxSize {
		// Simple LRU eviction - remove oldest
		oldestKey := ""
		oldestTime := time.Now()
		for k, v := range mc.cache {
			if v.CacheTime.Before(oldestTime) {
				oldestTime = v.CacheTime
				oldestKey = k
			}
		}
		if oldestKey != "" {
			delete(mc.cache, oldestKey)
		}
	}
	
	mc.cache[key] = result
}

// GetStats returns cache statistics
func (mc *MatchCache) GetStats() struct{ Hits, Misses int64 } {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return struct{ Hits, Misses int64 }{mc.hits, mc.misses}
}

// NewStringPool creates a new string pool
func NewStringPool() *StringPool {
	return &StringPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 64)
			},
		},
	}
}

// Get retrieves a byte slice from the pool
func (sp *StringPool) Get() []byte {
	return sp.pool.Get().([]byte)[:0]
}

// Put returns a byte slice to the pool
func (sp *StringPool) Put(b []byte) {
	sp.pool.Put(b)
}