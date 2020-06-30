package repository

import (
	"errors"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
	"time"
)

type mealRepo struct{}

func NewMealRepo() *mealRepo {
	return &mealRepo{}
}

// Looks for meal in db, if it doesn't exist returns nil
func (m mealRepo) Find(meal domain.Meal) error {
	result := config.DB.
		Where("catering_id = ? AND date = ?", meal.CateringID, meal.Date).
		Find(&meal)
	if result.RowsAffected != 0 {
		return errors.New("this meal already exist")
	}
	return nil
}

// Create meal entity
// returns new meal item and error
func (m mealRepo) Add(meal domain.Meal) (interface{}, error) {
	mealItem := config.DB.Create(&meal)
	if mealItem.Error != nil {
		return nil, mealItem.Error
	}

	return mealItem.Value, nil
}

// Returns list of meals withing provided date range
// Returns list of meals, total items if and error
func (m mealRepo) Get(limit int, dateQuery types.StartEndDateQuery, id string) ([]domain.Meal, int, error) {
	var meals []domain.Meal
	var total int

	startDate := dateQuery.StartDate
	endDate := dateQuery.EndDate

	if limit == 0 {
		limit = 10
	}

	config.DB.
		Where("catering_id = ? AND date BETWEEN ? and ?", id, startDate, endDate).
		Find(&meals).
		Count(&total)

	err := config.DB.
		Limit(limit).
		Where("catering_id = ? AND date BETWEEN ? and ?", id, startDate, endDate).
		Find(&meals).
		Error
	return meals, total, err
}

// Returns updated meals if exists
func (m mealRepo) Update(path types.PathMeal, meal domain.Meal) (error, int) {
	var mealModel domain.Meal

	t := 24 * time.Hour

	difference := meal.Date.Sub(time.Now().Truncate(t)).Hours()

	if difference < 0 {
		return errors.New("can't add meals to previous dates"), http.StatusBadRequest
	}

	if mealExist := config.DB.
		Where("catering_id = ? AND date = ?", path.ID, meal.Date).
		Find(&mealModel).RecordNotFound(); !mealExist {
		return errors.New("this date already exist"), http.StatusBadRequest
	}

	if resultSecond := config.DB.Model(&mealModel).
		Where("id = ?", path.MealID).
		Update(&meal); resultSecond.RowsAffected == 0 {
		if resultSecond.Error != nil {
			return errors.New(resultSecond.Error.Error()), http.StatusBadRequest
		}
		return errors.New("meal not found"), http.StatusNotFound
	}
	return nil, 0
}
