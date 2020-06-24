package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Dish struct {
	Base
	Name           string    `json:"name" gorm:"not null" binding:"required"`
	Weight         int       `json:"weight" gorm:"not null"binding:"required"`
	Price          int       `json:"price" gorm:"not null" binding:"required"`
	Desc           string    `json:"desc"`
	Images         string    `json:"images"`
	CateringID     uuid.UUID `json:"-"`
	DishCategoryID uuid.UUID `json:"categoryId"`
} //@name DishResponse

type DishUsecase interface {
	Add(c *gin.Context)
}

type DishRepository interface {
	Add(cateringId string, dish Dish) (Dish, error)
}
