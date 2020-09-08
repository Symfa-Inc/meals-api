package interfaces

import (
	uuid "github.com/satori/go.uuid"
)

// MealDish struct for DB
type MealDish struct {
	Base
	MealID uuid.UUID `json:"mealId"`
	DishID uuid.UUID `json:"dishId"`
} //@name MealDishRequest
