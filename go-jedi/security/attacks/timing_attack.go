package attacks

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"math"
	"sort"
	"time"
)

// TimingAttack represents a timing attack simulation
type TimingAttack struct {
	TargetFunction func([]byte, []byte) bool
	SampleSize     int
	Results        []TimingAttackResult
}

// TimingAttackResult represents the result of a timing attack
type TimingAttackResult struct {
	AttackType         string            `json:"attackType"`
	Success            bool              `json:"success"`
	ConfidenceLevel    float64           `json:"confidenceLevel"`
	StatisticalSignif  float64           `json:"statisticalSignificance"`
	TimingDifference   float64           `json:"timingDifferenceNs"`
	SampleSize         int               `json:"sampleSize"`
	ExecutionTime      time.Duration     `json:"executionTime"`
	SuccessRate        float64           `json:"successRate"`
	ErrorMessage       string            `json:"errorMessage,omitempty"`
	AttackMetadata     map[string]string `json:"attackMetadata"`
	TimingStatistics   TimingStatistics  `json:"timingStatistics"`
}

// TimingStatistics contains statistical analysis of timing measurements
type TimingStatistics struct {
	Mean           float64           `json:"mean"`
	StandardDev    float64           `json:"standardDeviation"`
	Variance       float64           `json:"variance"`
	Median         float64           `json:"median"`
	Min            float64           `json:"min"`
	Max            float64           `json:"max"`
	Quartiles      [3]float64        `json:"quartiles"`
	Outliers       int               `json:"outliers"`
	Distribution   []float64         `json:"distribution"`
	Histogram      map[string]int    `json:"histogram"`
}

// NewTimingAttack creates a new timing attack simulator
func NewTimingAttack(targetFunc func([]byte, []byte) bool, sampleSize int) *TimingAttack {
	return &TimingAttack{
		TargetFunction: targetFunc,
		SampleSize:     sampleSize,
		Results:        make([]TimingAttackResult, 0),
	}
}

// SimulatePasswordTimingAttack simulates timing attack on password verification
func (t *TimingAttack) SimulatePasswordTimingAttack() TimingAttackResult {
	start := time.Now()
	
	correctPassword := "correct_password_123"
	incorrectPassword := "incorrect_password"
	
	// Vulnerable password comparison function
	vulnerableCompare := func(password, input []byte) bool {
		if len(password) != len(input) {
			return false
		}
		
		// Vulnerable: early return on first mismatch
		for i := 0; i < len(password); i++ {
			if password[i] != input[i] {
				return false
			}
			// Add artificial delay to simulate processing
			time.Sleep(time.Nanosecond * 100)
		}
		return true
	}
	
	// Secure password comparison function
	secureCompare := func(password, input []byte) bool {
		return subtle.ConstantTimeCompare(password, input) == 1
	}
	
	// Test both vulnerable and secure implementations
	vulnerableResult := t.measurePasswordTimingAttack(vulnerableCompare, correctPassword, incorrectPassword)
	secureResult := t.measurePasswordTimingAttack(secureCompare, correctPassword, incorrectPassword)
	
	// Compare results
	attackSuccessful := vulnerableResult.Success && !secureResult.Success
	
	result := TimingAttackResult{
		AttackType:         "PasswordTimingAttack",
		Success:            attackSuccessful,
		ConfidenceLevel:    vulnerableResult.ConfidenceLevel,
		StatisticalSignif:  vulnerableResult.StatisticalSignif,
		TimingDifference:   vulnerableResult.TimingDifference,
		SampleSize:         t.SampleSize,
		ExecutionTime:      time.Since(start),
		SuccessRate:        calculateSuccessRate(vulnerableResult.Success, secureResult.Success),
		TimingStatistics:   vulnerableResult.TimingStatistics,
		AttackMetadata: map[string]string{
			"correct_password":     correctPassword,
			"vulnerable_attack":    fmt.Sprintf("%v", vulnerableResult.Success),
			"secure_attack":        fmt.Sprintf("%v", secureResult.Success),
			"timing_difference":    fmt.Sprintf("%.2f ns", vulnerableResult.TimingDifference),
		},
	}
	
	t.Results = append(t.Results, result)
	return result
}

// SimulateKeyComparisonAttack simulates timing attack on cryptographic key comparison
func (t *TimingAttack) SimulateKeyComparisonAttack() TimingAttackResult {
	start := time.Now()
	
	// Generate test keys
	correctKey := make([]byte, 32)
	rand.Read(correctKey)
	
	// Vulnerable key comparison
	vulnerableKeyCompare := func(key1, key2 []byte) bool {
		if len(key1) != len(key2) {
			return false
		}
		
		// Vulnerable: byte-by-byte comparison with early return
		for i := 0; i < len(key1); i++ {
			if key1[i] != key2[i] {
				return false
			}
			// Simulate processing delay
			time.Sleep(time.Nanosecond * 50)
		}
		return true
	}
	
	// Secure key comparison
	secureKeyCompare := func(key1, key2 []byte) bool {
		return subtle.ConstantTimeCompare(key1, key2) == 1
	}
	
	// Test timing attack on both implementations
	vulnerableResult := t.measureKeyTimingAttack(vulnerableKeyCompare, correctKey)
	secureResult := t.measureKeyTimingAttack(secureKeyCompare, correctKey)
	
	attackSuccessful := vulnerableResult.Success && !secureResult.Success
	
	result := TimingAttackResult{
		AttackType:         "KeyComparisonAttack",
		Success:            attackSuccessful,
		ConfidenceLevel:    vulnerableResult.ConfidenceLevel,
		StatisticalSignif:  vulnerableResult.StatisticalSignif,
		TimingDifference:   vulnerableResult.TimingDifference,
		SampleSize:         t.SampleSize,
		ExecutionTime:      time.Since(start),
		SuccessRate:        calculateSuccessRate(vulnerableResult.Success, secureResult.Success),
		TimingStatistics:   vulnerableResult.TimingStatistics,
		AttackMetadata: map[string]string{
			"key_length":           fmt.Sprintf("%d", len(correctKey)),
			"vulnerable_timing":    fmt.Sprintf("%.2f ns", vulnerableResult.TimingDifference),
			"secure_timing":        fmt.Sprintf("%.2f ns", secureResult.TimingDifference),
			"attack_effectiveness": fmt.Sprintf("%.2f%%", vulnerableResult.SuccessRate*100),
		},
	}
	
	t.Results = append(t.Results, result)
	return result
}

// SimulateHashComparisonAttack simulates timing attack on hash comparison
func (t *TimingAttack) SimulateHashComparisonAttack() TimingAttackResult {
	start := time.Now()
	
	// Generate test hashes (simulated)
	correctHash := "a1b2c3d4e5f6789012345678901234567890abcdef"
	
	// Vulnerable hash comparison
	vulnerableHashCompare := func(hash1, hash2 []byte) bool {
		h1, h2 := string(hash1), string(hash2)
		
		if len(h1) != len(h2) {
			return false
		}
		
		// Vulnerable: character-by-character comparison
		for i := 0; i < len(h1); i++ {
			if h1[i] != h2[i] {
				return false
			}
			// Simulate hash processing delay
			time.Sleep(time.Nanosecond * 75)
		}
		return true
	}
	
	// Secure hash comparison
	secureHashCompare := func(hash1, hash2 []byte) bool {
		return subtle.ConstantTimeCompare(hash1, hash2) == 1
	}
	
	// Test timing attack on both implementations
	vulnerableResult := t.measureHashTimingAttack(vulnerableHashCompare, correctHash)
	secureResult := t.measureHashTimingAttack(secureHashCompare, correctHash)
	
	attackSuccessful := vulnerableResult.Success && !secureResult.Success
	
	result := TimingAttackResult{
		AttackType:         "HashComparisonAttack",
		Success:            attackSuccessful,
		ConfidenceLevel:    vulnerableResult.ConfidenceLevel,
		StatisticalSignif:  vulnerableResult.StatisticalSignif,
		TimingDifference:   vulnerableResult.TimingDifference,
		SampleSize:         t.SampleSize,
		ExecutionTime:      time.Since(start),
		SuccessRate:        calculateSuccessRate(vulnerableResult.Success, secureResult.Success),
		TimingStatistics:   vulnerableResult.TimingStatistics,
		AttackMetadata: map[string]string{
			"hash_length":          fmt.Sprintf("%d", len(correctHash)),
			"vulnerable_exploited": fmt.Sprintf("%v", vulnerableResult.Success),
			"secure_protected":     fmt.Sprintf("%v", !secureResult.Success),
		},
	}
	
	t.Results = append(t.Results, result)
	return result
}

// SimulateRemoteTimingAttack simulates timing attack over network
func (t *TimingAttack) SimulateRemoteTimingAttack() TimingAttackResult {
	start := time.Now()
	
	// Simulate network delays and jitter
	networkDelay := func() time.Duration {
		// Random network delay between 1-10ms
		baseDelay := time.Millisecond * time.Duration(1+rand.Intn(10))
		
		// Add jitter (±20%)
		jitter := time.Duration(float64(baseDelay) * (rand.Float64()*0.4 - 0.2))
		
		return baseDelay + jitter
	}
	
	// Measure network timing variations
	var networkTimings []time.Duration
	
	for i := 0; i < t.SampleSize; i++ {
		start := time.Now()
		time.Sleep(networkDelay())
		
		// Simulate vulnerable operation
		vulnerableOp := func() {
			// Simulate timing-sensitive operation
			time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))
		}
		
		vulnerableOp()
		networkTimings = append(networkTimings, time.Since(start))
	}
	
	// Analyze timing patterns
	stats := calculateTimingStatistics(networkTimings)
	
	// Network timing attacks are harder but still possible
	// Check for detectable patterns despite network noise
	timingVariation := stats.StandardDev / stats.Mean
	attackSuccessful := timingVariation < 0.5 && stats.StandardDev > 1000 // Detectable pattern
	
	result := TimingAttackResult{
		AttackType:         "RemoteTimingAttack",
		Success:            attackSuccessful,
		ConfidenceLevel:    calculateConfidenceLevel(stats),
		StatisticalSignif:  calculateStatisticalSignificance(networkTimings),
		TimingDifference:   stats.StandardDev,
		SampleSize:         t.SampleSize,
		ExecutionTime:      time.Since(start),
		SuccessRate:        calculateRemoteSuccessRate(attackSuccessful, timingVariation),
		TimingStatistics:   stats,
		AttackMetadata: map[string]string{
			"network_jitter":    fmt.Sprintf("%.2f ms", stats.StandardDev/1000000),
			"timing_variation":  fmt.Sprintf("%.2f%%", timingVariation*100),
			"attack_difficulty": "high",
		},
	}
	
	t.Results = append(t.Results, result)
	return result
}

// Helper functions for timing attack measurements

func (t *TimingAttack) measurePasswordTimingAttack(compareFunc func([]byte, []byte) bool, correct, incorrect string) TimingAttackResult {
	var correctTimings, incorrectTimings []time.Duration
	
	// Measure timing for correct password
	for i := 0; i < t.SampleSize/2; i++ {
		start := time.Now()
		compareFunc([]byte(correct), []byte(correct))
		correctTimings = append(correctTimings, time.Since(start))
	}
	
	// Measure timing for incorrect password
	for i := 0; i < t.SampleSize/2; i++ {
		start := time.Now()
		compareFunc([]byte(correct), []byte(incorrect))
		incorrectTimings = append(incorrectTimings, time.Since(start))
	}
	
	// Statistical analysis
	correctStats := calculateTimingStatistics(correctTimings)
	incorrectStats := calculateTimingStatistics(incorrectTimings)
	
	// Calculate timing difference
	timingDifference := math.Abs(correctStats.Mean - incorrectStats.Mean)
	
	// T-test for statistical significance
	tStat := (correctStats.Mean - incorrectStats.Mean) / 
		math.Sqrt((correctStats.Variance/float64(len(correctTimings))) + 
			(incorrectStats.Variance/float64(len(incorrectTimings))))
	
	pValue := calculatePValue(tStat, len(correctTimings)+len(incorrectTimings)-2)
	
	// Attack is successful if timing difference is statistically significant
	attackSuccessful := pValue < 0.05 && timingDifference > 1000 // 1μs threshold
	
	return TimingAttackResult{
		Success:            attackSuccessful,
		ConfidenceLevel:    1.0 - pValue,
		StatisticalSignif:  pValue,
		TimingDifference:   timingDifference,
		SuccessRate:        calculateTimingSuccessRate(attackSuccessful, pValue),
		TimingStatistics:   correctStats,
	}
}

func (t *TimingAttack) measureKeyTimingAttack(compareFunc func([]byte, []byte) bool, correctKey []byte) TimingAttackResult {
	var matchTimings, mismatchTimings []time.Duration
	
	// Measure timing for matching keys
	for i := 0; i < t.SampleSize/2; i++ {
		testKey := make([]byte, len(correctKey))
		copy(testKey, correctKey)
		
		start := time.Now()
		compareFunc(correctKey, testKey)
		matchTimings = append(matchTimings, time.Since(start))
	}
	
	// Measure timing for mismatched keys (different at various positions)
	for i := 0; i < t.SampleSize/2; i++ {
		testKey := make([]byte, len(correctKey))
		copy(testKey, correctKey)
		
		// Modify key at different positions
		position := i % len(correctKey)
		testKey[position] = testKey[position] ^ 0xFF
		
		start := time.Now()
		compareFunc(correctKey, testKey)
		mismatchTimings = append(mismatchTimings, time.Since(start))
	}
	
	// Statistical analysis
	matchStats := calculateTimingStatistics(matchTimings)
	mismatchStats := calculateTimingStatistics(mismatchTimings)
	
	timingDifference := math.Abs(matchStats.Mean - mismatchStats.Mean)
	
	// T-test
	tStat := (matchStats.Mean - mismatchStats.Mean) / 
		math.Sqrt((matchStats.Variance/float64(len(matchTimings))) + 
			(mismatchStats.Variance/float64(len(mismatchTimings))))
	
	pValue := calculatePValue(tStat, len(matchTimings)+len(mismatchTimings)-2)
	
	attackSuccessful := pValue < 0.05 && timingDifference > 500 // 0.5μs threshold
	
	return TimingAttackResult{
		Success:            attackSuccessful,
		ConfidenceLevel:    1.0 - pValue,
		StatisticalSignif:  pValue,
		TimingDifference:   timingDifference,
		SuccessRate:        calculateTimingSuccessRate(attackSuccessful, pValue),
		TimingStatistics:   matchStats,
	}
}

func (t *TimingAttack) measureHashTimingAttack(compareFunc func([]byte, []byte) bool, correctHash string) TimingAttackResult {
	var timings []time.Duration
	
	// Generate test hashes with varying prefixes
	for i := 0; i < t.SampleSize; i++ {
		// Create hash with correct prefix of varying length
		prefixLength := i % len(correctHash)
		testHash := correctHash[:prefixLength] + generateRandomString(len(correctHash)-prefixLength)
		
		start := time.Now()
		compareFunc([]byte(correctHash), []byte(testHash))
		timings = append(timings, time.Since(start))
	}
	
	// Analyze timing patterns
	stats := calculateTimingStatistics(timings)
	
	// Check for correlation between prefix length and timing
	correlation := calculatePrefixTimingCorrelation(timings, correctHash)
	
	attackSuccessful := math.Abs(correlation) > 0.3 // Significant correlation
	
	return TimingAttackResult{
		Success:            attackSuccessful,
		ConfidenceLevel:    math.Abs(correlation),
		StatisticalSignif:  1.0 - math.Abs(correlation),
		TimingDifference:   stats.StandardDev,
		SuccessRate:        calculateTimingSuccessRate(attackSuccessful, 1.0-math.Abs(correlation)),
		TimingStatistics:   stats,
	}
}

// Statistical analysis functions

func calculateTimingStatistics(timings []time.Duration) TimingStatistics {
	if len(timings) == 0 {
		return TimingStatistics{}
	}
	
	// Convert to float64 for calculations
	values := make([]float64, len(timings))
	for i, t := range timings {
		values[i] = float64(t.Nanoseconds())
	}
	
	// Sort for percentile calculations
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	// Calculate statistics
	mean := calculateMean(values)
	variance := calculateVariance(values, mean)
	stdDev := math.Sqrt(variance)
	
	median := sorted[len(sorted)/2]
	min := sorted[0]
	max := sorted[len(sorted)-1]
	
	q1 := sorted[len(sorted)/4]
	q3 := sorted[3*len(sorted)/4]
	
	// Count outliers (values outside 1.5 * IQR)
	iqr := q3 - q1
	outliers := 0
	for _, v := range values {
		if v < q1-1.5*iqr || v > q3+1.5*iqr {
			outliers++
		}
	}
	
	// Create histogram
	histogram := createHistogram(values, 10)
	
	return TimingStatistics{
		Mean:         mean,
		StandardDev:  stdDev,
		Variance:     variance,
		Median:       median,
		Min:          min,
		Max:          max,
		Quartiles:    [3]float64{q1, median, q3},
		Outliers:     outliers,
		Distribution: sorted,
		Histogram:    histogram,
	}
}

func calculateMean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateVariance(values []float64, mean float64) float64 {
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	return sumSquares / float64(len(values)-1)
}

func calculatePValue(tStat float64, degreesOfFreedom int) float64 {
	// Simplified p-value calculation
	absTStat := math.Abs(tStat)
	if absTStat > 2.576 {
		return 0.01
	} else if absTStat > 1.96 {
		return 0.05
	} else if absTStat > 1.645 {
		return 0.10
	}
	return 0.20
}

func calculateConfidenceLevel(stats TimingStatistics) float64 {
	// Calculate confidence based on statistical properties
	cv := stats.StandardDev / stats.Mean // Coefficient of variation
	
	if cv < 0.1 {
		return 0.95
	} else if cv < 0.2 {
		return 0.90
	} else if cv < 0.3 {
		return 0.80
	}
	return 0.70
}

func calculateStatisticalSignificance(timings []time.Duration) float64 {
	// Simple statistical significance calculation
	if len(timings) < 2 {
		return 1.0
	}
	
	stats := calculateTimingStatistics(timings)
	cv := stats.StandardDev / stats.Mean
	
	if cv < 0.1 {
		return 0.01
	} else if cv < 0.2 {
		return 0.05
	} else if cv < 0.3 {
		return 0.10
	}
	return 0.20
}

func calculateSuccessRate(vulnerable, secure bool) float64 {
	if vulnerable && !secure {
		return 0.9 // 90% success rate
	} else if vulnerable && secure {
		return 0.5 // 50% success rate
	}
	return 0.1 // 10% success rate
}

func calculateTimingSuccessRate(successful bool, pValue float64) float64 {
	if successful {
		return 1.0 - pValue
	}
	return pValue
}

func calculateRemoteSuccessRate(successful bool, variation float64) float64 {
	if successful {
		return math.Max(0.1, 1.0-variation)
	}
	return 0.1
}

func calculatePrefixTimingCorrelation(timings []time.Duration, correctHash string) float64 {
	// Calculate correlation between prefix length and timing
	if len(timings) == 0 {
		return 0.0
	}
	
	// Simple correlation calculation
	var sumX, sumY, sumXY, sumX2, sumY2 float64
	n := float64(len(timings))
	
	for i, timing := range timings {
		x := float64(i % len(correctHash)) // Prefix length
		y := float64(timing.Nanoseconds())
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}
	
	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))
	
	if denominator == 0 {
		return 0.0
	}
	
	return numerator / denominator
}

func createHistogram(values []float64, bins int) map[string]int {
	if len(values) == 0 {
		return make(map[string]int)
	}
	
	min := values[0]
	max := values[0]
	
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	
	binSize := (max - min) / float64(bins)
	histogram := make(map[string]int)
	
	for _, v := range values {
		binIndex := int((v - min) / binSize)
		if binIndex >= bins {
			binIndex = bins - 1
		}
		
		binLabel := fmt.Sprintf("%.0f-%.0f", min+float64(binIndex)*binSize, min+float64(binIndex+1)*binSize)
		histogram[binLabel]++
	}
	
	return histogram
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	
	return string(result)
}

// RunAllTimingAttacks executes all timing attack simulations
func (t *TimingAttack) RunAllTimingAttacks() []TimingAttackResult {
	fmt.Println("Starting Timing Attack Simulations...")
	
	t.SimulatePasswordTimingAttack()
	t.SimulateKeyComparisonAttack()
	t.SimulateHashComparisonAttack()
	t.SimulateRemoteTimingAttack()
	
	return t.Results
}

// GenerateTimingReport generates a comprehensive timing attack report
func (t *TimingAttack) GenerateTimingReport() string {
	report := "=== Timing Attack Simulation Report ===\n\n"
	
	totalAttacks := len(t.Results)
	successfulAttacks := 0
	avgConfidence := 0.0
	
	for _, result := range t.Results {
		if result.Success {
			successfulAttacks++
		}
		avgConfidence += result.ConfidenceLevel
	}
	
	if totalAttacks > 0 {
		avgConfidence /= float64(totalAttacks)
	}
	
	report += fmt.Sprintf("Total timing attacks: %d\n", totalAttacks)
	report += fmt.Sprintf("Successful attacks: %d\n", successfulAttacks)
	report += fmt.Sprintf("Success rate: %.2f%%\n", float64(successfulAttacks)/float64(totalAttacks)*100)
	report += fmt.Sprintf("Average confidence: %.2f%%\n", avgConfidence*100)
	report += "\n=== Individual Attack Results ===\n\n"
	
	for _, result := range t.Results {
		report += fmt.Sprintf("Attack: %s\n", result.AttackType)
		report += fmt.Sprintf("  Success: %v\n", result.Success)
		report += fmt.Sprintf("  Confidence: %.2f%%\n", result.ConfidenceLevel*100)
		report += fmt.Sprintf("  Statistical significance: %.4f\n", result.StatisticalSignif)
		report += fmt.Sprintf("  Timing difference: %.2f ns\n", result.TimingDifference)
		report += fmt.Sprintf("  Sample size: %d\n", result.SampleSize)
		report += fmt.Sprintf("  Success rate: %.2f%%\n", result.SuccessRate*100)
		report += fmt.Sprintf("  Execution time: %v\n", result.ExecutionTime)
		report += "\n"
	}
	
	return report
}

// GetTimingResults returns all timing attack results
func (t *TimingAttack) GetTimingResults() []TimingAttackResult {
	return t.Results
}

// ClearTimingResults clears all timing attack results
func (t *TimingAttack) ClearTimingResults() {
	t.Results = make([]TimingAttackResult, 0)
}