package viewmodels

type GoogleAuthenticateRequest struct {
	Code string `json:"code"`
}

type GoogleAuthenticateResponse struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	UserID int64 `json:"user_id"`
}

type GetUserInfoResponse struct {
	UserID int64 `json:"user_id"`
	Name string `json:"user_name"`
	Email string `json:"email"`
	Photo string `json:"photo"`
	NumberOfPost int `json:"number_of_post"`
	NumberOfCollection int `json:"number_of_collection"`
}