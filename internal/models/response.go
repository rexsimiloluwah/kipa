package models

type SuccessResponse struct {
	Status  bool        `json:"status" swaggertype:"boolean"`
	Message string      `json:"message" swaggertype:"string"`
	Data    interface{} `json:"data,omitempty" swaggertype:"primitive,string"`
}

type ErrorResponse struct {
	Status bool   `json:"status" swaggertype:"boolean"`
	Error  string `json:"error" swaggertype:"string"`
}
