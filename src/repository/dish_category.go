package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/types"
	"time"
)

// CreateDishCategory creates dish category
// returns dish category and error
func CreateDishCategory(category models.DishCategory) (models.DishCategory, error) {
	result := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND deleted_at >  ?", category.CateringID, category.Name, time.Now()).
		Or("catering_id = ? AND name = ? AND deleted_at IS NULL", category.CateringID, category.Name).
		Find(&category)

	if result.RowsAffected != 0 {
		return models.DishCategory{}, errors.New("this category already exist")
	}

	err := config.DB.Create(&category).Error
	return category, err
}

// CreateDishCategoryTemporal creates a dish which will be deleted on next day
// returns dish category and error
func CreateDishCategoryTemporal(category models.DishCategory) (models.DishCategory, error) {
	trunc := time.Hour * 24
	deleteDate := time.Now().Truncate(trunc).AddDate(0, 0, 1)

	result := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND deleted_at >  ?", category.CateringID, category.Name, time.Now()).
		Find(&category)

	if result.RowsAffected != 0 {
		return models.DishCategory{}, errors.New("this category already exist")
	}

	err := config.DB.Create(&category).Update("deleted_at", deleteDate).Error
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
		Unscoped().
		Where("catering_id = ? AND (deleted_at > ? OR deleted_at IS NULL)", id, time.Now()).
		Find(&categories).
		Error

	return categories, err
}

// GetDishCategoryByKey returns single category item found by key
// and error if exists
func GetDishCategoryByKey(key, value, cateringId string) (models.DishCategory, error) {
	var dishCategory models.DishCategory
	err := config.DB.
		Where("catering_id = ? AND "+key+" = ?", cateringId, value).
		First(&dishCategory).Error
	return dishCategory, err
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
