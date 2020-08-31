package repository

import (
	"errors"
	"net/http"
	"time"

	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"github.com/Aiscom-LLC/meals-api/src/schemes/response"

	"github.com/jinzhu/gorm"
)

// MealRepo struct
type MealRepo struct{}

// NewMealRepo returns pointer to meal repository
// with all methods
func NewMealRepo() *MealRepo {
	return &MealRepo{}
}

// Find looks for meal in db, if it doesn't exist returns nil
func (m MealRepo) Find(meal *domain.Meal) error {
	if exist := config.DB.
		Where("catering_id = ? AND date = ?", meal.CateringID, meal.Date).
		Find(meal).RecordNotFound(); !exist {
		return errors.New("this meal already exist")
	}
	return nil
}

// Add create meal entity
// returns new meal item and error
func (m MealRepo) Add(meal *domain.Meal) error {
	if err := config.DB.Create(meal).Error; err != nil {
		return err
	}

	return nil
}

// Get returns list of meals, total items if and error
func (m MealRepo) Get(Date time.Time, id, clientID string) ([]response.GetMeal, int, error) {
	var meals []domain.Meal
	var mealsResponse []response.GetMeal

	if err := config.DB.
		Where("catering_id = ? AND client_id = ? AND date > ?", id, clientID, Date).
		Order("created_at").
		Find(&meals).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []response.GetMeal{}, http.StatusNotFound, err
		}
		return []response.GetMeal{}, http.StatusNotFound, err
	}

	for _, meal := range meals {
		var result []domain.Dish

		if err := config.DB.
			Unscoped().
			Model(&domain.Category{}).
			Select("categories.id as category_id, categories.deleted_at, d.*").
			Joins("left join dishes d on d.category_id = categories.id").
			Joins("left join meal_dishes md on md.dish_id = d.id").
			Joins("left join meals m on m.id = md.meal_id").
			Where("m.id = ? AND md.deleted_at IS NULL AND (categories.deleted_at > ? OR categories.deleted_at IS NULL)", meal.ID, Date).
			Order("d.created_at").
			Scan(&result).
			Error; err != nil {
			return []response.GetMeal{}, http.StatusNotFound, err
		}

		for i := range result {
			var imagesArray []domain.ImageArray
			config.DB.
				Model(&domain.Image{}).
				Select("images.path, images.id").
				Joins("left join image_dishes id on id.image_id = images.id").
				Joins("left join dishes d on id.dish_id = d.id").
				Where("d.id = ? AND id.deleted_at IS NULL", result[i].ID).
				Scan(&imagesArray)
			result[i].Images = imagesArray
		}

		mealDishes := response.GetMeal{
			MealID:  meal.MealID,
			Version: meal.Version,
			Person:  meal.Person,
			Date:    meal.CreatedAt.Format(time.RFC3339),
			Result:  result,
		}

		mealsResponse = append([]response.GetMeal{mealDishes}, mealsResponse...)
	}

	if mealsResponse == nil {
		//return nil, http.StatusNotFound, errors.New("meal for current day not found")
		mealsResponse = make([]response.GetMeal, 0)
	}

	return mealsResponse, 0, nil
}

// GetByKey get meal by provided key value arguments
// Returns meal, error and status code
func (m MealRepo) GetByKey(key, value string) (domain.Meal, int, error) {
	var meal domain.Meal

	if err := config.DB.
		Where(key+" = ?", value).
		First(&meal).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Meal{}, http.StatusNotFound, errors.New("meal with this date not found")
		}
		return domain.Meal{}, http.StatusBadRequest, err
	}

	return meal, 0, nil
}

// Get returns list of meals, total items if and error
func (m MealRepo) GetByRange(startDate time.Time, endDate time.Time, id, clientID string) ([]response.GetMeal, int, error) {
	var meals []domain.Meal
	var mealsResponse []response.GetMeal

	if err := config.DB.
		Where("catering_id = ? AND client_id = ? AND date >= ? AND date <= ?", id, clientID, startDate.UTC(), endDate.UTC()).
		Order("created_at").
		Find(&meals).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []response.GetMeal{}, http.StatusNotFound, err
		}
		return []response.GetMeal{}, http.StatusNotFound, err
	}

	for _, meal := range meals {
		var result []domain.Dish

		if err := config.DB.
			Unscoped().
			Model(&domain.Category{}).
			Select("categories.id as category_id, categories.deleted_at, d.*").
			Joins("left join dishes d on d.category_id = categories.id").
			Joins("left join meal_dishes md on md.dish_id = d.id").
			Joins("left join meals m on m.id = md.meal_id").
			Where("m.id = ? AND md.deleted_at IS NULL AND (categories.deleted_at > ? OR categories.deleted_at IS NULL)", meal.ID, startDate).
			Order("d.created_at").
			Scan(&result).
			Error; err != nil {
			return []response.GetMeal{}, http.StatusNotFound, err
		}

		for i := range result {
			var imagesArray []domain.ImageArray
			config.DB.
				Model(&domain.Image{}).
				Select("images.path, images.id").
				Joins("left join image_dishes id on id.image_id = images.id").
				Joins("left join dishes d on id.dish_id = d.id").
				Where("d.id = ? AND id.deleted_at IS NULL", result[i].ID).
				Scan(&imagesArray)
			result[i].Images = imagesArray
		}

		mealDishes := response.GetMeal{
			MealID:  meal.MealID,
			Version: meal.Version,
			Person:  meal.Person,
			Date:    meal.CreatedAt.Format(time.RFC3339),
			Result:  result,
		}

		mealsResponse = append([]response.GetMeal{mealDishes}, mealsResponse...)
	}

	if mealsResponse == nil {
		//return nil, http.StatusNotFound, errors.New("meal for current day not found")
		mealsResponse = make([]response.GetMeal, 0)
	}

	return mealsResponse, 0, nil
}