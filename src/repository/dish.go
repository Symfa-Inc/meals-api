package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"net/http"
	"time"
)

type dishRepo struct{}

func NewDishRepo() *dishRepo {
	return &dishRepo{}
}

// Creates new dish entity
// returns error or nil
func (d dishRepo) Add(cateringId string, dish domain.Dish) (domain.Dish, error) {
	var total int
	config.DB.
		Model(&domain.Dish{}).
		Where("catering_id = ? AND category_id = ?", cateringId, dish.CategoryID).
		Count(&total)

	if total >= 10 {
		return domain.Dish{}, errors.New("can't add more than 10 dishes for a single category")
	}

	if dishExist := config.DB.
		Where("catering_id = ? AND category_id = ? AND name = ?", cateringId, dish.CategoryID, dish.Name).
		Find(&dish).
		RecordNotFound(); !dishExist {
		return domain.Dish{}, errors.New("this dish already exist in that category")
	}

	if err := config.DB.Create(&dish).Error; err != nil {
		return domain.Dish{}, err
	}

	return dish, nil
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
func (d dishRepo) GetByKey(key, value, cateringId, categoryId string) (domain.Dish, error, int) {
	var dish domain.Dish
	if err := config.DB.
		Where("catering_id = ? AND category_id = ? AND "+key+" = ?", cateringId, categoryId, value).
		First(&dish).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Dish{}, errors.New("dish with that id not found"), http.StatusNotFound
		}
		return domain.Dish{}, err, http.StatusBadRequest
	}
	return dish, nil, 0
}

func (d dishRepo) FindById(cateringId, id string) (domain.Dish, error, int) {
	var dish domain.Dish
	if err := config.DB.
		Where("catering_id = ? AND id = ?", cateringId, id).
		First(&dish).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Dish{}, errors.New("dish with that id not found"), http.StatusNotFound
		}
		return domain.Dish{}, err, http.StatusBadRequest
	}
	return dish, nil, 0
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
		Unscoped().
		Where("id = ? AND (deleted_at > ? OR deleted_at IS NULL)", categoryId, time.Now()).
		Find(&domain.Category{}).
		RecordNotFound(); categoryNotExist {
		return nil, errors.New("category with that ID doesn't exist"), http.StatusNotFound
	}

	err := config.DB.
		Where("catering_id = ? AND category_id = ?", cateringId, categoryId).
		Find(&dishes).
		Error
	for i := range dishes {
		var imagesArray []domain.ImageArray
		config.DB.
			Model(&domain.Image{}).
			Select("images.path, images.id").
			Joins("left join image_dishes id on id.image_id = images.id").
			Joins("left join dishes d on id.dish_id = d.id").
			Where("d.id = ?", dishes[i].ID).
			Scan(&imagesArray)
		dishes[i].Images = imagesArray
	}

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
		Unscoped().
		Where("id = ? AND (deleted_at > ? OR deleted_at IS NULL)", dish.CategoryID, time.Now()).
		Find(&domain.Category{}).
		RecordNotFound(); categoryNotExist {
		return errors.New("dish category not found"), http.StatusNotFound
	}

	if result := config.DB.Model(&dish).
		Where("id = ? AND category_id = ?", path.DishID, dish.CategoryID).
		Update(&dish).RowsAffected; result == 0 {
		return errors.New("dish not found"), http.StatusNotFound
	}

	return nil, 0
}
