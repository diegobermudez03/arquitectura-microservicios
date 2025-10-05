package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"issuer/models"

	"github.com/gin-gonic/gin"
)

var countryMinAge = map[string]int{
	"US": 18,
	"CO": 14,
	"MX": 16,
	"CA" : 20,
}

var cardTypes = map[string]bool{
	"debit":   true,
	"credit":  true,
	"prepaid": true,
}

type Handlers struct {
	webhookURL string
}

func NewHandlers(webhookURL string) *Handlers {
	return &Handlers{webhookURL: webhookURL}
}

func (h *Handlers) IssueCard(c *gin.Context) {
	var req models.IssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received issue  request for: %v", req)

	// Validate country code
	minAge, exists := countryMinAge[req.CountryCode]
	if !exists {
		log.Printf("Country not eligible: %s", req.CountryCode)
		// Send webhook asynchronously for decline
		go h.processCardIssueAsync(req, &models.DeclineReason{Reason: "Country not eligible"})
		c.JSON(http.StatusOK, gin.H{"status": "request_received", "message": "Request is being processed"})
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
		// Send webhook asynchronously for decline
		go h.processCardIssueAsync(req, &models.DeclineReason{Reason: "User not eligible due to age"})
		c.JSON(http.StatusOK, gin.H{"status": "request_received", "message": "Request is being processed"})
		return
	}

	// Validate card type
	if _, ok := cardTypes[req.CardType]; !ok {
		log.Printf("Card type not eligible: %s", req.CardType)
		// Send webhook asynchronously for decline
		go h.processCardIssueAsync(req, &models.DeclineReason{Reason: "Card type not eligible"})
		c.JSON(http.StatusOK, gin.H{"status": "request_received", "message": "Request is being processed"})
		return
	}

	// Start async processing for successful case
	go h.processCardIssueAsync(req, nil)

	// Return immediately
	c.JSON(http.StatusOK, gin.H{"status": "request_received", "message": "Request is being processed"})
}

func (h *Handlers) processCardIssueAsync(req models.IssueRequest, declineReason *models.DeclineReason) {
	log.Printf("Starting async processing for request: %s", req.RequestUUID)
	// Simulate processing time (6 seconds)
	time.Sleep(6 * time.Second)

	// If we already have a decline reason, send it immediately
	if declineReason != nil {
		log.Printf("Sending decline webhook for request: %s", req.RequestUUID)
		h.sendWebhookResponse(req, declineReason, nil)
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
		CardType:   req.CardType,
	}

	log.Printf("Card generated successfully for %s %s", req.Name, req.Lastname)
	// Send webhook response
	log.Printf("Sending success webhook for request: %s", req.RequestUUID)
	h.sendWebhookResponse(req, nil, issuedCard)
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

func (h *Handlers) sendWebhookResponse(req models.IssueRequest, declineReason *models.DeclineReason, issuedCard *models.IssuedCard) {
	if h.webhookURL == "" {
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

	log.Printf("Sending webhook to: %s", h.webhookURL)
	resp, err := http.Post(h.webhookURL, "application/json", bytes.NewBuffer(jsonData))
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
