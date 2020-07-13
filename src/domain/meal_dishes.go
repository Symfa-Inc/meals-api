package domain

import (
	uuid "github.com/satori/go.uuid"
)

type MealDish struct {
	Base
	MealID uuid.UUID `json:"mealId"`
	DishID uuid.UUID `json:"dishId"`
} //@name MealDishRequest

type GetMealDish struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Weight       string `json:"weight,omitempty"`
	Price        string `json:"price,omitempty"`
	Images       string `json:"images,omitempty"`
	Desc         string `json:"desc,omitempty"`
	CategoryID   string `gorm:"column:category_id" json:"categoryId"`
	CategoryName string `gorm:"column:category_name" json:"categoryName"`
} //@name GetMealResponse

type MealDishRepository interface {
	Add(mealDish MealDish) error
	Delete(mealId string) error
}
