package viewmodels

type ClientRegisterRequest struct {
	//required: true
	//example: Andy
	Account  string `json:"account"`
	//required: true
	//example: pass123
	Password string `json:"password"`
}

type ClientRegisterResponse struct {
	Token string `json:"token"`
}
