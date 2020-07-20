package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
	"time"
)

type categoryRepo struct{}

func NewCategoryRepo() *categoryRepo {
	return &categoryRepo{}
}

// CreateCategory creates dish category
// returns dish category and error
func (dc categoryRepo) Add(category domain.Category) (domain.Category, error) {
	if exist := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND deleted_at >  ?", category.CateringID, category.Name, time.Now()).
		Or("catering_id = ? AND name = ? AND deleted_at IS NULL", category.CateringID, category.Name).
		Find(&category).RecordNotFound(); !exist {

		return domain.Category{}, errors.New("this category already exist")
	}

	err := config.DB.Create(&category).Error
	return category, err
}

// GetCategoriesDB returns list of categories of passed catering ID
// returns list of categories and error
func (dc categoryRepo) Get(id string) ([]domain.Category, error, int) {
	var categories []domain.Category
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

// GetCategoryByKey returns single category item found by key
// and error if exists
func (dc categoryRepo) GetByKey(key, value, cateringId string) (domain.Category, error) {
	var category domain.Category
	err := config.DB.
		Where("catering_id = ? AND "+key+" = ?", cateringId, value).
		First(&category).Error
	return category, err
}

// DeleteCategoryDB soft deletes reading from DB
// returns gorm.DB struct with methods
func (dc categoryRepo) Delete(path types.PathCategory) error {
	if categoryRows := config.DB.
		Unscoped().
		Model(&domain.Category{}).
		Where("catering_id = ? AND id = ?  AND (deleted_at > ? OR deleted_at IS NULL)", path.ID, path.CategoryID, time.Now()).
		Update("deleted_at", time.Now()).RowsAffected; categoryRows == 0 {
		return errors.New("category not found")
	}
	return nil
}

// UpdateCategoryDB checks if that name already exists in provided catering
// if its exists throws and error, if not updates the reading
func (dc categoryRepo) Update(path types.PathCategory, category domain.Category) (error, int) {
	var categoryModel domain.Category

	if categoryExist := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND (deleted_at > ? OR deleted_at IS NULL)", path.ID, category.Name, time.Now()).
		Find(&categoryModel).RecordNotFound(); !categoryExist {
		return errors.New("this category already exist"), http.StatusBadRequest
	}

	if resultSecond := config.DB.
		Unscoped().
		Model(&categoryModel).
		Where("id = ? AND (deleted_at > ? OR deleted_at IS NULL)", path.CategoryID, time.Now()).
		Update(&category); resultSecond.RowsAffected == 0 {
		return errors.New("category not found"), http.StatusNotFound
	}
	return nil, 0
}
