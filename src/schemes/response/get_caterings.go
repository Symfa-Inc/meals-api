package response

import (
	"go_api/src/domain"
)

// GetCaterings response scheme
type GetCaterings struct {
	Items []domain.Catering `json:"items"`
	Page  int               `json:"page"`
	Total int               `json:"total"`
} //@name GetCateringsResponse
