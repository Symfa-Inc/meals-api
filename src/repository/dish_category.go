package repository

import (
	"go_api/src/config"
	"go_api/src/models"
)

// CreateDishCategory creates dish category
// returns dish category and error
func CreateDishCategory(category models.DishCategory) (models.DishCategory, error) {
	err := config.DB.Create(&category).Error
	return category, err
}
