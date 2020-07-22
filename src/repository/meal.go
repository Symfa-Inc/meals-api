package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go_api/src/config"
	"go_api/src/domain"
	"net/http"
	"time"
)

// MealRepo struct
type MealRepo struct{}

// NewMealRepo returns pointer to meal repository
// with all methods
func NewMealRepo() *MealRepo {
	return &MealRepo{}
}

// Find looks for meal in db, if it doesn't exist returns nil
func (m MealRepo) Find(meal domain.Meal) error {
	if exist := config.DB.
		Where("catering_id = ? AND date = ?", meal.CateringID, meal.Date).
		Find(&meal).RecordNotFound(); !exist {
		return errors.New("this meal already exist")
	}
	return nil
}

// Add create meal entity
// returns new meal item and error
func (m MealRepo) Add(meal domain.Meal) (interface{}, error) {
	mealItem := config.DB.Create(&meal)

	if mealItem.Error != nil {
		return nil, mealItem.Error
	}

	return mealItem.Value, nil
}

// Get returns list of meals, total items if and error
func (m MealRepo) Get(mealDate time.Time, id string) ([]domain.Dish, uuid.UUID, int, error) {
	var meal domain.Meal
	var result []domain.Dish

	if err := config.DB.
		Where("catering_id = ? AND date = ?", id, mealDate).
		First(&meal).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []domain.Dish{}, uuid.Nil, http.StatusNotFound, err
		}
		return []domain.Dish{}, uuid.Nil, http.StatusBadRequest, err
	}

	err := config.DB.
		Unscoped().
		Model(&domain.Category{}).
		Select("categories.id as category_id, categories.deleted_at,d.*").
		Joins("left join dishes d on d.category_id = categories.id").
		Joins("left join meal_dishes md on md.dish_id = d.id").
		Joins("left join meals m on m.id = md.meal_id").
		Where("m.id = ? AND md.deleted_at IS NULL AND (categories.deleted_at > ? OR categories.deleted_at IS NULL)", meal.ID, time.Now()).
		Scan(&result).
		Error

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

	return result, meal.ID, http.StatusBadRequest, err
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
