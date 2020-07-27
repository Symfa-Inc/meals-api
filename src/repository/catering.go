package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// CateringRepo struct
type CateringRepo struct{}

// NewCateringRepo returns pointer
// to catering repository with all methods
func NewCateringRepo() *CateringRepo {
	return &CateringRepo{}
}

// Add creates catering in DB
// and error if exists
func (c CateringRepo) Add(catering *domain.Catering) error {

	if exist := config.DB.Where("name = ?", catering.Name).
		Find(catering).RowsAffected; exist != 0 {
		return errors.New("catering with that name already exist")
	}

	err := config.DB.Create(catering).Error

	if err != nil {
		return err
	}

	utils.AddDefaultCateringSchedules(catering.ID)

	return err
}

// Get returns list of caterings with pagination args
// and error if exists
func (c CateringRepo) Get(query types.PaginationQuery) ([]domain.Catering, int, error) {
	var caterings []domain.Catering
	var total int

	page := query.Page
	limit := query.Limit

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	config.DB.Find(&caterings).Count(&total)

	err := config.DB.
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&caterings).
		Order("created_at DESC").
		Error

	return caterings, total, err
}

// GetByKey returns single catering item found by key
// and error if exists
func (c CateringRepo) GetByKey(key, value string) (domain.Catering, error) {
	var catering domain.Catering
	err := config.DB.Where(key+" = ?", value).First(&catering).Error
	return catering, err
}

// Delete soft delete of catering with passed id
// returns error if exists
func (c CateringRepo) Delete(id string) error {
	if result := config.DB.Where("id = ?", id).
		Delete(&domain.Catering{}).RowsAffected; result == 0 {
		return errors.New("catering not found")
	}

	return nil
}

// Update updates catering with passed args
// returns updated catering struct and error if exists
func (c CateringRepo) Update(id string, catering domain.Catering) (int, error) {
	if nameExist := config.DB.Where("name = ?", catering.Name).
		Find(&catering).RowsAffected; nameExist != 0 {
		return http.StatusBadRequest, errors.New("catering with that name already exist")
	}

	if cateringExist := config.DB.Model(&catering).Where("id = ?", id).
		Update(&catering).RowsAffected; cateringExist == 0 {
		return http.StatusNotFound, errors.New("catering not found")
	}

	return 0, nil
}
