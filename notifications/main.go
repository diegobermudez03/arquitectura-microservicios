package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Register routes
	registerRoutes(r)

	// Start server
	log.Println("Starting Notifications Service on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// registerRoutes sets up all the API routes
func registerRoutes(r *gin.Engine) {
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// SSE endpoint for real-time notifications
	r.GET("/notifications/stream", StreamHandler)

	// Endpoint to send notifications
	r.POST("/notify", SendHandler)
}
