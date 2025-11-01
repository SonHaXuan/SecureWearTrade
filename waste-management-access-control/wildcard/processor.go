package wildcard

import (
	"fmt"
	"strings"
	"sync"
	"time"
	
	"../hibe"
	"../pattern"
)

// WildcardProcessor handles optimized wildcard pattern processing
type WildcardProcessor struct {
	PatternCache    map[string]*WildcardPattern
	OptimizedRules  map[string]*OptimizationRule
	MemoryOptimizer *WildcardMemoryManager
	Metrics         *WildcardMetrics
	mu              sync.RWMutex
}

// WildcardPattern represents an optimized wildcard pattern
type WildcardPattern struct {
	OriginalPattern   string
	Components        []WildcardComponent
	WildcardPositions []int
	StaticPositions   []int
	OptimizationLevel OptimizationLevel
	MemoryFootprint   int64
	CreatedAt         time.Time
}

// WildcardComponent represents a single component in a wildcard pattern
type WildcardComponent struct {
	Position    int
	Value       []byte
	IsWildcard  bool
	ComponentType ComponentType
	MemorySize  int
}

// ComponentType defines the type of waste-management URI component
type ComponentType int

const (
	TypeFacility ComponentType = iota
	TypeDepartment
	TypeBinLabel
	TypeBinID
	TypeDataType
	TypeAccessLevel
)

// OptimizationLevel defines the level of wildcard optimization applied
type OptimizationLevel int

const (
	OptimizationNone OptimizationLevel = iota
	OptimizationBasic
	OptimizationAdvanced
	OptimizationMaximum
)

// OptimizationRule defines rules for wildcard pattern optimization
type OptimizationRule struct {
	Pattern           string
	SkipComponents    []int
	MemoryReduction   float64
	SpeedImprovement  float64
	ApplicableTypes   []ComponentType
}

// WildcardMetrics tracks wildcard processing performance
type WildcardMetrics struct {
	TotalPatterns      int64
	OptimizedPatterns  int64
	MemorySaved        int64
	ProcessingTimeSaved time.Duration
	OptimizationRatio  float64
	AverageWildcards   float64
	mu                 sync.RWMutex
}

// WildcardMemoryManager optimizes memory usage for wildcard patterns
type WildcardMemoryManager struct {
	EmptySliceCache   map[int][]byte
	ComponentPools    map[ComponentType]*sync.Pool
	StringInternCache map[string][]byte
	mu                sync.RWMutex
}

// NewWildcardProcessor creates a new wildcard processor
func NewWildcardProcessor() *WildcardProcessor {
	wp := &WildcardProcessor{
		PatternCache:    make(map[string]*WildcardPattern),
		OptimizedRules:  make(map[string]*OptimizationRule),
		MemoryOptimizer: NewWildcardMemoryManager(),
		Metrics:         &WildcardMetrics{},
	}
	
	// Initialize optimization rules
	wp.initializeOptimizationRules()
	
	return wp
}

// ProcessWildcardPattern processes and optimizes a wildcard pattern
func (wp *WildcardProcessor) ProcessWildcardPattern(patternStr string) (*WildcardPattern, error) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	
	// Check cache first
	if cached, exists := wp.PatternCache[patternStr]; exists {
		return cached, nil
	}
	
	start := time.Now()
	
	// Parse the pattern
	pattern, err := wp.parseWildcardPattern(patternStr)
	if err != nil {
		return nil, err
	}
	
	// Apply optimizations
	wp.optimizePattern(pattern)
	
	// Cache the result
	wp.PatternCache[patternStr] = pattern
	
	// Update metrics
	wp.updateMetrics(pattern, time.Since(start))
	
	return pattern, nil
}

// parseWildcardPattern parses a pattern string into a WildcardPattern
func (wp *WildcardProcessor) parseWildcardPattern(patternStr string) (*WildcardPattern, error) {
	// Remove leading slash and split into components
	cleanPattern := strings.TrimPrefix(patternStr, "/")
	parts := strings.Split(cleanPattern, "/")
	
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid waste-management pattern: expected 6 components, got %d", len(parts))
	}
	
	pattern := &WildcardPattern{
		OriginalPattern:   patternStr,
		Components:        make([]WildcardComponent, len(parts)),
		WildcardPositions: []int{},
		StaticPositions:   []int{},
		OptimizationLevel: OptimizationNone,
		CreatedAt:         time.Now(),
	}
	
	totalMemory := int64(0)
	
	for i, part := range parts {
		component := WildcardComponent{
			Position:      i,
			ComponentType: wp.detectComponentType(i),
			IsWildcard:    part == "*",
		}
		
		if component.IsWildcard {
			// Wildcard component - no memory allocation needed
			component.Value = nil
			component.MemorySize = 0
			pattern.WildcardPositions = append(pattern.WildcardPositions, i)
		} else {
			// Static component - optimize memory allocation
			component.Value = wp.MemoryOptimizer.AllocateComponent(part, component.ComponentType)
			component.MemorySize = len(component.Value)
			pattern.StaticPositions = append(pattern.StaticPositions, i)
			totalMemory += int64(component.MemorySize)
		}
		
		pattern.Components[i] = component
	}
	
	pattern.MemoryFootprint = totalMemory
	
	return pattern, nil
}

// optimizePattern applies various optimizations to the wildcard pattern
func (wp *WildcardProcessor) optimizePattern(pattern *WildcardPattern) {
	// Calculate wildcard ratio
	wildcardCount := len(pattern.WildcardPositions)
	totalCount := len(pattern.Components)
	wildcardRatio := float64(wildcardCount) / float64(totalCount)
	
	// Apply optimization level based on wildcard ratio
	switch {
	case wildcardRatio >= 0.75:
		pattern.OptimizationLevel = OptimizationMaximum
		wp.applyMaximumOptimization(pattern)
	case wildcardRatio >= 0.50:
		pattern.OptimizationLevel = OptimizationAdvanced
		wp.applyAdvancedOptimization(pattern)
	case wildcardRatio >= 0.25:
		pattern.OptimizationLevel = OptimizationBasic
		wp.applyBasicOptimization(pattern)
	default:
		pattern.OptimizationLevel = OptimizationNone
	}
}

// applyMaximumOptimization applies maximum level optimizations
func (wp *WildcardProcessor) applyMaximumOptimization(pattern *WildcardPattern) {
	// Pre-allocate empty slices for wildcard positions
	for _, pos := range pattern.WildcardPositions {
		if pattern.Components[pos].Value == nil {
			pattern.Components[pos].Value = wp.MemoryOptimizer.GetEmptySlice(0)
		}
	}
	
	// Intern static component strings to reduce memory usage
	for _, pos := range pattern.StaticPositions {
		component := &pattern.Components[pos]
		component.Value = wp.MemoryOptimizer.InternString(string(component.Value))
	}
	
	// Calculate memory savings
	originalSize := int64(len(pattern.OriginalPattern))
	optimizedSize := pattern.MemoryFootprint
	memorySaved := originalSize - optimizedSize
	
	wp.Metrics.mu.Lock()
	wp.Metrics.MemorySaved += memorySaved
	wp.Metrics.mu.Unlock()
}

// applyAdvancedOptimization applies advanced level optimizations
func (wp *WildcardProcessor) applyAdvancedOptimization(pattern *WildcardPattern) {
	// Component type-specific optimizations
	for i := range pattern.Components {
		component := &pattern.Components[i]
		
		if !component.IsWildcard {
			switch component.ComponentType {
			case TypeFacility:
				// "facility" is always the same - use singleton
				component.Value = wp.MemoryOptimizer.GetSingleton("facility")
			case TypeBinLabel:
				// "bin" is always the same - use singleton
				component.Value = wp.MemoryOptimizer.GetSingleton("bin")
			case TypeDepartment:
				// Use department-specific pool
				component.Value = wp.MemoryOptimizer.GetFromPool(TypeDepartment, string(component.Value))
			case TypeDataType:
				// Use data type-specific pool
				component.Value = wp.MemoryOptimizer.GetFromPool(TypeDataType, string(component.Value))
			case TypeAccessLevel:
				// Use access level-specific pool
				component.Value = wp.MemoryOptimizer.GetFromPool(TypeAccessLevel, string(component.Value))
			}
		}
	}
}

// applyBasicOptimization applies basic level optimizations
func (wp *WildcardProcessor) applyBasicOptimization(pattern *WildcardPattern) {
	// Simple memory deduplication
	seen := make(map[string][]byte)
	
	for i := range pattern.Components {
		component := &pattern.Components[i]
		
		if !component.IsWildcard {
			value := string(component.Value)
			if cached, exists := seen[value]; exists {
				component.Value = cached
			} else {
				seen[value] = component.Value
			}
		}
	}
}

// detectComponentType determines component type based on position
func (wp *WildcardProcessor) detectComponentType(position int) ComponentType {
	switch position {
	case 0:
		return TypeFacility
	case 1:
		return TypeDepartment
	case 2:
		return TypeBinLabel
	case 3:
		return TypeBinID
	case 4:
		return TypeDataType
	case 5:
		return TypeAccessLevel
	default:
		return TypeDataType // Default fallback
	}
}

// GenerateHIBEPattern converts wildcard pattern to HIBE pattern
func (wp *WildcardProcessor) GenerateHIBEPattern(wildcardPattern *WildcardPattern) *hibe.WasteManagementPattern {
	components := make([]string, len(wildcardPattern.Components))
	wildcardMask := make([]bool, len(wildcardPattern.Components))
	
	activeDepth := 0
	for i, comp := range wildcardPattern.Components {
		if comp.IsWildcard {
			components[i] = ""  // Empty for wildcards
			wildcardMask[i] = true
		} else {
			components[i] = string(comp.Value)
			wildcardMask[i] = false
			activeDepth++
		}
	}
	
	return &hibe.WasteManagementPattern{
		Components:   components,
		WildcardMask: wildcardMask,
		Depth:        activeDepth,
		PatternType:  "wildcard-optimized",
		Facility:     wp.getComponentValue(wildcardPattern, 0),
		Department:   wp.getComponentValue(wildcardPattern, 1),
		BinID:    wp.getComponentValue(wildcardPattern, 3),
		DataType:     wp.getComponentValue(wildcardPattern, 4),
		AccessLevel:  wp.getComponentValue(wildcardPattern, 5),
	}
}

// getComponentValue safely gets component value
func (wp *WildcardProcessor) getComponentValue(pattern *WildcardPattern, position int) string {
	if position < len(pattern.Components) && !pattern.Components[position].IsWildcard {
		return string(pattern.Components[position].Value)
	}
	return "*"
}

// EstimateOptimizationGains estimates performance gains from wildcard optimization
func (wp *WildcardProcessor) EstimateOptimizationGains(pattern *WildcardPattern) OptimizationGains {
	wildcardCount := len(pattern.WildcardPositions)
	totalCount := len(pattern.Components)
	
	// Calculate estimated gains based on wildcard ratio
	wildcardRatio := float64(wildcardCount) / float64(totalCount)
	
	return OptimizationGains{
		MemoryReduction:    wildcardRatio * 0.6,  // Up to 60% memory reduction
		SpeedImprovement:   wildcardRatio * 0.8,  // Up to 80% speed improvement
		ComparisonsSaved:   wildcardCount,
		OptimizationLevel:  pattern.OptimizationLevel,
	}
}

// OptimizationGains represents estimated optimization benefits
type OptimizationGains struct {
	MemoryReduction   float64
	SpeedImprovement  float64
	ComparisonsSaved  int
	OptimizationLevel OptimizationLevel
}

// initializeOptimizationRules sets up predefined optimization rules
func (wp *WildcardProcessor) initializeOptimizationRules() {
	rules := []*OptimizationRule{
		{
			Pattern:          "/facility/*/bin/*/vitals/*",
			SkipComponents:   []int{1, 3, 5}, // Skip department, bin ID, access level
			MemoryReduction:  0.6,
			SpeedImprovement: 0.75,
			ApplicableTypes:  []ComponentType{TypeDepartment, TypeBinID, TypeAccessLevel},
		},
		{
			Pattern:          "/facility/*/bin/*/records/*",
			SkipComponents:   []int{1, 3, 5},
			MemoryReduction:  0.55,
			SpeedImprovement: 0.72,
			ApplicableTypes:  []ComponentType{TypeDepartment, TypeBinID, TypeAccessLevel},
		},
		{
			Pattern:          "/facility/cardiology/bin/*/vitals/*",
			SkipComponents:   []int{3, 5}, // Skip bin ID, access level
			MemoryReduction:  0.4,
			SpeedImprovement: 0.5,
			ApplicableTypes:  []ComponentType{TypeBinID, TypeAccessLevel},
		},
	}
	
	for _, rule := range rules {
		wp.OptimizedRules[rule.Pattern] = rule
	}
}

// updateMetrics updates wildcard processing metrics
func (wp *WildcardProcessor) updateMetrics(pattern *WildcardPattern, processingTime time.Duration) {
	wp.Metrics.mu.Lock()
	defer wp.Metrics.mu.Unlock()
	
	wp.Metrics.TotalPatterns++
	
	if pattern.OptimizationLevel > OptimizationNone {
		wp.Metrics.OptimizedPatterns++
	}
	
	wp.Metrics.ProcessingTimeSaved += processingTime
	
	// Update optimization ratio
	wp.Metrics.OptimizationRatio = float64(wp.Metrics.OptimizedPatterns) / float64(wp.Metrics.TotalPatterns)
	
	// Update average wildcards
	wildcardCount := float64(len(pattern.WildcardPositions))
	totalPatterns := float64(wp.Metrics.TotalPatterns)
	wp.Metrics.AverageWildcards = (wp.Metrics.AverageWildcards*(totalPatterns-1) + wildcardCount) / totalPatterns
}

// GetMetrics returns current wildcard processing metrics
func (wp *WildcardProcessor) GetMetrics() WildcardMetrics {
	wp.Metrics.mu.RLock()
	defer wp.Metrics.mu.RUnlock()
	return *wp.Metrics
}

// NewWildcardMemoryManager creates a new wildcard memory manager
func NewWildcardMemoryManager() *WildcardMemoryManager {
	wmm := &WildcardMemoryManager{
		EmptySliceCache:   make(map[int][]byte),
		ComponentPools:    make(map[ComponentType]*sync.Pool),
		StringInternCache: make(map[string][]byte),
	}
	
	// Initialize component pools
	for ct := TypeFacility; ct <= TypeAccessLevel; ct++ {
		wmm.ComponentPools[ct] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 32)
			},
		}
	}
	
	// Pre-allocate common empty slices
	wmm.EmptySliceCache[0] = []byte{}
	wmm.EmptySliceCache[8] = make([]byte, 0, 8)
	wmm.EmptySliceCache[16] = make([]byte, 0, 16)
	wmm.EmptySliceCache[32] = make([]byte, 0, 32)
	
	return wmm
}

// AllocateComponent allocates optimized memory for a component
func (wmm *WildcardMemoryManager) AllocateComponent(value string, componentType ComponentType) []byte {
	// Check if already interned
	if cached, exists := wmm.StringInternCache[value]; exists {
		return cached
	}
	
	// Allocate new byte slice
	result := make([]byte, len(value))
	copy(result, []byte(value))
	
	// Intern the string for future use
	wmm.StringInternCache[value] = result
	
	return result
}

// GetEmptySlice returns a pre-allocated empty slice
func (wmm *WildcardMemoryManager) GetEmptySlice(size int) []byte {
	wmm.mu.RLock()
	if cached, exists := wmm.EmptySliceCache[size]; exists {
		wmm.mu.RUnlock()
		return cached
	}
	wmm.mu.RUnlock()
	
	wmm.mu.Lock()
	defer wmm.mu.Unlock()
	
	// Double-check after acquiring write lock
	if cached, exists := wmm.EmptySliceCache[size]; exists {
		return cached
	}
	
	// Create new empty slice
	empty := make([]byte, 0, size)
	wmm.EmptySliceCache[size] = empty
	return empty
}

// InternString interns a string to reduce memory usage
func (wmm *WildcardMemoryManager) InternString(value string) []byte {
	wmm.mu.RLock()
	if cached, exists := wmm.StringInternCache[value]; exists {
		wmm.mu.RUnlock()
		return cached
	}
	wmm.mu.RUnlock()
	
	wmm.mu.Lock()
	defer wmm.mu.Unlock()
	
	// Double-check after acquiring write lock
	if cached, exists := wmm.StringInternCache[value]; exists {
		return cached
	}
	
	// Create new interned string
	bytes := make([]byte, len(value))
	copy(bytes, []byte(value))
	wmm.StringInternCache[value] = bytes
	return bytes
}

// GetSingleton returns a singleton byte slice for common values
func (wmm *WildcardMemoryManager) GetSingleton(value string) []byte {
	// Common waste-management singletons
	switch value {
	case "facility":
		return []byte("facility")
	case "bin":
		return []byte("bin")
	default:
		return wmm.InternString(value)
	}
}

// GetFromPool retrieves a byte slice from the appropriate component pool
func (wmm *WildcardMemoryManager) GetFromPool(componentType ComponentType, value string) []byte {
	if pool, exists := wmm.ComponentPools[componentType]; exists {
		bytes := pool.Get().([]byte)[:0]  // Reset length to 0
		bytes = append(bytes, []byte(value)...)
		return bytes
	}
	
	// Fallback to intern string
	return wmm.InternString(value)
}