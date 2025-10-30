package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type EncryptRequest struct {
	URI     string `json:"uri"`
	Message string `json:"message"`
}

type DecryptRequest struct {
	URI              string `json:"uri"`
	EncryptedMessage string `json:"encryptedMessage"`
}

type EncryptResponse struct {
	Success        bool    `json:"success"`
	Data           string  `json:"data,omitempty"`
	Error          string  `json:"error,omitempty"`
	ExecutionTime  int64   `json:"executionTime"`
	MemoryUsage    uint64  `json:"memoryUsage"`
	CPUPercentage  float64 `json:"cpuPercentage"`
	PowerUsage     float64 `json:"powerUsage"`
	EnergyJoules   float64 `json:"energyConsumptionJoules"`
}

type DecryptResponse struct {
	Success        bool    `json:"success"`
	Data           string  `json:"data,omitempty"`
	Error          string  `json:"error,omitempty"`
	ExecutionTime  int64   `json:"executionTime"`
	MemoryUsage    uint64  `json:"memoryUsage"`
	CPUPercentage  float64 `json:"cpuPercentage"`
	PowerUsage     float64 `json:"powerUsage"`
	EnergyJoules   float64 `json:"energyConsumptionJoules"`
}

type BaseResponse struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Time      time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Endpoints []string  `json:"endpoints"`
}

// Simple encryption function for demonstration
func encryptMessage(message string) (string, error) {
	key := sha256.Sum256([]byte("jedi-demo-key"))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// Create byte array of the size of block size + message length
	ciphertext := make([]byte, aes.BlockSize+len(message))
	iv := ciphertext[:aes.BlockSize]

	// Generate random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(message))

	// Return base64 encoded ciphertext
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Simple decryption function for demonstration
func decryptMessage(encryptedMessage string) (string, error) {
	key := sha256.Sum256([]byte("jedi-demo-key"))

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// Simulate system metrics
func getSystemMetrics() (uint64, float64, float64) {
	// Simulate memory usage, CPU percentage, and power usage
	memory := uint64(1024 * 1024 * 50) // 50MB
	cpu := 15.5                       // 15.5%
	power := 2.5                      // 2.5 Watts
	return memory, cpu, power
}

func main() {
	fmt.Println("🚀 JEDI Enhanced API Server Starting...")

	// Main endpoint handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := BaseResponse{
			Service:   "JEDI Enhanced API",
			Status:    "running",
			Time:      time.Now(),
			Message:   "JEDI container with encryption is working successfully!",
			Endpoints: []string{"/", "/health", "/encrypt", "/decrypt"},
		}

		json.NewEncoder(w).Encode(response)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := BaseResponse{
			Service:   "JEDI Enhanced Health Check",
			Status:    "healthy",
			Time:      time.Now(),
			Message:   "All systems operational including encryption!",
			Endpoints: []string{},
		}

		json.NewEncoder(w).Encode(response)
	})

	// Encrypt endpoint
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(EncryptResponse{
				Success: false,
				Error:   "Method not allowed. Use POST.",
			})
			return
		}

		var encryptReq EncryptRequest
		err := json.NewDecoder(r.Body).Decode(&encryptReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(EncryptResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid JSON: %v", err),
			})
			return
		}

		// Get system metrics before encryption
		memory, cpu, power := getSystemMetrics()

		// Perform encryption
		encryptedData, err := encryptMessage(encryptReq.Message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(EncryptResponse{
				Success: false,
				Error:   fmt.Sprintf("Encryption failed: %v", err),
			})
			return
		}

		// Calculate metrics
		executionTime := time.Since(startTime).Microseconds()
		executionTimeSeconds := float64(executionTime) / 1000000.0
		energyConsumption := power * executionTimeSeconds

		response := EncryptResponse{
			Success:        true,
			Data:           encryptedData,
			ExecutionTime:  executionTime,
			MemoryUsage:    memory,
			CPUPercentage:  cpu,
			PowerUsage:     power,
			EnergyJoules:   energyConsumption,
		}

		json.NewEncoder(w).Encode(response)
	})

	// Decrypt endpoint
	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(DecryptResponse{
				Success: false,
				Error:   "Method not allowed. Use POST.",
			})
			return
		}

		var decryptReq DecryptRequest
		err := json.NewDecoder(r.Body).Decode(&decryptReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DecryptResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid JSON: %v", err),
			})
			return
		}

		// Get system metrics before decryption
		memory, cpu, power := getSystemMetrics()

		// Perform decryption
		decryptedData, err := decryptMessage(decryptReq.EncryptedMessage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(DecryptResponse{
				Success: false,
				Error:   fmt.Sprintf("Decryption failed: %v", err),
			})
			return
		}

		// Calculate metrics
		executionTime := time.Since(startTime).Microseconds()
		executionTimeSeconds := float64(executionTime) / 1000000.0
		energyConsumption := power * executionTimeSeconds

		response := DecryptResponse{
			Success:        true,
			Data:           decryptedData,
			ExecutionTime:  executionTime,
			MemoryUsage:    memory,
			CPUPercentage:  cpu,
			PowerUsage:     power,
			EnergyJoules:   energyConsumption,
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("✅ JEDI Enhanced API Server started successfully!")
	fmt.Println("🌐 Server running on http://localhost:8080")
	fmt.Println("📊 Health check: http://localhost:8080/health")
	fmt.Println("🔐 Encrypt API: http://localhost:8080/encrypt")
	fmt.Println("🔓 Decrypt API: http://localhost:8080/decrypt")
	fmt.Println("🐳 Docker container with encryption is running!")

	// Start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}