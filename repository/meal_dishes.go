package repository

import (
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	uuid "github.com/satori/go.uuid"
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
		Where("id = ?", mealID).
		Delete(&domain.MealDish{}).
		Error; err != nil {
		return err
	}

	return nil
}

// FindByID return list of meal
func (md MealDishesRepo) FindByID(id uuid.UUID) ([]domain.MealDish, error) {
	var mealDish []domain.MealDish
	if err := config.DB.
		Select("*").
		Where("meal_id = ?", id).
		Error; err != nil {
		return []domain.MealDish{}, err
	}

	config.DB.
		Unscoped().
		Model(&domain.MealDish{}).
		Where("meal_id = ?", id).
		Find(&mealDish)

	return mealDish, nil
}

// Update updates dish for current meal
func (md MealDishesRepo) Update(meal domain.MealDish) error {
	if err := config.DB.
		Unscoped().
		Model(&domain.MealDish{}).
		Where("id = ?", meal.ID).
		Update(meal).Error; err != nil {
		return err
	}

	return nil
}
