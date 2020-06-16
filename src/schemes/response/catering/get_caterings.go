package catering

import "go_api/src/models"

type GetCaterings struct {
	Items []models.Catering `json:"items"`
	Page  int               `json:"page"`
	Total int               `json:"total"`
} //@name GetCateringsResponse
