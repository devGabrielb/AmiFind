package dtos

type CustomErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, msg string) *CustomErrorResponse {
	return &CustomErrorResponse{Code: code, Message: msg}
}
