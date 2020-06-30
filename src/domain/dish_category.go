package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/types"
)

type DishCategory struct {
	Base
	Name       string    `gorm:"type:varchar(30);not null" json:"name" binding:"required"`
	CateringID uuid.UUID `json:"-"`
} //@name DishCategoryResponse

type DishCategoryUsecase interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type DishCategoryRepository interface {
	Add(category DishCategory) (DishCategory, error)
	Get(id string) ([]DishCategory, error, int)
	GetByKey(id, value, cateringId string) (DishCategory, error)
	Delete(path types.PathDishCategory) error
	Update(path types.PathDishCategory, category DishCategory) (error, int)
}
