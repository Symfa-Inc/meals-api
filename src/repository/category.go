package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
	"time"
)

// CategoryRepo struct
type CategoryRepo struct{}

// NewCategoryRepo returns pointer to
// category repository with all methods
func NewCategoryRepo() *CategoryRepo {
	return &CategoryRepo{}
}

// Add creates dish category
// returns dish category and error
func (dc CategoryRepo) Add(category *domain.Category) error {
	if exist := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND deleted_at >  ?", category.CateringID, category.Name, time.Now()).
		Or("catering_id = ? AND name = ? AND deleted_at IS NULL", category.CateringID, category.Name).
		Find(category).RecordNotFound(); !exist {

		return errors.New("this category already exist")
	}

	err := config.DB.Create(category).Error
	return err
}

// Get returns list of categories of passed catering ID
// returns list of categories and error
func (dc CategoryRepo) Get(id string) ([]domain.Category, int, error) {
	var categories []domain.Category

	if cateringRows := config.DB.
		Where("id = ?", id).
		Find(&domain.Catering{}).RowsAffected; cateringRows == 0 {

		return nil, http.StatusNotFound, errors.New("catering with that ID is not found")
	}

	err := config.DB.
		Unscoped().
		Where("catering_id = ? AND (deleted_at > ? OR deleted_at IS NULL)", id, time.Now()).
		Find(&categories).
		Error

	return categories, 0, err
}

// GetByKey returns single category item found by key
// and error if exists
func (dc CategoryRepo) GetByKey(key, value, cateringID string) (domain.Category, error) {
	var category domain.Category
	err := config.DB.
		Where("catering_id = ? AND "+key+" = ?", cateringID, value).
		First(&category).Error
	return category, err
}

// Delete soft deletes reading from DB
// returns gorm.DB struct with methods
func (dc CategoryRepo) Delete(path types.PathCategory) (int, error) {
	if err := config.DB.
		Unscoped().
		Model(&domain.Category{}).
		Where("catering_id = ? AND id = ?  AND (deleted_at > ? OR deleted_at IS NULL)", path.ID, path.CategoryID, time.Now()).
		Update("deleted_at", time.Now()).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, errors.New("category not found")
		}
		return http.StatusBadRequest, err
	}

	return 0, nil
}

// Update checks if that name already exists in provided catering
// if its exists throws and error, if not updates the reading
func (dc CategoryRepo) Update(path types.PathCategory, category *domain.Category) (int, error) {
	if categoryRows := config.DB.
		Unscoped().
		Where("catering_id = ? AND name = ? AND (deleted_at > ? OR deleted_at IS NULL)", path.ID, category.Name, time.Now()).
		Find(&domain.Category{}).RowsAffected; categoryRows != 0 {
		return http.StatusBadRequest, errors.New("category with that name already exist")
	}

	if categoryNotExist := config.DB.
		Unscoped().
		Model(&domain.Category{}).
		Where("id = ? AND (deleted_at > ? OR deleted_at IS NULL)", path.CategoryID, time.Now()).
		Update(category); categoryNotExist.RowsAffected == 0 {
		return http.StatusNotFound, errors.New("category not found")
	}

	return 0, nil
}
