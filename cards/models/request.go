package models

// IssueRequest represents the request to issue a card
type IssueRequest struct {
	Name            string `json:"name"`
	Lastname        string `json:"last_name"`
	BirthDate       string `json:"birth_date"`
	CountryCode     string `json:"country_code"`
	CardType        string `json:"card_type"`
	SuscriptorToken string `json:"suscriptor_token"`
	RequestUUID     string `json:"request_uuid"`
}

// IssueCardRequest represents the request from frontend to issue a card
type IssueCardRequest struct {
	CardType  string `json:"card_type" binding:"required"`
	UserToken string `json:"user_token" binding:"required"`
}

// RequestData represents the data stored in Redis for a request
type RequestData struct {
	User      User   `json:"user"`
	CardType  string `json:"card_type"`
	UserToken string `json:"user_token"`
}
