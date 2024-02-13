package models

import "keeper/internal/utils"

type SuccessResponse struct {
	Status  bool        `json:"status" swaggertype:"boolean"`
	Message string      `json:"message" swaggertype:"string"`
	Data    interface{} `json:"data,omitempty" swaggertype:"array,object"`
}

type PaginatedSuccessResponse struct {
	Status   bool           `json:"status" swaggertype:"boolean"`
	Message  string         `json:"message" swaggertype:"string"`
	Data     interface{}    `json:"data,omitempty" swaggertype:"array,object"`
	PageInfo utils.PageInfo `json:"page_info,omitempty"`
}

type ErrorResponse struct {
	Status bool   `json:"status" swaggertype:"boolean"`
	Error  string `json:"error" swaggertype:"string"`
}
