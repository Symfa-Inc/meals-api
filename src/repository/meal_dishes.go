package repository

import (
	"go_api/src/config"
	"go_api/src/domain"
)

type mealDishesRepo struct{}

func NewMealDishesRepo() *mealDishesRepo {
	return &mealDishesRepo{}
}

func (md mealDishesRepo) Add(mealDish domain.MealDish) error {
	err := config.DB.
		Create(&mealDish).
		Error
	return err
}
