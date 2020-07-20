package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/types"
)

type Dish struct {
	Base
	Name       string       `json:"name" gorm:"not null" binding:"required"`
	Weight     int          `json:"weight" gorm:"not null"binding:"required"`
	Price      int          `json:"price" gorm:"not null" binding:"required"`
	Desc       string       `json:"desc"`
	Images     []ImageArray `json:"images"`
	CateringID uuid.UUID    `json:"-"`
	CategoryID uuid.UUID    `json:"categoryId,omitempty"`
} //@name DishRequest

type DishUsecase interface {
	Add(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

type DishRepository interface {
	Add(cateringId string, dish Dish) (Dish, error)
	Delete(path types.PathDish) error
	Get(cateringId, categoryId string) ([]Dish, error, int)
	GetByKey(key, value, cateringId, categoryId string) (Dish, error, int)
	Update(path types.PathDish, dish Dish) (error, int)
}
