package models

type SuscribeRequest struct {
	Name        string `json:"name"`
	CallbackURL string `json:"callback_url"`
}

type SuscribeResponse struct {
	SuscriptorToken string `json:"suscriptor_token"`
}

type ResponseRequest struct {
	RequestUUID     string `json:"request_uuid"`
	SuscriptorToken string `json:"suscriptor_token"`
}
