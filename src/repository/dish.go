package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/models"
)

func CreateDish(cateringId string, dish models.Dish) (models.Dish, error) {
	result := config.DB.
		Where("catering_id = ? AND dish_category_id = ? AND name = ?", cateringId, dish.DishCategoryID, dish.Name).
		Find(&dish)
	if result.RowsAffected != 0 {
		return models.Dish{}, errors.New("this dish already exist in that category")
	}

	err := config.DB.Create(&dish).Error
	if err != nil {
		return models.Dish{}, err
	}

	return dish, nil
}
