package meal

import "go_api/src/models"

type GetMealsModel struct {
	Items []models.Meal `json:"items"`
	Total int           `json:"total"`
} //@name GetMealsResponse
