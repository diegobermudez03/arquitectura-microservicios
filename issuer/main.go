package main

import (
	"issuer/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	webhookURL := os.Getenv("WEBHOOK_URL")
	if port == "" {
		port = "8080"
	}

	// Setup routes
	r := setupRoutes(handlers.NewHandlers(webhookURL))

	log.Printf("Starting Service on port %s", port)
	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(h *handlers.Handlers) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	v1 := r.Group("/v1")

	// Card issue endpoint
	v1.POST("/cards", h.IssueCard)

	return r
}
