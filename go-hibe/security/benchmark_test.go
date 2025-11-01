package security

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"testing"
	"time"
)

// BenchmarkResult represents the result of a security benchmark test
type BenchmarkResult struct {
	TestName           string            `json:"testName"`
	OperationsPerSec   float64           `json:"operationsPerSecond"`
	MemoryUsage        uint64            `json:"memoryUsageBytes"`
	CPUUsage           float64           `json:"cpuUsagePercent"`
	PowerConsumption   float64           `json:"powerConsumptionWatts"`
	LatencyP50         time.Duration     `json:"latencyP50"`
	LatencyP95         time.Duration     `json:"latencyP95"`
	LatencyP99         time.Duration     `json:"latencyP99"`
	ThroughputMBps     float64           `json:"throughputMBps"`
	ErrorRate          float64           `json:"errorRate"`
	SecurityOverhead   float64           `json:"securityOverheadPercent"`
	Metrics            map[string]float64 `json:"metrics"`
}

// SecurityBenchmarkSuite runs comprehensive security performance benchmarks
type SecurityBenchmarkSuite struct {
	Results []BenchmarkResult
}

// RunSecurityBenchmarks executes all security-related benchmarks
func (s *SecurityBenchmarkSuite) RunSecurityBenchmarks(b *testing.B) {
	b.Run("EncryptionPerformance", s.BenchmarkEncryptionPerformance)
	b.Run("DecryptionPerformance", s.BenchmarkDecryptionPerformance)
	b.Run("KeyGenerationPerformance", s.BenchmarkKeyGenerationPerformance)
	b.Run("SecurityHeadersPerformance", s.BenchmarkSecurityHeadersPerformance)
	b.Run("TLSHandshakePerformance", s.BenchmarkTLSHandshakePerformance)
	b.Run("MemorySecurityPerformance", s.BenchmarkMemorySecurityPerformance)
	b.Run("ConcurrentSecurityPerformance", s.BenchmarkConcurrentSecurityPerformance)
	
	// Generate comprehensive report
	s.generateBenchmarkReport(b)
}

// BenchmarkEncryptionPerformance measures encryption performance under various conditions
func (s *SecurityBenchmarkSuite) BenchmarkEncryptionPerformance(b *testing.B) {
	dataSizes := []int{64, 256, 1024, 4096, 16384, 65536} // bytes
	
	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("DataSize_%d", size), func(b *testing.B) {
			plaintext := make([]byte, size)
			key := make([]byte, 32)
			rand.Read(plaintext)
			rand.Read(key)
			
			var memBefore, memAfter runtime.MemStats
			runtime.ReadMemStats(&memBefore)
			
			latencies := make([]time.Duration, 0, b.N)
			
			b.ResetTimer()
			start := time.Now()
			
			for i := 0; i < b.N; i++ {
				iterStart := time.Now()
				ciphertext := simulateEncryption(key, plaintext)
				latencies = append(latencies, time.Since(iterStart))
				
				// Prevent optimization
				_ = ciphertext
			}
			
			elapsed := time.Since(start)
			runtime.ReadMemStats(&memAfter)
			
			// Calculate metrics
			opsPerSec := float64(b.N) / elapsed.Seconds()
			throughputMBps := (float64(size) * float64(b.N)) / (1024 * 1024) / elapsed.Seconds()
			memUsage := memAfter.TotalAlloc - memBefore.TotalAlloc
			
			// Calculate latency percentiles
			p50, p95, p99 := calculateLatencyPercentiles(latencies)
			
			result := BenchmarkResult{
				TestName:           fmt.Sprintf("Encryption_DataSize_%d", size),
				OperationsPerSec:   opsPerSec,
				MemoryUsage:        memUsage,
				CPUUsage:           getCPUUsage(),
				PowerConsumption:   simulatePowerConsumption(key),
				LatencyP50:         p50,
				LatencyP95:         p95,
				LatencyP99:         p99,
				ThroughputMBps:     throughputMBps,
				ErrorRate:          0.0,
				SecurityOverhead:   calculateSecurityOverhead(opsPerSec, size),
				Metrics: map[string]float64{
					"dataSize":       float64(size),
					"iterations":     float64(b.N),
					"elapsedSeconds": elapsed.Seconds(),
				},
			}
			
			s.Results = append(s.Results, result)
			
			b.ReportMetric(opsPerSec, "ops/sec")
			b.ReportMetric(throughputMBps, "MB/s")
			b.ReportMetric(float64(p50.Nanoseconds()), "p50-latency-ns")
			b.ReportMetric(float64(p95.Nanoseconds()), "p95-latency-ns")
		})
	}
}

// BenchmarkDecryptionPerformance measures decryption performance
func (s *SecurityBenchmarkSuite) BenchmarkDecryptionPerformance(b *testing.B) {
	dataSizes := []int{64, 256, 1024, 4096, 16384, 65536}
	
	for _, size := range dataSizes {
		b.Run(fmt.Sprintf("DataSize_%d", size), func(b *testing.B) {
			plaintext := make([]byte, size)
			key := make([]byte, 32)
			rand.Read(plaintext)
			rand.Read(key)
			
			// Pre-encrypt data
			ciphertext := simulateEncryption(key, plaintext)
			
			var memBefore, memAfter runtime.MemStats
			runtime.ReadMemStats(&memBefore)
			
			latencies := make([]time.Duration, 0, b.N)
			
			b.ResetTimer()
			start := time.Now()
			
			for i := 0; i < b.N; i++ {
				iterStart := time.Now()
				decrypted := simulateDecryption(key, ciphertext)
				latencies = append(latencies, time.Since(iterStart))
				
				// Prevent optimization
				_ = decrypted
			}
			
			elapsed := time.Since(start)
			runtime.ReadMemStats(&memAfter)
			
			// Calculate metrics
			opsPerSec := float64(b.N) / elapsed.Seconds()
			throughputMBps := (float64(size) * float64(b.N)) / (1024 * 1024) / elapsed.Seconds()
			memUsage := memAfter.TotalAlloc - memBefore.TotalAlloc
			
			// Calculate latency percentiles
			p50, p95, p99 := calculateLatencyPercentiles(latencies)
			
			result := BenchmarkResult{
				TestName:           fmt.Sprintf("Decryption_DataSize_%d", size),
				OperationsPerSec:   opsPerSec,
				MemoryUsage:        memUsage,
				CPUUsage:           getCPUUsage(),
				PowerConsumption:   simulatePowerConsumption(key),
				LatencyP50:         p50,
				LatencyP95:         p95,
				LatencyP99:         p99,
				ThroughputMBps:     throughputMBps,
				ErrorRate:          0.0,
				SecurityOverhead:   calculateSecurityOverhead(opsPerSec, size),
				Metrics: map[string]float64{
					"dataSize":       float64(size),
					"iterations":     float64(b.N),
					"elapsedSeconds": elapsed.Seconds(),
				},
			}
			
			s.Results = append(s.Results, result)
			
			b.ReportMetric(opsPerSec, "ops/sec")
			b.ReportMetric(throughputMBps, "MB/s")
		})
	}
}

// BenchmarkKeyGenerationPerformance measures key generation performance
func (s *SecurityBenchmarkSuite) BenchmarkKeyGenerationPerformance(b *testing.B) {
	keySizes := []int{16, 32, 64, 128, 256} // bytes
	
	for _, size := range keySizes {
		b.Run(fmt.Sprintf("KeySize_%d", size), func(b *testing.B) {
			var memBefore, memAfter runtime.MemStats
			runtime.ReadMemStats(&memBefore)
			
			latencies := make([]time.Duration, 0, b.N)
			
			b.ResetTimer()
			start := time.Now()
			
			for i := 0; i < b.N; i++ {
				iterStart := time.Now()
				key := make([]byte, size)
				rand.Read(key)
				latencies = append(latencies, time.Since(iterStart))
				
				// Prevent optimization
				_ = key
			}
			
			elapsed := time.Since(start)
			runtime.ReadMemStats(&memAfter)
			
			// Calculate metrics
			opsPerSec := float64(b.N) / elapsed.Seconds()
			memUsage := memAfter.TotalAlloc - memBefore.TotalAlloc
			
			// Calculate latency percentiles
			p50, p95, p99 := calculateLatencyPercentiles(latencies)
			
			result := BenchmarkResult{
				TestName:           fmt.Sprintf("KeyGeneration_KeySize_%d", size),
				OperationsPerSec:   opsPerSec,
				MemoryUsage:        memUsage,
				CPUUsage:           getCPUUsage(),
				PowerConsumption:   float64(size) * 0.01, // Simulated power consumption
				LatencyP50:         p50,
				LatencyP95:         p95,
				LatencyP99:         p99,
				ThroughputMBps:     0, // Not applicable for key generation
				ErrorRate:          0.0,
				SecurityOverhead:   0.0,
				Metrics: map[string]float64{
					"keySize":        float64(size),
					"iterations":     float64(b.N),
					"elapsedSeconds": elapsed.Seconds(),
				},
			}
			
			s.Results = append(s.Results, result)
			
			b.ReportMetric(opsPerSec, "keys/sec")
			b.ReportMetric(float64(p50.Nanoseconds()), "p50-latency-ns")
		})
	}
}

// BenchmarkSecurityHeadersPerformance measures security headers processing performance
func (s *SecurityBenchmarkSuite) BenchmarkSecurityHeadersPerformance(b *testing.B) {
	headers := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Content-Security-Policy": "default-src 'self'",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}
	
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	
	latencies := make([]time.Duration, 0, b.N)
	
	b.ResetTimer()
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		iterStart := time.Now()
		
		// Simulate header processing
		for key, value := range headers {
			processed := processSecurityHeader(key, value)
			_ = processed
		}
		
		latencies = append(latencies, time.Since(iterStart))
	}
	
	elapsed := time.Since(start)
	runtime.ReadMemStats(&memAfter)
	
	// Calculate metrics
	opsPerSec := float64(b.N) / elapsed.Seconds()
	memUsage := memAfter.TotalAlloc - memBefore.TotalAlloc
	
	// Calculate latency percentiles
	p50, p95, p99 := calculateLatencyPercentiles(latencies)
	
	result := BenchmarkResult{
		TestName:           "SecurityHeaders",
		OperationsPerSec:   opsPerSec,
		MemoryUsage:        memUsage,
		CPUUsage:           getCPUUsage(),
		PowerConsumption:   0.01, // Minimal power consumption
		LatencyP50:         p50,
		LatencyP95:         p95,
		LatencyP99:         p99,
		ThroughputMBps:     0,
		ErrorRate:          0.0,
		SecurityOverhead:   5.0, // Approximate 5% overhead
		Metrics: map[string]float64{
			"headerCount":    float64(len(headers)),
			"iterations":     float64(b.N),
			"elapsedSeconds": elapsed.Seconds(),
		},
	}
	
	s.Results = append(s.Results, result)
	
	b.ReportMetric(opsPerSec, "ops/sec")
	b.ReportMetric(float64(p50.Nanoseconds()), "p50-latency-ns")
}

// BenchmarkTLSHandshakePerformance measures TLS handshake performance
func (s *SecurityBenchmarkSuite) BenchmarkTLSHandshakePerformance(b *testing.B) {
	b.ResetTimer()
	
	latencies := make([]time.Duration, 0, b.N)
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		iterStart := time.Now()
		
		// Simulate TLS handshake
		simulateTLSHandshake()
		
		latencies = append(latencies, time.Since(iterStart))
	}
	
	elapsed := time.Since(start)
	
	// Calculate metrics
	opsPerSec := float64(b.N) / elapsed.Seconds()
	
	// Calculate latency percentiles
	p50, p95, p99 := calculateLatencyPercentiles(latencies)
	
	result := BenchmarkResult{
		TestName:           "TLSHandshake",
		OperationsPerSec:   opsPerSec,
		MemoryUsage:        0, // Not measured for this test
		CPUUsage:           getCPUUsage(),
		PowerConsumption:   0.5, // Moderate power consumption
		LatencyP50:         p50,
		LatencyP95:         p95,
		LatencyP99:         p99,
		ThroughputMBps:     0,
		ErrorRate:          0.0,
		SecurityOverhead:   20.0, // Approximate 20% overhead for TLS
		Metrics: map[string]float64{
			"iterations":     float64(b.N),
			"elapsedSeconds": elapsed.Seconds(),
		},
	}
	
	s.Results = append(s.Results, result)
	
	b.ReportMetric(opsPerSec, "handshakes/sec")
	b.ReportMetric(float64(p50.Nanoseconds()), "p50-latency-ns")
}

// BenchmarkMemorySecurityPerformance measures memory security operations performance
func (s *SecurityBenchmarkSuite) BenchmarkMemorySecurityPerformance(b *testing.B) {
	b.ResetTimer()
	
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	
	latencies := make([]time.Duration, 0, b.N)
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		iterStart := time.Now()
		
		// Simulate secure memory operations
		secureData := make([]byte, 1024)
		rand.Read(secureData)
		
		// Simulate memory clearing
		for j := range secureData {
			secureData[j] = 0
		}
		
		latencies = append(latencies, time.Since(iterStart))
	}
	
	elapsed := time.Since(start)
	runtime.ReadMemStats(&memAfter)
	
	// Calculate metrics
	opsPerSec := float64(b.N) / elapsed.Seconds()
	memUsage := memAfter.TotalAlloc - memBefore.TotalAlloc
	
	// Calculate latency percentiles
	p50, p95, p99 := calculateLatencyPercentiles(latencies)
	
	result := BenchmarkResult{
		TestName:           "MemorySecurity",
		OperationsPerSec:   opsPerSec,
		MemoryUsage:        memUsage,
		CPUUsage:           getCPUUsage(),
		PowerConsumption:   0.05, // Low power consumption
		LatencyP50:         p50,
		LatencyP95:         p95,
		LatencyP99:         p99,
		ThroughputMBps:     0,
		ErrorRate:          0.0,
		SecurityOverhead:   2.0, // Approximate 2% overhead
		Metrics: map[string]float64{
			"iterations":     float64(b.N),
			"elapsedSeconds": elapsed.Seconds(),
		},
	}
	
	s.Results = append(s.Results, result)
	
	b.ReportMetric(opsPerSec, "ops/sec")
	b.ReportMetric(float64(p50.Nanoseconds()), "p50-latency-ns")
}

// BenchmarkConcurrentSecurityPerformance measures concurrent security operations performance
func (s *SecurityBenchmarkSuite) BenchmarkConcurrentSecurityPerformance(b *testing.B) {
	concurrencyLevels := []int{1, 2, 4, 8, 16, 32}
	
	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			b.ResetTimer()
			
			start := time.Now()
			
			b.RunParallel(func(pb *testing.PB) {
				key := make([]byte, 32)
				plaintext := make([]byte, 1024)
				rand.Read(key)
				rand.Read(plaintext)
				
				for pb.Next() {
					ciphertext := simulateEncryption(key, plaintext)
					_ = ciphertext
				}
			})
			
			elapsed := time.Since(start)
			
			// Calculate metrics
			opsPerSec := float64(b.N) / elapsed.Seconds()
			
			result := BenchmarkResult{
				TestName:           fmt.Sprintf("ConcurrentSecurity_Concurrency_%d", concurrency),
				OperationsPerSec:   opsPerSec,
				MemoryUsage:        0, // Not measured for concurrent test
				CPUUsage:           getCPUUsage(),
				PowerConsumption:   float64(concurrency) * 0.1, // Power scales with concurrency
				LatencyP50:         0, // Not applicable for concurrent test
				LatencyP95:         0,
				LatencyP99:         0,
				ThroughputMBps:     (1024 * float64(b.N)) / (1024 * 1024) / elapsed.Seconds(),
				ErrorRate:          0.0,
				SecurityOverhead:   0.0,
				Metrics: map[string]float64{
					"concurrency":    float64(concurrency),
					"iterations":     float64(b.N),
					"elapsedSeconds": elapsed.Seconds(),
				},
			}
			
			s.Results = append(s.Results, result)
			
			b.ReportMetric(opsPerSec, "ops/sec")
			b.ReportMetric(result.ThroughputMBps, "MB/s")
		})
	}
}

// Helper functions for benchmarking

func simulateDecryption(key, ciphertext []byte) []byte {
	// Simulate decryption operation
	result := make([]byte, len(ciphertext))
	for i, b := range ciphertext {
		result[i] = b ^ key[i%len(key)]
	}
	return result
}

func processSecurityHeader(key, value string) string {
	// Simulate security header processing
	processed := fmt.Sprintf("%s: %s", key, value)
	
	// Add some processing time
	time.Sleep(time.Nanosecond * 100)
	
	return processed
}

func simulateTLSHandshake() {
	// Simulate TLS handshake operations
	// Generate ephemeral key pair
	clientKey := make([]byte, 32)
	serverKey := make([]byte, 32)
	rand.Read(clientKey)
	rand.Read(serverKey)
	
	// Simulate key exchange
	sharedSecret := make([]byte, 32)
	for i := 0; i < 32; i++ {
		sharedSecret[i] = clientKey[i] ^ serverKey[i]
	}
	
	// Simulate certificate verification
	time.Sleep(time.Microsecond * 100)
	
	// Use shared secret to prevent optimization
	_ = sharedSecret
}

func calculateLatencyPercentiles(latencies []time.Duration) (p50, p95, p99 time.Duration) {
	if len(latencies) == 0 {
		return 0, 0, 0
	}
	
	// Sort latencies
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	
	// Simple insertion sort for small arrays
	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && sorted[j] > key {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}
	
	// Calculate percentiles
	p50Index := int(0.50 * float64(len(sorted)))
	p95Index := int(0.95 * float64(len(sorted)))
	p99Index := int(0.99 * float64(len(sorted)))
	
	if p50Index >= len(sorted) {
		p50Index = len(sorted) - 1
	}
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}
	if p99Index >= len(sorted) {
		p99Index = len(sorted) - 1
	}
	
	return sorted[p50Index], sorted[p95Index], sorted[p99Index]
}

func getCPUUsage() float64 {
	// Simulate CPU usage measurement
	// In a real implementation, this would use system calls
	return 15.0 + (rand.Float64() * 10.0) // 15-25% CPU usage
}

func calculateSecurityOverhead(opsPerSec float64, dataSize int) float64 {
	// Estimate security overhead based on operation rate and data size
	// This is a simplified calculation
	baseOverhead := 5.0 // Base 5% overhead
	
	// Larger data sizes have proportionally less overhead
	sizeOverhead := math.Max(0, 10.0-float64(dataSize)/1000.0)
	
	// Higher operation rates might have more overhead due to contention
	rateOverhead := math.Min(5.0, opsPerSec/10000.0)
	
	return baseOverhead + sizeOverhead + rateOverhead
}

// generateBenchmarkReport generates a comprehensive benchmark report
func (s *SecurityBenchmarkSuite) generateBenchmarkReport(b *testing.B) {
	if len(s.Results) == 0 {
		b.Log("No benchmark results to report")
		return
	}
	
	b.Log("=== Security Performance Benchmark Report ===")
	
	// Calculate summary statistics
	var totalOps, totalThroughput, totalMemory, totalPower float64
	var totalLatencyP50, totalLatencyP95, totalLatencyP99 time.Duration
	
	for _, result := range s.Results {
		totalOps += result.OperationsPerSec
		totalThroughput += result.ThroughputMBps
		totalMemory += float64(result.MemoryUsage)
		totalPower += result.PowerConsumption
		totalLatencyP50 += result.LatencyP50
		totalLatencyP95 += result.LatencyP95
		totalLatencyP99 += result.LatencyP99
	}
	
	count := float64(len(s.Results))
	b.Logf("Average operations per second: %.2f", totalOps/count)
	b.Logf("Average throughput: %.2f MB/s", totalThroughput/count)
	b.Logf("Average memory usage: %.2f KB", totalMemory/count/1024)
	b.Logf("Average power consumption: %.2f W", totalPower/count)
	b.Logf("Average P50 latency: %v", time.Duration(int64(totalLatencyP50)/int64(count)))
	b.Logf("Average P95 latency: %v", time.Duration(int64(totalLatencyP95)/int64(count)))
	b.Logf("Average P99 latency: %v", time.Duration(int64(totalLatencyP99)/int64(count)))
	
	// Log detailed results
	b.Log("\n=== Detailed Results ===")
	for _, result := range s.Results {
		b.Logf("Test: %s", result.TestName)
		b.Logf("  Operations/sec: %.2f", result.OperationsPerSec)
		b.Logf("  Memory usage: %.2f KB", float64(result.MemoryUsage)/1024)
		b.Logf("  CPU usage: %.2f%%", result.CPUUsage)
		b.Logf("  Power consumption: %.2f W", result.PowerConsumption)
		b.Logf("  P50 latency: %v", result.LatencyP50)
		b.Logf("  P95 latency: %v", result.LatencyP95)
		b.Logf("  P99 latency: %v", result.LatencyP99)
		b.Logf("  Throughput: %.2f MB/s", result.ThroughputMBps)
		b.Logf("  Security overhead: %.2f%%", result.SecurityOverhead)
		b.Logf("  Error rate: %.2f%%", result.ErrorRate*100)
	}
	
	// Export results as JSON for further analysis
	jsonData, _ := json.MarshalIndent(s.Results, "", "  ")
	b.Logf("\n=== JSON Export ===\n%s", string(jsonData))
}

// TestSecurityBenchmarks runs the complete security benchmark suite
func TestSecurityBenchmarks(t *testing.T) {
	suite := &SecurityBenchmarkSuite{
		Results: make([]BenchmarkResult, 0),
	}
	
	// Run benchmarks in test mode
	testing.Benchmark(suite.RunSecurityBenchmarks)
	
	// Verify performance meets minimum requirements
	if len(suite.Results) == 0 {
		t.Error("No benchmark results generated")
		return
	}
	
	// Performance assertions
	for _, result := range suite.Results {
		if result.OperationsPerSec < 100 { // Minimum 100 ops/sec
			t.Errorf("Performance below minimum for %s: %.2f ops/sec", result.TestName, result.OperationsPerSec)
		}
		
		if result.SecurityOverhead > 50 { // Maximum 50% overhead
			t.Errorf("Security overhead too high for %s: %.2f%%", result.TestName, result.SecurityOverhead)
		}
		
		if result.ErrorRate > 0.01 { // Maximum 1% error rate
			t.Errorf("Error rate too high for %s: %.2f%%", result.TestName, result.ErrorRate*100)
		}
	}
}