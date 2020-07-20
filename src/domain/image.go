package domain

import "github.com/gin-gonic/gin"

type Image struct {
	Base
	Path     string  `json:"path,omitempty" binding:"required"`
	Category *string `json:"category" swaggerignore:"true"`
} // @name ImageResponse

type ImageArray struct {
	ID   string `json:"id" gorm:"column:id"`
	Path string `json:"path" gorm:"column:path"`
} //@name Image

type ImageRepository interface {
	GetByKey(key, value string) (Image, error)
	Delete(cateringId, imageId string) (error, int)
	Add(cateringId, dishId string, image Image) (Image, error, int)
	Get() ([]Image, error)
}

type ImageUsecase interface {
	Delete(c *gin.Context)
	Add(c *gin.Context)
	Get(c *gin.Context)
}
