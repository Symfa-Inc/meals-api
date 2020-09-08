package repository

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/interfaces"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

// DishRepo struct
type DishRepo struct{}

// NewDishRepo returns pointer to dish repository
// with all methods
func NewDishRepo() *DishRepo {
	return &DishRepo{}
}

// Add creates new dish entity
// returns error or nil
func (d DishRepo) Add(cateringID string, dish *interfaces.Dish) error {
	if dishExist := config.DB.
		Where("catering_id = ? AND category_id = ? AND name = ?", cateringID, dish.CategoryID, dish.Name).
		Find(dish).
		RecordNotFound(); !dishExist {
		return errors.New("this dish already exist in that category")
	}

	if err := config.DB.Create(dish).Error; err != nil {
		return err
	}

	return nil
}

// Delete soft delete of entity
// returns error or nil
func (d DishRepo) Delete(path url.PathDish) error {
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

// GetByKey get entity filtered by key and value
// returns entity and error or nil
func (d DishRepo) GetByKey(key, value, cateringID, categoryID string) (interfaces.Dish, int, error) {
	var dish interfaces.Dish

	if err := config.DB.
		Where("catering_id = ? AND category_id = ? AND "+key+" = ?", cateringID, categoryID, value).
		First(&dish).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return interfaces.Dish{}, http.StatusNotFound, errors.New("dish with that id not found")
		}
		return interfaces.Dish{}, http.StatusBadRequest, err
	}

	return dish, 0, nil
}

// FindByID finds dish by ID
// Returns Dish, err and status code
func (d DishRepo) FindByID(cateringID, id string) (interfaces.Dish, int, error) {
	var dish interfaces.Dish
	var imagesArray []interfaces.ImageArray

	if err := config.DB.
		Where("catering_id = ? AND id = ?", cateringID, id).
		First(&dish).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return interfaces.Dish{}, http.StatusNotFound, errors.New("dish with that id not found")
		}
		return interfaces.Dish{}, http.StatusBadRequest, err
	}

	config.DB.
		Model(&domain.Image{}).
		Select("images.path, images.id").
		Joins("left join image_dishes id on id.image_id = images.id").
		Joins("left join dishes d on id.dish_id = d.id").
		Where("d.id = ? AND id.deleted_at IS NULL", dish.ID).
		Scan(&imagesArray)

	dish.Images = imagesArray

	if len(dish.Images) == 0 {
		dish.Images = make([]interfaces.ImageArray, 0)
	}

	return dish, 0, nil
}

// Get list of dishes
// returns array of dishes and error or nil and status code
func (d DishRepo) Get(cateringID, categoryID string) ([]interfaces.Dish, int, error) {
	var dishes []interfaces.Dish

	if cateringNotExist := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		RecordNotFound(); cateringNotExist {
		return nil, http.StatusNotFound, errors.New("catering with that ID doesn't exist")
	}

	if categoryNotExist := config.DB.
		Unscoped().
		Where("id = ? AND (deleted_at > ? OR deleted_at IS NULL)", categoryID, time.Now()).
		Find(&domain.Category{}).
		RecordNotFound(); categoryNotExist {
		return nil, http.StatusNotFound, errors.New("category with that ID doesn't exist")
	}

	err := config.DB.
		Where("catering_id = ? AND category_id = ?", cateringID, categoryID).
		Find(&dishes).
		Error

	for i := range dishes {
		var imagesArray []interfaces.ImageArray
		config.DB.
			Model(&domain.Image{}).
			Select("images.path, images.id").
			Joins("left join image_dishes id on id.image_id = images.id").
			Joins("left join dishes d on id.dish_id = d.id").
			Where("d.id = ? AND id.deleted_at IS NULL", dishes[i].ID).
			Scan(&imagesArray)
		dishes[i].Images = imagesArray
	}

	return dishes, 0, err
}

// Update updates entity
// returns error or nil and status code
func (d DishRepo) Update(path url.PathDish, dish interfaces.Dish) (int, error) {
	if cateringNotExist := config.DB.
		Where("id = ?", path.CateringID).
		Find(&domain.Catering{}).
		RecordNotFound(); cateringNotExist {
		return http.StatusNotFound, errors.New("catering not found")
	}

	if categoryNotExist := config.DB.
		Unscoped().
		Where("id = ? AND (deleted_at > ? OR deleted_at IS NULL)", dish.CategoryID, time.Now()).
		Find(&domain.Category{}).
		RecordNotFound(); categoryNotExist {
		return http.StatusNotFound, errors.New("dish category not found")
	}

	if result := config.DB.Model(&dish).
		Where("id = ? AND category_id = ?", path.DishID, dish.CategoryID).
		Update(&dish).RowsAffected; result == 0 {
		return http.StatusNotFound, errors.New("dish not found")
	}

	return 0, nil
}
