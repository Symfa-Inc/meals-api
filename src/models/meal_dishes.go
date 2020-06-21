package models

import uuid "github.com/satori/go.uuid"

type MealDish struct {
	Base
	MealID uuid.UUID `json:"mealId"`
	DishID uuid.UUID `json:"dishId"`
} //@name MealDishResponse
