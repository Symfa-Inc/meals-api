package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
)

type dishRepo struct{}

func NewDishRepo() *dishRepo {
	return &dishRepo{}
}

// Creates new dish entity
// returns error or nil
func (d dishRepo) Add(cateringId string, dish domain.Dish) error {
	var total int
	config.DB.
		Model(&domain.Dish{}).
		Where("catering_id = ? AND dish_category_id = ?", cateringId, dish.DishCategoryID).
		Count(&total)

	if total >= 10 {
		return errors.New("can't add more than 10 dishes for a single category")
	}

	if dishExist := config.DB.
		Where("catering_id = ? AND dish_category_id = ? AND name = ?", cateringId, dish.DishCategoryID, dish.Name).
		Find(&dish).
		RecordNotFound(); !dishExist {
		return errors.New("this dish already exist in that category")
	}

	if err := config.DB.Create(&dish).Error; err != nil {
		return err
	}
	return nil
}

// Soft delete of entity
// returns error or nil
func (d dishRepo) Delete(path types.PathDish) error {
	if cateringNotExist := config.DB.Where("id = ?", path.CateringID).
		Find(&domain.Catering{}).RecordNotFound(); cateringNotExist {
		return errors.New("catering with that ID doesn't exist")
	}

	if rows := config.DB.Where("catering_id = ? AND id = ?", path.CateringID, path.DishID).
		Delete(&domain.Dish{}).RowsAffected; rows == 0 {
		return errors.New("dish not found")
	}

	return nil
}

// Get entity filtered by key and value
// returns entity and error or nil
func (d dishRepo) GetByKey(key, value, cateringId, categoryId string) (domain.Dish, error) {
	var dish domain.Dish
	err := config.DB.
		Where("catering_id = ? and dish_category_id = ? AND "+key+" = ?", cateringId, categoryId, value).
		First(&dish).Error
	return dish, err
}

// Get list of dishes
// returns array of dishes and error or nil and status code
func (d dishRepo) Get(cateringId, categoryId string) ([]domain.Dish, error, int) {
	var dishes []domain.Dish

	if cateringNotExist := config.DB.
		Where("id = ?", cateringId).
		Find(&domain.Catering{}).
		RecordNotFound(); cateringNotExist {
		return nil, errors.New("catering with that ID doesn't exist"), http.StatusNotFound
	}

	if categoryNotExist := config.DB.
		Where("id = ?", categoryId).
		Find(&domain.DishCategory{}).
		RecordNotFound(); categoryNotExist {
		return nil, errors.New("category with that ID doesn't exist"), http.StatusNotFound
	}

	err := config.DB.
		Where("catering_id = ? AND dish_category_id = ?", cateringId, categoryId).
		Find(&dishes).
		Error
	return dishes, err, 0
}

// Updates entity
// returns error or nil and status code
func (d dishRepo) Update(path types.PathDish, dish domain.Dish) (error, int) {
	if cateringNotExist := config.DB.
		Where("id = ?", path.CateringID).
		Find(&domain.Catering{}).
		RecordNotFound(); cateringNotExist {
		return errors.New("catering not found"), http.StatusNotFound
	}

	if categoryNotExist := config.DB.
		Where("id = ?", dish.DishCategoryID).
		Find(&domain.DishCategory{}).
		RecordNotFound(); categoryNotExist {
		return errors.New("dish category not found"), http.StatusNotFound
	}

	if result := config.DB.Model(&dish).
		Where("id = ? AND dish_category_id = ?", path.DishID, dish.DishCategoryID).
		Update(&dish).RowsAffected; result == 0 {
		return errors.New("dish not found"), http.StatusNotFound
	}

	return nil, 0
}
