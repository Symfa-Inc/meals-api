package repository

import (
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/types"
)

// CreateCatering creates catering in DB
// and error if exists
func CreateCatering(catering models.Catering) (models.Catering, error) {
	err := config.DB.Create(&catering).Error
	return catering, err
}

// GetCateringsDB returns list of caterings with pagination args
// and error if exists
func GetCateringsDB(query types.PaginationQuery) ([]models.Catering, int, error) {
	var caterings []models.Catering
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
func GetCateringByKey(key, value string) (models.Catering, error) {
	var catering models.Catering
	err := config.DB.Where(key+" = ?", value).First(&catering).Error
	return catering, err
}

// DeleteCateringDB soft delete of catering with passed id
// returns error if exists
func DeleteCateringDB(id string) *gorm.DB {
	return config.DB.Where("id = ?", id).Delete(&models.Catering{})
}

// UpdateCateringDB updates catering with passed args
// returns updated catering struct and error if exists
func UpdateCateringDB(id string, catering models.Catering) *gorm.DB {
	return config.DB.Model(&catering).Where("id = ?", id).Update(&catering)
}
