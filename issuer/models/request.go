package models

type IssueRequest struct {
	Name            string `json:"name" binding:"required"`
	Lastname        string `json:"lastname" binding:"required"`
	BirthDate       string `json:"birthDate" binding:"required"`
	CountryCode     string `json:"countryCode" binding:"required"`
	CardType        string `json:"cardType" binding:"required"`
	SuscriptorToken string `json:"suscriptorToken" binding:"required"`
	RequestUUID     string `json:"requestUUID" binding:"required"`
}
