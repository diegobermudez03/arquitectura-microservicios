package models

import "time"

// DeclineReason represents the reason for card decline
type DeclineReason struct {
Reason string `json:"reason"`
}

// IssuedCard represents an issued card
type IssuedCard struct {
PAN        string `json:"pan"`
CVV        string `json:"cvv"`
ExpiryDate string `json:"expiry_date"`
CardType   string `json:"card_type"`
}

// IssuerResponse represents the response from the issuer
type IssuerResponse struct {
DeclineReason   *DeclineReason `json:"decline_reason,omitempty"`
IssuedCard      *IssuedCard    `json:"issued_card,omitempty"`
RequestUUID     string         `json:"request_uuid"`
SuscriptorToken string         `json:"suscriptor_token"`
Status          string         `json:"status"`
}

// WebhookEvent represents the webhook event structure
type WebhookEvent struct {
ID        string                 `json:"id"`
Type      string                 `json:"type"`
Timestamp time.Time              `json:"timestamp"`
Source    string                 `json:"source"`
Data      IssuerResponse         `json:"data"`
Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NotificationRequest represents the request to notifications service
type NotificationRequest struct {
UserToken      string         `json:"user_token"`
IssuerResponse IssuerResponse `json:"issuer_response"`
}
