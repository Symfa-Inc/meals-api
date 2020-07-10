package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/schemes/response"
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
	if exist := config.DB.
		Where("catering_id = ? AND date = ?", meal.CateringID, meal.Date).
		Find(&meal).RecordNotFound(); !exist {
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
func (m mealRepo) Get(mealId, id string) (map[string][]interface{}, error, int) {
	var meal domain.Meal
	var result []response.GetMealsModel
	var categoriesArray []response.GetMealsModel

	if err := config.DB.
		Where("catering_id = ? AND id = ?", id, mealId).
		First(&meal).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return map[string][]interface{}{}, errors.New(err.Error()), http.StatusNotFound
		}
		return map[string][]interface{}{}, errors.New(err.Error()), http.StatusBadRequest
	}

	mealsMap := make(map[string][]interface{})

	err := config.DB.
		Model(&domain.Category{}).
		Select("categories.name as category_name, categories.id as category_id, d.*").
		Joins("left join dishes d on d.category_id = categories.id").
		Joins("left join meal_dishes md on md.dish_id = d.id").
		Joins("left join meals m on m.id = md.meal_id").
		Where("m.id = ?", meal.ID).
		Scan(&result).
		Error

	config.DB.
		Model(&domain.Category{}).
		Select("categories.name as category_name, categories.id as category_id").
		Where("categories.catering_id = ?", id).
		Scan(&categoriesArray)

	for i := range categoriesArray {
		result = append(result, categoriesArray[i])
	}

	for _, dish := range result {
		mealsMap[dish.CategoryName] = append(mealsMap[dish.CategoryName], dish.DishStruct)
	}

	for key, element := range mealsMap {
		element = element[:len(element)-1]
		mealsMap[key] = element
	}

	return mealsMap, err, http.StatusBadRequest
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

func (m mealRepo) GetByKey(key, value string) (domain.Meal, error) {
	var meal domain.Meal
	err := config.DB.
		Where(key+"= ?", value).
		First(&meal).Error
	return meal, err
}
