package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"time"
)

type dishCategoryRepo struct{}

func NewDishCategoryRepo() *dishCategoryRepo {
	return &dishCategoryRepo{}
}

// CreateDishCategory creates dish category
// returns dish category and error
func (dc dishCategoryRepo) Add(category domain.DishCategory) (domain.DishCategory, error) {
	result := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND deleted_at >  ?", category.CateringID, category.Name, time.Now()).
		Or("catering_id = ? AND name = ? AND deleted_at IS NULL", category.CateringID, category.Name).
		Find(&category)

	if result.RowsAffected != 0 {
		return domain.DishCategory{}, errors.New("this category already exist")
	}

	err := config.DB.Create(&category).Error
	return category, err
}

// GetDishCategoriesDB returns list of categories of passed catering ID
// returns list of categories and error
func (dc dishCategoryRepo) Get(id string) ([]domain.DishCategory, error) {
	var categories []domain.DishCategory
	catering := config.DB.
		Where("id = ?", id).
		Find(&domain.Catering{})

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
func (dc dishCategoryRepo) GetByKey(key, value, cateringId string) (domain.DishCategory, error) {
	var dishCategory domain.DishCategory
	err := config.DB.
		Where("catering_id = ? AND "+key+" = ?", cateringId, value).
		First(&dishCategory).Error
	return dishCategory, err
}

// DeleteDishCategoryDB soft deletes reading from DB
// returns gorm.DB struct with methods
func (dc dishCategoryRepo) Delete(path types.PathDishCategory) error {
	result := config.DB.
		Where("catering_id = ? AND id = ?", path.ID, path.CategoryID).
		Delete(&domain.DishCategory{})
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

// UpdateDishCategoryDB checks if that name already exists in provided catering
// if its exists throws and error, if not updates the reading
func (dc dishCategoryRepo) Update(path types.PathDishCategory, category domain.DishCategory) error {
	var categoryModel domain.DishCategory

	result := config.DB.
		Where("catering_id = ? AND name = ?", path.ID, category.Name).
		Find(&categoryModel)

	if result.RowsAffected != 0 {
		return errors.New("this category already exist")
	}

	resultSecond := config.DB.Model(&categoryModel).Where("id = ?", path.CategoryID).Update(&category)
	if resultSecond.RowsAffected == 0 {
		if resultSecond.Error != nil {
			return errors.New(resultSecond.Error.Error())
		}
		return errors.New("category not found")
	}
	return nil
}
