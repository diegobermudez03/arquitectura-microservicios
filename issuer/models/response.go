package models

type DeclineReason struct {
	Reason string `json:"reason"`
}

type IssuedCard struct {
	PAN        string `json:"pan"`
	CVV        string `json:"cvv"`
	ExpiryDate string `json:"expiry_date"`
}

type WebhookResponse struct {
	DeclineReason   *DeclineReason `json:"decline_reason,omitempty"`
	IssuedCard      *IssuedCard    `json:"issued_card,omitempty"`
	RequestUUID     string         `json:"request_uuid"`
	SuscriptorToken string         `json:"suscriptor_token"`
	Status          string         `json:"status"`
}
