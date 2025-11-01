package security

import (
	"crypto/rand"
	"fmt"
	"math"
	"runtime"
	"sort"
	"testing"
	"time"
)

// SideChannelTestResult represents the result of a side-channel attack test
type SideChannelTestResult struct {
	TestName           string        `json:"testName"`
	AttackSuccessful   bool          `json:"attackSuccessful"`
	TimingVariation    float64       `json:"timingVariation"`
	MemoryLeakage      bool          `json:"memoryLeakage"`
	PowerVariation     float64       `json:"powerVariation"`
	StatisticalSignif  float64       `json:"statisticalSignificance"`
	ConfidenceInterval [2]float64    `json:"confidenceInterval"`
	SampleSize         int           `json:"sampleSize"`
	ErrorMessage       string        `json:"errorMessage,omitempty"`
}

// SideChannelExperiment conducts side-channel attack experiments
type SideChannelExperiment struct {
	TestIterations int
	SampleSize     int
	Results        []SideChannelTestResult
}

// TestSideChannelResistance performs comprehensive side-channel attack resistance testing
func TestSideChannelResistance(t *testing.T) {
	experiment := &SideChannelExperiment{
		TestIterations: 1000,
		SampleSize:     10000,
		Results:        make([]SideChannelTestResult, 0),
	}

	// Test 1: Timing Attack Test
	t.Run("TimingAttack", func(t *testing.T) {
		result := experiment.testTimingAttack()
		experiment.Results = append(experiment.Results, result)
		
		if result.AttackSuccessful {
			t.Errorf("Timing attack successful: variation %.2f ns", result.TimingVariation)
		}
	})

	// Test 2: Memory Access Pattern Test
	t.Run("MemoryAccessPattern", func(t *testing.T) {
		result := experiment.testMemoryAccessPattern()
		experiment.Results = append(experiment.Results, result)
		
		if result.MemoryLeakage {
			t.Errorf("Memory access pattern leakage detected")
		}
	})

	// Test 3: Cache Timing Attack Test
	t.Run("CacheTimingAttack", func(t *testing.T) {
		result := experiment.testCacheTimingAttack()
		experiment.Results = append(experiment.Results, result)
		
		if result.AttackSuccessful {
			t.Errorf("Cache timing attack successful")
		}
	})

	// Test 4: Power Analysis Attack Test
	t.Run("PowerAnalysisAttack", func(t *testing.T) {
		result := experiment.testPowerAnalysisAttack()
		experiment.Results = append(experiment.Results, result)
		
		if result.AttackSuccessful {
			t.Errorf("Power analysis attack successful: variation %.2f", result.PowerVariation)
		}
	})

	// Test 5: Branch Prediction Attack Test
	t.Run("BranchPredictionAttack", func(t *testing.T) {
		result := experiment.testBranchPredictionAttack()
		experiment.Results = append(experiment.Results, result)
		
		if result.AttackSuccessful {
			t.Errorf("Branch prediction attack successful")
		}
	})

	// Generate statistical analysis
	experiment.generateStatisticalAnalysis(t)
}

// testTimingAttack tests for timing attack vulnerabilities
func (e *SideChannelExperiment) testTimingAttack() SideChannelTestResult {
	// Test different key sizes and measure timing variations
	keySize1 := 128 // bits
	keySize2 := 256 // bits
	
	var timings1, timings2 []time.Duration
	
	// Measure timing for different key sizes
	for i := 0; i < e.SampleSize; i++ {
		// Generate random keys
		key1 := make([]byte, keySize1/8)
		key2 := make([]byte, keySize2/8)
		rand.Read(key1)
		rand.Read(key2)
		
		// Measure encryption time for key1
		start := time.Now()
		simulateEncryption(key1, []byte("test message"))
		timings1 = append(timings1, time.Since(start))
		
		// Measure encryption time for key2
		start = time.Now()
		simulateEncryption(key2, []byte("test message"))
		timings2 = append(timings2, time.Since(start))
	}
	
	// Calculate statistical significance
	mean1 := calculateMean(timings1)
	mean2 := calculateMean(timings2)
	std1 := calculateStdDev(timings1, mean1)
	std2 := calculateStdDev(timings2, mean2)
	
	// T-test for significant difference
	tStat := (mean1 - mean2) / math.Sqrt((std1*std1)/float64(len(timings1)) + (std2*std2)/float64(len(timings2)))
	pValue := calculatePValue(tStat, len(timings1)+len(timings2)-2)
	
	timingVariation := math.Abs(mean1 - mean2)
	attackSuccessful := pValue < 0.05 && timingVariation > 1000 // 1μs threshold
	
	return SideChannelTestResult{
		TestName:           "TimingAttack",
		AttackSuccessful:   attackSuccessful,
		TimingVariation:    timingVariation,
		MemoryLeakage:      false,
		PowerVariation:     0,
		StatisticalSignif:  pValue,
		ConfidenceInterval: [2]float64{mean1 - 1.96*std1, mean1 + 1.96*std1},
		SampleSize:         e.SampleSize,
	}
}

// testMemoryAccessPattern tests for memory access pattern leakage
func (e *SideChannelExperiment) testMemoryAccessPattern() SideChannelTestResult {
	// Test memory allocation patterns for different operations
	var memStats1, memStats2 []runtime.MemStats
	
	// Collect memory stats for different operations
	for i := 0; i < e.SampleSize/10; i++ {
		// Operation 1: Small data encryption
		runtime.GC()
		var ms1 runtime.MemStats
		runtime.ReadMemStats(&ms1)
		
		smallData := make([]byte, 64)
		rand.Read(smallData)
		simulateEncryption(smallData[:16], smallData[16:])
		
		runtime.ReadMemStats(&ms1)
		memStats1 = append(memStats1, ms1)
		
		// Operation 2: Large data encryption
		runtime.GC()
		var ms2 runtime.MemStats
		runtime.ReadMemStats(&ms2)
		
		largeData := make([]byte, 4096)
		rand.Read(largeData)
		simulateEncryption(largeData[:16], largeData[16:])
		
		runtime.ReadMemStats(&ms2)
		memStats2 = append(memStats2, ms2)
	}
	
	// Analyze memory allocation patterns
	memoryLeakage := analyzeMemoryPatterns(memStats1, memStats2)
	
	return SideChannelTestResult{
		TestName:           "MemoryAccessPattern",
		AttackSuccessful:   memoryLeakage,
		TimingVariation:    0,
		MemoryLeakage:      memoryLeakage,
		PowerVariation:     0,
		StatisticalSignif:  0.05,
		ConfidenceInterval: [2]float64{0, 1},
		SampleSize:         e.SampleSize / 10,
	}
}

// testCacheTimingAttack tests for cache timing attack vulnerabilities
func (e *SideChannelExperiment) testCacheTimingAttack() SideChannelTestResult {
	// Simulate cache timing attack using memory access patterns
	cacheLineSize := 64 // bytes
	testData := make([]byte, cacheLineSize*1024) // 64KB test data
	rand.Read(testData)
	
	var accessTimes []time.Duration
	
	// Measure cache access times
	for i := 0; i < e.SampleSize; i++ {
		// Flush cache (simulate)
		runtime.GC()
		
		// Access memory location
		offset := (i % 1024) * cacheLineSize
		start := time.Now()
		_ = testData[offset] // Memory access
		accessTimes = append(accessTimes, time.Since(start))
	}
	
	// Analyze timing patterns
	mean := calculateMean(accessTimes)
	stdDev := calculateStdDev(accessTimes, mean)
	
	// Check for consistent timing patterns (indication of cache attack success)
	variation := stdDev / mean
	attackSuccessful := variation < 0.1 // Low variation indicates predictable cache behavior
	
	return SideChannelTestResult{
		TestName:           "CacheTimingAttack",
		AttackSuccessful:   attackSuccessful,
		TimingVariation:    variation,
		MemoryLeakage:      false,
		PowerVariation:     0,
		StatisticalSignif:  0.05,
		ConfidenceInterval: [2]float64{mean - 1.96*stdDev, mean + 1.96*stdDev},
		SampleSize:         e.SampleSize,
	}
}

// testPowerAnalysisAttack tests for power analysis attack vulnerabilities
func (e *SideChannelExperiment) testPowerAnalysisAttack() SideChannelTestResult {
	// Simulate power consumption analysis
	// In a real implementation, this would interface with power measurement hardware
	
	var powerConsumption1, powerConsumption2 []float64
	
	// Measure power consumption for different operations
	for i := 0; i < e.SampleSize; i++ {
		// Operation 1: Encryption with key containing many 1s
		key1 := make([]byte, 16)
		for j := range key1 {
			key1[j] = 0xFF // All bits set
		}
		power1 := simulatePowerConsumption(key1)
		powerConsumption1 = append(powerConsumption1, power1)
		
		// Operation 2: Encryption with key containing many 0s
		key2 := make([]byte, 16)
		for j := range key2 {
			key2[j] = 0x00 // All bits clear
		}
		power2 := simulatePowerConsumption(key2)
		powerConsumption2 = append(powerConsumption2, power2)
	}
	
	// Statistical analysis
	mean1 := calculateMeanFloat64(powerConsumption1)
	mean2 := calculateMeanFloat64(powerConsumption2)
	std1 := calculateStdDevFloat64(powerConsumption1, mean1)
	std2 := calculateStdDevFloat64(powerConsumption2, mean2)
	
	// T-test for significant difference
	tStat := (mean1 - mean2) / math.Sqrt((std1*std1)/float64(len(powerConsumption1)) + (std2*std2)/float64(len(powerConsumption2)))
	pValue := calculatePValue(tStat, len(powerConsumption1)+len(powerConsumption2)-2)
	
	powerVariation := math.Abs(mean1 - mean2)
	attackSuccessful := pValue < 0.05 && powerVariation > 0.1 // 10% threshold
	
	return SideChannelTestResult{
		TestName:           "PowerAnalysisAttack",
		AttackSuccessful:   attackSuccessful,
		TimingVariation:    0,
		MemoryLeakage:      false,
		PowerVariation:     powerVariation,
		StatisticalSignif:  pValue,
		ConfidenceInterval: [2]float64{mean1 - 1.96*std1, mean1 + 1.96*std1},
		SampleSize:         e.SampleSize,
	}
}

// testBranchPredictionAttack tests for branch prediction attack vulnerabilities
func (e *SideChannelExperiment) testBranchPredictionAttack() SideChannelTestResult {
	// Test branch prediction patterns in cryptographic operations
	var timingsBranch1, timingsBranch2 []time.Duration
	
	for i := 0; i < e.SampleSize; i++ {
		// Generate predictable vs unpredictable branch patterns
		data1 := make([]byte, 16)
		data2 := make([]byte, 16)
		
		// Predictable pattern (all same values)
		for j := range data1 {
			data1[j] = 0x55
		}
		
		// Unpredictable pattern (random values)
		rand.Read(data2)
		
		// Measure timing for predictable branches
		start := time.Now()
		simulateBranchIntensiveOperation(data1)
		timingsBranch1 = append(timingsBranch1, time.Since(start))
		
		// Measure timing for unpredictable branches
		start = time.Now()
		simulateBranchIntensiveOperation(data2)
		timingsBranch2 = append(timingsBranch2, time.Since(start))
	}
	
	// Statistical analysis
	mean1 := calculateMean(timingsBranch1)
	mean2 := calculateMean(timingsBranch2)
	std1 := calculateStdDev(timingsBranch1, mean1)
	std2 := calculateStdDev(timingsBranch2, mean2)
	
	// T-test for significant difference
	tStat := (mean1 - mean2) / math.Sqrt((std1*std1)/float64(len(timingsBranch1)) + (std2*std2)/float64(len(timingsBranch2)))
	pValue := calculatePValue(tStat, len(timingsBranch1)+len(timingsBranch2)-2)
	
	timingVariation := math.Abs(mean1 - mean2)
	attackSuccessful := pValue < 0.05 && timingVariation > 500 // 0.5μs threshold
	
	return SideChannelTestResult{
		TestName:           "BranchPredictionAttack",
		AttackSuccessful:   attackSuccessful,
		TimingVariation:    timingVariation,
		MemoryLeakage:      false,
		PowerVariation:     0,
		StatisticalSignif:  pValue,
		ConfidenceInterval: [2]float64{mean1 - 1.96*std1, mean1 + 1.96*std1},
		SampleSize:         e.SampleSize,
	}
}

// generateStatisticalAnalysis generates comprehensive statistical analysis
func (e *SideChannelExperiment) generateStatisticalAnalysis(t *testing.T) {
	totalTests := len(e.Results)
	successfulAttacks := 0
	
	for _, result := range e.Results {
		if result.AttackSuccessful {
			successfulAttacks++
		}
	}
	
	successRate := float64(successfulAttacks) / float64(totalTests)
	
	t.Logf("=== Side-Channel Attack Resistance Test Results ===")
	t.Logf("Total tests: %d", totalTests)
	t.Logf("Successful attacks: %d", successfulAttacks)
	t.Logf("Success rate: %.2f%%", successRate*100)
	
	// Log detailed results for each test
	for _, result := range e.Results {
		t.Logf("Test: %s", result.TestName)
		t.Logf("  Attack successful: %v", result.AttackSuccessful)
		t.Logf("  Timing variation: %.2f ns", result.TimingVariation)
		t.Logf("  Memory leakage: %v", result.MemoryLeakage)
		t.Logf("  Power variation: %.2f", result.PowerVariation)
		t.Logf("  Statistical significance (p-value): %.4f", result.StatisticalSignif)
		t.Logf("  Sample size: %d", result.SampleSize)
		t.Logf("  Confidence interval: [%.2f, %.2f]", result.ConfidenceInterval[0], result.ConfidenceInterval[1])
	}
	
	// Assert that attack success rate is below acceptable threshold
	if successRate > 0.10 { // 10% threshold for side-channel attacks
		t.Errorf("Side-channel attack success rate too high: %.2f%% (threshold: 10%%)", successRate*100)
	}
}

// Helper functions for simulation and statistical analysis

func simulateEncryption(key, plaintext []byte) []byte {
	// Simulate encryption operation with timing variation
	result := make([]byte, len(plaintext))
	for i, b := range plaintext {
		// Simple XOR with timing variation based on key
		result[i] = b ^ key[i%len(key)]
		
		// Add artificial delay based on key bit count (simulating side-channel vulnerability)
		if popcount(key[i%len(key)]) > 4 {
			time.Sleep(time.Nanosecond * 10)
		}
	}
	return result
}

func simulatePowerConsumption(key []byte) float64 {
	// Simulate power consumption based on key properties
	basePower := 1.0
	
	// Power consumption correlates with number of 1s in key
	for _, b := range key {
		basePower += float64(popcount(b)) * 0.01
	}
	
	// Add random noise
	noise := (rand.Float64() - 0.5) * 0.05
	return basePower + noise
}

func simulateBranchIntensiveOperation(data []byte) {
	// Simulate operation with many branches
	result := 0
	for _, b := range data {
		if b&0x01 != 0 {
			result += int(b)
		}
		if b&0x02 != 0 {
			result *= 2
		}
		if b&0x04 != 0 {
			result ^= 0xFF
		}
		if b&0x08 != 0 {
			result = result << 1
		}
	}
	// Use result to prevent optimization
	_ = result
}

func popcount(b byte) int {
	count := 0
	for b != 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func analyzeMemoryPatterns(stats1, stats2 []runtime.MemStats) bool {
	// Analyze memory allocation patterns for differences
	if len(stats1) == 0 || len(stats2) == 0 {
		return false
	}
	
	// Simple analysis: check if memory usage patterns are significantly different
	var allocs1, allocs2 []float64
	for _, s := range stats1 {
		allocs1 = append(allocs1, float64(s.Alloc))
	}
	for _, s := range stats2 {
		allocs2 = append(allocs2, float64(s.Alloc))
	}
	
	mean1 := calculateMeanFloat64(allocs1)
	mean2 := calculateMeanFloat64(allocs2)
	
	// If memory usage patterns are significantly different, it indicates potential leakage
	return math.Abs(mean1-mean2) > 1000 // 1KB threshold
}

func calculateMean(data []time.Duration) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := time.Duration(0)
	for _, d := range data {
		sum += d
	}
	return float64(sum) / float64(len(data))
}

func calculateMeanFloat64(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, d := range data {
		sum += d
	}
	return sum / float64(len(data))
}

func calculateStdDev(data []time.Duration, mean float64) float64 {
	if len(data) <= 1 {
		return 0
	}
	
	sumSquares := 0.0
	for _, d := range data {
		diff := float64(d) - mean
		sumSquares += diff * diff
	}
	
	variance := sumSquares / float64(len(data)-1)
	return math.Sqrt(variance)
}

func calculateStdDevFloat64(data []float64, mean float64) float64 {
	if len(data) <= 1 {
		return 0
	}
	
	sumSquares := 0.0
	for _, d := range data {
		diff := d - mean
		sumSquares += diff * diff
	}
	
	variance := sumSquares / float64(len(data)-1)
	return math.Sqrt(variance)
}

func calculatePValue(tStat float64, degreesOfFreedom int) float64 {
	// Simplified p-value calculation (in practice, use proper statistical library)
	// This is a rough approximation
	absTStat := math.Abs(tStat)
	if absTStat > 2.576 {
		return 0.01 // p < 0.01
	} else if absTStat > 1.96 {
		return 0.05 // p < 0.05
	} else if absTStat > 1.645 {
		return 0.10 // p < 0.10
	}
	return 0.20 // p > 0.10
}

// BenchmarkSideChannelResistance provides performance benchmarks for side-channel resistance
func BenchmarkSideChannelResistance(b *testing.B) {
	key := make([]byte, 16)
	plaintext := make([]byte, 64)
	rand.Read(key)
	rand.Read(plaintext)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		simulateEncryption(key, plaintext)
	}
}