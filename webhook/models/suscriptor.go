package models

type SuscribeRequest struct {
	Name        string `json:"name" binding:"required"`
	CallbackURL string `json:"callback_url" binding:"required"`
}

type SuscribeResponse struct {
	SuscriptorToken string `json:"suscriptor_token"`
}

type ResponseRequest struct {
	RequestUUID     string `json:"request_uuid"`
	SuscriptorToken string `json:"suscriptor_token"`
}
