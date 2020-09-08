package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"time"

	"github.com/gin-gonic/gin"
)

// MealAPI is meal interface for API
type MealAPI interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
}

// MealRepository is meal interface for repository
type MealRepository interface {
	Find(meal *domain.Meal) error
	Add(meal *domain.Meal) error
	Get(mealDate time.Time, id, clientID string) ([]models.GetMeal, int, error)
	GetByKey(key, value string) (domain.Meal, int, error)
	GetByRange(startDate time.Time, endDate time.Time, id, clientID string) ([]models.GetMeal, int, error)
}

// MealService is meal interface for service
type MealService interface {
	Add(path url.PathClient, body models.AddMeal, user interface{}) ([]models.GetMeal, int, error)
}
