package response

import (
	"github.com/Aiscom-LLC/meals-api/src/domain"
)

// GetCaterings response scheme
type GetCaterings struct {
	Items []domain.Catering `json:"items"`
	Page  int               `json:"page"`
	Total int               `json:"total"`
} //@name GetCateringsResponse
