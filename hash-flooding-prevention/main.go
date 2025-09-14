package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"./prevention"
	"../ipfs-blockchain-binding"
)

func main() {
	fmt.Println("=== IPFS-Blockchain Cryptographic Binding & Hash Flooding Prevention System ===")
	
	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "binding-demo":
			runCryptographicBindingDemo()
		case "flooding-demo":
			runHashFloodingPreventionDemo()
		case "integrated-demo":
			runIntegratedSystemDemo()
		case "service":
			runMultiTierService()
		default:
			printUsage()
		}
	} else {
		runIntegratedSystemDemo()
	}
}

// runCryptographicBindingDemo demonstrates the IPFS-blockchain binding system
func runCryptographicBindingDemo() {
	fmt.Println("\n🔗 Running Cryptographic Binding Demonstration...")
	
	// Initialize components
	ethConnector, err := binding.NewEthereumConnector(
		"http://localhost:8545",                           // Ethereum node URL
		"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // Private key
		"0x5FbDB2315678afecb367f032d93F642f64180aa3",         // Contract address
		big.NewInt(31337),                                 // Chain ID (localhost)
	)
	if err != nil {
		log.Printf("Warning: Could not connect to Ethereum node: %v", err)
		log.Println("Continuing with demonstration using mock data...")
	}
	
	ipfsConnector := binding.NewIPFSConnector("http://localhost:5001")
	cryptoBinding := binding.NewCryptographicBinding(ethConnector, ipfsConnector)
	
	// Run real-world example
	err = cryptoBinding.RealWorldExample()
	if err != nil {
		log.Printf("Demo completed with simulated data due to: %v", err)
	}
}

// runHashFloodingPreventionDemo demonstrates hash flooding attack prevention
func runHashFloodingPreventionDemo() {
	fmt.Println("\n🛡️  Running Hash Flooding Prevention Demonstration...")
	
	// Initialize prevention system
	config := &prevention.PreventionConfig{
		EnableTieredLimiting:        true,
		EnableFalsePositiveTracking: true,
		MaxClientsTracked:          10000,
		CleanupInterval:            5 * time.Minute,
		MetricsRetentionPeriod:     24 * time.Hour,
	}
	
	floodPrevention := prevention.NewHashFloodingPrevention(config)
	
	// Demonstrate different attack scenarios
	demonstrateAttackScenarios(floodPrevention)
	
	// Show false positive analysis
	floodPrevention.PrintDetailedReport()
}

// runIntegratedSystemDemo demonstrates the complete integrated system
func runIntegratedSystemDemo() {
	fmt.Println("\n🚀 Running Integrated System Demonstration...")
	
	// Initialize all components
	config := &prevention.PreventionConfig{
		EnableTieredLimiting:        true,
		EnableFalsePositiveTracking: true,
		MaxClientsTracked:          10000,
		CleanupInterval:            5 * time.Minute,
		MetricsRetentionPeriod:     24 * time.Hour,
	}
	
	floodPrevention := prevention.NewHashFloodingPrevention(config)
	
	// Demonstrate cryptographic binding
	fmt.Println("\n=== Part 1: Cryptographic Binding ===")
	runCryptographicBindingDemo()
	
	// Demonstrate hash flooding prevention
	fmt.Println("\n=== Part 2: Hash Flooding Prevention ===")
	demonstrateAttackScenarios(floodPrevention)
	
	// Show integrated workflow
	fmt.Println("\n=== Part 3: Integrated Workflow ===")
	demonstrateIntegratedWorkflow(floodPrevention)
	
	// Final system report
	fmt.Println("\n=== System Performance Report ===")
	floodPrevention.PrintDetailedReport()
}

// runMultiTierService runs the HTTP API service
func runMultiTierService() {
	fmt.Println("\n🌐 Starting Multi-Tier Rate Limiting HTTP Service...")
	
	// Initialize prevention system
	config := &prevention.PreventionConfig{
		EnableTieredLimiting:        true,
		EnableFalsePositiveTracking: true,
		MaxClientsTracked:          10000,
		CleanupInterval:            5 * time.Minute,
		MetricsRetentionPeriod:     24 * time.Hour,
	}
	
	floodPrevention := prevention.NewHashFloodingPrevention(config)
	
	// Initialize service
	serviceConfig := &prevention.ServiceConfig{
		Port:                    8080,
		EnableMetrics:          true,
		EnableAdvancedAnalytics: true,
		RequestTimeout:         30 * time.Second,
		MaxConcurrentRequests:  1000,
	}
	
	service := prevention.NewMultiTierRateLimitingService(floodPrevention, serviceConfig)
	
	// Start service in goroutine
	go func() {
		if err := service.Start(); err != nil {
			log.Printf("Service error: %v", err)
		}
	}()
	
	// Demonstrate API usage
	time.Sleep(2 * time.Second)
	service.DemonstrateRealWorldScenario()
	
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := service.Stop(ctx); err != nil {
		log.Printf("Service shutdown error: %v", err)
	}
	
	fmt.Println("Service stopped successfully")
}

// demonstrateAttackScenarios shows various attack scenarios and prevention
func demonstrateAttackScenarios(floodPrevention *prevention.HashFloodingPrevention) {
	fmt.Println("\n--- Attack Scenario Simulations ---")
	
	// Scenario 1: Basic legitimate user
	fmt.Println("\nScenario 1: Legitimate Healthcare Professional")
	testClient("doctor_0x742d35Cc", big.NewInt(5000000000), 10, floodPrevention) // 5 gwei, 10 requests
	
	// Scenario 2: Research institution with higher gas fees
	fmt.Println("\nScenario 2: Research Institution (Premium Tier)")
	testClient("research_0x8b2c9f", big.NewInt(30000000000), 50, floodPrevention) // 30 gwei, 50 requests
	
	// Scenario 3: Potential attacker with low gas fees
	fmt.Println("\nScenario 3: Potential Attacker (Low Gas Fees)")
	testClient("attacker_0x123456", big.NewInt(1000000000), 150, floodPrevention) // 1 gwei, 150 requests
	
	// Scenario 4: Enterprise with high gas fees
	fmt.Println("\nScenario 4: Enterprise Hospital (Platinum Tier)")
	testClient("hospital_0x7c3e9a", big.NewInt(150000000000), 1000, floodPrevention) // 150 gwei, 1000 requests
	
	// Scenario 5: Burst attack simulation
	fmt.Println("\nScenario 5: Burst Attack Simulation")
	simulateBurstAttack("burst_attacker", big.NewInt(2000000000), floodPrevention)
}

// testClient simulates requests from a specific client
func testClient(clientID string, gasFeePaid *big.Int, requestCount int, floodPrevention *prevention.HashFloodingPrevention) {
	fmt.Printf("Testing client: %s with %s gwei for %d requests\n", 
		clientID, weiToGwei(gasFeePaid), requestCount)
	
	allowed := 0
	blocked := 0
	
	for i := 0; i < requestCount; i++ {
		result, err := floodPrevention.ValidateHashRequest(clientID, gasFeePaid, "test_request")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		
		if result.Allowed {
			allowed++
		} else {
			blocked++
			if i < 5 || i == requestCount-1 { // Show first 5 and last rejection
				fmt.Printf("  Request %d blocked: %s\n", i+1, result.RejectionReason)
			}
		}
		
		// Small delay to simulate real requests
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Printf("Results: %d allowed, %d blocked (%.1f%% blocked)\n", 
		allowed, blocked, float64(blocked)/float64(requestCount)*100)
}

// simulateBurstAttack simulates a burst attack pattern
func simulateBurstAttack(clientID string, gasFeePaid *big.Int, floodPrevention *prevention.HashFloodingPrevention) {
	fmt.Printf("Simulating burst attack from: %s\n", clientID)
	
	// Phase 1: Normal requests
	fmt.Println("  Phase 1: Normal request pattern...")
	for i := 0; i < 10; i++ {
		floodPrevention.ValidateHashRequest(clientID, gasFeePaid, "normal_request")
		time.Sleep(100 * time.Millisecond)
	}
	
	// Phase 2: Sudden burst
	fmt.Println("  Phase 2: Sudden burst attack...")
	burstAllowed := 0
	burstBlocked := 0
	
	for i := 0; i < 600; i++ { // Burst of 600 requests
		result, _ := floodPrevention.ValidateHashRequest(clientID, gasFeePaid, "burst_request")
		if result.Allowed {
			burstAllowed++
		} else {
			burstBlocked++
		}
		
		// No delay for burst simulation
		if i%100 == 0 {
			fmt.Printf("    Burst progress: %d requests, %d blocked\n", i+1, burstBlocked)
		}
	}
	
	fmt.Printf("  Burst results: %d allowed, %d blocked\n", burstAllowed, burstBlocked)
	
	// Phase 3: Post-burst cooldown testing
	fmt.Println("  Phase 3: Testing cooldown period...")
	for i := 0; i < 5; i++ {
		result, _ := floodPrevention.ValidateHashRequest(clientID, gasFeePaid, "post_burst_request")
		fmt.Printf("    Post-burst request %d: %t (cooldown: %v)\n", 
			i+1, result.Allowed, result.CooldownRemaining)
		time.Sleep(time.Second)
	}
}

// demonstrateIntegratedWorkflow shows the complete workflow
func demonstrateIntegratedWorkflow(floodPrevention *prevention.HashFloodingPrevention) {
	fmt.Println("Demonstrating integrated IPFS-blockchain binding with hash flooding prevention...")
	
	// Simulate healthcare data upload workflow
	workflows := []struct {
		clientID     string
		patientID    string
		gasFeePaid   *big.Int
		dataSize     string
		description  string
	}{
		{
			clientID:    "doctor_cardiology_0x742d35",
			patientID:   "patient_12345",
			gasFeePaid:  big.NewInt(15000000000), // 15 gwei
			dataSize:    "2.5MB ECG data",
			description: "Cardiologist uploading patient ECG data",
		},
		{
			clientID:    "research_genomics_0x8b2c9f",
			patientID:   "research_cohort_001",
			gasFeePaid:  big.NewInt(45000000000), // 45 gwei
			dataSize:    "150MB genetic data",
			description: "Genomics research lab uploading genetic analysis",
		},
		{
			clientID:    "emergency_dept_0x7c3e9a",
			patientID:   "emergency_67890",
			gasFeePaid:  big.NewInt(200000000000), // 200 gwei
			dataSize:    "5MB critical vitals",
			description: "Emergency department uploading critical patient data",
		},
	}
	
	for i, workflow := range workflows {
		fmt.Printf("\n--- Workflow %d: %s ---\n", i+1, workflow.description)
		fmt.Printf("Client: %s\n", workflow.clientID)
		fmt.Printf("Patient: %s\n", workflow.patientID)
		fmt.Printf("Data: %s\n", workflow.dataSize)
		fmt.Printf("Gas Fee: %s gwei\n", weiToGwei(workflow.gasFeePaid))
		
		// Step 1: Hash flooding prevention check
		result, err := floodPrevention.ValidateHashRequest(workflow.clientID, workflow.gasFeePaid, "healthcare_data_upload")
		if err != nil {
			fmt.Printf("Error in flood prevention: %v\n", err)
			continue
		}
		
		if !result.Allowed {
			fmt.Printf("❌ Upload blocked by flood prevention: %s\n", result.RejectionReason)
			continue
		}
		
		fmt.Printf("✅ Flood prevention passed - Tier: %s\n", result.CurrentTier)
		
		// Step 2: Simulate cryptographic binding process
		fmt.Println("🔗 Creating cryptographic binding...")
		
		// In real implementation, this would involve actual IPFS upload and blockchain transaction
		mockIPFSHash := fmt.Sprintf("Qm%s...%s", workflow.patientID[:8], workflow.clientID[len(workflow.clientID)-8:])
		mockTxHash := fmt.Sprintf("0x%s...%s", workflow.clientID[2:10], workflow.patientID)
		
		fmt.Printf("   IPFS Hash: %s\n", mockIPFSHash)
		fmt.Printf("   Transaction Hash: %s\n", mockTxHash)
		
		// Step 3: Simulate binding verification
		bindingHash := fmt.Sprintf("0x%x", []byte(workflow.clientID+workflow.patientID))
		fmt.Printf("   Binding Hash: %s\n", bindingHash)
		fmt.Printf("✅ Cryptographic binding created successfully\n")
		
		// Step 4: Access control validation
		fmt.Println("🔒 Validating access control permissions...")
		fmt.Printf("   Department: %s access granted\n", extractDepartment(workflow.clientID))
		fmt.Printf("   Data Type: Healthcare records\n")
		fmt.Printf("   Access Level: %s\n", determineAccessLevel(workflow.gasFeePaid))
		fmt.Printf("✅ Access control validation passed\n")
		
		fmt.Printf("🎉 Workflow completed successfully!\n")
	}
}

// Helper functions
func weiToGwei(wei *big.Int) string {
	gwei := new(big.Int).Div(wei, big.NewInt(1000000000))
	return gwei.String()
}

func extractDepartment(clientID string) string {
	if contains(clientID, "cardiology") {
		return "Cardiology"
	} else if contains(clientID, "research") {
		return "Research"
	} else if contains(clientID, "emergency") {
		return "Emergency"
	}
	return "General"
}

func determineAccessLevel(gasFeePaid *big.Int) string {
	gweiAmount := new(big.Int).Div(gasFeePaid, big.NewInt(1000000000)).Int64()
	
	if gweiAmount >= 100 {
		return "Critical/Real-time"
	} else if gweiAmount >= 25 {
		return "High Priority"
	} else if gweiAmount >= 10 {
		return "Standard"
	}
	return "Basic"
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr ||
		   len(str) >= len(substr) && str[len(str)-len(substr):] == substr
}

func printUsage() {
	fmt.Println("Usage: go run main.go [command]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  binding-demo      - Run cryptographic binding demonstration")
	fmt.Println("  flooding-demo     - Run hash flooding prevention demonstration")
	fmt.Println("  integrated-demo   - Run complete integrated system demo (default)")
	fmt.Println("  service          - Start HTTP API service")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run main.go")
	fmt.Println("  go run main.go integrated-demo")
	fmt.Println("  go run main.go service")
	fmt.Println("  go run main.go flooding-demo")
	fmt.Println("")
	fmt.Println("API Endpoints (when running service):")
	fmt.Println("  POST http://localhost:8080/validate-hash")
	fmt.Println("  GET  http://localhost:8080/service-tiers")
	fmt.Println("  GET  http://localhost:8080/system-metrics")
	fmt.Println("  GET  http://localhost:8080/false-positive-analysis")
	fmt.Println("  GET  http://localhost:8080/health")
}