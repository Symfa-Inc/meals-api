package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
)

type cateringRepo struct{}

func NewCateringRepo() *cateringRepo {
	return &cateringRepo{}
}

// CreateCatering creates catering in DB
// and error if exists
func (c cateringRepo) Add(catering domain.Catering) error {
	if exist := config.DB.Where("name = ?", catering.Name).
		Find(&catering).RowsAffected; exist != 0 {
		return errors.New("catering with that name already exist")
	}
	return config.DB.Create(&catering).Error
}

// GetCateringsDB returns list of caterings with pagination args
// and error if exists
func (c cateringRepo) Get(query types.PaginationQuery) ([]domain.Catering, int, error) {
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
		Error

	return caterings, total, err
}

// GetCateringByKey returns single catering item found by key
// and error if exists
func (c cateringRepo) GetByKey(key, value string) (domain.Catering, error) {
	var catering domain.Catering
	err := config.DB.Where(key+" = ?", value).First(&catering).Error
	return catering, err
}

// DeleteCateringDB soft delete of catering with passed id
// returns error if exists
func (c cateringRepo) Delete(id string) error {
	if result := config.DB.Where("id = ?", id).
		Delete(&domain.Catering{}).RowsAffected; result == 0 {
		return errors.New("catering not found")
	}

	return nil
}

// UpdateCateringDB updates catering with passed args
// returns updated catering struct and error if exists
func (c cateringRepo) Update(id string, catering domain.Catering) (error, int) {
	if cateringExist := config.DB.Where("id = ?", id).
		Find(&domain.Catering{}).RowsAffected; cateringExist == 0 {
		return errors.New("catering not found"), http.StatusNotFound
	}

	if nameExist := config.DB.Where("name = ?", catering.Name).
		Find(&catering).RowsAffected; nameExist != 0 {
		return errors.New("catering with that name already exist"), http.StatusBadRequest
	}
	return config.DB.Model(&catering).Where("id = ?", id).Update(&catering).Error, 0
}
