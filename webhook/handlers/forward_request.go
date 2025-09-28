package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ForwardRequestHandler struct{}

func NewForwardRequestHandler() *ForwardRequestHandler {
	return &ForwardRequestHandler{}
}

func (h *ForwardRequestHandler) HandleForwardRequest(c *gin.Context) {
	// Get the issuer URL from environment
	issuerURL := os.Getenv("ISSUER_URL")
	if issuerURL == "" {
		log.Printf("ISSUER_URL environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ISSUER_URL not configured"})
		return
	}

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Forward the request to the issuer
	resp, err := http.Post(issuerURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error forwarding request to issuer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request to issuer"})
		return
	}
	defer resp.Body.Close()

	// Read the response from the issuer
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading issuer response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read issuer response"})
		return
	}

	// Set the same status code and headers
	c.Status(resp.StatusCode)
	c.Header("Content-Type", resp.Header.Get("Content-Type"))

	// Try to parse as JSON and return as JSON, otherwise return as string
	var jsonResponse interface{}
	if err := json.Unmarshal(responseBody, &jsonResponse); err == nil {
		c.JSON(resp.StatusCode, jsonResponse)
	} else {
		c.String(resp.StatusCode, string(responseBody))
	}

	log.Printf("Request forwarded to issuer, status: %d", resp.StatusCode)
}
