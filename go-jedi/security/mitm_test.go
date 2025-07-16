package security

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// MITMTestResult represents the result of a MITM attack test
type MITMTestResult struct {
	TestName           string        `json:"testName"`
	AttackSuccessful   bool          `json:"attackSuccessful"`
	DataIntercepted    bool          `json:"dataIntercepted"`
	DataModified       bool          `json:"dataModified"`
	ExecutionTime      time.Duration `json:"executionTime"`
	ConfidenceInterval float64       `json:"confidenceInterval"`
	TLSBypassAttempt   bool          `json:"tlsBypassAttempt"`
	CertificateValid   bool          `json:"certificateValid"`
	ErrorMessage       string        `json:"errorMessage,omitempty"`
}

// MITMExperiment conducts experimental validation of MITM attack resistance
type MITMExperiment struct {
	TestIterations int
	SampleSize     int
	TargetEndpoint string
	Results        []MITMTestResult
}

// RunMITMResistanceTests performs comprehensive MITM attack resistance testing
func TestMITMResistance(t *testing.T) {
	experiment := &MITMExperiment{
		TestIterations: 100,
		SampleSize:     1000,
		TargetEndpoint: "/encrypt",
		Results:        make([]MITMTestResult, 0),
	}

	// Test 1: Certificate Validation Test
	t.Run("CertificateValidation", func(t *testing.T) {
		result := experiment.testCertificateValidation()
		experiment.Results = append(experiment.Results, result)
		
		if result.AttackSuccessful {
			t.Errorf("Certificate validation bypassed: %s", result.ErrorMessage)
		}
	})

	// Test 2: TLS Downgrade Attack Test
	t.Run("TLSDowngradeAttack", func(t *testing.T) {
		result := experiment.testTLSDowngradeAttack()
		experiment.Results = append(experiment.Results, result)
		
		if result.AttackSuccessful {
			t.Errorf("TLS downgrade attack succeeded: %s", result.ErrorMessage)
		}
	})

	// Test 3: Data Interception Test
	t.Run("DataInterception", func(t *testing.T) {
		result := experiment.testDataInterception()
		experiment.Results = append(experiment.Results, result)
		
		if result.DataIntercepted {
			t.Errorf("Data interception successful: %s", result.ErrorMessage)
		}
	})

	// Test 4: Message Modification Test
	t.Run("MessageModification", func(t *testing.T) {
		result := experiment.testMessageModification()
		experiment.Results = append(experiment.Results, result)
		
		if result.DataModified {
			t.Errorf("Message modification attack succeeded: %s", result.ErrorMessage)
		}
	})

	// Test 5: Session Hijacking Test
	t.Run("SessionHijacking", func(t *testing.T) {
		result := experiment.testSessionHijacking()
		experiment.Results = append(experiment.Results, result)
		
		if result.AttackSuccessful {
			t.Errorf("Session hijacking attack succeeded: %s", result.ErrorMessage)
		}
	})

	// Generate statistical analysis
	experiment.generateStatisticalAnalysis(t)
}

// testCertificateValidation tests certificate validation bypass attempts
func (e *MITMExperiment) testCertificateValidation() MITMTestResult {
	start := time.Now()
	
	// Create a test server with self-signed certificate
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"data": "encrypted_data"})
	}))
	defer server.Close()

	// Create client that accepts invalid certificates (simulating MITM)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // This simulates a vulnerable client
			},
		},
	}

	// Test if we can bypass certificate validation
	resp, err := client.Get(server.URL + e.TargetEndpoint)
	attackSuccessful := err == nil && resp.StatusCode == http.StatusOK
	
	// In a secure implementation, this should fail with certificate validation
	// The test should verify that proper certificate pinning is implemented
	
	return MITMTestResult{
		TestName:           "CertificateValidation",
		AttackSuccessful:   attackSuccessful,
		DataIntercepted:    attackSuccessful,
		DataModified:       false,
		ExecutionTime:      time.Since(start),
		ConfidenceInterval: 0.95,
		TLSBypassAttempt:   true,
		CertificateValid:   false,
		ErrorMessage:       formatError(err),
	}
}

// testTLSDowngradeAttack tests TLS protocol downgrade attacks
func (e *MITMExperiment) testTLSDowngradeAttack() MITMTestResult {
	start := time.Now()
	
	// Create server with weak TLS configuration
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "vulnerable"})
	}))
	defer server.Close()

	// Attempt to force TLS 1.0 or 1.1 (should be rejected)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS11, // Try to force weak TLS
			},
		},
	}

	resp, err := client.Get(server.URL + e.TargetEndpoint)
	attackSuccessful := err == nil && resp.StatusCode == http.StatusOK

	return MITMTestResult{
		TestName:           "TLSDowngradeAttack",
		AttackSuccessful:   attackSuccessful,
		DataIntercepted:    false,
		DataModified:       false,
		ExecutionTime:      time.Since(start),
		ConfidenceInterval: 0.95,
		TLSBypassAttempt:   true,
		CertificateValid:   true,
		ErrorMessage:       formatError(err),
	}
}

// testDataInterception tests ability to intercept encrypted data
func (e *MITMExperiment) testDataInterception() MITMTestResult {
	start := time.Now()
	
	// Create a proxy that intercepts traffic
	var interceptedData []byte
	interceptor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and store the body
		body, _ := io.ReadAll(r.Body)
		interceptedData = body
		
		// Forward to actual server (simulating transparent proxy)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"intercepted": "true"})
	}))
	defer interceptor.Close()

	// Generate test payload
	testPayload := generateTestPayload(1024)
	
	// Send request through interceptor
	resp, err := http.Post(interceptor.URL+e.TargetEndpoint, "application/json", bytes.NewBuffer(testPayload))
	
	dataIntercepted := len(interceptedData) > 0
	attackSuccessful := err == nil && resp.StatusCode == http.StatusOK && dataIntercepted

	return MITMTestResult{
		TestName:           "DataInterception",
		AttackSuccessful:   attackSuccessful,
		DataIntercepted:    dataIntercepted,
		DataModified:       false,
		ExecutionTime:      time.Since(start),
		ConfidenceInterval: 0.95,
		TLSBypassAttempt:   false,
		CertificateValid:   true,
		ErrorMessage:       formatError(err),
	}
}

// testMessageModification tests ability to modify messages in transit
func (e *MITMExperiment) testMessageModification() MITMTestResult {
	start := time.Now()
	
	originalMessage := "original_secure_message"
	modifiedMessage := "modified_malicious_message"
	var actualMessage string
	
	// Create server that logs received messages
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		actualMessage = string(body)
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"received": actualMessage})
	}))
	defer server.Close()

	// Create MITM proxy that modifies messages
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read original message
		body, _ := io.ReadAll(r.Body)
		
		// Modify message (simulating MITM attack)
		modifiedBody := strings.ReplaceAll(string(body), originalMessage, modifiedMessage)
		
		// Forward modified message
		client := &http.Client{}
		req, _ := http.NewRequest(r.Method, server.URL+r.URL.Path, strings.NewReader(modifiedBody))
		resp, _ := client.Do(req)
		
		// Return response
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}))
	defer proxy.Close()

	// Send original message through proxy
	_, err := http.Post(proxy.URL+e.TargetEndpoint, "application/json", strings.NewReader(originalMessage))
	
	dataModified := actualMessage == modifiedMessage
	attackSuccessful := err == nil && dataModified

	return MITMTestResult{
		TestName:           "MessageModification",
		AttackSuccessful:   attackSuccessful,
		DataIntercepted:    true,
		DataModified:       dataModified,
		ExecutionTime:      time.Since(start),
		ConfidenceInterval: 0.95,
		TLSBypassAttempt:   false,
		CertificateValid:   true,
		ErrorMessage:       formatError(err),
	}
}

// testSessionHijacking tests session hijacking resistance
func (e *MITMExperiment) testSessionHijacking() MITMTestResult {
	start := time.Now()
	
	// Create server with session management
	sessions := make(map[string]bool)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		if sessions[sessionID] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"access": "granted"})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"access": "denied"})
		}
	}))
	defer server.Close()

	// Create legitimate session
	legitimateSessionID := "legitimate_session_123"
	sessions[legitimateSessionID] = true
	
	// Attempt to hijack session with different ID
	hijackedSessionID := "hijacked_session_456"
	
	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL+e.TargetEndpoint, nil)
	req.Header.Set("X-Session-ID", hijackedSessionID)
	
	resp, err := client.Do(req)
	attackSuccessful := err == nil && resp.StatusCode == http.StatusOK

	return MITMTestResult{
		TestName:           "SessionHijacking",
		AttackSuccessful:   attackSuccessful,
		DataIntercepted:    false,
		DataModified:       false,
		ExecutionTime:      time.Since(start),
		ConfidenceInterval: 0.95,
		TLSBypassAttempt:   false,
		CertificateValid:   true,
		ErrorMessage:       formatError(err),
	}
}

// generateStatisticalAnalysis generates statistical analysis of test results
func (e *MITMExperiment) generateStatisticalAnalysis(t *testing.T) {
	totalTests := len(e.Results)
	successfulAttacks := 0
	totalTime := time.Duration(0)
	
	for _, result := range e.Results {
		if result.AttackSuccessful {
			successfulAttacks++
		}
		totalTime += result.ExecutionTime
	}
	
	successRate := float64(successfulAttacks) / float64(totalTests)
	avgTime := totalTime / time.Duration(totalTests)
	
	// Calculate confidence interval for success rate
	confidenceInterval := calculateConfidenceInterval(successRate, totalTests, 0.95)
	
	t.Logf("=== MITM Attack Resistance Test Results ===")
	t.Logf("Total tests: %d", totalTests)
	t.Logf("Successful attacks: %d", successfulAttacks)
	t.Logf("Success rate: %.2f%%", successRate*100)
	t.Logf("Average execution time: %v", avgTime)
	t.Logf("95%% Confidence interval: [%.2f%%, %.2f%%]", 
		confidenceInterval[0]*100, confidenceInterval[1]*100)
	
	// Log individual test results
	for _, result := range e.Results {
		t.Logf("Test: %s, Attack successful: %v, Time: %v", 
			result.TestName, result.AttackSuccessful, result.ExecutionTime)
	}
	
	// Assert that attack success rate is below acceptable threshold
	if successRate > 0.05 { // 5% threshold
		t.Errorf("MITM attack success rate too high: %.2f%% (threshold: 5%%)", successRate*100)
	}
}

// Helper functions
func generateTestPayload(size int) []byte {
	payload := make([]byte, size)
	rand.Read(payload)
	return payload
}

func formatError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func calculateConfidenceInterval(proportion float64, sampleSize int, confidence float64) [2]float64 {
	// Simple normal approximation for confidence interval
	z := 1.96 // For 95% confidence
	se := (proportion * (1 - proportion)) / float64(sampleSize)
	margin := z * se
	
	return [2]float64{
		proportion - margin,
		proportion + margin,
	}
}

// BenchmarkMITMResistance provides performance benchmarks for MITM resistance
func BenchmarkMITMResistance(b *testing.B) {
	experiment := &MITMExperiment{
		TestIterations: b.N,
		SampleSize:     1000,
		TargetEndpoint: "/encrypt",
		Results:        make([]MITMTestResult, 0),
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result := experiment.testCertificateValidation()
		experiment.Results = append(experiment.Results, result)
	}
}