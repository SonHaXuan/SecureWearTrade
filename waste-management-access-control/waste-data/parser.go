package waste-management

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	
	"../hibe"
	"../pattern"
	"../wildcard"
)

// WasteManagementURIParser provides comprehensive waste-management URI parsing and validation
type WasteManagementURIParser struct {
	validators      map[ComponentPosition]*regexp.Regexp
	allowedValues   map[ComponentPosition][]string
	parseCache      map[string]*WasteManagementURI
	cacheSize       int
	cacheHits       int64
	cacheMisses     int64
	mu              sync.RWMutex
	
	// Performance optimizations
	fastValidators  map[ComponentPosition]func(string) bool
	compiledPatterns map[string]*pattern.CompiledPattern
}

// ComponentPosition represents the position of components in waste-management URIs
type ComponentPosition int

const (
	PosFacility ComponentPosition = iota
	PosDepartment
	PosBinLabel
	PosBinID
	PosDataType
	PosAccessLevel
)

// WasteManagementURI represents a parsed and validated waste-management URI
type WasteManagementURI struct {
	// Original URI
	OriginalURI string
	
	// Parsed components
	Facility     string
	Department   string
	BinLabel string
	BinID    string
	DataType     string
	AccessLevel  string
	
	// Metadata
	IsWildcard      []bool
	ComponentCount  int
	ValidationLevel ValidationLevel
	ParsedAt        time.Time
	
	// Derived information
	BinIDNumeric int64
	DepartmentType   DepartmentType
	DataTypeCategory DataTypeCategory
	AccessPriority   AccessPriority
}

// ValidationLevel defines the level of URI validation performed
type ValidationLevel int

const (
	ValidationNone ValidationLevel = iota
	ValidationBasic
	ValidationStrict
	ValidationWasteManagement
)

// DepartmentType categorizes waste departments
type DepartmentType int

const (
	DeptUnknown DepartmentType = iota
	DeptSpecialist
	DeptEmergency
	DeptGeneral
	DeptDiagnostic
)

// DataTypeCategory categorizes waste-management data types
type DataTypeCategory int

const (
	DataUnknown DataTypeCategory = iota
	DataVital
	DataRecord
	DataImaging
	DataLaboratory
)

// AccessPriority defines access priority levels
type AccessPriority int

const (
	PriorityLow AccessPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// ParsingError represents errors during URI parsing
type ParsingError struct {
	URI       string
	Component string
	Position  int
	Reason    string
	Timestamp time.Time
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("parsing error at position %d (%s): %s in URI '%s'",
		e.Position, e.Component, e.Reason, e.URI)
}

// NewWasteManagementURIParser creates a new waste-management URI parser
func NewWasteManagementURIParser(cacheSize int) *WasteManagementURIParser {
	parser := &WasteManagementURIParser{
		validators:       make(map[ComponentPosition]*regexp.Regexp),
		allowedValues:    make(map[ComponentPosition][]string),
		parseCache:       make(map[string]*WasteManagementURI),
		cacheSize:        cacheSize,
		fastValidators:   make(map[ComponentPosition]func(string) bool),
		compiledPatterns: make(map[string]*pattern.CompiledPattern),
	}
	
	parser.initializeValidators()
	parser.initializeAllowedValues()
	parser.initializeFastValidators()
	
	return parser
}

// ParseWasteManagementURI parses and validates a waste-management URI
func (p *WasteManagementURIParser) ParseWasteManagementURI(uri string) (*WasteManagementURI, error) {
	// Check cache first
	p.mu.RLock()
	if cached, exists := p.parseCache[uri]; exists {
		p.cacheHits++
		p.mu.RUnlock()
		return cached, nil
	}
	p.cacheMisses++
	p.mu.RUnlock()
	
	// Perform actual parsing
	waste-managementURI, err := p.parseURI(uri)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	p.cacheResult(uri, waste-managementURI)
	
	return waste-managementURI, nil
}

// parseURI performs the actual URI parsing
func (p *WasteManagementURIParser) parseURI(uri string) (*WasteManagementURI, error) {
	// Basic format validation
	if !strings.HasPrefix(uri, "/facility/") {
		return nil, &ParsingError{
			URI:       uri,
			Component: "prefix",
			Position:  0,
			Reason:    "URI must start with '/facility/'",
			Timestamp: time.Now(),
		}
	}
	
	// Split URI into components
	components := p.splitURI(uri)
	if len(components) != 6 {
		return nil, &ParsingError{
			URI:       uri,
			Component: "structure",
			Position:  -1,
			Reason:    fmt.Sprintf("expected 6 components, got %d", len(components)),
			Timestamp: time.Now(),
		}
	}
	
	// Create waste-management URI object
	waste-managementURI := &WasteManagementURI{
		OriginalURI:     uri,
		ComponentCount:  len(components),
		IsWildcard:      make([]bool, len(components)),
		ParsedAt:        time.Now(),
		ValidationLevel: ValidationWasteManagement,
	}
	
	// Parse and validate each component
	if err := p.parseAndValidateComponents(waste-managementURI, components); err != nil {
		return nil, err
	}
	
	// Derive additional metadata
	p.deriveMetadata(waste-managementURI)
	
	return waste-managementURI, nil
}

// splitURI efficiently splits URI into components
func (p *WasteManagementURIParser) splitURI(uri string) []string {
	// Remove leading slash and split
	if len(uri) > 0 && uri[0] == '/' {
		uri = uri[1:]
	}
	return strings.Split(uri, "/")
}

// parseAndValidateComponents parses and validates individual URI components
func (p *WasteManagementURIParser) parseAndValidateComponents(waste-managementURI *WasteManagementURI, components []string) error {
	componentMap := []struct {
		position ComponentPosition
		setter   func(string)
	}{
		{PosFacility, func(s string) { waste-managementURI.Facility = s }},
		{PosDepartment, func(s string) { waste-managementURI.Department = s }},
		{PosBinLabel, func(s string) { waste-managementURI.BinLabel = s }},
		{PosBinID, func(s string) { waste-managementURI.BinID = s }},
		{PosDataType, func(s string) { waste-managementURI.DataType = s }},
		{PosAccessLevel, func(s string) { waste-managementURI.AccessLevel = s }},
	}
	
	for i, component := range components {
		if i >= len(componentMap) {
			break
		}
		
		pos := componentMap[i].position
		
		// Check if component is wildcard
		if component == "*" {
			waste-managementURI.IsWildcard[i] = true
			componentMap[i].setter("*")
			continue
		}
		
		// Validate component using fast validators first
		if fastValidator, exists := p.fastValidators[pos]; exists {
			if !fastValidator(component) {
				return &ParsingError{
					URI:       waste-managementURI.OriginalURI,
					Component: component,
					Position:  i,
					Reason:    fmt.Sprintf("invalid %s", p.getComponentName(pos)),
					Timestamp: time.Now(),
				}
			}
		}
		
		// Additional regex validation for complex patterns
		if validator, exists := p.validators[pos]; exists {
			if !validator.MatchString(component) {
				return &ParsingError{
					URI:       waste-managementURI.OriginalURI,
					Component: component,
					Position:  i,
					Reason:    fmt.Sprintf("invalid format for %s", p.getComponentName(pos)),
					Timestamp: time.Now(),
				}
			}
		}
		
		waste-managementURI.IsWildcard[i] = false
		componentMap[i].setter(component)
	}
	
	return nil
}

// deriveMetadata derives additional metadata from parsed components
func (p *WasteManagementURIParser) deriveMetadata(waste-managementURI *WasteManagementURI) {
	// Parse bin ID as numeric if possible
	if waste-managementURI.BinID != "*" {
		if id, err := strconv.ParseInt(waste-managementURI.BinID, 10, 64); err == nil {
			waste-managementURI.BinIDNumeric = id
		}
	}
	
	// Determine department type
	waste-managementURI.DepartmentType = p.getDepartmentType(waste-managementURI.Department)
	
	// Determine data type category
	waste-managementURI.DataTypeCategory = p.getDataTypeCategory(waste-managementURI.DataType)
	
	// Determine access priority
	waste-managementURI.AccessPriority = p.getAccessPriority(waste-managementURI.AccessLevel)
}

// getDepartmentType categorizes department type
func (p *WasteManagementURIParser) getDepartmentType(department string) DepartmentType {
	switch department {
	case "cardiology", "neurology", "oncology":
		return DeptSpecialist
	case "emergency":
		return DeptEmergency
	case "general":
		return DeptGeneral
	case "imaging", "radiology":
		return DeptDiagnostic
	case "*":
		return DeptUnknown
	default:
		return DeptGeneral
	}
}

// getDataTypeCategory categorizes data type
func (p *WasteManagementURIParser) getDataTypeCategory(dataType string) DataTypeCategory {
	switch dataType {
	case "vitals":
		return DataVital
	case "records":
		return DataRecord
	case "imaging":
		return DataImaging
	case "labs":
		return DataLaboratory
	case "*":
		return DataUnknown
	default:
		return DataRecord
	}
}

// getAccessPriority determines access priority
func (p *WasteManagementURIParser) getAccessPriority(accessLevel string) AccessPriority {
	switch accessLevel {
	case "critical":
		return PriorityCritical
	case "realtime":
		return PriorityHigh
	case "routine":
		return PriorityNormal
	case "historical":
		return PriorityLow
	case "*":
		return PriorityNormal
	default:
		return PriorityNormal
	}
}

// GenerateHIBEPattern converts parsed URI to HIBE pattern
func (p *WasteManagementURIParser) GenerateHIBEPattern(waste-managementURI *WasteManagementURI) *hibe.WasteManagementPattern {
	components := []string{
		waste-managementURI.Facility,
		waste-managementURI.Department,
		waste-managementURI.BinLabel,
		waste-managementURI.BinID,
		waste-managementURI.DataType,
		waste-managementURI.AccessLevel,
	}
	
	// Count active (non-wildcard) components
	activeDepth := 0
	for _, isWildcard := range waste-managementURI.IsWildcard {
		if !isWildcard {
			activeDepth++
		}
	}
	
	return &hibe.WasteManagementPattern{
		Components:   components,
		WildcardMask: waste-managementURI.IsWildcard,
		Depth:        activeDepth,
		PatternType:  "waste-management-parsed",
		Facility:     waste-managementURI.Facility,
		Department:   waste-managementURI.Department,
		BinID:    waste-managementURI.BinID,
		DataType:     waste-managementURI.DataType,
		AccessLevel:  waste-managementURI.AccessLevel,
	}
}

// GenerateWildcardPattern converts parsed URI to wildcard pattern
func (p *WasteManagementURIParser) GenerateWildcardPattern(waste-managementURI *WasteManagementURI) (*wildcard.WildcardPattern, error) {
	wildcardProcessor := wildcard.NewWildcardProcessor()
	return wildcardProcessor.ProcessWildcardPattern(waste-managementURI.OriginalURI)
}

// ValidateAccessPermissions validates if URI allows specific access patterns
func (p *WasteManagementURIParser) ValidateAccessPermissions(waste-managementURI *WasteManagementURI, requiredAccess *AccessRequirements) bool {
	// Check department access
	if requiredAccess.DepartmentRequired != "" && 
	   waste-managementURI.Department != "*" && 
	   waste-managementURI.Department != requiredAccess.DepartmentRequired {
		return false
	}
	
	// Check bin access
	if requiredAccess.BinIDRequired != "" && 
	   waste-managementURI.BinID != "*" && 
	   waste-managementURI.BinID != requiredAccess.BinIDRequired {
		return false
	}
	
	// Check data type access
	if requiredAccess.DataTypeRequired != "" && 
	   waste-managementURI.DataType != "*" && 
	   waste-managementURI.DataType != requiredAccess.DataTypeRequired {
		return false
	}
	
	// Check access level priority
	if waste-managementURI.AccessPriority < requiredAccess.MinAccessPriority {
		return false
	}
	
	return true
}

// AccessRequirements defines required access parameters
type AccessRequirements struct {
	DepartmentRequired  string
	BinIDRequired   string
	DataTypeRequired    string
	MinAccessPriority   AccessPriority
}

// EstimateOptimizationPotential estimates optimization potential for the URI
func (p *WasteManagementURIParser) EstimateOptimizationPotential(waste-managementURI *WasteManagementURI) OptimizationPotential {
	wildcardCount := 0
	for _, isWildcard := range waste-managementURI.IsWildcard {
		if isWildcard {
			wildcardCount++
		}
	}
	
	wildcardRatio := float64(wildcardCount) / float64(len(waste-managementURI.IsWildcard))
	
	return OptimizationPotential{
		WildcardCount:       wildcardCount,
		WildcardRatio:       wildcardRatio,
		MemoryReduction:     wildcardRatio * 0.6,  // Up to 60%
		SpeedImprovement:    wildcardRatio * 0.8,  // Up to 80%
		PatternComplexity:   p.calculatePatternComplexity(waste-managementURI),
		OptimizationScore:   wildcardRatio * 100,
	}
}

// OptimizationPotential represents optimization potential metrics
type OptimizationPotential struct {
	WildcardCount     int
	WildcardRatio     float64
	MemoryReduction   float64
	SpeedImprovement  float64
	PatternComplexity float64
	OptimizationScore float64
}

// calculatePatternComplexity calculates pattern complexity score
func (p *WasteManagementURIParser) calculatePatternComplexity(waste-managementURI *WasteManagementURI) float64 {
	complexity := 0.0
	
	// Base complexity from component types
	if waste-managementURI.Department != "*" {
		complexity += 1.0
	}
	if waste-managementURI.BinID != "*" {
		complexity += 2.0  // Bin ID is more specific
	}
	if waste-managementURI.DataType != "*" {
		complexity += 1.5
	}
	if waste-managementURI.AccessLevel != "*" {
		complexity += 1.0
	}
	
	// Adjust for department type
	switch waste-managementURI.DepartmentType {
	case DeptEmergency:
		complexity += 0.5  // Emergency access is more complex
	case DeptSpecialist:
		complexity += 0.3
	}
	
	// Adjust for access priority
	switch waste-managementURI.AccessPriority {
	case PriorityCritical:
		complexity += 0.5
	case PriorityHigh:
		complexity += 0.3
	}
	
	return complexity
}

// GetParsingStats returns parsing performance statistics
func (p *WasteManagementURIParser) GetParsingStats() ParsingStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	total := p.cacheHits + p.cacheMisses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(p.cacheHits) / float64(total) * 100
	}
	
	return ParsingStats{
		TotalParses:   total,
		CacheHits:     p.cacheHits,
		CacheMisses:   p.cacheMisses,
		CacheHitRate:  hitRate,
		CacheSize:     len(p.parseCache),
		MaxCacheSize:  p.cacheSize,
	}
}

// ParsingStats contains parsing performance statistics
type ParsingStats struct {
	TotalParses   int64
	CacheHits     int64
	CacheMisses   int64
	CacheHitRate  float64
	CacheSize     int
	MaxCacheSize  int
}

// cacheResult caches a parsing result
func (p *WasteManagementURIParser) cacheResult(uri string, result *WasteManagementURI) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Simple LRU eviction if cache is full
	if len(p.parseCache) >= p.cacheSize {
		// Remove oldest entry (simplified LRU)
		oldestURI := ""
		oldestTime := time.Now()
		for u, cached := range p.parseCache {
			if cached.ParsedAt.Before(oldestTime) {
				oldestTime = cached.ParsedAt
				oldestURI = u
			}
		}
		if oldestURI != "" {
			delete(p.parseCache, oldestURI)
		}
	}
	
	p.parseCache[uri] = result
}

// initializeValidators sets up regex validators for each component
func (p *WasteManagementURIParser) initializeValidators() {
	p.validators[PosFacility] = regexp.MustCompile(`^(facility|\*)$`)
	p.validators[PosDepartment] = regexp.MustCompile(`^(cardiology|neurology|oncology|emergency|general|imaging|radiology|\*)$`)
	p.validators[PosBinLabel] = regexp.MustCompile(`^(bin|\*)$`)
	p.validators[PosBinID] = regexp.MustCompile(`^(\d{1,10}|\*)$`)
	p.validators[PosDataType] = regexp.MustCompile(`^(vitals|records|imaging|labs|\*)$`)
	p.validators[PosAccessLevel] = regexp.MustCompile(`^(realtime|historical|critical|routine|\*)$`)
}

// initializeAllowedValues sets up allowed values for each component
func (p *WasteManagementURIParser) initializeAllowedValues() {
	p.allowedValues[PosFacility] = []string{"facility", "*"}
	p.allowedValues[PosDepartment] = []string{"cardiology", "neurology", "oncology", "emergency", "general", "imaging", "radiology", "*"}
	p.allowedValues[PosBinLabel] = []string{"bin", "*"}
	p.allowedValues[PosDataType] = []string{"vitals", "records", "imaging", "labs", "*"}
	p.allowedValues[PosAccessLevel] = []string{"realtime", "historical", "critical", "routine", "*"}
}

// initializeFastValidators sets up fast validation functions
func (p *WasteManagementURIParser) initializeFastValidators() {
	// Facility validator (always "facility" or "*")
	p.fastValidators[PosFacility] = func(s string) bool {
		return s == "facility" || s == "*"
	}
	
	// Department validator
	p.fastValidators[PosDepartment] = func(s string) bool {
		switch s {
		case "cardiology", "neurology", "oncology", "emergency", "general", "imaging", "radiology", "*":
			return true
		default:
			return false
		}
	}
	
	// Bin label validator (always "bin" or "*")
	p.fastValidators[PosBinLabel] = func(s string) bool {
		return s == "bin" || s == "*"
	}
	
	// Bin ID validator (numeric or "*")
	p.fastValidators[PosBinID] = func(s string) bool {
		if s == "*" {
			return true
		}
		if len(s) == 0 || len(s) > 10 {
			return false
		}
		for _, r := range s {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	}
	
	// Data type validator
	p.fastValidators[PosDataType] = func(s string) bool {
		switch s {
		case "vitals", "records", "imaging", "labs", "*":
			return true
		default:
			return false
		}
	}
	
	// Access level validator
	p.fastValidators[PosAccessLevel] = func(s string) bool {
		switch s {
		case "realtime", "historical", "critical", "routine", "*":
			return true
		default:
			return false
		}
	}
}

// getComponentName returns human-readable component name
func (p *WasteManagementURIParser) getComponentName(pos ComponentPosition) string {
	switch pos {
	case PosFacility:
		return "facility"
	case PosDepartment:
		return "department"
	case PosBinLabel:
		return "bin label"
	case PosBinID:
		return "bin ID"
	case PosDataType:
		return "data type"
	case PosAccessLevel:
		return "access level"
	default:
		return "unknown"
	}
}

// IsValidWasteManagementURI performs quick validation without full parsing
func (p *WasteManagementURIParser) IsValidWasteManagementURI(uri string) bool {
	// Quick format check
	if !strings.HasPrefix(uri, "/facility/") {
		return false
	}
	
	components := p.splitURI(uri)
	if len(components) != 6 {
		return false
	}
	
	// Quick validation of each component
	for i, component := range components {
		if i >= len(p.fastValidators) {
			continue
		}
		
		pos := ComponentPosition(i)
		if validator, exists := p.fastValidators[pos]; exists {
			if !validator(component) {
				return false
			}
		}
	}
	
	return true
}

// GetOptimizedParsingRecommendations provides optimization recommendations
func (p *WasteManagementURIParser) GetOptimizedParsingRecommendations(waste-managementURI *WasteManagementURI) []OptimizationRecommendation {
	var recommendations []OptimizationRecommendation
	
	potential := p.EstimateOptimizationPotential(waste-managementURI)
	
	if potential.WildcardRatio > 0.5 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Type:        "wildcard-optimization",
			Description: "High wildcard ratio detected - apply maximum wildcard optimization",
			Priority:    PriorityHigh,
			EstimatedGain: fmt.Sprintf("%.1f%% speed improvement, %.1f%% memory reduction", 
				potential.SpeedImprovement*100, potential.MemoryReduction*100),
		})
	}
	
	if waste-managementURI.AccessPriority == PriorityCritical {
		recommendations = append(recommendations, OptimizationRecommendation{
			Type:        "critical-path-optimization",
			Description: "Critical access detected - enable priority processing",
			Priority:    PriorityCritical,
			EstimatedGain: "Reduced latency for critical waste-management operations",
		})
	}
	
	if potential.PatternComplexity < 2.0 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Type:        "simple-pattern-optimization",
			Description: "Simple pattern detected - use lightweight processing",
			Priority:    PriorityNormal,
			EstimatedGain: "15-25% reduced processing overhead",
		})
	}
	
	return recommendations
}

// OptimizationRecommendation represents an optimization recommendation
type OptimizationRecommendation struct {
	Type          string
	Description   string
	Priority      AccessPriority
	EstimatedGain string
}