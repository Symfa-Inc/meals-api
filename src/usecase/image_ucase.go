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
// @Success 204 "Successfully added default image"
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
		err, code := imageRepo.AddDefault(path.ID, query.DishId, parsedID)
		if err != nil {
			utils.CreateError(code, err.Error(), c)
			return
		}
		c.Status(http.StatusNoContent)
		return
	} else {
		file, err := c.FormFile("image")
		if err != nil {
			utils.CreateError(http.StatusBadRequest, "form file error:"+err.Error(), c)
			return
		}

		filename := filepath.Base(file.Filename)
		dir, err := os.Getwd()
		imagePath := dir + "/src/static/images/" + filename
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}

		image := domain.Image{
			Path: filename,
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
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/images/{imageId} [delete]
func (i image) Delete(c *gin.Context) {
	var path types.PathImage

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err, code := imageRepo.Delete(path.CateringID, path.ImageId); err != nil {
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
