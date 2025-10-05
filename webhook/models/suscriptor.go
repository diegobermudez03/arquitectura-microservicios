package models

type SuscribeRequest struct {
	Name        string `json:"name"`
	CallbackURL string `json:"callback_url"`
}

type SuscribeResponse struct {
	SuscriptorToken string `json:"suscriptor_token"`
}
