package domain

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/types"
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
	Get(limit int, dateQuery types.StartEndDateQuery, id string) ([]Meal, int, error)
	Update(path types.PathMeal, meal Meal) error
}
