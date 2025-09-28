package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"cards/internal"
	"cards/models"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	redisService    *internal.RedisService
	postgresService *internal.PostgresService
}

func NewWebhookHandler(redisService *internal.RedisService, postgresService *internal.PostgresService) *WebhookHandler {
	return &WebhookHandler{
		redisService:    redisService,
		postgresService: postgresService,
	}
}

func (h *WebhookHandler) Webhook(c *gin.Context) {
	var webhookEvent models.WebhookEvent
	if err := c.ShouldBindJSON(&webhookEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	response := webhookEvent.Data

	requestData, err := h.redisService.GetRequest(ctx, response.RequestUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	userToken := findUserTokenByRequest(ctx, h.redisService, response.RequestUUID)
	if userToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User token not found"})
		return
	}

	// Get user from database
	userRecord, err := h.postgresService.GetUserByToken(userToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in database"})
		return
	}

	// Process based on issuance result
	if response.IssuedCard != nil {
		// Successful issuance - store in issued_cards table
		issuedCardRecord := h.postgresService.CreateIssuedCardRecord(
			userRecord.ID,
			userToken,
			response,
		)

		if err := h.postgresService.StoreIssuedCard(issuedCardRecord); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store issued card"})
			return
		}
	} else if response.DeclineReason != nil {
		// Failed attempt - store in failed_attempts table
		failedAttemptRecord := h.postgresService.CreateFailedAttemptRecord(
			userRecord.ID,
			userToken,
			requestData.CardType,
			response,
		)

		if err := h.postgresService.StoreFailedAttempt(failedAttemptRecord); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store failed attempt"})
			return
		}
	}

	// Clean up Redis request
	if err := h.redisService.DeleteRequest(ctx, response.RequestUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}

	// Send notification
	notificationRequest := models.NotificationRequest{
		UserToken:      userToken,
		IssuerResponse: response,
	}

	notificationsURL := os.Getenv("NOTIFICATIONS_URL")
	if notificationsURL != "" {
		notificationJSON, err := json.Marshal(notificationRequest)
		if err == nil {
			http.Post(notificationsURL, "application/json", bytes.NewBuffer(notificationJSON))
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func findUserTokenByRequest(ctx context.Context, redisService *internal.RedisService, requestUUID string) string {
	requestData, err := redisService.GetRequest(ctx, requestUUID)
	if err != nil {
		return ""
	}

	userToken := ""
	keys, err := redisService.GetAllUserKeys(ctx)
	if err != nil {
		return ""
	}

	for _, key := range keys {
		user, err := redisService.GetUser(ctx, key[5:])
		if err == nil && user.Name == requestData.User.Name && user.Lastname == requestData.User.Lastname {
			userToken = key[5:]
			break
		}
	}

	return userToken
}
