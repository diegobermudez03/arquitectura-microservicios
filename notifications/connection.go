package main

import (
	"log"
	"notifications/models"
	"sync"
)

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
