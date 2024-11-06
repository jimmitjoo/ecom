package handlers

// ErrorResponse representerar ett API-fel
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Ogiltig förfrågan"`
}

// SuccessResponse representerar ett lyckat API-svar
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
}
