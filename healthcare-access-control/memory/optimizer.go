package memory

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// MemoryOptimizer provides advanced memory allocation optimizations
type MemoryOptimizer struct {
	StringPool      *OptimizedStringPool
	ComponentPool   *HealthcareComponentPool
	KeyPool         *CryptoKeyPool
	AllocationStats *AllocationStats
	GCController    *GCController
	MemoryMonitor   *MemoryMonitor
}

// OptimizedStringPool manages string allocations with size-based pooling
type OptimizedStringPool struct {
	pools    map[int]*sync.Pool  // Size-based pools
	maxSize  int
	hitCount int64
	missCount int64
	mu       sync.RWMutex
}

// HealthcareComponentPool manages healthcare-specific component allocations
type HealthcareComponentPool struct {
	hospitalPool    *sync.Pool
	departmentPool  *sync.Pool
	patientPool     *sync.Pool
	patientIDPool   *sync.Pool
	dataTypePool    *sync.Pool
	accessLevelPool *sync.Pool
	wildcardPool    *sync.Pool
	
	// Pre-allocated common values
	commonValues map[string][]byte
	mu           sync.RWMutex
}

// CryptoKeyPool manages cryptographic key component allocations
type CryptoKeyPool struct {
	bigIntPool     *sync.Pool
	componentPool  *sync.Pool
	privateKeyPool *sync.Pool
	
	// Reusable buffers for cryptographic operations
	hashBuffers   *sync.Pool
	tempBuffers   *sync.Pool
}

// AllocationStats tracks memory allocation statistics
type AllocationStats struct {
	TotalAllocations   int64
	TotalDeallocations int64
	CurrentUsage       int64
	PeakUsage          int64
	PoolHits           int64
	PoolMisses         int64
	
	// Size distribution
	SmallAllocations   int64  // < 64 bytes
	MediumAllocations  int64  // 64-512 bytes
	LargeAllocations   int64  // > 512 bytes
	
	// Time tracking
	StartTime          time.Time
	LastGC             time.Time
	GCCount            int64
	
	mu sync.RWMutex
}

// GCController manages garbage collection optimization
type GCController struct {
	gcPercent       int
	lastGC          time.Time
	gcThreshold     int64
	currentMemory   int64
	autoTuning      bool
	mu              sync.RWMutex
}

// MemoryMonitor provides real-time memory monitoring
type MemoryMonitor struct {
	isMonitoring  bool
	interval      time.Duration
	callbacks     []MemoryCallback
	stats         *runtime.MemStats
	mu            sync.RWMutex
}

// MemoryCallback defines callback function for memory events
type MemoryCallback func(stats *AllocationStats)

// OptimizedPattern represents a memory-optimized healthcare pattern
type OptimizedPattern struct {
	Components [][]byte
	Size       int
	Timestamp  time.Time
}

// NewMemoryOptimizer creates a new memory optimizer
func NewMemoryOptimizer() *MemoryOptimizer {
	return &MemoryOptimizer{
		StringPool:      NewOptimizedStringPool(),
		ComponentPool:   NewHealthcareComponentPool(),
		KeyPool:         NewCryptoKeyPool(),
		AllocationStats: NewAllocationStats(),
		GCController:    NewGCController(),
		MemoryMonitor:   NewMemoryMonitor(),
	}
}

// AllocateHealthcarePattern allocates memory-optimized healthcare pattern
func (mo *MemoryOptimizer) AllocateHealthcarePattern(components []string, isWildcard []bool) *OptimizedPattern {
	pattern := &OptimizedPattern{
		Components: make([][]byte, len(components)),
		Size:       0,
		Timestamp:  time.Now(),
	}
	
	for i, component := range components {
		if isWildcard[i] {
			// Use pre-allocated empty slice for wildcards - zero memory allocation
			pattern.Components[i] = mo.ComponentPool.GetWildcard()
		} else {
			// Allocate optimized component based on type
			pattern.Components[i] = mo.allocateOptimizedComponent(component, i)
			pattern.Size += len(pattern.Components[i])
		}
	}
	
	mo.AllocationStats.recordAllocation(int64(pattern.Size))
	return pattern
}

// allocateOptimizedComponent allocates memory for a specific component type
func (mo *MemoryOptimizer) allocateOptimizedComponent(component string, position int) []byte {
	switch position {
	case 0: // hospital
		return mo.ComponentPool.GetHospital(component)
	case 1: // department
		return mo.ComponentPool.GetDepartment(component)
	case 2: // patient label
		return mo.ComponentPool.GetPatient(component)
	case 3: // patient ID
		return mo.ComponentPool.GetPatientID(component)
	case 4: // data type
		return mo.ComponentPool.GetDataType(component)
	case 5: // access level
		return mo.ComponentPool.GetAccessLevel(component)
	default:
		return mo.StringPool.Get(component)
	}
}

// ReleasePattern returns pattern memory to pools
func (mo *MemoryOptimizer) ReleasePattern(pattern *OptimizedPattern) {
	if pattern == nil {
		return
	}
	
	for i, component := range pattern.Components {
		if component != nil && len(component) > 0 {
			mo.releaseComponent(component, i)
		}
	}
	
	mo.AllocationStats.recordDeallocation(int64(pattern.Size))
}

// releaseComponent returns component memory to appropriate pool
func (mo *MemoryOptimizer) releaseComponent(component []byte, position int) {
	switch position {
	case 0: // hospital
		mo.ComponentPool.ReleaseHospital(component)
	case 1: // department
		mo.ComponentPool.ReleaseDepartment(component)
	case 2: // patient label
		mo.ComponentPool.ReleasePatient(component)
	case 3: // patient ID
		mo.ComponentPool.ReleasePatientID(component)
	case 4: // data type
		mo.ComponentPool.ReleaseDataType(component)
	case 5: // access level
		mo.ComponentPool.ReleaseAccessLevel(component)
	default:
		mo.StringPool.Release(component)
	}
}

// OptimizeMemoryUsage performs comprehensive memory optimization
func (mo *MemoryOptimizer) OptimizeMemoryUsage() {
	// Force garbage collection if memory usage is high
	if mo.shouldTriggerGC() {
		runtime.GC()
		mo.GCController.recordGC()
	}
	
	// Optimize pool sizes based on usage patterns
	mo.optimizePoolSizes()
	
	// Clean up expired allocations
	mo.cleanupExpiredAllocations()
}

// shouldTriggerGC determines if garbage collection should be triggered
func (mo *MemoryOptimizer) shouldTriggerGC() bool {
	stats := mo.AllocationStats.GetSnapshot()
	
	// Trigger GC if:
	// 1. Current usage exceeds threshold
	// 2. Time since last GC exceeds interval
	// 3. Memory growth rate is high
	
	memoryThreshold := stats.PeakUsage * 8 / 10 // 80% of peak
	timeSinceGC := time.Since(mo.GCController.lastGC)
	
	return stats.CurrentUsage > memoryThreshold || timeSinceGC > 30*time.Second
}

// optimizePoolSizes adjusts pool sizes based on usage patterns
func (mo *MemoryOptimizer) optimizePoolSizes() {
	stats := mo.AllocationStats.GetSnapshot()
	
	// Adjust string pool sizes based on allocation distribution
	if stats.SmallAllocations > stats.MediumAllocations*2 {
		mo.StringPool.expandSmallPools()
	}
	
	if stats.LargeAllocations > stats.SmallAllocations {
		mo.StringPool.expandLargePools()
	}
}

// cleanupExpiredAllocations removes old unused allocations
func (mo *MemoryOptimizer) cleanupExpiredAllocations() {
	// Clean up component pools
	mo.ComponentPool.cleanup()
	
	// Clean up string pools
	mo.StringPool.cleanup()
	
	// Clean up key pools
	mo.KeyPool.cleanup()
}

// GetCurrentUsage returns current memory usage in bytes
func (mo *MemoryOptimizer) GetCurrentUsage() int64 {
	return atomic.LoadInt64(&mo.AllocationStats.CurrentUsage)
}

// GetOptimizationStats returns memory optimization statistics
func (mo *MemoryOptimizer) GetOptimizationStats() OptimizationStats {
	stats := mo.AllocationStats.GetSnapshot()
	poolHitRate := float64(stats.PoolHits) / float64(stats.PoolHits+stats.PoolMisses) * 100
	
	return OptimizationStats{
		TotalAllocations:   stats.TotalAllocations,
		CurrentUsage:       stats.CurrentUsage,
		PeakUsage:         stats.PeakUsage,
		PoolHitRate:       poolHitRate,
		MemoryEfficiency:  mo.calculateMemoryEfficiency(),
		GCCount:           stats.GCCount,
		LastGC:            stats.LastGC,
	}
}

// OptimizationStats contains memory optimization statistics
type OptimizationStats struct {
	TotalAllocations  int64
	CurrentUsage      int64
	PeakUsage         int64
	PoolHitRate       float64
	MemoryEfficiency  float64
	GCCount           int64
	LastGC            time.Time
}

// calculateMemoryEfficiency calculates memory usage efficiency
func (mo *MemoryOptimizer) calculateMemoryEfficiency() float64 {
	stats := mo.AllocationStats.GetSnapshot()
	
	if stats.PeakUsage == 0 {
		return 100.0
	}
	
	// Efficiency = (Peak - Current) / Peak * 100
	return float64(stats.PeakUsage-stats.CurrentUsage) / float64(stats.PeakUsage) * 100
}

// NewOptimizedStringPool creates a new optimized string pool
func NewOptimizedStringPool() *OptimizedStringPool {
	pool := &OptimizedStringPool{
		pools:   make(map[int]*sync.Pool),
		maxSize: 1024,
	}
	
	// Initialize size-based pools
	sizes := []int{8, 16, 32, 64, 128, 256, 512, 1024}
	for _, size := range sizes {
		pool.pools[size] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, size)
			},
		}
	}
	
	return pool
}

// Get retrieves an optimized byte slice for the given string
func (sp *OptimizedStringPool) Get(s string) []byte {
	length := len(s)
	poolSize := sp.getPoolSize(length)
	
	if pool, exists := sp.pools[poolSize]; exists {
		bytes := pool.Get().([]byte)[:0]  // Reset length
		bytes = append(bytes, s...)
		atomic.AddInt64(&sp.hitCount, 1)
		return bytes
	}
	
	// Fallback to direct allocation
	atomic.AddInt64(&sp.missCount, 1)
	return []byte(s)
}

// Release returns a byte slice to the appropriate pool
func (sp *OptimizedStringPool) Release(bytes []byte) {
	if bytes == nil || cap(bytes) > sp.maxSize {
		return
	}
	
	poolSize := sp.getPoolSize(cap(bytes))
	if pool, exists := sp.pools[poolSize]; exists {
		pool.Put(bytes[:0])  // Reset length but keep capacity
	}
}

// getPoolSize returns the appropriate pool size for a given length
func (sp *OptimizedStringPool) getPoolSize(length int) int {
	sizes := []int{8, 16, 32, 64, 128, 256, 512, 1024}
	for _, size := range sizes {
		if length <= size {
			return size
		}
	}
	return sp.maxSize
}

// expandSmallPools increases capacity of small pools
func (sp *OptimizedStringPool) expandSmallPools() {
	// Implementation for expanding small pools based on usage
	sp.mu.Lock()
	defer sp.mu.Unlock()
	
	// Add more small size pools if needed
	if _, exists := sp.pools[4]; !exists {
		sp.pools[4] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 4)
			},
		}
	}
}

// expandLargePools increases capacity of large pools
func (sp *OptimizedStringPool) expandLargePools() {
	// Implementation for expanding large pools based on usage
	sp.mu.Lock()
	defer sp.mu.Unlock()
	
	// Add larger size pools if needed
	if _, exists := sp.pools[2048]; !exists {
		sp.pools[2048] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 2048)
			},
		}
		sp.maxSize = 2048
	}
}

// cleanup removes unused pool entries
func (sp *OptimizedStringPool) cleanup() {
	// Implementation for cleaning up unused pool entries
	sp.mu.Lock()
	defer sp.mu.Unlock()
	
	// Remove pools that haven't been used recently
	// This is a simplified cleanup - in practice, you'd track usage timestamps
}

// NewHealthcareComponentPool creates a new healthcare component pool
func NewHealthcareComponentPool() *HealthcareComponentPool {
	pool := &HealthcareComponentPool{
		commonValues: make(map[string][]byte),
	}
	
	// Initialize component-specific pools
	pool.hospitalPool = &sync.Pool{New: func() interface{} { return make([]byte, 0, 16) }}
	pool.departmentPool = &sync.Pool{New: func() interface{} { return make([]byte, 0, 32) }}
	pool.patientPool = &sync.Pool{New: func() interface{} { return make([]byte, 0, 16) }}
	pool.patientIDPool = &sync.Pool{New: func() interface{} { return make([]byte, 0, 16) }}
	pool.dataTypePool = &sync.Pool{New: func() interface{} { return make([]byte, 0, 16) }}
	pool.accessLevelPool = &sync.Pool{New: func() interface{} { return make([]byte, 0, 16) }}
	pool.wildcardPool = &sync.Pool{New: func() interface{} { return []byte{} }}
	
	// Pre-allocate common healthcare values
	pool.preAllocateCommonValues()
	
	return pool
}

// preAllocateCommonValues pre-allocates commonly used healthcare values
func (hcp *HealthcareComponentPool) preAllocateCommonValues() {
	commonValues := map[string][]byte{
		"hospital":    []byte("hospital"),
		"patient":     []byte("patient"),
		"cardiology":  []byte("cardiology"),
		"neurology":   []byte("neurology"),
		"oncology":    []byte("oncology"),
		"emergency":   []byte("emergency"),
		"general":     []byte("general"),
		"vitals":      []byte("vitals"),
		"records":     []byte("records"),
		"imaging":     []byte("imaging"),
		"labs":        []byte("labs"),
		"realtime":    []byte("realtime"),
		"historical":  []byte("historical"),
		"critical":    []byte("critical"),
		"routine":     []byte("routine"),
	}
	
	hcp.mu.Lock()
	for k, v := range commonValues {
		hcp.commonValues[k] = v
	}
	hcp.mu.Unlock()
}

// Component getter methods with optimized allocation
func (hcp *HealthcareComponentPool) GetHospital(value string) []byte {
	return hcp.getCommonOrAllocate(value, hcp.hospitalPool)
}

func (hcp *HealthcareComponentPool) GetDepartment(value string) []byte {
	return hcp.getCommonOrAllocate(value, hcp.departmentPool)
}

func (hcp *HealthcareComponentPool) GetPatient(value string) []byte {
	return hcp.getCommonOrAllocate(value, hcp.patientPool)
}

func (hcp *HealthcareComponentPool) GetPatientID(value string) []byte {
	return hcp.getCommonOrAllocate(value, hcp.patientIDPool)
}

func (hcp *HealthcareComponentPool) GetDataType(value string) []byte {
	return hcp.getCommonOrAllocate(value, hcp.dataTypePool)
}

func (hcp *HealthcareComponentPool) GetAccessLevel(value string) []byte {
	return hcp.getCommonOrAllocate(value, hcp.accessLevelPool)
}

func (hcp *HealthcareComponentPool) GetWildcard() []byte {
	return hcp.wildcardPool.Get().([]byte)
}

// getCommonOrAllocate returns pre-allocated common value or allocates from pool
func (hcp *HealthcareComponentPool) getCommonOrAllocate(value string, pool *sync.Pool) []byte {
	// Check for common values first
	hcp.mu.RLock()
	if common, exists := hcp.commonValues[value]; exists {
		hcp.mu.RUnlock()
		return common
	}
	hcp.mu.RUnlock()
	
	// Allocate from pool
	bytes := pool.Get().([]byte)[:0]
	bytes = append(bytes, value...)
	return bytes
}

// Component release methods
func (hcp *HealthcareComponentPool) ReleaseHospital(bytes []byte) {
	if !hcp.isCommonValue(bytes) {
		hcp.hospitalPool.Put(bytes[:0])
	}
}

func (hcp *HealthcareComponentPool) ReleaseDepartment(bytes []byte) {
	if !hcp.isCommonValue(bytes) {
		hcp.departmentPool.Put(bytes[:0])
	}
}

func (hcp *HealthcareComponentPool) ReleasePatient(bytes []byte) {
	if !hcp.isCommonValue(bytes) {
		hcp.patientPool.Put(bytes[:0])
	}
}

func (hcp *HealthcareComponentPool) ReleasePatientID(bytes []byte) {
	if !hcp.isCommonValue(bytes) {
		hcp.patientIDPool.Put(bytes[:0])
	}
}

func (hcp *HealthcareComponentPool) ReleaseDataType(bytes []byte) {
	if !hcp.isCommonValue(bytes) {
		hcp.dataTypePool.Put(bytes[:0])
	}
}

func (hcp *HealthcareComponentPool) ReleaseAccessLevel(bytes []byte) {
	if !hcp.isCommonValue(bytes) {
		hcp.accessLevelPool.Put(bytes[:0])
	}
}

// isCommonValue checks if bytes represent a pre-allocated common value
func (hcp *HealthcareComponentPool) isCommonValue(bytes []byte) bool {
	hcp.mu.RLock()
	defer hcp.mu.RUnlock()
	
	for _, common := range hcp.commonValues {
		if len(bytes) == len(common) && 
		   (*(*string)(unsafe.Pointer(&bytes))) == (*(*string)(unsafe.Pointer(&common))) {
			return true
		}
	}
	return false
}

// cleanup removes unused allocations
func (hcp *HealthcareComponentPool) cleanup() {
	// Healthcare component pools are typically long-lived
	// Cleanup would involve removing least recently used entries
}

// NewCryptoKeyPool creates a new crypto key pool
func NewCryptoKeyPool() *CryptoKeyPool {
	return &CryptoKeyPool{
		bigIntPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 256) // Typical crypto key size
			},
		},
		componentPool: &sync.Pool{
			New: func() interface{} {
				return make([][]byte, 0, 8)
			},
		},
		privateKeyPool: &sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{})
			},
		},
		hashBuffers: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 64) // SHA256 size
			},
		},
		tempBuffers: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 512)
			},
		},
	}
}

// Key pool methods for cryptographic operations
func (ckp *CryptoKeyPool) GetBigIntBuffer() []byte {
	return ckp.bigIntPool.Get().([]byte)[:0]
}

func (ckp *CryptoKeyPool) ReleaseBigIntBuffer(buf []byte) {
	ckp.bigIntPool.Put(buf[:0])
}

func (ckp *CryptoKeyPool) GetHashBuffer() []byte {
	return ckp.hashBuffers.Get().([]byte)[:0]
}

func (ckp *CryptoKeyPool) ReleaseHashBuffer(buf []byte) {
	ckp.hashBuffers.Put(buf[:0])
}

func (ckp *CryptoKeyPool) cleanup() {
	// Crypto key cleanup - remove expired keys and buffers
}

// NewAllocationStats creates a new allocation statistics tracker
func NewAllocationStats() *AllocationStats {
	return &AllocationStats{
		StartTime: time.Now(),
	}
}

// recordAllocation records a memory allocation
func (as *AllocationStats) recordAllocation(size int64) {
	as.mu.Lock()
	defer as.mu.Unlock()
	
	atomic.AddInt64(&as.TotalAllocations, 1)
	atomic.AddInt64(&as.CurrentUsage, size)
	
	if as.CurrentUsage > as.PeakUsage {
		as.PeakUsage = as.CurrentUsage
	}
	
	// Categorize by size
	switch {
	case size < 64:
		atomic.AddInt64(&as.SmallAllocations, 1)
	case size < 512:
		atomic.AddInt64(&as.MediumAllocations, 1)
	default:
		atomic.AddInt64(&as.LargeAllocations, 1)
	}
}

// recordDeallocation records a memory deallocation
func (as *AllocationStats) recordDeallocation(size int64) {
	atomic.AddInt64(&as.TotalDeallocations, 1)
	atomic.AddInt64(&as.CurrentUsage, -size)
}

// GetSnapshot returns a snapshot of current statistics
func (as *AllocationStats) GetSnapshot() AllocationStats {
	as.mu.RLock()
	defer as.mu.RUnlock()
	
	return AllocationStats{
		TotalAllocations:   atomic.LoadInt64(&as.TotalAllocations),
		TotalDeallocations: atomic.LoadInt64(&as.TotalDeallocations),
		CurrentUsage:       atomic.LoadInt64(&as.CurrentUsage),
		PeakUsage:         as.PeakUsage,
		PoolHits:          atomic.LoadInt64(&as.PoolHits),
		PoolMisses:        atomic.LoadInt64(&as.PoolMisses),
		SmallAllocations:  atomic.LoadInt64(&as.SmallAllocations),
		MediumAllocations: atomic.LoadInt64(&as.MediumAllocations),
		LargeAllocations:  atomic.LoadInt64(&as.LargeAllocations),
		StartTime:         as.StartTime,
		LastGC:            as.LastGC,
		GCCount:           atomic.LoadInt64(&as.GCCount),
	}
}

// NewGCController creates a new garbage collection controller
func NewGCController() *GCController {
	return &GCController{
		gcPercent:   100,  // Default Go GC target
		gcThreshold: 1024 * 1024 * 10, // 10MB threshold
		autoTuning:  true,
	}
}

// recordGC records a garbage collection event
func (gc *GCController) recordGC() {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	
	gc.lastGC = time.Now()
	
	// Auto-tune GC based on memory patterns
	if gc.autoTuning {
		gc.tuneGCPercent()
	}
}

// tuneGCPercent automatically adjusts GC percentage based on usage
func (gc *GCController) tuneGCPercent() {
	// Simple auto-tuning algorithm
	// In practice, this would be more sophisticated
	if gc.currentMemory > gc.gcThreshold*2 {
		gc.gcPercent = 50  // More aggressive GC
	} else if gc.currentMemory < gc.gcThreshold/2 {
		gc.gcPercent = 200 // Less aggressive GC
	} else {
		gc.gcPercent = 100 // Default
	}
}

// NewMemoryMonitor creates a new memory monitor
func NewMemoryMonitor() *MemoryMonitor {
	return &MemoryMonitor{
		interval: 10 * time.Second,
		stats:    &runtime.MemStats{},
	}
}

// StartMonitoring begins memory monitoring
func (mm *MemoryMonitor) StartMonitoring() {
	mm.mu.Lock()
	if mm.isMonitoring {
		mm.mu.Unlock()
		return
	}
	mm.isMonitoring = true
	mm.mu.Unlock()
	
	go mm.monitorLoop()
}

// StopMonitoring stops memory monitoring
func (mm *MemoryMonitor) StopMonitoring() {
	mm.mu.Lock()
	mm.isMonitoring = false
	mm.mu.Unlock()
}

// monitorLoop runs the monitoring loop
func (mm *MemoryMonitor) monitorLoop() {
	ticker := time.NewTicker(mm.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			mm.mu.RLock()
			if !mm.isMonitoring {
				mm.mu.RUnlock()
				return
			}
			callbacks := mm.callbacks
			mm.mu.RUnlock()
			
			// Collect memory stats and notify callbacks
			runtime.ReadMemStats(mm.stats)
			
			// Convert to our allocation stats format
			allocStats := &AllocationStats{
				CurrentUsage: int64(mm.stats.HeapInuse),
				PeakUsage:    int64(mm.stats.HeapSys),
				GCCount:      int64(mm.stats.NumGC),
			}
			
			for _, callback := range callbacks {
				go callback(allocStats)
			}
		}
	}
}

// AddCallback adds a memory monitoring callback
func (mm *MemoryMonitor) AddCallback(callback MemoryCallback) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.callbacks = append(mm.callbacks, callback)
}