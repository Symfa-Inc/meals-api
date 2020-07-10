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

type MealsResult struct {
	Category
	Dish
}

type MealRepository interface {
	Find(meal Meal) error
	Add(meal Meal) (interface{}, error)
	Get(mealId, id string) (map[string][]interface{}, error, int)
	Update(path types.PathMeal, meal Meal) (error, int)
	GetByKey(key, value string) (Meal, error)
}
