package meal

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/models"
	"go_api/src/schemes/request/meal"
	"go_api/src/types"
	"time"
)

// Looks for meal in db, if it doesn't exist returns nil
func FindMealDB(meal models.Meal) error {
	result := config.DB.
		Debug().
		Where("catering_id = ? AND date = ?", meal.CateringID, meal.Date).
		Find(&meal)
	if result.RowsAffected != 0 {
		return errors.New("this meal already exist")
	}
	return nil
}

// Create meal entity
// returns new meal item and error
func CreateMealDB(meal models.Meal) (*gorm.DB, error) {
	mealItem := config.DB.Create(&meal)
	if mealItem.Error != nil {
		return nil, mealItem.Error
	}

	return mealItem, nil
}

// Returns list of meals withing provided date range
// Returns list of meals, total items if and error
func GetMealsDB(limit int, dateQuery types.StartEndDateQuery, id string) ([]models.Meal, int, error) {
	var meals []models.Meal
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
func UpdateMealDB(id string, meal meal.UpdateMealRequest) (*gorm.DB, error) {
	var mealModel models.Meal

	t := 24 * time.Hour

	difference := meal.Date.Sub(time.Now().Truncate(t)).Hours()

	if difference < 0 {
		return nil, errors.New("can't add meals to previous dates")
	}

	result := config.DB.
		Debug().
		Where("catering_id = ? AND date = ?", meal.CateringID, meal.Date).
		Find(&mealModel)

	if result.RowsAffected != 0 {
		return nil, errors.New("this date already exist")
	}

	return config.DB.Debug().Model(&mealModel).Where("id = ?", id).Update(&meal), nil
}
