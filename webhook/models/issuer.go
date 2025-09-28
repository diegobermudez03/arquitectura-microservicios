package models

import "time"

type DeclineReason struct {
	Reason string `json:"reason"`
}

type IssuedCard struct {
	PAN        string `json:"pan"`
	CVV        string `json:"cvv"`
	ExpiryDate string `json:"expiry_date"`
	CardType   string `json:"card_type"`
}

type IssuerResponse struct {
	DeclineReason   *DeclineReason `json:"decline_reason,omitempty"`
	IssuedCard      *IssuedCard    `json:"issued_card,omitempty"`
	RequestUUID     string         `json:"request_uuid"`
	SuscriptorToken string         `json:"suscriptor_token"`
	Status          string         `json:"status"`
}

// Webhook Event Structure
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Data      IssuerResponse         `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
