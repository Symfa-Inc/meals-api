package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/domain"
	"net/http"
	"os"
)

// ImageRepo struct
type ImageRepo struct{}

// NewImageRepo returns pointer to image repository
// with all methods
func NewImageRepo() *ImageRepo {
	return &ImageRepo{}
}

// GetByKey returns image struct and error by provided key and value
func (i ImageRepo) GetByKey(key, value string) (domain.Image, error) {
	var image domain.Image
	err := config.DB.
		Where(key+"= ?", value).
		First(&image).Error
	return image, err
}

// Add image for provided dish id, and also adds it in imageDish table
// Returns image struct, error and status code
func (i ImageRepo) Add(cateringID, dishID string, image *domain.Image) (int, error) {
	if err := config.DB.
		Where("id = ? AND catering_id = ?", dishID, cateringID).
		Find(&domain.Dish{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, err
		}
		return http.StatusBadRequest, err
	}

	if err := config.DB.Create(image).Error; err != nil {
		return http.StatusBadRequest, err
	}

	parsedDishID, _ := uuid.FromString(dishID)
	imageDish := domain.ImageDish{
		ImageID: image.ID,
		DishID:  parsedDishID,
	}

	if err := config.DB.Create(&imageDish).Error; err != nil {
		return http.StatusBadRequest, err
	}

	return 0, nil
}

// AddDefault adds default image for provided dish
// id and only creates imageDish column
// Returns error and status code
func (i ImageRepo) AddDefault(cateringID, dishID string, imageID uuid.UUID) (domain.Image, int, error) {
	var image domain.Image

	if err := config.DB.
		Where("id = ? AND catering_id = ?", dishID, cateringID).
		Find(&domain.Dish{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Image{}, http.StatusNotFound, err
		}
		return domain.Image{}, http.StatusBadRequest, err
	}

	if err := config.DB.
		Where("id = ?", imageID).
		Find(&image).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Image{}, http.StatusNotFound, err
		}
		return domain.Image{}, http.StatusBadRequest, err
	}

	parsedDishID, _ := uuid.FromString(dishID)

	if rows := config.DB.
		Where("image_id = ? AND dish_id = ?", imageID, dishID).
		Find(&domain.ImageDish{}).
		RowsAffected; rows != 0 {
		return domain.Image{}, http.StatusBadRequest, errors.New("can't add the same default image to the dish")
	}

	imageDish := domain.ImageDish{
		ImageID: imageID,
		DishID:  parsedDishID,
	}

	if err := config.DB.Create(&imageDish).Error; err != nil {
		return domain.Image{}, http.StatusBadRequest, err
	}

	return image, 0, nil
}

// Delete deletes image from imageDish table
// Returns error and status code
func (i ImageRepo) Delete(cateringID, imageID, dishID string) (int, error) {
	var imageToDelete domain.Image
	if err := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		Error; err != nil {

		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, err
		}

		return http.StatusBadRequest, err
	}

	if rows := config.DB.
		Where("image_id = ? AND dish_id = ?", imageID, dishID).
		Delete(&domain.ImageDish{}).RowsAffected; rows == 0 {
		return http.StatusNotFound, errors.New("image or dish with that ID not found")
	}

	if imageExist := config.DB.
		Where("id = ? AND category IS NULL", imageID).
		Find(&imageToDelete).
		Delete(&domain.Image{}).RowsAffected; imageExist != 0 {
		dir, _ := os.Getwd()
		imagePath := dir + "/src/static/images/" + imageToDelete.Path

		if err := os.Remove(imagePath); err != nil {
			return http.StatusBadRequest, err
		}

		return 0, nil
	}

	return 0, nil
}

// Get return list of default images and error
func (i ImageRepo) Get() ([]domain.Image, error) {
	var images []domain.Image

	if err := config.DB.
		Where("category IS NOT NULL").
		Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

// UpdateDishImage updates already existing image in ImageDish table
// Doesn't delete or change previous image in image table
func (i ImageRepo) UpdateDishImage(cateringID, imageID, dishID string, image *domain.Image) (int, error) {
	if err := config.DB.
		Where("id = ?", cateringID).
		Find(&domain.Catering{}).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, err
		}
		return http.StatusBadRequest, err
	}

	if err := config.DB.
		Where("image_id = ? AND dish_id = ?", imageID, dishID).
		Find(&domain.ImageDish{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, err
		}
		return http.StatusBadRequest, err
	}

	if err := config.DB.
		Create(image).
		Error; err != nil {
		return http.StatusNotFound, err
	}

	if err := config.DB.
		Model(&domain.ImageDish{}).
		Where("image_id = ? AND dish_id = ?", imageID, dishID).
		Update("image_id", image.ID).Error; err != nil {
		return http.StatusBadRequest, err
	}

	return 0, nil
}
