package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Meal struct {
	Base
	Date       time.Time `json:"date,omitempty" binding:"required"`
	CateringID uuid.UUID `json:"-"`
} // @name MealsResponse

type MealUsecase interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

type MealRepository interface {
	Find(meal Meal) error
	Add(meal Meal) (interface{}, error)
	Get(mealDate time.Time, id string) ([]GetMealDish, uuid.UUID, error, int)
	GetByKey(key, value string) (Meal, error, int)
}
