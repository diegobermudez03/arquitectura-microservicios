package models

type IssueRequest struct {
	Name            string `json:"name"`
	Lastname        string `json:"last_name"`
	BirthDate       string `json:"birth_date"`
	CountryCode     string `json:"country_code"`
	CardType        string `json:"card_type"`
	SuscriptorToken string `json:"suscriptor_token"`
	RequestUUID     string `json:"request_uuid"`
}
