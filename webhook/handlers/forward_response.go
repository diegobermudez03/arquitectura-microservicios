package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"webhook/internal"
	"webhook/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// Parse the issuer response
	var issuerResponse models.IssuerResponse
	if err := json.Unmarshal(body, &issuerResponse); err != nil {
		log.Printf("Error parsing issuer response: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issuer response format"})
		return
	}

	// Validate required fields
	if issuerResponse.SuscriptorToken == "" {
		log.Printf("Missing suscriptor_token in issuer response")
		c.JSON(http.StatusBadRequest, gin.H{"error": "suscriptor_token is required"})
		return
	}

	// Look up suscriptor info in Redis
	suscriptor, err := h.redisService.GetSuscriptor(issuerResponse.SuscriptorToken)
	if err != nil {
		log.Printf("Error retrieving suscriptor %s: %v", issuerResponse.SuscriptorToken, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Suscriptor not found"})
		return
	}

	callbackURL, exists := suscriptor["callback_url"]
	if !exists || callbackURL == "" {
		log.Printf("No callback URL found for suscriptor %s", issuerResponse.SuscriptorToken)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No callback URL configured"})
		return
	}

	// Create webhook event
	webhookEvent := h.createWebhookEvent(issuerResponse, suscriptor["name"])

	// Forward the webhook event to the suscriptor
	if err := h.forwardWebhookEvent(callbackURL, webhookEvent); err != nil {
		log.Printf("Error forwarding webhook event to %s: %v", callbackURL, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward webhook event"})
		return
	}

	log.Printf("Webhook event successfully forwarded to %s for request %s", callbackURL, issuerResponse.RequestUUID)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":  "Webhook event forwarded successfully",
		"event_id": webhookEvent.ID,
		"status":   issuerResponse.Status,
	})
}

func (h *ForwardResponseHandler) createWebhookEvent(issuerResponse models.IssuerResponse, suscriptorName string) models.WebhookEvent {
	// Generate unique event ID
	eventID := uuid.New().String()

	// Determine event type based on status
	eventType := "card.issued"
	if issuerResponse.Status == "declined" {
		eventType = "card.declined"
	}

	// Create metadata
	metadata := map[string]interface{}{
		"suscriptor_name": suscriptorName,
		"request_uuid":    issuerResponse.RequestUUID,
	}

	// Add card-specific metadata if card was issued
	if issuerResponse.IssuedCard != nil {
		metadata["card_type"] = issuerResponse.IssuedCard.CardType
		metadata["has_pan"] = issuerResponse.IssuedCard.PAN != ""
	}

	// Add decline reason if declined
	if issuerResponse.DeclineReason != nil {
		metadata["decline_reason"] = issuerResponse.DeclineReason.Reason
	}

	return models.WebhookEvent{
		ID:        eventID,
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    "webhook-service",
		Data:      issuerResponse,
		Metadata:  metadata,
	}
}

func (h *ForwardResponseHandler) forwardWebhookEvent(callbackURL string, event models.WebhookEvent) error {
	// Marshal the webhook event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook event: %w", err)
	}

	// Forward the webhook event to the suscriptor
	resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(eventJSON))
	if err != nil {
		return fmt.Errorf("failed to send webhook event: %w", err)
	}
	defer resp.Body.Close()

	// Log the result
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Webhook event %s successfully delivered to %s, status: %d", event.ID, callbackURL, resp.StatusCode)
	} else {
		log.Printf("Webhook event %s delivery to %s returned non-2xx status: %d", event.ID, callbackURL, resp.StatusCode)
	}

	return nil
}
