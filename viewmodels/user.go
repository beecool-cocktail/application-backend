package viewmodels

type GoogleAuthenticateRequest struct {
	Code string `json:"code"`
}

type GoogleAuthenticateResponse struct {
	Token string `json:"token"`
}