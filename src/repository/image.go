package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go_api/src/config"
	"go_api/src/domain"
	"net/http"
)

type imageRepo struct{}

func NewImageRepo() *imageRepo {
	return &imageRepo{}
}

func (i imageRepo) GetByKey(key, value string) (domain.Image, error) {
	var image domain.Image
	err := config.DB.
		Where(key+"= ?", value).
		First(&image).Error
	return image, err
}

func (i imageRepo) Add(cateringId, dishId string, image domain.Image) (domain.Image, error, int) {
	if err := config.DB.
		Where("id = ? AND catering_id = ?", dishId, cateringId).
		Find(&domain.Dish{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Image{}, err, http.StatusNotFound
		}
		return domain.Image{}, err, http.StatusBadRequest
	}

	if err := config.DB.Create(&image).Error; err != nil {
		return domain.Image{}, err, http.StatusBadRequest
	}

	parsedDishId, _ := uuid.FromString(dishId)
	imageDish := domain.ImageDish{
		ImageID: image.ID,
		DishID:  parsedDishId,
	}

	if err := config.DB.Create(&imageDish).Error; err != nil {
		return domain.Image{}, err, http.StatusBadRequest
	}

	return image, nil, 0
}

func (i imageRepo) Delete(cateringId, imageId string) (error, int) {
	if err := config.DB.
		Where("id = ?", cateringId).
		Find(&domain.Catering{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return err, http.StatusNotFound
		}
		return err, http.StatusBadRequest
	}

	if rows := config.DB.
		Where("image_id =  ?", imageId).
		Delete(&domain.ImageDish{}).RowsAffected; rows == 0 {
		return errors.New("image with that ID not found"), http.StatusNotFound
	}

	config.DB.
		Where("id =  ?", imageId).
		Delete(&domain.Image{})

	return nil, 0
}

func (i imageRepo) Get() ([]domain.Image, error) {
	var images []domain.Image
	if err := config.DB.
		Where("category IS NOT NULL").
		Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
