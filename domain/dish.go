package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Dish struct used in DB
type Dish struct {
	Base
	Name       string       `json:"name" gorm:"type:not null" binding:"required"`
	Weight     float32      `json:"weight" gorm:"type:not null" binding:"required"`
	Price      float32      `json:"price" gorm:"type:not null" binding:"required"`
	Desc       string       `json:"desc"`
	Images     []ImageArray `json:"images"`
	CateringID uuid.UUID    `json:"-"`
	CategoryID uuid.UUID    `json:"categoryId,omitempty"`
} //@name DishRequest

// DishUsecase is dish interface for usecase
type DishUsecase interface {
	Add(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

// DishRepository is dish interface for repository
type DishRepository interface {
	Add(cateringID string, dish *Dish) error
	Delete(path url.PathDish) error
	Get(cateringID, categoryID string) ([]Dish, int, error)
	GetByKey(key, value, cateringID, categoryID string) (Dish, int, error)
	Update(path url.PathDish, dish Dish) (int, error)
}
