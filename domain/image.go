package domain

import (
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
	ID   string `json:"id" gorm:"type:column:id"`
	Path string `json:"path" gorm:"type:column:path"`
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

// ImageUsecase is image interface for usecase
type ImageUsecase interface {
	Delete(c *gin.Context)
	Add(c *gin.Context)
	Get(c *gin.Context)
}
