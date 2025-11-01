package attacks

import (
	"crypto/rand"
	"fmt"
	"math"
	"sort"
	"time"
)

// PowerAnalysisAttack represents a power analysis attack simulation
type PowerAnalysisAttack struct {
	SampleSize     int
	DeviceType     string
	PowerModel     PowerModel
	Results        []PowerAnalysisResult
}

// PowerAnalysisResult represents the result of a power analysis attack
type PowerAnalysisResult struct {
	AttackType        string              `json:"attackType"`
	Success           bool                `json:"success"`
	ConfidenceLevel   float64             `json:"confidenceLevel"`
	KeyRecoveryRate   float64             `json:"keyRecoveryRate"`
	PowerDifference   float64             `json:"powerDifferenceWatts"`
	SampleSize        int                 `json:"sampleSize"`
	ExecutionTime     time.Duration       `json:"executionTime"`
	ErrorMessage      string              `json:"errorMessage,omitempty"`
	AttackMetadata    map[string]string   `json:"attackMetadata"`
	PowerStatistics   PowerStatistics     `json:"powerStatistics"`
	KeyBits           []KeyBitRecovery    `json:"keyBits"`
}

// PowerStatistics contains statistical analysis of power measurements
type PowerStatistics struct {
	Mean              float64           `json:"mean"`
	StandardDev       float64           `json:"standardDeviation"`
	Variance          float64           `json:"variance"`
	Median            float64           `json:"median"`
	Min               float64           `json:"min"`
	Max               float64           `json:"max"`
	SNR               float64           `json:"signalToNoiseRatio"`
	Correlation       float64           `json:"correlation"`
	PowerTrace        []float64         `json:"powerTrace"`
	Frequencies       []float64         `json:"frequencies"`
	Histogram         map[string]int    `json:"histogram"`
}

// KeyBitRecovery represents recovery information for individual key bits
type KeyBitRecovery struct {
	BitPosition   int     `json:"bitPosition"`
	RecoveredBit  int     `json:"recoveredBit"`
	Confidence    float64 `json:"confidence"`
	PowerDiff     float64 `json:"powerDifference"`
	Attempts      int     `json:"attempts"`
}

// PowerModel simulates power consumption patterns
type PowerModel struct {
	BasePower        float64 `json:"basePower"`
	BitPowerFactor   float64 `json:"bitPowerFactor"`
	NoiseLevel       float64 `json:"noiseLevel"`
	FrequencyFactor  float64 `json:"frequencyFactor"`
	TemperatureFactor float64 `json:"temperatureFactor"`
}

// NewPowerAnalysisAttack creates a new power analysis attack simulator
func NewPowerAnalysisAttack(sampleSize int, deviceType string) *PowerAnalysisAttack {
	var powerModel PowerModel
	
	switch deviceType {
	case "mobile":
		powerModel = PowerModel{
			BasePower:        1.5,  // 1.5W base power
			BitPowerFactor:   0.01, // 10mW per bit operation
			NoiseLevel:       0.05, // 5% noise
			FrequencyFactor:  0.02, // 2% variation per GHz
			TemperatureFactor: 0.01, // 1% per degree C
		}
	case "iot":
		powerModel = PowerModel{
			BasePower:        0.1,   // 100mW base power
			BitPowerFactor:   0.001, // 1mW per bit operation
			NoiseLevel:       0.02,  // 2% noise
			FrequencyFactor:  0.01,  // 1% variation per GHz
			TemperatureFactor: 0.005, // 0.5% per degree C
		}
	case "server":
		powerModel = PowerModel{
			BasePower:        50.0,  // 50W base power
			BitPowerFactor:   0.1,   // 100mW per bit operation
			NoiseLevel:       0.1,   // 10% noise
			FrequencyFactor:  0.05,  // 5% variation per GHz
			TemperatureFactor: 0.02,  // 2% per degree C
		}
	default:
		powerModel = PowerModel{
			BasePower:        1.0,
			BitPowerFactor:   0.01,
			NoiseLevel:       0.05,
			FrequencyFactor:  0.02,
			TemperatureFactor: 0.01,
		}
	}
	
	return &PowerAnalysisAttack{
		SampleSize:     sampleSize,
		DeviceType:     deviceType,
		PowerModel:     powerModel,
		Results:        make([]PowerAnalysisResult, 0),
	}
}

// SimulateSimplePowerAnalysis simulates Simple Power Analysis (SPA) attack
func (p *PowerAnalysisAttack) SimulateSimplePowerAnalysis() PowerAnalysisResult {
	start := time.Now()
	
	// Generate test key
	testKey := make([]byte, 16) // 128-bit key
	rand.Read(testKey)
	
	// Simulate power traces for different operations
	var powerTraces [][]float64
	keyBits := make([]KeyBitRecovery, 0)
	
	// For each bit in the key
	for byteIndex := 0; byteIndex < len(testKey); byteIndex++ {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			bit := (testKey[byteIndex] >> bitIndex) & 1
			
			// Simulate power consumption for this bit
			powerTrace := p.generatePowerTrace(bit, "bit_operation")
			powerTraces = append(powerTraces, powerTrace)
			
			// Analyze power trace to recover bit
			recoveredBit, confidence := p.analyzePowerTrace(powerTrace, bit)
			
			keyBits = append(keyBits, KeyBitRecovery{
				BitPosition:  byteIndex*8 + bitIndex,
				RecoveredBit: recoveredBit,
				Confidence:   confidence,
				PowerDiff:    p.calculatePowerDifference(powerTrace),
				Attempts:     1,
			})
		}
	}
	
	// Calculate overall statistics
	stats := p.calculatePowerStatistics(powerTraces)
	keyRecoveryRate := p.calculateKeyRecoveryRate(keyBits)
	
	// SPA attack is successful if we can recover a significant portion of the key
	attackSuccessful := keyRecoveryRate > 0.7 // 70% threshold
	
	result := PowerAnalysisResult{
		AttackType:        "SimplePowerAnalysis",
		Success:           attackSuccessful,
		ConfidenceLevel:   p.calculateOverallConfidence(keyBits),
		KeyRecoveryRate:   keyRecoveryRate,
		PowerDifference:   stats.StandardDev,
		SampleSize:        p.SampleSize,
		ExecutionTime:     time.Since(start),
		PowerStatistics:   stats,
		KeyBits:           keyBits,
		AttackMetadata: map[string]string{
			"key_length":        fmt.Sprintf("%d", len(testKey)*8),
			"device_type":       p.DeviceType,
			"base_power":        fmt.Sprintf("%.2f W", p.PowerModel.BasePower),
			"noise_level":       fmt.Sprintf("%.2f%%", p.PowerModel.NoiseLevel*100),
			"recovery_rate":     fmt.Sprintf("%.2f%%", keyRecoveryRate*100),
		},
	}
	
	p.Results = append(p.Results, result)
	return result
}

// SimulateDifferentialPowerAnalysis simulates Differential Power Analysis (DPA) attack
func (p *PowerAnalysisAttack) SimulateDifferentialPowerAnalysis() PowerAnalysisResult {
	start := time.Now()
	
	// Generate multiple plaintexts and a fixed key
	key := make([]byte, 16)
	rand.Read(key)
	
	var powerTraces [][]float64
	var plaintexts [][]byte
	keyBits := make([]KeyBitRecovery, 0)
	
	// Collect power traces for different plaintexts
	for i := 0; i < p.SampleSize; i++ {
		plaintext := make([]byte, 16)
		rand.Read(plaintext)
		plaintexts = append(plaintexts, plaintext)
		
		// Simulate encryption operation
		powerTrace := p.simulateEncryptionPowerTrace(plaintext, key)
		powerTraces = append(powerTraces, powerTrace)
	}
	
	// DPA attack: correlate power consumption with key hypotheses
	for byteIndex := 0; byteIndex < len(key); byteIndex++ {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			bitPosition := byteIndex*8 + bitIndex
			
			// Test all possible key bit values (0 and 1)
			correlations := make([]float64, 2)
			
			for keyBitHypothesis := 0; keyBitHypothesis < 2; keyBitHypothesis++ {
				// Calculate correlation between power traces and key hypothesis
				correlation := p.calculateDPACorrelation(powerTraces, plaintexts, byteIndex, bitIndex, keyBitHypothesis)
				correlations[keyBitHypothesis] = math.Abs(correlation)
			}
			
			// Choose the hypothesis with higher correlation
			recoveredBit := 0
			if correlations[1] > correlations[0] {
				recoveredBit = 1
			}
			
			confidence := math.Abs(correlations[1] - correlations[0])
			actualBit := (key[byteIndex] >> bitIndex) & 1
			
			keyBits = append(keyBits, KeyBitRecovery{
				BitPosition:  bitPosition,
				RecoveredBit: recoveredBit,
				Confidence:   confidence,
				PowerDiff:    correlations[recoveredBit],
				Attempts:     p.SampleSize,
			})
		}
	}
	
	// Calculate statistics
	stats := p.calculatePowerStatistics(powerTraces)
	keyRecoveryRate := p.calculateKeyRecoveryRate(keyBits)
	
	// DPA attack is successful if correlation is strong enough
	attackSuccessful := keyRecoveryRate > 0.6 && stats.Correlation > 0.3
	
	result := PowerAnalysisResult{
		AttackType:        "DifferentialPowerAnalysis",
		Success:           attackSuccessful,
		ConfidenceLevel:   p.calculateOverallConfidence(keyBits),
		KeyRecoveryRate:   keyRecoveryRate,
		PowerDifference:   stats.StandardDev,
		SampleSize:        p.SampleSize,
		ExecutionTime:     time.Since(start),
		PowerStatistics:   stats,
		KeyBits:           keyBits,
		AttackMetadata: map[string]string{
			"attack_method":     "correlation",
			"correlation":       fmt.Sprintf("%.3f", stats.Correlation),
			"sample_size":       fmt.Sprintf("%d", p.SampleSize),
			"key_recovery":      fmt.Sprintf("%.2f%%", keyRecoveryRate*100),
		},
	}
	
	p.Results = append(p.Results, result)
	return result
}

// SimulateCorrelationPowerAnalysis simulates Correlation Power Analysis (CPA) attack
func (p *PowerAnalysisAttack) SimulateCorrelationPowerAnalysis() PowerAnalysisResult {
	start := time.Now()
	
	// Generate test data
	key := make([]byte, 16)
	rand.Read(key)
	
	var powerTraces [][]float64
	var plaintexts [][]byte
	keyBits := make([]KeyBitRecovery, 0)
	
	// Collect power traces
	for i := 0; i < p.SampleSize; i++ {
		plaintext := make([]byte, 16)
		rand.Read(plaintext)
		plaintexts = append(plaintexts, plaintext)
		
		// Simulate S-box operations (critical for CPA)
		powerTrace := p.simulateSBoxPowerTrace(plaintext, key)
		powerTraces = append(powerTraces, powerTrace)
	}
	
	// CPA attack: use Hamming weight model
	for byteIndex := 0; byteIndex < len(key); byteIndex++ {
		bestCorrelation := 0.0
		bestKeyByte := 0
		
		// Test all possible key byte values
		for keyHypothesis := 0; keyHypothesis < 256; keyHypothesis++ {
			// Calculate predicted power consumption using Hamming weight model
			var predictions []float64
			for _, plaintext := range plaintexts {
				sboxInput := plaintext[byteIndex] ^ byte(keyHypothesis)
				sboxOutput := p.simulateSBox(sboxInput)
				hammingWeight := p.calculateHammingWeight(sboxOutput)
				predictions = append(predictions, float64(hammingWeight))
			}
			
			// Calculate correlation between predictions and actual power traces
			correlation := p.calculateCorrelation(predictions, powerTraces, byteIndex)
			
			if math.Abs(correlation) > math.Abs(bestCorrelation) {
				bestCorrelation = correlation
				bestKeyByte = keyHypothesis
			}
		}
		
		// Convert byte to individual bits
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			recoveredBit := (bestKeyByte >> bitIndex) & 1
			actualBit := (key[byteIndex] >> bitIndex) & 1
			
			keyBits = append(keyBits, KeyBitRecovery{
				BitPosition:  byteIndex*8 + bitIndex,
				RecoveredBit: recoveredBit,
				Confidence:   math.Abs(bestCorrelation),
				PowerDiff:    bestCorrelation,
				Attempts:     256, // Tested all byte values
			})
		}
	}
	
	// Calculate statistics
	stats := p.calculatePowerStatistics(powerTraces)
	keyRecoveryRate := p.calculateKeyRecoveryRate(keyBits)
	
	// CPA attack is successful if correlation is strong
	attackSuccessful := keyRecoveryRate > 0.8 && stats.Correlation > 0.5
	
	result := PowerAnalysisResult{
		AttackType:        "CorrelationPowerAnalysis",
		Success:           attackSuccessful,
		ConfidenceLevel:   p.calculateOverallConfidence(keyBits),
		KeyRecoveryRate:   keyRecoveryRate,
		PowerDifference:   stats.StandardDev,
		SampleSize:        p.SampleSize,
		ExecutionTime:     time.Since(start),
		PowerStatistics:   stats,
		KeyBits:           keyBits,
		AttackMetadata: map[string]string{
			"model":             "hamming_weight",
			"correlation":       fmt.Sprintf("%.3f", stats.Correlation),
			"attack_complexity": "high",
			"key_recovery":      fmt.Sprintf("%.2f%%", keyRecoveryRate*100),
		},
	}
	
	p.Results = append(p.Results, result)
	return result
}

// SimulateElectromagneticAnalysis simulates Electromagnetic Analysis (EMA) attack
func (p *PowerAnalysisAttack) SimulateElectromagneticAnalysis() PowerAnalysisResult {
	start := time.Now()
	
	// Generate test data
	key := make([]byte, 16)
	rand.Read(key)
	
	var emTraces [][]float64
	keyBits := make([]KeyBitRecovery, 0)
	
	// Collect electromagnetic traces
	for i := 0; i < p.SampleSize; i++ {
		plaintext := make([]byte, 16)
		rand.Read(plaintext)
		
		// Simulate EM emissions during encryption
		emTrace := p.simulateEMTrace(plaintext, key)
		emTraces = append(emTraces, emTrace)
	}
	
	// Analyze EM traces for key recovery
	for byteIndex := 0; byteIndex < len(key); byteIndex++ {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			bitPosition := byteIndex*8 + bitIndex
			
			// Analyze EM signature for this bit
			emSignature := p.analyzeEMSignature(emTraces, byteIndex, bitIndex)
			recoveredBit := p.extractBitFromEMSignature(emSignature)
			confidence := p.calculateEMConfidence(emSignature)
			
			keyBits = append(keyBits, KeyBitRecovery{
				BitPosition:  bitPosition,
				RecoveredBit: recoveredBit,
				Confidence:   confidence,
				PowerDiff:    emSignature.Amplitude,
				Attempts:     p.SampleSize,
			})
		}
	}
	
	// Calculate statistics
	stats := p.calculateEMStatistics(emTraces)
	keyRecoveryRate := p.calculateKeyRecoveryRate(keyBits)
	
	// EMA attack success depends on EM leakage
	attackSuccessful := keyRecoveryRate > 0.5 && stats.SNR > 2.0
	
	result := PowerAnalysisResult{
		AttackType:        "ElectromagneticAnalysis",
		Success:           attackSuccessful,
		ConfidenceLevel:   p.calculateOverallConfidence(keyBits),
		KeyRecoveryRate:   keyRecoveryRate,
		PowerDifference:   stats.StandardDev,
		SampleSize:        p.SampleSize,
		ExecutionTime:     time.Since(start),
		PowerStatistics:   stats,
		KeyBits:           keyBits,
		AttackMetadata: map[string]string{
			"em_frequency":      "2.4GHz",
			"snr":              fmt.Sprintf("%.2f dB", stats.SNR),
			"leakage_type":     "data_dependent",
			"key_recovery":     fmt.Sprintf("%.2f%%", keyRecoveryRate*100),
		},
	}
	
	p.Results = append(p.Results, result)
	return result
}

// Helper functions for power analysis

func (p *PowerAnalysisAttack) generatePowerTrace(bit int, operation string) []float64 {
	traceLength := 1000 // Sample points
	trace := make([]float64, traceLength)
	
	for i := 0; i < traceLength; i++ {
		basePower := p.PowerModel.BasePower
		
		// Add bit-dependent power consumption
		if bit == 1 {
			basePower += p.PowerModel.BitPowerFactor
		}
		
		// Add noise
		noise := (rand.Float64() - 0.5) * p.PowerModel.NoiseLevel * basePower
		
		// Add frequency-dependent variations
		freqVariation := math.Sin(float64(i)/100.0) * p.PowerModel.FrequencyFactor * basePower
		
		trace[i] = basePower + noise + freqVariation
	}
	
	return trace
}

func (p *PowerAnalysisAttack) simulateEncryptionPowerTrace(plaintext, key []byte) []float64 {
	traceLength := 1000
	trace := make([]float64, traceLength)
	
	for i := 0; i < traceLength; i++ {
		power := p.PowerModel.BasePower
		
		// Simulate power consumption based on data being processed
		if i < len(plaintext) {
			// XOR operation power
			xorResult := plaintext[i%len(plaintext)] ^ key[i%len(key)]
			hammingWeight := p.calculateHammingWeight(xorResult)
			power += float64(hammingWeight) * p.PowerModel.BitPowerFactor
		}
		
		// Add noise
		noise := (rand.Float64() - 0.5) * p.PowerModel.NoiseLevel * power
		trace[i] = power + noise
	}
	
	return trace
}

func (p *PowerAnalysisAttack) simulateSBoxPowerTrace(plaintext, key []byte) []float64 {
	traceLength := 256 // One sample per S-box operation
	trace := make([]float64, traceLength)
	
	for i := 0; i < traceLength; i++ {
		if i < len(plaintext) {
			sboxInput := plaintext[i%len(plaintext)] ^ key[i%len(key)]
			sboxOutput := p.simulateSBox(sboxInput)
			hammingWeight := p.calculateHammingWeight(sboxOutput)
			
			power := p.PowerModel.BasePower + float64(hammingWeight)*p.PowerModel.BitPowerFactor
			noise := (rand.Float64() - 0.5) * p.PowerModel.NoiseLevel * power
			trace[i] = power + noise
		} else {
			trace[i] = p.PowerModel.BasePower
		}
	}
	
	return trace
}

func (p *PowerAnalysisAttack) simulateEMTrace(plaintext, key []byte) []float64 {
	traceLength := 1000
	trace := make([]float64, traceLength)
	
	for i := 0; i < traceLength; i++ {
		// EM emissions are similar to power but with different characteristics
		basePower := p.PowerModel.BasePower * 0.1 // EM is typically weaker
		
		if i < len(plaintext) {
			dataValue := plaintext[i%len(plaintext)] ^ key[i%len(key)]
			hammingWeight := p.calculateHammingWeight(dataValue)
			basePower += float64(hammingWeight) * p.PowerModel.BitPowerFactor * 0.5
		}
		
		// EM has different noise characteristics
		noise := (rand.Float64() - 0.5) * p.PowerModel.NoiseLevel * basePower * 2.0
		
		// Add frequency-dependent EM emissions
		emFreq := math.Sin(float64(i)/50.0) * basePower * 0.1
		
		trace[i] = basePower + noise + emFreq
	}
	
	return trace
}

func (p *PowerAnalysisAttack) simulateSBox(input byte) byte {
	// Simplified S-box simulation (not a real cryptographic S-box)
	sbox := []byte{
		0x63, 0x7C, 0x77, 0x7B, 0xF2, 0x6B, 0x6F, 0xC5,
		0x30, 0x01, 0x67, 0x2B, 0xFE, 0xD7, 0xAB, 0x76,
		// ... (simplified, would normally be 256 bytes)
	}
	
	return sbox[input%16] // Simplified lookup
}

func (p *PowerAnalysisAttack) calculateHammingWeight(value byte) int {
	weight := 0
	for i := 0; i < 8; i++ {
		if (value>>i)&1 == 1 {
			weight++
		}
	}
	return weight
}

func (p *PowerAnalysisAttack) analyzePowerTrace(trace []float64, actualBit int) (int, float64) {
	// Simple analysis: check if power is above or below threshold
	avgPower := 0.0
	for _, power := range trace {
		avgPower += power
	}
	avgPower /= float64(len(trace))
	
	threshold := p.PowerModel.BasePower + p.PowerModel.BitPowerFactor/2
	
	recoveredBit := 0
	if avgPower > threshold {
		recoveredBit = 1
	}
	
	// Calculate confidence based on distance from threshold
	confidence := math.Abs(avgPower-threshold) / (p.PowerModel.BitPowerFactor / 2)
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	return recoveredBit, confidence
}

func (p *PowerAnalysisAttack) calculatePowerDifference(trace []float64) float64 {
	if len(trace) < 2 {
		return 0.0
	}
	
	max := trace[0]
	min := trace[0]
	
	for _, power := range trace {
		if power > max {
			max = power
		}
		if power < min {
			min = power
		}
	}
	
	return max - min
}

func (p *PowerAnalysisAttack) calculateDPACorrelation(traces [][]float64, plaintexts [][]byte, byteIndex, bitIndex, keyBitHypothesis int) float64 {
	if len(traces) != len(plaintexts) || len(traces) == 0 {
		return 0.0
	}
	
	var group0, group1 []float64
	
	for i, plaintext := range plaintexts {
		if byteIndex < len(plaintext) {
			// Predict intermediate value
			intermediate := plaintext[byteIndex] ^ byte(keyBitHypothesis<<bitIndex)
			bit := (intermediate >> bitIndex) & 1
			
			// Extract relevant power sample
			if bitIndex < len(traces[i]) {
				power := traces[i][bitIndex]
				if bit == 0 {
					group0 = append(group0, power)
				} else {
					group1 = append(group1, power)
				}
			}
		}
	}
	
	if len(group0) == 0 || len(group1) == 0 {
		return 0.0
	}
	
	// Calculate mean difference
	mean0 := calculateMeanFloat64(group0)
	mean1 := calculateMeanFloat64(group1)
	
	return mean1 - mean0
}

func (p *PowerAnalysisAttack) calculateCorrelation(predictions []float64, traces [][]float64, byteIndex int) float64 {
	if len(predictions) != len(traces) || len(predictions) == 0 {
		return 0.0
	}
	
	// Extract power samples for the relevant time point
	var powerSamples []float64
	for _, trace := range traces {
		if byteIndex < len(trace) {
			powerSamples = append(powerSamples, trace[byteIndex])
		}
	}
	
	if len(powerSamples) != len(predictions) {
		return 0.0
	}
	
	// Calculate Pearson correlation coefficient
	return calculatePearsonCorrelation(predictions, powerSamples)
}

func (p *PowerAnalysisAttack) calculatePowerStatistics(traces [][]float64) PowerStatistics {
	if len(traces) == 0 {
		return PowerStatistics{}
	}
	
	// Flatten all traces into one array
	var allPowers []float64
	for _, trace := range traces {
		allPowers = append(allPowers, trace...)
	}
	
	if len(allPowers) == 0 {
		return PowerStatistics{}
	}
	
	// Calculate statistics
	mean := calculateMeanFloat64(allPowers)
	variance := calculateVarianceFloat64(allPowers, mean)
	stdDev := math.Sqrt(variance)
	
	// Sort for median calculation
	sorted := make([]float64, len(allPowers))
	copy(sorted, allPowers)
	sort.Float64s(sorted)
	
	median := sorted[len(sorted)/2]
	min := sorted[0]
	max := sorted[len(sorted)-1]
	
	// Calculate SNR (simplified)
	snr := mean / stdDev
	
	// Calculate correlation (simplified - correlation with first trace)
	correlation := 0.0
	if len(traces) > 1 {
		correlation = calculatePearsonCorrelation(traces[0], traces[1])
	}
	
	return PowerStatistics{
		Mean:        mean,
		StandardDev: stdDev,
		Variance:    variance,
		Median:      median,
		Min:         min,
		Max:         max,
		SNR:         snr,
		Correlation: correlation,
		PowerTrace:  traces[0], // First trace as example
	}
}

func (p *PowerAnalysisAttack) calculateEMStatistics(traces [][]float64) PowerStatistics {
	// Similar to power statistics but with EM-specific calculations
	stats := p.calculatePowerStatistics(traces)
	
	// EM-specific SNR calculation
	if len(traces) > 0 && len(traces[0]) > 0 {
		// Calculate SNR based on signal peaks
		maxSignal := 0.0
		for _, trace := range traces {
			for _, value := range trace {
				if math.Abs(value) > maxSignal {
					maxSignal = math.Abs(value)
				}
			}
		}
		
		stats.SNR = 20 * math.Log10(maxSignal/stats.StandardDev) // SNR in dB
	}
	
	return stats
}

func (p *PowerAnalysisAttack) calculateKeyRecoveryRate(keyBits []KeyBitRecovery) float64 {
	if len(keyBits) == 0 {
		return 0.0
	}
	
	// For simulation, assume we know the correct key
	// In practice, this would be determined by comparing with known test vectors
	correctBits := 0
	for _, keyBit := range keyBits {
		// Simplified: assume recovery is correct if confidence > 0.5
		if keyBit.Confidence > 0.5 {
			correctBits++
		}
	}
	
	return float64(correctBits) / float64(len(keyBits))
}

func (p *PowerAnalysisAttack) calculateOverallConfidence(keyBits []KeyBitRecovery) float64 {
	if len(keyBits) == 0 {
		return 0.0
	}
	
	totalConfidence := 0.0
	for _, keyBit := range keyBits {
		totalConfidence += keyBit.Confidence
	}
	
	return totalConfidence / float64(len(keyBits))
}

// EM-specific analysis functions

type EMSignature struct {
	Frequency float64
	Amplitude float64
	Phase     float64
}

func (p *PowerAnalysisAttack) analyzeEMSignature(traces [][]float64, byteIndex, bitIndex int) EMSignature {
	// Simplified EM signature analysis
	avgAmplitude := 0.0
	maxAmplitude := 0.0
	
	for _, trace := range traces {
		if byteIndex < len(trace) {
			amplitude := math.Abs(trace[byteIndex])
			avgAmplitude += amplitude
			if amplitude > maxAmplitude {
				maxAmplitude = amplitude
			}
		}
	}
	
	if len(traces) > 0 {
		avgAmplitude /= float64(len(traces))
	}
	
	return EMSignature{
		Frequency: 2.4e9, // 2.4 GHz
		Amplitude: avgAmplitude,
		Phase:     float64(bitIndex) * math.Pi / 8.0, // Phase based on bit position
	}
}

func (p *PowerAnalysisAttack) extractBitFromEMSignature(signature EMSignature) int {
	// Extract bit based on amplitude threshold
	threshold := p.PowerModel.BasePower * 0.05 // 5% of base power
	
	if signature.Amplitude > threshold {
		return 1
	}
	return 0
}

func (p *PowerAnalysisAttack) calculateEMConfidence(signature EMSignature) float64 {
	// Confidence based on amplitude and frequency characteristics
	amplitudeConfidence := signature.Amplitude / (p.PowerModel.BasePower * 0.1)
	if amplitudeConfidence > 1.0 {
		amplitudeConfidence = 1.0
	}
	
	return amplitudeConfidence
}

// Utility functions

func calculateMeanFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateVarianceFloat64(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0.0
	}
	
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	
	return sumSquares / float64(len(values)-1)
}

func calculatePearsonCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0.0
	}
	
	meanX := calculateMeanFloat64(x)
	meanY := calculateMeanFloat64(y)
	
	var sumXY, sumX2, sumY2 float64
	
	for i := 0; i < len(x); i++ {
		dx := x[i] - meanX
		dy := y[i] - meanY
		sumXY += dx * dy
		sumX2 += dx * dx
		sumY2 += dy * dy
	}
	
	denominator := math.Sqrt(sumX2 * sumY2)
	if denominator == 0 {
		return 0.0
	}
	
	return sumXY / denominator
}

// Main attack execution functions

func (p *PowerAnalysisAttack) RunAllPowerAttacks() []PowerAnalysisResult {
	fmt.Println("Starting Power Analysis Attack Simulations...")
	
	p.SimulateSimplePowerAnalysis()
	p.SimulateDifferentialPowerAnalysis()
	p.SimulateCorrelationPowerAnalysis()
	p.SimulateElectromagneticAnalysis()
	
	return p.Results
}

func (p *PowerAnalysisAttack) GeneratePowerReport() string {
	report := "=== Power Analysis Attack Simulation Report ===\n\n"
	
	totalAttacks := len(p.Results)
	successfulAttacks := 0
	avgKeyRecovery := 0.0
	avgConfidence := 0.0
	
	for _, result := range p.Results {
		if result.Success {
			successfulAttacks++
		}
		avgKeyRecovery += result.KeyRecoveryRate
		avgConfidence += result.ConfidenceLevel
	}
	
	if totalAttacks > 0 {
		avgKeyRecovery /= float64(totalAttacks)
		avgConfidence /= float64(totalAttacks)
	}
	
	report += fmt.Sprintf("Device type: %s\n", p.DeviceType)
	report += fmt.Sprintf("Total attacks: %d\n", totalAttacks)
	report += fmt.Sprintf("Successful attacks: %d\n", successfulAttacks)
	report += fmt.Sprintf("Success rate: %.2f%%\n", float64(successfulAttacks)/float64(totalAttacks)*100)
	report += fmt.Sprintf("Average key recovery rate: %.2f%%\n", avgKeyRecovery*100)
	report += fmt.Sprintf("Average confidence: %.2f%%\n", avgConfidence*100)
	report += fmt.Sprintf("Base power: %.2f W\n", p.PowerModel.BasePower)
	report += fmt.Sprintf("Noise level: %.2f%%\n", p.PowerModel.NoiseLevel*100)
	report += "\n=== Individual Attack Results ===\n\n"
	
	for _, result := range p.Results {
		report += fmt.Sprintf("Attack: %s\n", result.AttackType)
		report += fmt.Sprintf("  Success: %v\n", result.Success)
		report += fmt.Sprintf("  Key recovery rate: %.2f%%\n", result.KeyRecoveryRate*100)
		report += fmt.Sprintf("  Confidence level: %.2f%%\n", result.ConfidenceLevel*100)
		report += fmt.Sprintf("  Power difference: %.4f W\n", result.PowerDifference)
		report += fmt.Sprintf("  Sample size: %d\n", result.SampleSize)
		report += fmt.Sprintf("  Execution time: %v\n", result.ExecutionTime)
		report += fmt.Sprintf("  SNR: %.2f dB\n", result.PowerStatistics.SNR)
		report += "\n"
	}
	
	return report
}

func (p *PowerAnalysisAttack) GetPowerResults() []PowerAnalysisResult {
	return p.Results
}

func (p *PowerAnalysisAttack) ClearPowerResults() {
	p.Results = make([]PowerAnalysisResult, 0)
}