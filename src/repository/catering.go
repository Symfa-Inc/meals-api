package repository

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
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
func (c CateringRepo) Get(cateringID string, query types.PaginationQuery) ([]domain.Catering, int, error) {
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

	if cateringID != "" {
		config.DB.
			Find(&caterings).
			Where("id = ?", cateringID).
			Count(&total)

		err := config.DB.
			Limit(limit).
			Offset((page-1)*limit).
			Where("id = ?", cateringID).
			Order("created_at DESC").
			Find(&caterings).
			Error
		return caterings, total, err
	}
	config.DB.
		Find(&caterings).
		Count(&total)

	err := config.DB.
		Limit(limit).
		Offset((page - 1) * limit).
		Order("created_at DESC").
		Find(&caterings).
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
	var cateringUsers []domain.CateringUser
	if cateringExist := config.DB.
		Where("id = ?", id).
		Delete(&domain.Catering{}).
		RowsAffected; cateringExist == 0 {
		return errors.New("catering not found")
	}

	config.DB.
		Where("catering_id = ?", id).
		Find(&cateringUsers)

	for i := range cateringUsers {
		config.DB.
			Model(&domain.User{}).
			Where("id = ?", cateringUsers[i].UserID).
			Update(map[string]interface{}{
				"status":     types.StatusTypesEnum.Deleted,
				"deleted_at": time.Now(),
			})
	}
	return nil
}

// Update updates catering with passed args
// returns updated catering struct and error if exists
func (c CateringRepo) Update(id string, catering domain.Catering) (int, error) {
	if cateringExist := config.DB.
		Where("name = ? AND id = ?", catering.Name, id).
		Find(&catering).
		RowsAffected; cateringExist == 0 {
		if nameExist := config.DB.
			Where("name = ?", catering.Name).
			Find(&catering).
			RowsAffected; nameExist != 0 {
			return http.StatusBadRequest, errors.New("catering with that name already exist")
		}
	}

	if cateringExist := config.DB.Model(&catering).Where("id = ?", id).
		Update(&catering).RowsAffected; cateringExist == 0 {
		return http.StatusNotFound, errors.New("catering not found")
	}

	return 0, nil
}
