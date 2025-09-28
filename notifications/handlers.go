package main

import (
	"encoding/json"
	"log"
	"net/http"

	"notifications/models"

	"github.com/gin-gonic/gin"
)

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
