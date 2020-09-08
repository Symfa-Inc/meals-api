package domain

import (
	uuid "github.com/satori/go.uuid"
)

// MealDish struct for DB
type MealDish struct {
	Base
	MealID uuid.UUID `json:"mealId"`
	DishID uuid.UUID `json:"dishId"`
} //@name MealDishRequest

// MealDishRepository is mealDish interface for repository
type MealDishRepository interface {
	Add(mealDish MealDish) error
	Delete(mealID string) error
}
