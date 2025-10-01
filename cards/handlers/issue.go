package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cards/internal"
	"cards/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IssueHandler struct {
	redisService    *internal.RedisService
	postgresService *internal.PostgresService
}

func NewIssueHandler(redisService *internal.RedisService, postgresService *internal.PostgresService) *IssueHandler {
	return &IssueHandler{
		redisService:    redisService,
		postgresService: postgresService,
	}
}

func (h *IssueHandler) Issue(c *gin.Context) {
	var req models.IssueCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received card issue request for: %v", req)
	ctx := context.Background()

	user, err := h.redisService.GetUser(ctx, req.UserToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	requestUUID := uuid.New().String()

	issueRequest := models.IssueRequest{
		Name:            user.Name,
		Lastname:        user.Lastname,
		BirthDate:       user.BirthDate,
		CountryCode:     user.CountryCode,
		CardType:        req.CardType,
		SuscriptorToken: os.Getenv("SUSCRIPTOR_TOKEN"),
		RequestUUID:     requestUUID,
	}

	requestData := models.RequestData{
		User:      *user,
		CardType:  req.CardType,
		UserToken: req.UserToken,
	}

	if err := h.redisService.StoreRequest(ctx, requestUUID, requestData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store request"})
		return
	}

	log.Printf("Storing request in Redis for user UUID and request UUID: %s and %s", req.UserToken, requestUUID)
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WEBHOOK_URL not configured"})
		return
	}

	requestJSON, err := json.Marshal(issueRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to webhook"})
		return
	}
	defer resp.Body.Close()

	c.Status(http.StatusAccepted)
}
