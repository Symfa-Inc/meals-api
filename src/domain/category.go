package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/types"
)

type Category struct {
	Base
	Name       string    `gorm:"type:varchar(30);not null" json:"name" binding:"required"`
	CateringID uuid.UUID `json:"-"`
} //@name CategoryResponse

type CategoryUsecase interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type CategoryRepository interface {
	Add(category Category) (Category, error)
	Get(id string) ([]Category, error, int)
	GetByKey(id, value, cateringId string) (Category, error)
	Delete(path types.PathCategory) error
	Update(path types.PathCategory, category Category) (error, int)
}
