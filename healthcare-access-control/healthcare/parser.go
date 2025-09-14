package healthcare

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

// HealthcareURIParser provides comprehensive healthcare URI parsing and validation
type HealthcareURIParser struct {
	validators      map[ComponentPosition]*regexp.Regexp
	allowedValues   map[ComponentPosition][]string
	parseCache      map[string]*HealthcareURI
	cacheSize       int
	cacheHits       int64
	cacheMisses     int64
	mu              sync.RWMutex
	
	// Performance optimizations
	fastValidators  map[ComponentPosition]func(string) bool
	compiledPatterns map[string]*pattern.CompiledPattern
}

// ComponentPosition represents the position of components in healthcare URIs
type ComponentPosition int

const (
	PosHospital ComponentPosition = iota
	PosDepartment
	PosPatientLabel
	PosPatientID
	PosDataType
	PosAccessLevel
)

// HealthcareURI represents a parsed and validated healthcare URI
type HealthcareURI struct {
	// Original URI
	OriginalURI string
	
	// Parsed components
	Hospital     string
	Department   string
	PatientLabel string
	PatientID    string
	DataType     string
	AccessLevel  string
	
	// Metadata
	IsWildcard      []bool
	ComponentCount  int
	ValidationLevel ValidationLevel
	ParsedAt        time.Time
	
	// Derived information
	PatientIDNumeric int64
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
	ValidationHealthcare
)

// DepartmentType categorizes medical departments
type DepartmentType int

const (
	DeptUnknown DepartmentType = iota
	DeptSpecialist
	DeptEmergency
	DeptGeneral
	DeptDiagnostic
)

// DataTypeCategory categorizes healthcare data types
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

// NewHealthcareURIParser creates a new healthcare URI parser
func NewHealthcareURIParser(cacheSize int) *HealthcareURIParser {
	parser := &HealthcareURIParser{
		validators:       make(map[ComponentPosition]*regexp.Regexp),
		allowedValues:    make(map[ComponentPosition][]string),
		parseCache:       make(map[string]*HealthcareURI),
		cacheSize:        cacheSize,
		fastValidators:   make(map[ComponentPosition]func(string) bool),
		compiledPatterns: make(map[string]*pattern.CompiledPattern),
	}
	
	parser.initializeValidators()
	parser.initializeAllowedValues()
	parser.initializeFastValidators()
	
	return parser
}

// ParseHealthcareURI parses and validates a healthcare URI
func (p *HealthcareURIParser) ParseHealthcareURI(uri string) (*HealthcareURI, error) {
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
	healthcareURI, err := p.parseURI(uri)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	p.cacheResult(uri, healthcareURI)
	
	return healthcareURI, nil
}

// parseURI performs the actual URI parsing
func (p *HealthcareURIParser) parseURI(uri string) (*HealthcareURI, error) {
	// Basic format validation
	if !strings.HasPrefix(uri, "/hospital/") {
		return nil, &ParsingError{
			URI:       uri,
			Component: "prefix",
			Position:  0,
			Reason:    "URI must start with '/hospital/'",
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
	
	// Create healthcare URI object
	healthcareURI := &HealthcareURI{
		OriginalURI:     uri,
		ComponentCount:  len(components),
		IsWildcard:      make([]bool, len(components)),
		ParsedAt:        time.Now(),
		ValidationLevel: ValidationHealthcare,
	}
	
	// Parse and validate each component
	if err := p.parseAndValidateComponents(healthcareURI, components); err != nil {
		return nil, err
	}
	
	// Derive additional metadata
	p.deriveMetadata(healthcareURI)
	
	return healthcareURI, nil
}

// splitURI efficiently splits URI into components
func (p *HealthcareURIParser) splitURI(uri string) []string {
	// Remove leading slash and split
	if len(uri) > 0 && uri[0] == '/' {
		uri = uri[1:]
	}
	return strings.Split(uri, "/")
}

// parseAndValidateComponents parses and validates individual URI components
func (p *HealthcareURIParser) parseAndValidateComponents(healthcareURI *HealthcareURI, components []string) error {
	componentMap := []struct {
		position ComponentPosition
		setter   func(string)
	}{
		{PosHospital, func(s string) { healthcareURI.Hospital = s }},
		{PosDepartment, func(s string) { healthcareURI.Department = s }},
		{PosPatientLabel, func(s string) { healthcareURI.PatientLabel = s }},
		{PosPatientID, func(s string) { healthcareURI.PatientID = s }},
		{PosDataType, func(s string) { healthcareURI.DataType = s }},
		{PosAccessLevel, func(s string) { healthcareURI.AccessLevel = s }},
	}
	
	for i, component := range components {
		if i >= len(componentMap) {
			break
		}
		
		pos := componentMap[i].position
		
		// Check if component is wildcard
		if component == "*" {
			healthcareURI.IsWildcard[i] = true
			componentMap[i].setter("*")
			continue
		}
		
		// Validate component using fast validators first
		if fastValidator, exists := p.fastValidators[pos]; exists {
			if !fastValidator(component) {
				return &ParsingError{
					URI:       healthcareURI.OriginalURI,
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
					URI:       healthcareURI.OriginalURI,
					Component: component,
					Position:  i,
					Reason:    fmt.Sprintf("invalid format for %s", p.getComponentName(pos)),
					Timestamp: time.Now(),
				}
			}
		}
		
		healthcareURI.IsWildcard[i] = false
		componentMap[i].setter(component)
	}
	
	return nil
}

// deriveMetadata derives additional metadata from parsed components
func (p *HealthcareURIParser) deriveMetadata(healthcareURI *HealthcareURI) {
	// Parse patient ID as numeric if possible
	if healthcareURI.PatientID != "*" {
		if id, err := strconv.ParseInt(healthcareURI.PatientID, 10, 64); err == nil {
			healthcareURI.PatientIDNumeric = id
		}
	}
	
	// Determine department type
	healthcareURI.DepartmentType = p.getDepartmentType(healthcareURI.Department)
	
	// Determine data type category
	healthcareURI.DataTypeCategory = p.getDataTypeCategory(healthcareURI.DataType)
	
	// Determine access priority
	healthcareURI.AccessPriority = p.getAccessPriority(healthcareURI.AccessLevel)
}

// getDepartmentType categorizes department type
func (p *HealthcareURIParser) getDepartmentType(department string) DepartmentType {
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
func (p *HealthcareURIParser) getDataTypeCategory(dataType string) DataTypeCategory {
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
func (p *HealthcareURIParser) getAccessPriority(accessLevel string) AccessPriority {
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
func (p *HealthcareURIParser) GenerateHIBEPattern(healthcareURI *HealthcareURI) *hibe.HealthcarePattern {
	components := []string{
		healthcareURI.Hospital,
		healthcareURI.Department,
		healthcareURI.PatientLabel,
		healthcareURI.PatientID,
		healthcareURI.DataType,
		healthcareURI.AccessLevel,
	}
	
	// Count active (non-wildcard) components
	activeDepth := 0
	for _, isWildcard := range healthcareURI.IsWildcard {
		if !isWildcard {
			activeDepth++
		}
	}
	
	return &hibe.HealthcarePattern{
		Components:   components,
		WildcardMask: healthcareURI.IsWildcard,
		Depth:        activeDepth,
		PatternType:  "healthcare-parsed",
		Hospital:     healthcareURI.Hospital,
		Department:   healthcareURI.Department,
		PatientID:    healthcareURI.PatientID,
		DataType:     healthcareURI.DataType,
		AccessLevel:  healthcareURI.AccessLevel,
	}
}

// GenerateWildcardPattern converts parsed URI to wildcard pattern
func (p *HealthcareURIParser) GenerateWildcardPattern(healthcareURI *HealthcareURI) (*wildcard.WildcardPattern, error) {
	wildcardProcessor := wildcard.NewWildcardProcessor()
	return wildcardProcessor.ProcessWildcardPattern(healthcareURI.OriginalURI)
}

// ValidateAccessPermissions validates if URI allows specific access patterns
func (p *HealthcareURIParser) ValidateAccessPermissions(healthcareURI *HealthcareURI, requiredAccess *AccessRequirements) bool {
	// Check department access
	if requiredAccess.DepartmentRequired != "" && 
	   healthcareURI.Department != "*" && 
	   healthcareURI.Department != requiredAccess.DepartmentRequired {
		return false
	}
	
	// Check patient access
	if requiredAccess.PatientIDRequired != "" && 
	   healthcareURI.PatientID != "*" && 
	   healthcareURI.PatientID != requiredAccess.PatientIDRequired {
		return false
	}
	
	// Check data type access
	if requiredAccess.DataTypeRequired != "" && 
	   healthcareURI.DataType != "*" && 
	   healthcareURI.DataType != requiredAccess.DataTypeRequired {
		return false
	}
	
	// Check access level priority
	if healthcareURI.AccessPriority < requiredAccess.MinAccessPriority {
		return false
	}
	
	return true
}

// AccessRequirements defines required access parameters
type AccessRequirements struct {
	DepartmentRequired  string
	PatientIDRequired   string
	DataTypeRequired    string
	MinAccessPriority   AccessPriority
}

// EstimateOptimizationPotential estimates optimization potential for the URI
func (p *HealthcareURIParser) EstimateOptimizationPotential(healthcareURI *HealthcareURI) OptimizationPotential {
	wildcardCount := 0
	for _, isWildcard := range healthcareURI.IsWildcard {
		if isWildcard {
			wildcardCount++
		}
	}
	
	wildcardRatio := float64(wildcardCount) / float64(len(healthcareURI.IsWildcard))
	
	return OptimizationPotential{
		WildcardCount:       wildcardCount,
		WildcardRatio:       wildcardRatio,
		MemoryReduction:     wildcardRatio * 0.6,  // Up to 60%
		SpeedImprovement:    wildcardRatio * 0.8,  // Up to 80%
		PatternComplexity:   p.calculatePatternComplexity(healthcareURI),
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
func (p *HealthcareURIParser) calculatePatternComplexity(healthcareURI *HealthcareURI) float64 {
	complexity := 0.0
	
	// Base complexity from component types
	if healthcareURI.Department != "*" {
		complexity += 1.0
	}
	if healthcareURI.PatientID != "*" {
		complexity += 2.0  // Patient ID is more specific
	}
	if healthcareURI.DataType != "*" {
		complexity += 1.5
	}
	if healthcareURI.AccessLevel != "*" {
		complexity += 1.0
	}
	
	// Adjust for department type
	switch healthcareURI.DepartmentType {
	case DeptEmergency:
		complexity += 0.5  // Emergency access is more complex
	case DeptSpecialist:
		complexity += 0.3
	}
	
	// Adjust for access priority
	switch healthcareURI.AccessPriority {
	case PriorityCritical:
		complexity += 0.5
	case PriorityHigh:
		complexity += 0.3
	}
	
	return complexity
}

// GetParsingStats returns parsing performance statistics
func (p *HealthcareURIParser) GetParsingStats() ParsingStats {
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
func (p *HealthcareURIParser) cacheResult(uri string, result *HealthcareURI) {
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
func (p *HealthcareURIParser) initializeValidators() {
	p.validators[PosHospital] = regexp.MustCompile(`^(hospital|\*)$`)
	p.validators[PosDepartment] = regexp.MustCompile(`^(cardiology|neurology|oncology|emergency|general|imaging|radiology|\*)$`)
	p.validators[PosPatientLabel] = regexp.MustCompile(`^(patient|\*)$`)
	p.validators[PosPatientID] = regexp.MustCompile(`^(\d{1,10}|\*)$`)
	p.validators[PosDataType] = regexp.MustCompile(`^(vitals|records|imaging|labs|\*)$`)
	p.validators[PosAccessLevel] = regexp.MustCompile(`^(realtime|historical|critical|routine|\*)$`)
}

// initializeAllowedValues sets up allowed values for each component
func (p *HealthcareURIParser) initializeAllowedValues() {
	p.allowedValues[PosHospital] = []string{"hospital", "*"}
	p.allowedValues[PosDepartment] = []string{"cardiology", "neurology", "oncology", "emergency", "general", "imaging", "radiology", "*"}
	p.allowedValues[PosPatientLabel] = []string{"patient", "*"}
	p.allowedValues[PosDataType] = []string{"vitals", "records", "imaging", "labs", "*"}
	p.allowedValues[PosAccessLevel] = []string{"realtime", "historical", "critical", "routine", "*"}
}

// initializeFastValidators sets up fast validation functions
func (p *HealthcareURIParser) initializeFastValidators() {
	// Hospital validator (always "hospital" or "*")
	p.fastValidators[PosHospital] = func(s string) bool {
		return s == "hospital" || s == "*"
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
	
	// Patient label validator (always "patient" or "*")
	p.fastValidators[PosPatientLabel] = func(s string) bool {
		return s == "patient" || s == "*"
	}
	
	// Patient ID validator (numeric or "*")
	p.fastValidators[PosPatientID] = func(s string) bool {
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
func (p *HealthcareURIParser) getComponentName(pos ComponentPosition) string {
	switch pos {
	case PosHospital:
		return "hospital"
	case PosDepartment:
		return "department"
	case PosPatientLabel:
		return "patient label"
	case PosPatientID:
		return "patient ID"
	case PosDataType:
		return "data type"
	case PosAccessLevel:
		return "access level"
	default:
		return "unknown"
	}
}

// IsValidHealthcareURI performs quick validation without full parsing
func (p *HealthcareURIParser) IsValidHealthcareURI(uri string) bool {
	// Quick format check
	if !strings.HasPrefix(uri, "/hospital/") {
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
func (p *HealthcareURIParser) GetOptimizedParsingRecommendations(healthcareURI *HealthcareURI) []OptimizationRecommendation {
	var recommendations []OptimizationRecommendation
	
	potential := p.EstimateOptimizationPotential(healthcareURI)
	
	if potential.WildcardRatio > 0.5 {
		recommendations = append(recommendations, OptimizationRecommendation{
			Type:        "wildcard-optimization",
			Description: "High wildcard ratio detected - apply maximum wildcard optimization",
			Priority:    PriorityHigh,
			EstimatedGain: fmt.Sprintf("%.1f%% speed improvement, %.1f%% memory reduction", 
				potential.SpeedImprovement*100, potential.MemoryReduction*100),
		})
	}
	
	if healthcareURI.AccessPriority == PriorityCritical {
		recommendations = append(recommendations, OptimizationRecommendation{
			Type:        "critical-path-optimization",
			Description: "Critical access detected - enable priority processing",
			Priority:    PriorityCritical,
			EstimatedGain: "Reduced latency for critical healthcare operations",
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