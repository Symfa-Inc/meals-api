package api

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Image struct
type Image struct{}

// NewImage returns pointer to image struct
// with all methods
func NewImage() *Image {
	return &Image{}
}

var imageService = services.NewImageService()
var imageRepo = repository.NewImageRepo()

// Add adds image for dish with provided ID
// @Summary Add image for certain dish
// @Tags catering dishes
// @Produce json
// @Param id path string false "Catering ID"
// @Param dishId query string false "Dish ID"
// @Param image formData file false  "Image File"
// @Param id formData string false  "id of default image"
// @Success 201 {object} domain.Image "Image model"
// @Success 200 {object} domain.Image "default image"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/dishes/{dishId}/images [post]
func (i Image) Add(c *gin.Context) {
	var path url.PathDish

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	image, code, err := imageService.Add(c, path)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusCreated, image)
}

// Delete deletes image from dish
// @Summary Soft delete
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param imageId path string true "Image ID"
// @Param dishId path string true "Dish ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/dishes/{dishId}/images/{imageId} [delete]
func (i Image) Delete(c *gin.Context) {
	var path url.PathImageDish

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if code, err := imageRepo.Delete(path.CateringID, path.ImageID, path.DishID); err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// Get return list of default images
// @Summary Returns list of images
// @Tags catering images
// @Produce json
// @Success 200 {array} swagger.GetImages "List of images"
// @Failure 400 {object} Error "Error"
// @Router /images [get]
func (i Image) Get(c *gin.Context) {
	images, err := imageRepo.Get()

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	c.JSON(http.StatusOK, images)
}

// Update updates image for dish with provided ID
// @Summary Updates image for certain dish
// @Tags catering dishes
// @Produce json
// @Param id path string false "Catering ID"
// @Param imageId path string false "Image ID"
// @Param dishId path string false "Dish ID"
// @Param image formData file false  "Image File"
// @Success 201 {object} domain.Image "Image model"
// @Success 200 {object} domain.Image "default image"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Not Found"
// @Router /caterings/{id}/dishes/{dishId}/images/{imageId} [put]
func (i Image) Update(c *gin.Context) {
	var path url.PathImageDish

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	image, code, err := imageService.Update(c, path)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, image)
}
