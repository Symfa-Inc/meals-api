package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
	"time"
)

type dishCategoryRepo struct{}

func NewDishCategoryRepo() *dishCategoryRepo {
	return &dishCategoryRepo{}
}

// CreateDishCategory creates dish category
// returns dish category and error
func (dc dishCategoryRepo) Add(category domain.DishCategory) (domain.DishCategory, error) {
	if exist := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND deleted_at >  ?", category.CateringID, category.Name, time.Now()).
		Or("catering_id = ? AND name = ? AND deleted_at IS NULL", category.CateringID, category.Name).
		Find(&category).RecordNotFound(); !exist {

		return domain.DishCategory{}, errors.New("this category already exist")
	}

	err := config.DB.Create(&category).Error
	return category, err
}

// GetDishCategoriesDB returns list of categories of passed catering ID
// returns list of categories and error
func (dc dishCategoryRepo) Get(id string) ([]domain.DishCategory, error, int) {
	var categories []domain.DishCategory
	if cateringRows := config.DB.
		Where("id = ?", id).
		Find(&domain.Catering{}).RowsAffected; cateringRows == 0 {

		return nil, errors.New("catering with that ID is not found"), http.StatusNotFound
	}

	err := config.DB.
		Unscoped().
		Where("catering_id = ? AND (deleted_at > ? OR deleted_at IS NULL)", id, time.Now()).
		Find(&categories).
		Error

	return categories, err, 0
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
	if dishCategoryRows := config.DB.
		Where("catering_id = ? AND id = ?", path.ID, path.CategoryID).
		Delete(&domain.DishCategory{}).RowsAffected; dishCategoryRows == 0 {
		return errors.New("category not found")
	}
	return nil
}

// UpdateDishCategoryDB checks if that name already exists in provided catering
// if its exists throws and error, if not updates the reading
func (dc dishCategoryRepo) Update(path types.PathDishCategory, category domain.DishCategory) (error, int) {
	var categoryModel domain.DishCategory

	if categoryExist := config.DB.
		Where("catering_id = ? AND name = ?", path.ID, category.Name).
		Find(&categoryModel).RecordNotFound(); !categoryExist {

		return errors.New("this category already exist"), http.StatusBadRequest
	}

	if resultSecond := config.DB.Model(&categoryModel).
		Where("id = ?", path.CategoryID).
		Update(&category); resultSecond.RowsAffected == 0 {
		return errors.New("category not found"), http.StatusNotFound
	}
	return nil, 0
}
