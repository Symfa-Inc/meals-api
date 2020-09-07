package models

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	uuid "github.com/satori/go.uuid"
)

// GetMeal struct response
type GetMeal struct {
	Version string        `json:"version"`
	MealID  uuid.UUID     `json:"mealId"`
	Date    string        `json:"date"`
	Person  string        `json:"person"`
	Result  []domain.Dish `json:"dishes"`
} //@name GetMealsResponse
