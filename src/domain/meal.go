package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"time"
)

// Meal struct for DB
type Meal struct {
	Base
	Date       time.Time `json:"date,omitempty" binding:"required"`
	CateringID uuid.UUID `json:"-"`
} // @name MealsResponse

// MealUsecase is meal interface for usecase
type MealUsecase interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

// MealRepository is meal interface for repository
type MealRepository interface {
	Find(meal Meal) error
	Add(meal Meal) (interface{}, error)
	Get(mealDate time.Time, id string) ([]Dish, uuid.UUID, int, error)
	GetByKey(key, value string) (Meal, int, error)
}
