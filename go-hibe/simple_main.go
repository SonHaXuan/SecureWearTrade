package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 HIBE API Server Starting...")

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
			"message":   "HIBE API is running!",
		})
	})

	// Basic info endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "HIBE API",
			"status":  "running",
			"version": "1.0.0",
			"endpoints": []string{
				"/health",
				"/",
			},
		})
	})

	fmt.Println("✅ HIBE API Server started successfully!")
	fmt.Println("🌐 Server running on http://localhost:8080")
	fmt.Println("📊 Health check: http://localhost:8080/health")

	// Start server
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
	}
}