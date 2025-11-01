package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Service  string    `json:"service"`
	Status   string    `json:"status"`
	Time     time.Time `json:"timestamp"`
	Message  string    `json:"message"`
	Endpoints []string `json:"endpoints"`
}

func main() {
	fmt.Println("🚀 HIBE API Server Starting...")

	// Simple HTTP server without Gin
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := Response{
			Service: "HIBE API",
			Status:  "running",
			Time:    time.Now(),
			Message: "HIBE container is working successfully!",
			Endpoints: []string{
				"/",
				"/health",
			},
		}

		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := Response{
			Service:  "HIBE Health Check",
			Status:   "healthy",
			Time:     time.Now(),
			Message:  "All systems operational!",
			Endpoints: []string{},
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("✅ HIBE API Server started successfully!")
	fmt.Println("🌐 Server running on http://localhost:8080")
	fmt.Println("📊 Health check: http://localhost:8080/health")
	fmt.Println("🐳 Docker container is running!")

	// Start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}