package repository

import (
	"errors"
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
		Where("catering_id = ? AND client_id = ? AND name = ? AND (deleted_at >  ? OR deleted_at IS NULL)",
			category.CateringID, category.ClientID, category.Name, category.DeletedAt).
		Find(category).RecordNotFound(); !exist {

		return errors.New("this category already exist")
	}

	err := config.DB.Create(category).Error
	return err
}

// Get returns list of categories of passed catering ID
// returns list of categories and error
func (dc CategoryRepo) Get(cateringID, clientID, date string) ([]domain.Category, int, error) {
	var categories []domain.Category

	if cateringRows := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).RowsAffected; cateringRows == 0 {

		return nil, http.StatusNotFound, errors.New("catering with that ID is not found")
	}

	err := config.DB.
		Unscoped().
		Where("catering_id = ? AND client_id = ? AND (deleted_at > ? OR deleted_at IS NULL)",
			cateringID, clientID, date).
		Order("created_at").
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
	result := config.DB.
		Unscoped().
		Model(&domain.Category{}).
		Where("catering_id = ? AND id = ? AND client_id = ? AND (deleted_at > ? OR deleted_at IS NULL)", path.ID, path.CategoryID, path.ClientID, time.Now()).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return http.StatusBadRequest, result.Error
	}

	if result.RowsAffected == 0 {
		return http.StatusNotFound, errors.New("category not found")
	}

	return 0, nil
}

// Update checks if that name already exists in provided catering
// if its exists throws and error, if not updates the reading
func (dc CategoryRepo) Update(path types.PathCategory, category *domain.Category) (int, error) {
	if categoryExist := config.DB.
		Where("catering_id = ? AND name = ? AND id = ? AND (deleted_at > ? OR deleted_at IS NULL)",
			path.ID, category.Name, path.CategoryID, time.Now()).
		Find(&category).
		RowsAffected; categoryExist == 0 {
		if nameExist := config.DB.
			Debug().
			Where("catering_id = ? AND client_id = ? AND name = ?", path.ID, path.ClientID, category.Name).
			Find(&category).
			RowsAffected; nameExist != 0 {
			return http.StatusBadRequest, errors.New("category with that name already exist")
		}
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
