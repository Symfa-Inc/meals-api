package response

import (
	"go_api/src/domain"
)

type GetCaterings struct {
	Items []domain.Catering `json:"items"`
	Page  int               `json:"page"`
	Total int               `json:"total"`
} //@name GetCateringsResponse
