package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"webhook/internal"
	"webhook/models"

	"github.com/gin-gonic/gin"
)

type SuscribeHandler struct {
	redisService *internal.RedisService
}

func NewSuscribeHandler(redisService *internal.RedisService) *SuscribeHandler {
	return &SuscribeHandler{
		redisService: redisService,
	}
}

func (h *SuscribeHandler) HandleSuscribe(c *gin.Context) {
	var req models.SuscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate secure token
	token, err := generateSecureToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Store suscriptor info in Redis
	suscriptor := map[string]string{
		"name":         req.Name,
		"callback_url": req.CallbackURL,
	}

	if err := h.redisService.StoreSuscriptor(token, suscriptor); err != nil {
		log.Printf("Error storing suscriptor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store suscriptor"})
		return
	}

	log.Printf("New suscriptor registered: %s with token: %s", req.Name, token)

	response := models.SuscribeResponse{
		SuscriptorToken: token,
	}

	c.JSON(http.StatusOK, response)
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
