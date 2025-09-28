package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"cards/internal"
	"cards/models"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	redisService    *internal.RedisService
	postgresService *internal.PostgresService
}

func NewRegisterHandler(redisService *internal.RedisService, postgresService *internal.PostgresService) *RegisterHandler {
	return &RegisterHandler{
		redisService:    redisService,
		postgresService: postgresService,
	}
}

func (h *RegisterHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := generateRandomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user := models.User{
		Name:        req.Name,
		Lastname:    req.Lastname,
		BirthDate:   req.BirthDate,
		CountryCode: req.CountryCode,
	}

	ctx := context.Background()

	// Store user in Redis
	if err := h.redisService.StoreUser(ctx, token, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user in Redis"})
		return
	}

	// Store user in PostgreSQL
	_, err = h.postgresService.StoreUser(token, user)
	if err != nil {
		// If PostgreSQL fails, we should still return success since Redis worked
		// In production, you might want to handle this differently
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user in database"})
		return
	}

	response := models.RegisterResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
