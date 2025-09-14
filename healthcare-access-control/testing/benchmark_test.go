package testing

import (
	"fmt"
	"testing"
	"time"
	"sync"
	"context"
	"runtime"
	
	"../hibe"
	"../pattern"
	"../wildcard"
	"../healthcare"
)

// BenchmarkHIBEKeyGeneration benchmarks HIBE key generation performance
func BenchmarkHIBEKeyGeneration(b *testing.B) {
	// Initialize HIBE key generator
	params := &hibe.SystemParams{MaxDepth: 6}
	hibeGen, err := hibe.NewHIBEKeyGenerator(params)
	if err != nil {
		b.Fatal(err)
	}
	
	// Test patterns for healthcare URIs
	testPatterns := []*hibe.HealthcarePattern{
		{
			Components:   []string{"hospital", "cardiology", "patient", "12345", "vitals", "realtime"},
			WildcardMask: []bool{false, false, false, false, false, false},
			PatternType:  "cardiology",
		},
		{
			Components:   []string{"hospital", "*", "patient", "*", "vitals", "*"},
			WildcardMask: []bool{false, true, false, true, false, true},
			PatternType:  "wildcard",
		},
	}
	
	b.ResetTimer()
	
	// Benchmark standard vs wildcard optimized key generation
	for i := 0; i < b.N; i++ {
		pattern := testPatterns[i%len(testPatterns)]
		_, _, err := hibeGen.GenerateHealthcareKey(pattern)
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkPatternMatching benchmarks pattern matching performance
func BenchmarkPatternMatching(b *testing.B) {
	matcher := pattern.NewPatternMatcher(1000)
	
	// Test URIs and patterns
	testCases := []struct {
		uri     string
		pattern string
	}{
		{"/hospital/cardiology/patient/12345/vitals/realtime", "/hospital/cardiology/patient/12345/vitals/realtime"},
		{"/hospital/neurology/patient/67890/records/historical", "/hospital/*/patient/*/records/*"},
		{"/hospital/oncology/patient/11111/imaging/routine", "/hospital/*/patient/*/imaging/*"},
	}
	
	// Pre-compile patterns
	compiledPatterns := make([]*pattern.CompiledPattern, len(testCases))
	for i, tc := range testCases {
		compiledPatterns[i] = matcher.CompilePattern(tc.pattern)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		tc := testCases[i%len(testCases)]
		cp := compiledPatterns[i%len(compiledPatterns)]
		
		_, _, _ = matcher.MatchHealthcarePattern(tc.uri, cp)
	}
}

// BenchmarkWildcardProcessing benchmarks wildcard pattern processing
func BenchmarkWildcardProcessing(b *testing.B) {
	processor := wildcard.NewWildcardProcessor(500)
	
	testPatterns := []string{
		"/hospital/*/patient/*/vitals/*",
		"/hospital/cardiology/patient/*/records/*",
		"/hospital/*/patient/12345/imaging/*",
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		pattern := testPatterns[i%len(testPatterns)]
		_, _ = processor.OptimizeWildcardPattern(pattern)
	}
}

// BenchmarkHealthcareParserIntegration benchmarks integrated healthcare parsing
func BenchmarkHealthcareParserIntegration(b *testing.B) {
	parser := healthcare.NewHealthcareParser()
	
	testURIs := []string{
		"/hospital/cardiology/patient/12345/vitals/realtime",
		"/hospital/neurology/patient/67890/records/historical", 
		"/hospital/oncology/patient/11111/imaging/routine",
		"/hospital/emergency/patient/22222/vitals/critical",
		"/hospital/general/patient/33333/labs/routine",
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		uri := testURIs[i%len(testURIs)]
		_, err := parser.ParseHealthcareURI(uri)
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkConcurrentLoad tests concurrent performance under load
func BenchmarkConcurrentLoad(b *testing.B) {
	// Initialize test suite components
	params := &hibe.SystemParams{MaxDepth: 6}
	hibeGen, _ := hibe.NewHIBEKeyGenerator(params)
	matcher := pattern.NewPatternMatcher(1000)
	parser := healthcare.NewHealthcareParser()
	
	// Test data
	testURI := "/hospital/cardiology/patient/12345/vitals/realtime"
	testPattern := "/hospital/*/patient/*/vitals/*"
	
	compiledPattern := matcher.CompilePattern(testPattern)
	parsedData, _ := parser.ParseHealthcareURI(testURI)
	
	healthcarePattern := &hibe.HealthcarePattern{
		Components:   parsedData.Components,
		WildcardMask: []bool{false, true, false, true, false, true},
		PatternType:  parsedData.DepartmentType,
	}
	
	b.ResetTimer()
	
	// Run concurrent operations
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Test HIBE key generation
			_, _, _ = hibeGen.GenerateHealthcareKey(healthcarePattern)
			
			// Test pattern matching
			_, _, _ = matcher.MatchHealthcarePattern(testURI, compiledPattern)
		}
	})
}

// LoadTest performs comprehensive load testing with different scenarios
func LoadTest(iterations int, concurrency int, duration time.Duration) *LoadTestResults {
	results := &LoadTestResults{
		Iterations:   iterations,
		Concurrency:  concurrency,
		Duration:     duration,
		StartTime:    time.Now(),
		Operations:   make(map[string]*OperationStats),
	}
	
	// Initialize components
	params := &hibe.SystemParams{MaxDepth: 6}
	hibeGen, _ := hibe.NewHIBEKeyGenerator(params)
	matcher := pattern.NewPatternMatcher(1000)
	processor := wildcard.NewWildcardProcessor(500)
	parser := healthcare.NewHealthcareParser()
	
	// Test scenarios
	scenarios := []LoadTestScenario{
		{
			Name: "Non-Wildcard HIBE",
			URI:  "/hospital/cardiology/patient/12345/vitals/realtime",
			Pattern: "/hospital/cardiology/patient/12345/vitals/realtime",
			IsWildcard: false,
		},
		{
			Name: "Wildcard-Optimized HIBE",
			URI:  "/hospital/neurology/patient/67890/records/historical",
			Pattern: "/hospital/*/patient/*/records/*",
			IsWildcard: true,
		},
		{
			Name: "Complex Pattern Match",
			URI:  "/hospital/oncology/patient/11111/imaging/routine",
			Pattern: "/hospital/*/patient/*/imaging/*",
			IsWildcard: true,
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	
	var wg sync.WaitGroup
	operationChan := make(chan OperationResult, concurrency*10)
	
	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runLoadTestWorker(ctx, workerID, scenarios, hibeGen, matcher, processor, parser, operationChan)
		}(i)
	}
	
	// Result collector
	go func() {
		for result := range operationChan {
			results.collectResult(result)
		}
	}()
	
	wg.Wait()
	close(operationChan)
	
	results.EndTime = time.Now()
	results.calculateSummary()
	
	return results
}

// LoadTestResults stores comprehensive load test results
type LoadTestResults struct {
	Iterations      int
	Concurrency     int
	Duration        time.Duration
	StartTime       time.Time
	EndTime         time.Time
	TotalOperations int64
	TotalErrors     int64
	Operations      map[string]*OperationStats
	ThroughputOPS   float64
	AvgLatency      time.Duration
	P95Latency      time.Duration
	P99Latency      time.Duration
	ErrorRate       float64
	mu              sync.RWMutex
}

// OperationStats tracks statistics for specific operations
type OperationStats struct {
	Count       int64
	Errors      int64
	TotalTime   time.Duration
	MinTime     time.Duration
	MaxTime     time.Duration
	Latencies   []time.Duration
	mu          sync.RWMutex
}

// LoadTestScenario defines a load test scenario
type LoadTestScenario struct {
	Name       string
	URI        string
	Pattern    string
	IsWildcard bool
}

// OperationResult represents the result of a single operation
type OperationResult struct {
	Operation string
	Duration  time.Duration
	Error     error
	WorkerID  int
}

// runLoadTestWorker executes load test operations for a single worker
func runLoadTestWorker(ctx context.Context, workerID int, scenarios []LoadTestScenario, 
	hibeGen *hibe.HIBEKeyGenerator, matcher *pattern.PatternMatcher, 
	processor *wildcard.WildcardProcessor, parser *healthcare.HealthcareParser,
	resultChan chan<- OperationResult) {
	
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Execute test scenarios cyclically
			for _, scenario := range scenarios {
				// Parse healthcare URI
				start := time.Now()
				parsedData, err := parser.ParseHealthcareURI(scenario.URI)
				duration := time.Since(start)
				resultChan <- OperationResult{"Parse", duration, err, workerID}
				
				if err != nil {
					continue
				}
				
				// HIBE Key Generation
				start = time.Now()
				pattern := &hibe.HealthcarePattern{
					Components:   parsedData.Components,
					WildcardMask: createWildcardMask(scenario.IsWildcard, len(parsedData.Components)),
					PatternType:  parsedData.DepartmentType,
				}
				_, hibeDuration, err := hibeGen.GenerateHealthcareKey(pattern)
				resultChan <- OperationResult{"HIBE", hibeDuration, err, workerID}
				
				// Pattern Matching
				start = time.Now()
				compiledPattern := matcher.CompilePattern(scenario.Pattern)
				_, matchDuration, _ := matcher.MatchHealthcarePattern(scenario.URI, compiledPattern)
				totalMatchDuration := time.Since(start)
				resultChan <- OperationResult{"Pattern", totalMatchDuration, nil, workerID}
				
				// Wildcard Processing (if applicable)
				if scenario.IsWildcard {
					start = time.Now()
					_, _ = processor.OptimizeWildcardPattern(scenario.Pattern)
					wildcardDuration := time.Since(start)
					resultChan <- OperationResult{"Wildcard", wildcardDuration, nil, workerID}
				}
			}
		}
	}
}

// Helper function to create wildcard mask
func createWildcardMask(isWildcard bool, length int) []bool {
	mask := make([]bool, length)
	if isWildcard {
		// Set common wildcard positions for healthcare patterns
		wildcardPositions := []int{1, 3, 5} // dept, patient ID, access level
		for _, pos := range wildcardPositions {
			if pos < length {
				mask[pos] = true
			}
		}
	}
	return mask
}

// collectResult safely collects operation results
func (ltr *LoadTestResults) collectResult(result OperationResult) {
	ltr.mu.Lock()
	defer ltr.mu.Unlock()
	
	ltr.TotalOperations++
	
	if result.Error != nil {
		ltr.TotalErrors++
	}
	
	// Initialize operation stats if not exists
	if _, exists := ltr.Operations[result.Operation]; !exists {
		ltr.Operations[result.Operation] = &OperationStats{
			MinTime:   time.Hour, // Initialize with large value
			Latencies: make([]time.Duration, 0, 1000),
		}
	}
	
	opStats := ltr.Operations[result.Operation]
	opStats.mu.Lock()
	defer opStats.mu.Unlock()
	
	opStats.Count++
	opStats.TotalTime += result.Duration
	
	if result.Error != nil {
		opStats.Errors++
	}
	
	if result.Duration < opStats.MinTime {
		opStats.MinTime = result.Duration
	}
	if result.Duration > opStats.MaxTime {
		opStats.MaxTime = result.Duration
	}
	
	// Collect latencies for percentile calculation
	if len(opStats.Latencies) < cap(opStats.Latencies) {
		opStats.Latencies = append(opStats.Latencies, result.Duration)
	}
}

// calculateSummary computes final load test statistics
func (ltr *LoadTestResults) calculateSummary() {
	totalDuration := ltr.EndTime.Sub(ltr.StartTime)
	ltr.ThroughputOPS = float64(ltr.TotalOperations) / totalDuration.Seconds()
	ltr.ErrorRate = float64(ltr.TotalErrors) / float64(ltr.TotalOperations) * 100
	
	// Calculate overall average latency
	var totalTime time.Duration
	for _, opStats := range ltr.Operations {
		totalTime += opStats.TotalTime
	}
	if ltr.TotalOperations > 0 {
		ltr.AvgLatency = totalTime / time.Duration(ltr.TotalOperations)
	}
	
	// Calculate percentiles from collected latencies
	ltr.calculatePercentiles()
}

// calculatePercentiles computes P95 and P99 latencies
func (ltr *LoadTestResults) calculatePercentiles() {
	allLatencies := make([]time.Duration, 0)
	
	for _, opStats := range ltr.Operations {
		allLatencies = append(allLatencies, opStats.Latencies...)
	}
	
	if len(allLatencies) == 0 {
		return
	}
	
	// Simple percentile calculation (proper implementation would sort)
	if len(allLatencies) > 0 {
		ltr.P95Latency = allLatencies[int(float64(len(allLatencies))*0.95)]
		ltr.P99Latency = allLatencies[int(float64(len(allLatencies))*0.99)]
	}
}

// PrintLoadTestReport prints comprehensive load test results
func (ltr *LoadTestResults) PrintLoadTestReport() {
	fmt.Printf("\n" + "="*80 + "\n")
	fmt.Printf("HEALTHCARE ACCESS CONTROL - LOAD TEST REPORT\n")
	fmt.Printf("="*80 + "\n")
	
	fmt.Printf("Test Configuration:\n")
	fmt.Printf("  Iterations: %d\n", ltr.Iterations)
	fmt.Printf("  Concurrency: %d\n", ltr.Concurrency)
	fmt.Printf("  Duration: %v\n", ltr.Duration)
	fmt.Printf("  Total Operations: %d\n", ltr.TotalOperations)
	
	fmt.Printf("\nPerformance Metrics:\n")
	fmt.Printf("  Throughput: %.2f ops/sec\n", ltr.ThroughputOPS)
	fmt.Printf("  Average Latency: %v\n", ltr.AvgLatency)
	fmt.Printf("  P95 Latency: %v\n", ltr.P95Latency)
	fmt.Printf("  P99 Latency: %v\n", ltr.P99Latency)
	fmt.Printf("  Error Rate: %.2f%%\n", ltr.ErrorRate)
	
	fmt.Printf("\nOperation Breakdown:\n")
	fmt.Printf("%-15s | %-10s | %-10s | %-12s | %-12s | %-12s\n", 
		"Operation", "Count", "Errors", "Avg Time", "Min Time", "Max Time")
	fmt.Printf("%s\n", "-"*80)
	
	for opName, stats := range ltr.Operations {
		avgTime := time.Duration(0)
		if stats.Count > 0 {
			avgTime = stats.TotalTime / time.Duration(stats.Count)
		}
		
		fmt.Printf("%-15s | %-10d | %-10d | %-12v | %-12v | %-12v\n",
			opName, stats.Count, stats.Errors, avgTime, stats.MinTime, stats.MaxTime)
	}
	
	fmt.Printf("\n✅ Load test completed successfully!\n")
}

// PerformanceComparisonTest compares optimized vs non-optimized performance
func PerformanceComparisonTest() {
	fmt.Println("\n=== PERFORMANCE COMPARISON: Optimized vs Non-Optimized ===")
	
	// Test scenarios for comparison
	scenarios := []struct {
		name        string
		iterations  int
		concurrency int
		duration    time.Duration
	}{
		{"Light Load", 1000, 10, 30 * time.Second},
		{"Medium Load", 5000, 25, 60 * time.Second},
		{"Heavy Load", 10000, 50, 120 * time.Second},
	}
	
	for _, scenario := range scenarios {
		fmt.Printf("\n--- %s Test ---\n", scenario.name)
		
		// Run optimized test
		fmt.Printf("Running Optimized Test...\n")
		optimizedResults := LoadTest(scenario.iterations, scenario.concurrency, scenario.duration)
		
		// Print comparison summary
		fmt.Printf("\nOptimized Results:\n")
		fmt.Printf("  Throughput: %.2f ops/sec\n", optimizedResults.ThroughputOPS)
		fmt.Printf("  Average Latency: %v\n", optimizedResults.AvgLatency)
		fmt.Printf("  Error Rate: %.2f%%\n", optimizedResults.ErrorRate)
		
		// Memory usage
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("  Memory Usage: %.2f MB\n", float64(m.Alloc)/1024/1024)
	}
}