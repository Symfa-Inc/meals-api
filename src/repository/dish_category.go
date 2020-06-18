package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/types"
)

// CreateDishCategory creates dish category
// returns dish category and error
func CreateDishCategory(category models.DishCategory) (models.DishCategory, error) {
	result := config.DB.
		Where("catering_id = ? AND name = ?", category.CateringID, category.Name).
		Find(&category)

	if result.RowsAffected != 0 {
		return models.DishCategory{}, errors.New("this category already exist")
	}

	err := config.DB.Create(&category).Error
	return category, err
}

// GetDishCategoriesDB returns list of categories of passed catering ID
// returns list of categories and error
func GetDishCategoriesDB(id string) ([]models.DishCategory, error) {
	var categories []models.DishCategory
	catering := config.DB.
		Where("id = ?", id).
		Find(&models.Catering{})

	if catering.RowsAffected == 0 {
		return nil, errors.New("catering with that ID is not found")
	}

	err := config.DB.
		Where("catering_id = ?", id).
		Find(&categories).
		Error

	return categories, err
}

// DeleteDishCategoryDB soft deletes reading from DB
// returns gorm.DB struct with methods
func DeleteDishCategoryDB(path types.PathDishCategory) *gorm.DB {
	return config.DB.
		Where("catering_id = ? AND id = ?", path.ID, path.CategoryID).
		Delete(&models.DishCategory{})
}

// UpdateDishCategoryDB checks if that name already exists in provided catering
// if its exists throws and error, if not updates the reading
func UpdateDishCategoryDB(path types.PathDishCategory, category models.DishCategory) (*gorm.DB, error) {
	var categoryModel models.DishCategory

	result := config.DB.
		Where("catering_id = ? AND name = ?", path.ID, category.Name).
		Find(&categoryModel)

	if result.RowsAffected != 0 {
		return nil, errors.New("this category already exist")
	}

	return config.DB.Model(&categoryModel).Where("id = ?", path.CategoryID).Update(&category), nil
}
