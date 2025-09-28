package main

import (
	"issuer/handlers"

	"github.com/gin-gonic/gin"
)

func setupRoutes() *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Card issue endpoint
	r.POST("/issue", handlers.IssueCard)

	return r
}
