package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
)

type cateringRepo struct{}

func NewCateringRepo() *cateringRepo {
	return &cateringRepo{}
}

// CreateCatering creates catering in DB
// and error if exists
func (c cateringRepo) Add(catering domain.Catering) (domain.Catering, error) {
	err := config.DB.Create(&catering).Error
	return catering, err
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
	result := config.DB.Where("id = ?", id).Delete(&domain.Catering{})

	if result.RowsAffected == 0 {
		return errors.New("catering not found")
	}

	return nil
}

// UpdateCateringDB updates catering with passed args
// returns updated catering struct and error if exists
func (c cateringRepo) Update(id string, catering domain.Catering) error {
	result := config.DB.Model(&catering).Where("id = ?", id).Update(&catering)

	if result.RowsAffected == 0 {
		if result.Error != nil {
			return errors.New(result.Error.Error())
		}
		return errors.New("catering not found")
	}

	return nil
}
