package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Image struct for DB
type Image struct {
	Base
	Path     string  `json:"path,omitempty" binding:"required"`
	Category *string `json:"category" swaggerignore:"true"`
} // @name ImageResponse

// ImageArray struct
type ImageArray struct {
	ID   string `json:"id" gorm:"column:id"`
	Path string `json:"path" gorm:"column:path"`
} //@name Image

// ImageRepository is image interface for repository
type ImageRepository interface {
	GetByKey(key, value string) (Image, error)
	Delete(cateringID, imageID, dishID string) (int, error)
	Add(cateringID, dishID string, image *Image) (int, error)
	AddDefault(cateringID, dishID string, imageID uuid.UUID) (Image, int, error)
	UpdateDishImage(cateringID, imageID, dishID string, image *Image) (int, error)
	Get() ([]Image, error)
}

// ImageAPI is image interface for API
type ImageAPI interface {
	Delete(c *gin.Context)
	Add(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

// ImageService is image interface for API
type ImageService interface {
	Add(c *gin.Context, path url.PathDish) (Image, int, error)
	Update(c *gin.Context, path url.PathImageDish) (Image, int, error)
}
