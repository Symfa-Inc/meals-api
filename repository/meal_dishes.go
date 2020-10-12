package repository

import (
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
)

// MealDishesRepo struct
type MealDishesRepo struct{}

// NewMealDishesRepo returns pointer to mealDish repository
// with all methods
func NewMealDishesRepo() *MealDishesRepo {
	return &MealDishesRepo{}
}

// Add mealDish, returns err or nil
func (md MealDishesRepo) Add(mealDish domain.MealDish) error {
	err := config.DB.
		Create(&mealDish).
		Error
	return err
}

// Delete soft deletes mealDish, returns err or nil
func (md MealDishesRepo) Delete(mealID string) error {
	if err := config.DB.
		Where("meal_id = ?", mealID).
		Delete(&domain.MealDish{}).
		Error; err != nil {
		return err
	}

	return nil
}
