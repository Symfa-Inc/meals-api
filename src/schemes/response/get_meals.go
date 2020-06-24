package response

import (
	"go_api/src/domain"
)

type GetMealsModel struct {
	Items []domain.Meal `json:"items"`
	Total int           `json:"total"`
} //@name GetMealsResponse
