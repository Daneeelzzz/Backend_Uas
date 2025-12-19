package helper

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"` // Tambahan Meta
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalData   int `json:"total_data"`
	Limit       int `json:"limit"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}