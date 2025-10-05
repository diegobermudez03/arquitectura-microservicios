package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notifications/models"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS to allow all connections
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
	}))

	// Register routes
	registerRoutes(r)

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	fmt.Println("running on port", port)
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

// StreamHandler handles SSE connections for real-time notifications
func StreamHandler(c *gin.Context) {
	// Get user token from query parameter
	userToken := c.Query("user_token")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_token is required"})
		return
	}

	log.Printf("New SSE connection request for user: %s", userToken)

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Create a channel for this connection
	notificationChan := connManager.AddConnection(userToken)

	// Ensure connection is cleaned up when done
	defer func() {
		connManager.RemoveConnection(userToken)
		log.Printf("SSE connection closed for user: %s", userToken)
	}()

	// Send initial connection event
	c.String(http.StatusOK, "data: {\"status\":\"connected\"}\n\n")
	c.Writer.Flush()

	// Wait for notification
	select {
	case notification := <-notificationChan:
		// Send the notification as JSON
		jsonData, err := json.Marshal(notification)
		if err != nil {
			log.Printf("Error marshaling notification: %v", err)
			c.String(http.StatusInternalServerError, "data: {\"error\":\"internal server error\"}\n\n")
			return
		}

		// Send SSE formatted data
		c.String(http.StatusOK, "data: %s\n\n", string(jsonData))
		c.Writer.Flush()
		log.Printf("Notification sent via SSE to user: %s", userToken)

	case <-c.Request.Context().Done():
		// Client disconnected
		log.Printf("Client disconnected for user: %s", userToken)
		return
	}
}

// SendHandler handles notification sending requests
func SendHandler(c *gin.Context) {
	var req models.NotificationRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Validate required fields
	if req.UserToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_token is required"})
		return
	}

	log.Printf("Received notification request for user: %s", req.UserToken)

	// Try to send notification
	success := connManager.SendNotification(req.UserToken, req.IssuerResponse)

	if success {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Notification sent successfully",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not connected or notification failed",
		})
	}
}

// ConnectionManager manages active SSE connections
type ConnectionManager struct {
	connections map[string]chan models.IssuerResponse
	mutex       sync.RWMutex
}

// Global connection manager instance
var connManager = &ConnectionManager{
	connections: make(map[string]chan models.IssuerResponse),
}

// AddConnection adds a new connection for a user
func (cm *ConnectionManager) AddConnection(userToken string) chan models.IssuerResponse {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Close existing connection if any
	if existingChan, exists := cm.connections[userToken]; exists {
		close(existingChan)
	}

	// Create new channel for this user
	ch := make(chan models.IssuerResponse, 1)
	cm.connections[userToken] = ch
	log.Printf("Connection added for user: %s", userToken)
	return ch
}

// SendNotification sends a notification to a user
func (cm *ConnectionManager) SendNotification(userToken string, response models.IssuerResponse) bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if ch, exists := cm.connections[userToken]; exists {
		select {
		case ch <- response:
			log.Printf("Notification sent to user: %s", userToken)
			return true
		default:
			log.Printf("Failed to send notification to user: %s (channel full)", userToken)
			return false
		}
	}

	log.Printf("No active connection found for user: %s", userToken)
	return false
}

// RemoveConnection removes a connection for a user
func (cm *ConnectionManager) RemoveConnection(userToken string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if ch, exists := cm.connections[userToken]; exists {
		close(ch)
		delete(cm.connections, userToken)
		log.Printf("Connection removed for user: %s", userToken)
	}
}
