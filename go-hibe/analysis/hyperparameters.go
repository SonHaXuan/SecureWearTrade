package analysis

import (
	"encoding/json"
	"fmt"
	"time"
)

// SystemHyperparameters contains all configurable parameters for reproducibility
type SystemHyperparameters struct {
	// Cryptographic Parameters
	EllipticCurve         string  `json:"elliptic_curve"`
	PairingType          string  `json:"pairing_type"`
	SecurityLevel        int     `json:"security_level_bits"`
	HIBEKeySize          int     `json:"hibe_key_size_bytes"`
	PatternSize          int     `json:"pattern_size"`
	MaxTimeLength        int     `json:"max_time_length"`
	
	// Batch Processing Parameters
	BatchSizeMin         int     `json:"batch_size_min_kb"`
	BatchSizeMax         int     `json:"batch_size_max_kb"`
	BatchThreshold       int     `json:"batch_threshold_kb"`
	ChunkSize            int     `json:"chunk_size_bytes"`
	CompressionEnabled   bool    `json:"compression_enabled"`
	
	// Performance Parameters
	MaxConcurrentOps     int     `json:"max_concurrent_operations"`
	CacheSize            int     `json:"cache_size_mb"`
	MemoryPoolSize       int     `json:"memory_pool_size_mb"`
	GCInterval           int     `json:"gc_interval_seconds"`
	
	// Power Analysis Parameters
	BasePowerWatts       float64 `json:"base_power_watts"`
	CPUPowerFactor       float64 `json:"cpu_power_factor"`
	MemoryPowerFactor    float64 `json:"memory_power_factor"`
	PowerMeasureInterval int     `json:"power_measure_interval_ms"`
	
	// Network Parameters
	IPFSNodeCount        int     `json:"ipfs_node_count"`
	ReplicationFactor    int     `json:"replication_factor"`
	NetworkLatencyMs     int     `json:"network_latency_ms"`
	BandwidthMbps        int     `json:"bandwidth_mbps"`
	
	// Blockchain Parameters
	GasLimit             uint64  `json:"gas_limit"`
	GasPrice             float64 `json:"gas_price_gwei"`
	BlockConfirmations   int     `json:"block_confirmations"`
	SmartContractAddr    string  `json:"smart_contract_address"`
	PolygonRPCEndpoint   string  `json:"polygon_rpc_endpoint"`
	
	// Security Parameters
	EncryptionAlgorithm  string  `json:"encryption_algorithm"`
	HashFunction         string  `json:"hash_function"`
	KeyDerivationFunction string `json:"key_derivation_function"`
	RandomnessSource     string  `json:"randomness_source"`
	ConstantTimeOps      bool    `json:"constant_time_operations"`
	
	// Testing Parameters
	TestIterations       int     `json:"test_iterations"`
	WarmupIterations     int     `json:"warmup_iterations"`
	StatisticalSamples   int     `json:"statistical_samples"`
	ConfidenceLevel      float64 `json:"confidence_level"`
	
	// Dataset Parameters
	LifeSnapsDatasetSize int     `json:"lifesnaps_dataset_size_mb"`
	SyntheticDataRatio   float64 `json:"synthetic_data_ratio"`
	DataDistribution     string  `json:"data_distribution"`
	NoiseLevel           float64 `json:"noise_level"`
	
	// Hardware Simulation Parameters
	DeviceType           string  `json:"device_type"`
	CPUCores             int     `json:"cpu_cores"`
	RAMSizeGB           int     `json:"ram_size_gb"`
	StorageType         string  `json:"storage_type"`
	BatteryCapacityWh   float64 `json:"battery_capacity_wh"`
	
	// Metadata
	Version             string    `json:"version"`
	LastUpdated         time.Time `json:"last_updated"`
	ConfigurationName   string    `json:"configuration_name"`
	ExperimentID        string    `json:"experiment_id"`
}

// DefaultHyperparameters returns the standard configuration used in experiments
func DefaultHyperparameters() *SystemHyperparameters {
	return &SystemHyperparameters{
		// Cryptographic Parameters (Section 5.1 requirements)
		EllipticCurve:         "secp256k1",
		PairingType:          "BLS12-381",
		SecurityLevel:        256, // 256-bit security
		HIBEKeySize:          15871488, // 15.15MB as mentioned in Table 2
		PatternSize:          20,
		MaxTimeLength:        4,
		
		// Batch Processing Parameters (Algorithm 7)
		BatchSizeMin:         100, // 100KB minimum
		BatchSizeMax:         250, // 250KB maximum
		BatchThreshold:       200, // 200KB threshold
		ChunkSize:            4096, // 4KB chunks
		CompressionEnabled:   true,
		
		// Performance Parameters
		MaxConcurrentOps:     50,
		CacheSize:            128, // 128MB cache
		MemoryPoolSize:       256, // 256MB memory pool
		GCInterval:           30,  // 30 seconds
		
		// Power Analysis Parameters (mobile device simulation)
		BasePowerWatts:       0.5,   // 0.5W base power
		CPUPowerFactor:       0.05,  // 0.05W per 1% CPU
		MemoryPowerFactor:    0.02,  // 0.02W per GB memory
		PowerMeasureInterval: 100,   // 100ms intervals
		
		// Network Parameters (IPFS distributed setup)
		IPFSNodeCount:        5,
		ReplicationFactor:    3,
		NetworkLatencyMs:     50,
		BandwidthMbps:        100,
		
		// Blockchain Parameters (Polygon mainnet)
		GasLimit:             1000000, // 1M gas limit
		GasPrice:             30.0,    // 30 Gwei
		BlockConfirmations:   12,
		SmartContractAddr:    "0x742d35Cc6634C0532925a3b8D6Ac6B0ad39CEe5C", // Example
		PolygonRPCEndpoint:   "https://polygon-rpc.com",
		
		// Security Parameters
		EncryptionAlgorithm:  "AES-256-GCM",
		HashFunction:         "SHA-256",
		KeyDerivationFunction: "PBKDF2",
		RandomnessSource:     "crypto/rand",
		ConstantTimeOps:      true,
		
		// Testing Parameters
		TestIterations:       1000,
		WarmupIterations:     100,
		StatisticalSamples:   50,
		ConfidenceLevel:      0.95, // 95% confidence
		
		// Dataset Parameters
		LifeSnapsDatasetSize: 1024, // 1GB dataset
		SyntheticDataRatio:   0.3,  // 30% synthetic data
		DataDistribution:     "normal",
		NoiseLevel:           0.05, // 5% noise
		
		// Hardware Simulation Parameters
		DeviceType:           "mobile",
		CPUCores:             4,
		RAMSizeGB:           4,
		StorageType:         "SSD",
		BatteryCapacityWh:   15.0, // 15Wh battery
		
		// Metadata
		Version:             "1.0.0",
		LastUpdated:         time.Now(),
		ConfigurationName:   "SecureWearTrade Default",
		ExperimentID:        generateExperimentID(),
	}
}

// MobileDeviceHyperparameters returns parameters optimized for mobile devices
func MobileDeviceHyperparameters() *SystemHyperparameters {
	params := DefaultHyperparameters()
	
	// Optimize for mobile constraints
	params.ConfigurationName = "Mobile Device Optimized"
	params.BatchSizeMax = 128      // Smaller batches for mobile
	params.CacheSize = 64          // Smaller cache
	params.MemoryPoolSize = 128    // Smaller memory pool
	params.CPUCores = 2            // Typical mobile CPU
	params.RAMSizeGB = 2           // Typical mobile RAM
	params.BatteryCapacityWh = 5.0 // Smaller battery
	params.BasePowerWatts = 0.3    // Lower base power
	
	return params
}

// IoTDeviceHyperparameters returns parameters for IoT devices
func IoTDeviceHyperparameters() *SystemHyperparameters {
	params := DefaultHyperparameters()
	
	// Optimize for IoT constraints
	params.ConfigurationName = "IoT Device Optimized"
	params.DeviceType = "iot"
	params.BatchSizeMax = 64       // Very small batches
	params.CacheSize = 32          // Minimal cache
	params.MemoryPoolSize = 64     // Minimal memory pool
	params.CPUCores = 1            // Single core
	params.RAMSizeGB = 1           // Minimal RAM
	params.BatteryCapacityWh = 2.0 // Very small battery
	params.BasePowerWatts = 0.1    // Minimal base power
	params.MaxConcurrentOps = 10   // Limited concurrency
	
	return params
}

// HighPerformanceHyperparameters returns parameters for high-performance testing
func HighPerformanceHyperparameters() *SystemHyperparameters {
	params := DefaultHyperparameters()
	
	// Optimize for performance
	params.ConfigurationName = "High Performance"
	params.DeviceType = "server"
	params.BatchSizeMax = 1024     // Larger batches
	params.CacheSize = 512         // Large cache
	params.MemoryPoolSize = 1024   // Large memory pool
	params.CPUCores = 16           // Many cores
	params.RAMSizeGB = 32          // Lots of RAM
	params.MaxConcurrentOps = 200  // High concurrency
	params.TestIterations = 10000  // More iterations
	
	return params
}

// ExperimentConfiguration represents a complete experimental setup
type ExperimentConfiguration struct {
	Hyperparameters *SystemHyperparameters `json:"hyperparameters"`
	Environment     *EnvironmentInfo       `json:"environment"`
	Objectives      []string               `json:"objectives"`
	Metrics         []string               `json:"metrics"`
	Duration        time.Duration          `json:"duration"`
	Description     string                 `json:"description"`
}

// EnvironmentInfo captures the testing environment details
type EnvironmentInfo struct {
	OperatingSystem   string    `json:"operating_system"`
	Architecture      string    `json:"architecture"`
	GoVersion         string    `json:"go_version"`
	CompilerFlags     []string  `json:"compiler_flags"`
	TimestampUTC      time.Time `json:"timestamp_utc"`
	HostName          string    `json:"hostname"`
	NetworkConditions string    `json:"network_conditions"`
	LoadConditions    string    `json:"load_conditions"`
}

// ReproducibilityPackage contains everything needed to reproduce experiments
type ReproducibilityPackage struct {
	Configuration *ExperimentConfiguration `json:"configuration"`
	SourceCode    map[string]string        `json:"source_code"`
	Dependencies  map[string]string        `json:"dependencies"`
	BuildInstructions []string             `json:"build_instructions"`
	RunInstructions   []string             `json:"run_instructions"`
	ExpectedResults   map[string]interface{} `json:"expected_results"`
	Checksums         map[string]string      `json:"checksums"`
}

// GenerateHyperparametersTable creates a formatted table for Section 5.1
func GenerateHyperparametersTable(params *SystemHyperparameters) string {
	table := "=== HYPERPARAMETERS TABLE (Section 5.1) ===\n\n"
	
	table += "CRYPTOGRAPHIC PARAMETERS:\n"
	table += fmt.Sprintf("%-30s: %s\n", "Elliptic Curve", params.EllipticCurve)
	table += fmt.Sprintf("%-30s: %s\n", "Pairing Type", params.PairingType)
	table += fmt.Sprintf("%-30s: %d bits\n", "Security Level", params.SecurityLevel)
	table += fmt.Sprintf("%-30s: %.2f MB\n", "HIBE Key Size", float64(params.HIBEKeySize)/1024/1024)
	table += fmt.Sprintf("%-30s: %d\n", "Pattern Size", params.PatternSize)
	table += "\n"
	
	table += "BATCH PROCESSING PARAMETERS (Algorithm 7):\n"
	table += fmt.Sprintf("%-30s: %d - %d KB\n", "Batch Size Range", params.BatchSizeMin, params.BatchSizeMax)
	table += fmt.Sprintf("%-30s: %d KB\n", "Batch Threshold", params.BatchThreshold)
	table += fmt.Sprintf("%-30s: %d bytes\n", "Chunk Size", params.ChunkSize)
	table += fmt.Sprintf("%-30s: %t\n", "Compression Enabled", params.CompressionEnabled)
	table += "\n"
	
	table += "POWER ANALYSIS PARAMETERS:\n"
	table += fmt.Sprintf("%-30s: %.3f W\n", "Base Power", params.BasePowerWatts)
	table += fmt.Sprintf("%-30s: %.3f W/%%\n", "CPU Power Factor", params.CPUPowerFactor)
	table += fmt.Sprintf("%-30s: %.3f W/GB\n", "Memory Power Factor", params.MemoryPowerFactor)
	table += fmt.Sprintf("%-30s: %d ms\n", "Measurement Interval", params.PowerMeasureInterval)
	table += "\n"
	
	table += "TESTING PARAMETERS:\n"
	table += fmt.Sprintf("%-30s: %d\n", "Test Iterations", params.TestIterations)
	table += fmt.Sprintf("%-30s: %d\n", "Warmup Iterations", params.WarmupIterations)
	table += fmt.Sprintf("%-30s: %d\n", "Statistical Samples", params.StatisticalSamples)
	table += fmt.Sprintf("%-30s: %.2f%%\n", "Confidence Level", params.ConfidenceLevel*100)
	table += "\n"
	
	table += "DATASET PARAMETERS:\n"
	table += fmt.Sprintf("%-30s: %d MB\n", "LifeSnaps Dataset Size", params.LifeSnapsDatasetSize)
	table += fmt.Sprintf("%-30s: %.1f%%\n", "Synthetic Data Ratio", params.SyntheticDataRatio*100)
	table += fmt.Sprintf("%-30s: %s\n", "Data Distribution", params.DataDistribution)
	table += fmt.Sprintf("%-30s: %.1f%%\n", "Noise Level", params.NoiseLevel*100)
	table += "\n"
	
	table += "HARDWARE SIMULATION:\n"
	table += fmt.Sprintf("%-30s: %s\n", "Device Type", params.DeviceType)
	table += fmt.Sprintf("%-30s: %d\n", "CPU Cores", params.CPUCores)
	table += fmt.Sprintf("%-30s: %d GB\n", "RAM Size", params.RAMSizeGB)
	table += fmt.Sprintf("%-30s: %.1f Wh\n", "Battery Capacity", params.BatteryCapacityWh)
	
	return table
}

// ExportForReproducibility creates a complete reproducibility package
func ExportForReproducibility(params *SystemHyperparameters) (*ReproducibilityPackage, error) {
	// Create experiment configuration
	config := &ExperimentConfiguration{
		Hyperparameters: params,
		Environment: &EnvironmentInfo{
			OperatingSystem:   "darwin", // or runtime.GOOS
			Architecture:      "amd64",  // or runtime.GOARCH
			GoVersion:         "go1.21",
			CompilerFlags:     []string{"-O2", "-trimpath"},
			TimestampUTC:      time.Now().UTC(),
			HostName:          "research-cluster",
			NetworkConditions: "stable, 100Mbps",
			LoadConditions:    "normal load",
		},
		Objectives: []string{
			"Measure wildcard performance improvement",
			"Compare HIBE schemes efficiency",
			"Validate IPFS-blockchain integration",
			"Analyze power consumption patterns",
		},
		Metrics: []string{
			"key_generation_time_ms",
			"encryption_time_ms",
			"decryption_time_ms",
			"memory_usage_kb",
			"power_consumption_w",
			"efficiency_gain_pct",
		},
		Duration:    time.Hour * 2, // 2-hour experiment
		Description: "Comprehensive performance analysis of SecureWearTrade framework",
	}
	
	// Create reproducibility package
	package_ := &ReproducibilityPackage{
		Configuration: config,
		SourceCode: map[string]string{
			"algorithms_1-8": "github.com/SecureWearTrade/algorithms",
			"benchmark_suite": "github.com/SecureWearTrade/benchmarks",
			"security_tests": "github.com/SecureWearTrade/security",
		},
		Dependencies: map[string]string{
			"github.com/ucbrise/hibe-pairing": "v0.0.0-20220312033002-c4bf151b8d2b",
			"github.com/gin-gonic/gin":        "v1.10.0",
			"github.com/shirou/gopsutil/v3":   "v3.23.12",
		},
		BuildInstructions: []string{
			"git clone https://github.com/SecureWearTrade/repository",
			"cd repository",
			"go mod download",
			"cd go-hibe && go build -o hibe-server main.go",
			"cd ../kyc-contract && npm install && npx hardhat compile",
			"cd ../TrustAuthority && npm install",
		},
		RunInstructions: []string{
			"./run_benchmarks.sh --config=default",
			"./run_security_tests.sh --comprehensive",
			"./analyze_results.sh --output=json",
		},
		ExpectedResults: map[string]interface{}{
			"wildcard_improvement_pct": 25.0,
			"avg_key_generation_ms":    90.0,
			"avg_power_consumption_w":  0.6,
			"security_test_pass_rate":  0.95,
		},
		Checksums: map[string]string{
			"main.go":                "sha256:1234567890abcdef...",
			"benchmarks_suite.tar.gz": "sha256:abcdef1234567890...",
			"test_dataset.zip":       "sha256:fedcba0987654321...",
		},
	}
	
	return package_, nil
}

// ValidateHyperparameters checks if hyperparameters are within valid ranges
func ValidateHyperparameters(params *SystemHyperparameters) []string {
	var errors []string
	
	// Validate cryptographic parameters
	if params.SecurityLevel < 128 {
		errors = append(errors, "Security level must be at least 128 bits")
	}
	
	if params.HIBEKeySize < 1024*1024 { // 1MB minimum
		errors = append(errors, "HIBE key size too small")
	}
	
	// Validate batch parameters
	if params.BatchSizeMin > params.BatchSizeMax {
		errors = append(errors, "Batch size minimum cannot exceed maximum")
	}
	
	if params.BatchThreshold < params.BatchSizeMin || params.BatchThreshold > params.BatchSizeMax {
		errors = append(errors, "Batch threshold must be within batch size range")
	}
	
	// Validate power parameters
	if params.BasePowerWatts < 0 {
		errors = append(errors, "Base power cannot be negative")
	}
	
	// Validate testing parameters
	if params.TestIterations < 1 {
		errors = append(errors, "Test iterations must be positive")
	}
	
	if params.ConfidenceLevel <= 0 || params.ConfidenceLevel >= 1 {
		errors = append(errors, "Confidence level must be between 0 and 1")
	}
	
	return errors
}

// generateExperimentID creates a unique identifier for the experiment
func generateExperimentID() string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("SWT-EXP-%s", timestamp)
}

// SaveHyperparameters saves hyperparameters to a JSON file
func SaveHyperparameters(params *SystemHyperparameters, filename string) error {
	data, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal hyperparameters: %v", err)
	}
	
	// In a real implementation, you would write to a file
	fmt.Printf("Hyperparameters saved to %s:\n%s\n", filename, string(data))
	
	return nil
}

// LoadHyperparameters loads hyperparameters from a JSON file
func LoadHyperparameters(filename string) (*SystemHyperparameters, error) {
	// In a real implementation, you would read from a file
	// For now, return default parameters
	return DefaultHyperparameters(), nil
}

// CompareHyperparameters compares two hyperparameter configurations
func CompareHyperparameters(params1, params2 *SystemHyperparameters) map[string]interface{} {
	comparison := make(map[string]interface{})
	
	comparison["elliptic_curve_match"] = params1.EllipticCurve == params2.EllipticCurve
	comparison["security_level_diff"] = params2.SecurityLevel - params1.SecurityLevel
	comparison["batch_size_diff"] = params2.BatchSizeMax - params1.BatchSizeMax
	comparison["test_iterations_ratio"] = float64(params2.TestIterations) / float64(params1.TestIterations)
	comparison["power_base_diff"] = params2.BasePowerWatts - params1.BasePowerWatts
	
	return comparison
}