package services

import (
	"github.com/Aiscom-LLC/meals-api/api/types"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"path/filepath"
)

// ImageService struct
type ImageService struct{}

// NewImageService returns pointer to image struct
// with all methods
func NewImageService() *ImageService {
	return &ImageService{}
}

var imageRepo = repository.NewImageRepo()

func (i *ImageService) Add(c *gin.Context, path types.PathDish) (domain.Image, int, error) {
	id := c.PostForm("id")

	if id != "" {
		parsedID, err := uuid.FromString(id)

		if err != nil {
			return domain.Image{}, http.StatusBadRequest, err
		}

		image, code, err := imageRepo.AddDefault(path.CateringID, path.DishID, parsedID)

		return image, code, err
	}

	file, err := c.FormFile("image")

	if err != nil {
		return domain.Image{}, http.StatusBadRequest, err
	}

	dir, _ := os.Getwd()
	ext := filepath.Ext(file.Filename)
	randomString := utils.GenerateString(10)

	imagePath := dir + "/src/static/images/" + randomString + ext

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		return domain.Image{}, http.StatusBadRequest, err
	}

	image := &domain.Image{
		Path: "/" + randomString + ext,
	}

	code, err := imageRepo.Add(path.CateringID, path.DishID, image)

	return *image, code, err
}

func (i *ImageService) Update(c *gin.Context, path types.PathImageDish) (domain.Image, int, error) {
	file, err := c.FormFile("image")

	if err != nil {
		return domain.Image{}, http.StatusBadRequest, err
	}

	dir, _ := os.Getwd()
	ext := filepath.Ext(file.Filename)
	randomString := utils.GenerateString(10)

	imagePath := dir + "/src/static/images/" + randomString + ext

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		return domain.Image{}, http.StatusBadRequest, err
	}

	image := &domain.Image{
		Path: "/" + randomString + ext,
	}

	code, err := imageRepo.UpdateDishImage(path.CateringID, path.ImageID, path.DishID, image)

	return *image, code, err
}
