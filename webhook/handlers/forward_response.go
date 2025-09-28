package handlers

import (
"bytes"
"encoding/json"
"io"
"log"
"net/http"

"github.com/gin-gonic/gin"
"webhook/internal"
)

type ForwardResponseHandler struct {
redisService *internal.RedisService
}

func NewForwardResponseHandler(redisService *internal.RedisService) *ForwardResponseHandler {
return &ForwardResponseHandler{
redisService: redisService,
}
}

func (h *ForwardResponseHandler) HandleForwardResponse(c *gin.Context) {
// Read the request body
body, err := io.ReadAll(c.Request.Body)
if err != nil {
log.Printf("Error reading request body: %v", err)
c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
return
}

// Parse the JSON to extract suscriptor_token
var requestData map[string]interface{}
if err := json.Unmarshal(body, &requestData); err != nil {
log.Printf("Error parsing JSON: %v", err)
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
return
}

// Extract suscriptor_token
token, ok := requestData["suscriptor_token"].(string)
if !ok || token == "" {
log.Printf("Missing or invalid suscriptor_token")
c.JSON(http.StatusBadRequest, gin.H{"error": "suscriptor_token is required"})
return
}

// Look up suscriptor info in Redis
suscriptor, err := h.redisService.GetSuscriptor(token)
if err != nil {
log.Printf("Error retrieving suscriptor %s: %v", token, err)
c.JSON(http.StatusNotFound, gin.H{"error": "Suscriptor not found"})
return
}

callbackURL, exists := suscriptor["callback_url"]
if !exists || callbackURL == "" {
log.Printf("No callback URL found for suscriptor %s", token)
c.JSON(http.StatusInternalServerError, gin.H{"error": "No callback URL configured"})
return
}

// Forward the response to the suscriptor callback URL
resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(body))
if err != nil {
log.Printf("Error forwarding response to %s: %v", callbackURL, err)
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward response to callback"})
return
}
defer resp.Body.Close()

// Log the result
if resp.StatusCode >= 200 && resp.StatusCode < 300 {
log.Printf("Response successfully forwarded to %s, status: %d", callbackURL, resp.StatusCode)
} else {
log.Printf("Callback to %s returned non-2xx status: %d", callbackURL, resp.StatusCode)
}

// Return success response
c.JSON(http.StatusOK, gin.H{
"message": "Response forwarded successfully",
"status":  resp.StatusCode,
})
}
