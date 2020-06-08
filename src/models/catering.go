package models

import (
	"go_api/src/config"
	"go_api/src/types"
)

// Catering model
type Catering struct {
	Base
	Name string `gorm:"type:varchar(30);unique;not null" json:"name,omitempty" binding:"required"`
}

func CreateCatering(catering Catering) (Catering, error) {
	err := config.DB.Create(&catering).Error
	return catering, err
}

func GetCateringsDB(query types.PaginationQuery) ([]Catering, error) {
	var caterings []Catering
	page := query.Page
	limit := query.Limit
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	err := config.DB.Limit(limit).Offset((page - 1) * limit).Find(&caterings).Error
	return caterings, err
}

func GetCateringByKey(key, value string) (Catering, error) {
	var catering Catering
	err := config.DB.Where(key+" = ?", value).First(&catering).Error
	return catering, err
}

func DeleteCateringDB(id string) error {
	return config.DB.Where("id = ?", id).Delete(&Catering{}).Error
}

func UpdateCateringDB(id string, catering Catering) (Catering, error) {
	err := config.DB.Model(&catering).Where("id = ?", id).Update(&catering).Error
	return catering, err
}