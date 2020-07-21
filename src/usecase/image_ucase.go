package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"os"
	"path/filepath"
)

type image struct{}

func NewImage() *image {
	return &image{}
}

var imageRepo = repository.NewImageRepo()

// AddImage adds image for dish with provided ID
// @Summary Add image for certain dish
// @Tags catering images
// @Produce json
// @Param id path string false "Catering ID"
// @Param dishId query string false "Dish ID"
// @Param image formData file false  "Image File"
// @Param id formData string false  "id of default image"
// @Success 201 {object} domain.Image "Image model"
// @Success 200 {object} domain.Image "default image"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/images [post]
func (i image) Add(c *gin.Context) {
	var path types.PathId
	var query types.DishIdQuery

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	id := c.PostForm("id")
	if id != "" {
		parsedID, err := uuid.FromString(id)
		if err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
		image, err, code := imageRepo.AddDefault(path.ID, query.DishId, parsedID)
		if err != nil {
			utils.CreateError(code, err.Error(), c)
			return
		}
		c.JSON(http.StatusOK, image)
		return
	} else {
		file, err := c.FormFile("image")
		if err != nil {
			utils.CreateError(http.StatusBadRequest, "form file error:"+err.Error(), c)
			return
		}

		dir, _ := os.Getwd()
		ext := filepath.Ext(file.Filename)
		randomString := utils.GenerateString(10)

		imagePath := dir + "/src/static/images/" + randomString + ext

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}

		image := domain.Image{
			Path: "/" + randomString + ext,
		}

		result, err, code := imageRepo.Add(path.ID, query.DishId, image)
		if err != nil {
			utils.CreateError(code, err.Error(), c)
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

// DeleteImage soft delete of image
// @Summary Soft delete
// @Tags catering images
// @Produce json
// @Param id path string true "Catering ID"
// @Param imageId path string true "Image ID"
// @Param dishId path string true "Dish ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/images/{imageId}/dish/{dishId} [delete]
func (i image) Delete(c *gin.Context) {
	var path types.PathImageDish

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err, code := imageRepo.Delete(path.CateringID, path.ImageId, path.DishId); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetImages return list of default images
// @Summary Returns list of images
// @Tags catering images
// @Produce json
// @Success 200 {array} response.GetImages "List of images"
// @Failure 400 {object} types.Error "Error"
// @Router /images [get]
func (i image) Get(c *gin.Context) {
	images, err := imageRepo.Get()
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, images)
}

// Update Updates image for dish with provided ID
// @Summary Updates image for certain dish
// @Tags catering images
// @Produce json
// @Param id path string false "Catering ID"
// @Param imageId path string false "Image ID"
// @Param dishId path string false "Dish ID"
// @Param image formData file false  "Image File"
// @Success 201 {object} domain.Image "Image model"
// @Success 200 {object} domain.Image "default image"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/images/{imageId}/dish/{dishId} [put]
func (i image) Update(c *gin.Context) {
	var path types.PathImageDish

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.CreateError(http.StatusBadRequest, "form file error:"+err.Error(), c)
		return
	}

	dir, _ := os.Getwd()
	ext := filepath.Ext(file.Filename)
	randomString := utils.GenerateString(10)

	imagePath := dir + "/src/static/images/" + randomString + ext

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	image := domain.Image{
		Path: "/" + randomString + ext,
	}

	imageResult, err, code := imageRepo.UpdateDishImage(path.CateringID, path.ImageId, path.DishId, image)
	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, imageResult)
}
