package viewmodels

// swagger:model
type ResponseData struct {
	ErrorCode    string      `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	Data         interface{} `json:"data"`
}
