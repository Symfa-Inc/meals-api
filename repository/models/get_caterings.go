package models

import "github.com/Aiscom-LLC/meals-api/domain"

// GetCaterings response scheme
type GetCaterings struct {
	Items []domain.Catering `json:"items"`
	Page  int               `json:"page"`
	Total int               `json:"total"`
} //@name GetCateringsResponse
