package handlers

import (
	"net/http"

	"cards/internal"

	"github.com/gin-gonic/gin"
)

type CardsHandler struct {
	postgresService *internal.PostgresService
}

func NewCardsHandler(postgresService *internal.PostgresService) *CardsHandler {
	return &CardsHandler{
		postgresService: postgresService,
	}
}

// GetCardsByCitizenID handles GET /v1/:citizen_id/cards
func (h *CardsHandler) GetCardsByCitizenID(c *gin.Context) {
	citizenID := c.Param("citizen_id")

	// Validate that citizen_id contains only digits
	if !isDigitsOnly(citizenID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Citizen ID must contain only digits"})
		return
	}

	// Get cards from database
	fullCards, err := h.postgresService.GetCardsByCitizenID(citizenID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Citizen not found or no cards available"})
		return
	}

	c.JSON(http.StatusOK, fullCards)
}
