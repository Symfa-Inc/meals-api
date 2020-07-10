package domain

import (
	uuid "github.com/satori/go.uuid"
)

type MealDish struct {
	Base
	MealID uuid.UUID `json:"mealId"`
	DishID uuid.UUID `json:"dishId"`
} //@name MealDishRequest

type MealDishRepository interface {
	Add(mealDish MealDish) error
}
