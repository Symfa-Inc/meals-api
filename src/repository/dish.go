package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
)

type dishRepo struct{}

func NewDishRepo() *dishRepo {
	return &dishRepo{}
}

func (d dishRepo) Add(cateringId string, dish domain.Dish) (domain.Dish, error) {
	result := config.DB.
		Where("catering_id = ? AND dish_category_id = ? AND name = ?", cateringId, dish.DishCategoryID, dish.Name).
		Find(&dish)
	if result.RowsAffected != 0 {
		return domain.Dish{}, errors.New("this dish already exist in that category")
	}

	err := config.DB.Create(&dish).Error
	if err != nil {
		return domain.Dish{}, err
	}

	return dish, nil
}
