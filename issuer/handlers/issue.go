package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"issuer/models"

	"github.com/gin-gonic/gin"
)

var countryMinAge = map[string]int{
	"US": 18,
	"CO": 14,
	"DE": 16,
}

func IssueCard(c *gin.Context) {
	var req models.IssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing card issue request for: %s %s", req.Name, req.Lastname)

	// Validate country code
	minAge, exists := countryMinAge[req.CountryCode]
	if !exists {
		log.Printf("Country not eligible: %s", req.CountryCode)
		sendWebhookResponse(req, &models.DeclineReason{Reason: "Country not eligible"}, nil)
		c.JSON(http.StatusOK, gin.H{"status": "declined", "reason": "Country not eligible"})
		return
	}

	// Parse birth date and calculate age
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		log.Printf("Error parsing birth date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format"})
		return
	}

	age := int(time.Since(birthDate).Hours() / 24 / 365)
	if age < minAge {
		log.Printf("User not eligible due to age: %d < %d", age, minAge)
		sendWebhookResponse(req, &models.DeclineReason{Reason: "User not eligible due to age"}, nil)
		c.JSON(http.StatusOK, gin.H{"status": "declined", "reason": "User not eligible due to age"})
		return
	}

	// Generate card details
	pan := generatePAN()
	cvv := generateCVV()
	expiryDate := generateExpiryDate()

	issuedCard := &models.IssuedCard{
		PAN:        pan,
		CVV:        cvv,
		ExpiryDate: expiryDate,
	}

	log.Printf("Card generated successfully for %s %s", req.Name, req.Lastname)

	// Simulate processing time
	time.Sleep(3 * time.Second)

	// Send webhook response
	sendWebhookResponse(req, nil, issuedCard)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"card":   issuedCard,
	})
}

func generatePAN() string {
	// Generate 16-digit PAN starting with "4242"
	rand.Seed(time.Now().UnixNano())
	pan := "4242"
	for i := 0; i < 12; i++ {
		pan += strconv.Itoa(rand.Intn(10))
	}
	return pan
}

func generateCVV() string {
	rand.Seed(time.Now().UnixNano())
	cvv := ""
	for i := 0; i < 3; i++ {
		cvv += strconv.Itoa(rand.Intn(10))
	}
	return cvv
}

func generateExpiryDate() string {
	now := time.Now()
	// 6 years from now plus up to 30 random days
	years := 6
	randomDays := rand.Intn(31) // 0-30 days
	expiry := now.AddDate(years, 0, randomDays)
	return expiry.Format("2006-01-02")
}

func sendWebhookResponse(req models.IssueRequest, declineReason *models.DeclineReason, issuedCard *models.IssuedCard) {
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Printf("WEBHOOK_URL not set, skipping webhook call")
		return
	}

	response := models.WebhookResponse{
		DeclineReason:   declineReason,
		IssuedCard:      issuedCard,
		RequestUUID:     req.RequestUUID,
		SuscriptorToken: req.SuscriptorToken,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling webhook response: %v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending webhook: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading webhook response: %v", err)
		return
	}

	log.Printf("Webhook sent successfully. Response: %s", string(body))
}
