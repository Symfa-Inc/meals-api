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

func (md mealDishesRepo) Delete(mealId string) error {
	if err := config.DB.
		Where("meal_id = ?", mealId).
		Delete(&domain.MealDish{}).
		Error; err != nil {
		return err
	}

	return nil
}
