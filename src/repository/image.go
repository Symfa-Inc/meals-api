package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go_api/src/config"
	"go_api/src/domain"
	"net/http"
	"os"
)

type imageRepo struct{}

func NewImageRepo() *imageRepo {
	return &imageRepo{}
}

// Returns image struct and error by provided key and value
func (i imageRepo) GetByKey(key, value string) (domain.Image, error) {
	var image domain.Image
	err := config.DB.
		Where(key+"= ?", value).
		First(&image).Error
	return image, err
}

// Adds image for provided dish id, and also adds it in imageDish table
// Returns image struct, error and status code
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

// Adds default image for provided dish id and only creates imageDish column
// Returns error and status code
func (i imageRepo) AddDefault(cateringId, dishId string, imageId uuid.UUID) (domain.Image, error, int) {
	var image domain.Image
	if err := config.DB.
		Where("id = ? AND catering_id = ?", dishId, cateringId).
		Find(&domain.Dish{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Image{}, err, http.StatusNotFound
		}
		return domain.Image{}, err, http.StatusBadRequest
	}

	if err := config.DB.
		Where("id = ?", imageId).
		Find(&image).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Image{}, err, http.StatusNotFound
		}
		return domain.Image{}, err, http.StatusBadRequest
	}

	parsedDishId, _ := uuid.FromString(dishId)
	if rows := config.DB.
		Where("image_id = ? AND dish_id = ?", imageId, dishId).
		Find(&domain.ImageDish{}).
		RowsAffected; rows != 0 {
		return domain.Image{}, errors.New("can't add the same default image to the dish"), http.StatusBadRequest
	}

	imageDish := domain.ImageDish{
		ImageID: imageId,
		DishID:  parsedDishId,
	}

	if err := config.DB.Create(&imageDish).Error; err != nil {
		return domain.Image{}, err, http.StatusBadRequest
	}

	return image, nil, 0
}

// Soft delete of image
// Deletes image from imageDish table and also from images table
// Returns error and status code
func (i imageRepo) Delete(cateringId, imageId, dishId string) (error, int) {
	var imageToDelete domain.Image
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
		Where("image_id = ? AND dish_id = ?", imageId, dishId).
		Delete(&domain.ImageDish{}).RowsAffected; rows == 0 {
		return errors.New("image or dish with that ID not found"), http.StatusNotFound
	}

	if imageExist := config.DB.
		Where("id = ? AND category IS NULL", imageId).
		Find(&imageToDelete).
		Delete(&domain.Image{}).RowsAffected; imageExist != 0 {
		dir, _ := os.Getwd()
		imagePath := dir + "/src/static/images/" + imageToDelete.Path

		if err := os.Remove(imagePath); err != nil {
			return err, http.StatusBadRequest
		}

		return nil, 0
	}

	return nil, 0
}

// Return list of default images and error
func (i imageRepo) Get() ([]domain.Image, error) {
	var images []domain.Image
	if err := config.DB.
		Where("category IS NOT NULL").
		Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (i imageRepo) UpdateDishImage(cateringId, imageId, dishId string, image domain.Image) (domain.Image, error, int) {
	if err := config.DB.
		Where("id = ?", cateringId).
		Find(&domain.Catering{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Image{}, err, http.StatusNotFound
		}
		return domain.Image{}, err, http.StatusBadRequest
	}

	if err := config.DB.
		Where("image_id = ? AND dish_id = ?", imageId, dishId).
		Find(&domain.ImageDish{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Image{}, err, http.StatusNotFound
		}
		return domain.Image{}, err, http.StatusBadRequest
	}

	if err := config.DB.
		Create(&image).
		Error; err != nil {
		return domain.Image{}, err, http.StatusNotFound
	}

	if err := config.DB.
		Model(&domain.ImageDish{}).
		Where("image_id = ? AND dish_id = ?", imageId, dishId).
		Update("image_id", image.ID).Error; err != nil {
		return domain.Image{}, err, http.StatusBadRequest
	}

	//dir, _ := os.Getwd()
	//imagePath := dir + "/src/static/images/" + prevImage.Path
	//
	//if err := os.Remove(imagePath); err != nil {
	//	return domain.Image{}, err, http.StatusBadRequest
	//}

	return image, nil, 0
}
