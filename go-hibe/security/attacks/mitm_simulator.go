package attacks

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

// MITMSimulator simulates various Man-in-the-Middle attacks
type MITMSimulator struct {
	TargetHost     string
	TargetPort     int
	ProxyPort      int
	CertificateCA  *x509.Certificate
	PrivateKey     interface{}
	AttackResults  []MITMAttackResult
}

// MITMAttackResult represents the result of a MITM attack simulation
type MITMAttackResult struct {
	AttackType         string            `json:"attackType"`
	Success            bool              `json:"success"`
	DataIntercepted    bool              `json:"dataIntercepted"`
	DataModified       bool              `json:"dataModified"`
	BytesIntercepted   int               `json:"bytesIntercepted"`
	ExecutionTime      time.Duration     `json:"executionTime"`
	ErrorMessage       string            `json:"errorMessage,omitempty"`
	InterceptedData    string            `json:"interceptedData,omitempty"`
	ModifiedData       string            `json:"modifiedData,omitempty"`
	AttackMetadata     map[string]string `json:"attackMetadata"`
}

// NewMITMSimulator creates a new MITM simulator instance
func NewMITMSimulator(targetHost string, targetPort int) *MITMSimulator {
	return &MITMSimulator{
		TargetHost:    targetHost,
		TargetPort:    targetPort,
		ProxyPort:     8888,
		AttackResults: make([]MITMAttackResult, 0),
	}
}

// SimulateCertificateSubstitution simulates certificate substitution attack
func (m *MITMSimulator) SimulateCertificateSubstitution() MITMAttackResult {
	start := time.Now()
	
	// Generate a fake certificate for the target domain
	fakeCert, fakeKey, err := m.generateFakeCertificate(m.TargetHost)
	if err != nil {
		return MITMAttackResult{
			AttackType:    "CertificateSubstitution",
			Success:       false,
			ErrorMessage:  err.Error(),
			ExecutionTime: time.Since(start),
		}
	}
	
	// Create a fake server using the fake certificate
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intercept and log the request
		body, _ := io.ReadAll(r.Body)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":      "intercepted",
			"original_host": r.Host,
			"method":      r.Method,
			"path":        r.URL.Path,
		})
		
		// Log intercepted data
		fmt.Printf("Intercepted request: %s %s\n", r.Method, r.URL.Path)
		fmt.Printf("Intercepted body: %s\n", string(body))
	}))
	
	// Configure server with fake certificate
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{fakeCert.Raw},
				PrivateKey:  fakeKey,
			},
		},
	}
	
	server.StartTLS()
	defer server.Close()
	
	// Test if a client would accept the fake certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Simulating vulnerable client
			},
		},
	}
	
	// Attempt to connect to the fake server
	resp, err := client.Get(server.URL + "/test")
	success := err == nil && resp.StatusCode == http.StatusOK
	
	var interceptedData string
	if success {
		body, _ := io.ReadAll(resp.Body)
		interceptedData = string(body)
		resp.Body.Close()
	}
	
	result := MITMAttackResult{
		AttackType:         "CertificateSubstitution",
		Success:            success,
		DataIntercepted:    success,
		DataModified:       false,
		BytesIntercepted:   len(interceptedData),
		ExecutionTime:      time.Since(start),
		InterceptedData:    interceptedData,
		AttackMetadata: map[string]string{
			"fake_cert_subject": fakeCert.Subject.String(),
			"server_url":        server.URL,
		},
	}
	
	if err != nil {
		result.ErrorMessage = err.Error()
	}
	
	m.AttackResults = append(m.AttackResults, result)
	return result
}

// SimulateSSLStripping simulates SSL stripping attack
func (m *MITMSimulator) SimulateSSLStripping() MITMAttackResult {
	start := time.Now()
	
	// Create a proxy that strips SSL/TLS
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intercept HTTPS requests and serve HTTP responses
		if r.URL.Scheme == "https" || r.Header.Get("X-Forwarded-Proto") == "https" {
			// Strip SSL - redirect to HTTP
			httpURL := strings.Replace(r.URL.String(), "https://", "http://", 1)
			
			// Log the attack
			fmt.Printf("SSL Stripping: %s -> %s\n", r.URL.String(), httpURL)
			
			// Serve modified response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":      "ssl_stripped",
				"original_url": r.URL.String(),
				"modified_url": httpURL,
				"attack":       "successful",
			})
			
			return
		}
		
		// Forward regular HTTP requests
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Normal HTTP response"))
	}))
	defer proxy.Close()
	
	// Test the SSL stripping attack
	client := &http.Client{}
	req, _ := http.NewRequest("GET", proxy.URL+"/secure-endpoint", nil)
	req.Header.Set("X-Forwarded-Proto", "https") // Simulate HTTPS request
	
	resp, err := client.Do(req)
	success := err == nil && resp.StatusCode == http.StatusOK
	
	var interceptedData string
	if success {
		body, _ := io.ReadAll(resp.Body)
		interceptedData = string(body)
		resp.Body.Close()
	}
	
	result := MITMAttackResult{
		AttackType:         "SSLStripping",
		Success:            success,
		DataIntercepted:    success,
		DataModified:       true,
		BytesIntercepted:   len(interceptedData),
		ExecutionTime:      time.Since(start),
		InterceptedData:    interceptedData,
		AttackMetadata: map[string]string{
			"proxy_url":      proxy.URL,
			"attack_method":  "protocol_downgrade",
		},
	}
	
	if err != nil {
		result.ErrorMessage = err.Error()
	}
	
	m.AttackResults = append(m.AttackResults, result)
	return result
}

// SimulateTrafficInterception simulates traffic interception and modification
func (m *MITMSimulator) SimulateTrafficInterception() MITMAttackResult {
	start := time.Now()
	
	originalMessage := "sensitive_user_data"
	modifiedMessage := "modified_by_attacker"
	var interceptedData, modifiedData string
	
	// Create intercepting proxy
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read original request
		body, _ := io.ReadAll(r.Body)
		interceptedData = string(body)
		
		// Modify the request
		if strings.Contains(interceptedData, originalMessage) {
			modifiedData = strings.Replace(interceptedData, originalMessage, modifiedMessage, -1)
		} else {
			modifiedData = interceptedData + "_modified"
		}
		
		// Log the interception
		fmt.Printf("Traffic Interception:\n")
		fmt.Printf("  Original: %s\n", interceptedData)
		fmt.Printf("  Modified: %s\n", modifiedData)
		
		// Send modified response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":           "intercepted",
			"original_data":    interceptedData,
			"modified_data":    modifiedData,
			"attack_success":   "true",
		})
	}))
	defer proxy.Close()
	
	// Send request through proxy
	client := &http.Client{}
	req, _ := http.NewRequest("POST", proxy.URL+"/api/data", strings.NewReader(originalMessage))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	success := err == nil && resp.StatusCode == http.StatusOK
	
	var responseData string
	if success {
		body, _ := io.ReadAll(resp.Body)
		responseData = string(body)
		resp.Body.Close()
	}
	
	result := MITMAttackResult{
		AttackType:         "TrafficInterception",
		Success:            success,
		DataIntercepted:    len(interceptedData) > 0,
		DataModified:       interceptedData != modifiedData,
		BytesIntercepted:   len(interceptedData),
		ExecutionTime:      time.Since(start),
		InterceptedData:    responseData,
		ModifiedData:       modifiedData,
		AttackMetadata: map[string]string{
			"proxy_url":        proxy.URL,
			"original_message": originalMessage,
			"modified_message": modifiedMessage,
		},
	}
	
	if err != nil {
		result.ErrorMessage = err.Error()
	}
	
	m.AttackResults = append(m.AttackResults, result)
	return result
}

// SimulateSessionHijacking simulates session hijacking attack
func (m *MITMSimulator) SimulateSessionHijacking() MITMAttackResult {
	start := time.Now()
	
	capturedSessions := make(map[string]string)
	
	// Create session hijacking proxy
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture session cookies
		for _, cookie := range r.Cookies() {
			if strings.Contains(cookie.Name, "session") || strings.Contains(cookie.Name, "auth") {
				capturedSessions[cookie.Name] = cookie.Value
				fmt.Printf("Captured session: %s=%s\n", cookie.Name, cookie.Value)
			}
		}
		
		// Capture Authorization headers
		if auth := r.Header.Get("Authorization"); auth != "" {
			capturedSessions["Authorization"] = auth
			fmt.Printf("Captured auth header: %s\n", auth)
		}
		
		// Respond with captured session info
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":           "session_hijacked",
			"captured_sessions": capturedSessions,
			"attack_success":   len(capturedSessions) > 0,
		})
	}))
	defer proxy.Close()
	
	// Test session hijacking
	client := &http.Client{}
	req, _ := http.NewRequest("GET", proxy.URL+"/secure-area", nil)
	
	// Add session cookie
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "user123_session_token",
	})
	
	// Add authorization header
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
	
	resp, err := client.Do(req)
	success := err == nil && resp.StatusCode == http.StatusOK && len(capturedSessions) > 0
	
	var responseData string
	if success {
		body, _ := io.ReadAll(resp.Body)
		responseData = string(body)
		resp.Body.Close()
	}
	
	result := MITMAttackResult{
		AttackType:         "SessionHijacking",
		Success:            success,
		DataIntercepted:    len(capturedSessions) > 0,
		DataModified:       false,
		BytesIntercepted:   len(responseData),
		ExecutionTime:      time.Since(start),
		InterceptedData:    responseData,
		AttackMetadata: map[string]string{
			"sessions_captured": fmt.Sprintf("%d", len(capturedSessions)),
			"proxy_url":         proxy.URL,
		},
	}
	
	if err != nil {
		result.ErrorMessage = err.Error()
	}
	
	m.AttackResults = append(m.AttackResults, result)
	return result
}

// SimulateDNSSpoofing simulates DNS spoofing attack
func (m *MITMSimulator) SimulateDNSSpoofing() MITMAttackResult {
	start := time.Now()
	
	// Create fake DNS server
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate malicious server response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":      "dns_spoofed",
			"fake_server": "true",
			"real_host":   m.TargetHost,
			"attack":      "successful",
		})
	}))
	defer fakeServer.Close()
	
	// Test DNS spoofing by directing client to fake server
	client := &http.Client{}
	
	// In a real attack, this would involve DNS manipulation
	// Here we simulate by directly connecting to the fake server
	req, _ := http.NewRequest("GET", fakeServer.URL+"/api/login", nil)
	req.Host = m.TargetHost // Spoof the host header
	
	resp, err := client.Do(req)
	success := err == nil && resp.StatusCode == http.StatusOK
	
	var responseData string
	if success {
		body, _ := io.ReadAll(resp.Body)
		responseData = string(body)
		resp.Body.Close()
	}
	
	result := MITMAttackResult{
		AttackType:         "DNSSpoofing",
		Success:            success,
		DataIntercepted:    success,
		DataModified:       true,
		BytesIntercepted:   len(responseData),
		ExecutionTime:      time.Since(start),
		InterceptedData:    responseData,
		AttackMetadata: map[string]string{
			"fake_server_url": fakeServer.URL,
			"spoofed_host":    m.TargetHost,
			"attack_method":   "server_substitution",
		},
	}
	
	if err != nil {
		result.ErrorMessage = err.Error()
	}
	
	m.AttackResults = append(m.AttackResults, result)
	return result
}

// RunAllAttacks executes all MITM attack simulations
func (m *MITMSimulator) RunAllAttacks() []MITMAttackResult {
	fmt.Println("Starting MITM Attack Simulations...")
	
	// Run all attack simulations
	m.SimulateCertificateSubstitution()
	m.SimulateSSLStripping()
	m.SimulateTrafficInterception()
	m.SimulateSessionHijacking()
	m.SimulateDNSSpoofing()
	
	return m.AttackResults
}

// GenerateReport generates a comprehensive report of all attacks
func (m *MITMSimulator) GenerateReport() string {
	report := "=== MITM Attack Simulation Report ===\n\n"
	
	totalAttacks := len(m.AttackResults)
	successfulAttacks := 0
	totalBytesIntercepted := 0
	
	for _, result := range m.AttackResults {
		if result.Success {
			successfulAttacks++
		}
		totalBytesIntercepted += result.BytesIntercepted
	}
	
	report += fmt.Sprintf("Total attacks simulated: %d\n", totalAttacks)
	report += fmt.Sprintf("Successful attacks: %d\n", successfulAttacks)
	report += fmt.Sprintf("Success rate: %.2f%%\n", float64(successfulAttacks)/float64(totalAttacks)*100)
	report += fmt.Sprintf("Total bytes intercepted: %d\n", totalBytesIntercepted)
	report += "\n=== Individual Attack Results ===\n\n"
	
	for _, result := range m.AttackResults {
		report += fmt.Sprintf("Attack: %s\n", result.AttackType)
		report += fmt.Sprintf("  Success: %v\n", result.Success)
		report += fmt.Sprintf("  Data intercepted: %v\n", result.DataIntercepted)
		report += fmt.Sprintf("  Data modified: %v\n", result.DataModified)
		report += fmt.Sprintf("  Bytes intercepted: %d\n", result.BytesIntercepted)
		report += fmt.Sprintf("  Execution time: %v\n", result.ExecutionTime)
		
		if result.ErrorMessage != "" {
			report += fmt.Sprintf("  Error: %s\n", result.ErrorMessage)
		}
		
		report += "\n"
	}
	
	return report
}

// Helper function to generate fake certificates
func (m *MITMSimulator) generateFakeCertificate(hostname string) (*x509.Certificate, interface{}, error) {
	// Generate private key
	priv, err := tls.LoadX509KeyPair("", "") // This would normally load a real key
	if err != nil {
		// Generate a temporary key for simulation
		// In practice, you'd use crypto/rsa or crypto/ecdsa
		return nil, nil, fmt.Errorf("key generation not implemented for simulation")
	}
	
	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Fake CA"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{hostname},
		BasicConstraintsValid: true,
	}
	
	// For simulation purposes, return a minimal certificate structure
	return &template, priv.PrivateKey, nil
}

// GetAttackResults returns all attack results
func (m *MITMSimulator) GetAttackResults() []MITMAttackResult {
	return m.AttackResults
}

// ClearResults clears all attack results
func (m *MITMSimulator) ClearResults() {
	m.AttackResults = make([]MITMAttackResult, 0)
}