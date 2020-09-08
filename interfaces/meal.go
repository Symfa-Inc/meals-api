package interfaces

import (
	"github.com/Aiscom-LLC/meals-api/api/url"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// MealBase struct without deletedAt
type MealBase struct {
	ID        uuid.UUID  `gorm:"type:uuid;" json:"id"`
	DeletedAt *time.Time `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

// BeforeCreate func which generates uuid v4 for each inserted row
func (base *MealBase) BeforeCreate(scope *gorm.Scope) error {
	uuidv4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidv4)
}

// Meal struct for DB
type Meal struct {
	MealBase
	CreatedAt  time.Time `json:"createdAt"`
	Date       time.Time `json:"date,omitempty" binding:"required"`
	CateringID uuid.UUID `json:"-"`
	ClientID   uuid.UUID `json:"-"`
	MealID     uuid.UUID `json:"mealId"`
	Version    string    `json:"version"`
	Person     string    `json:"person"`
} // @name MealsResponse

// MealAPI is meal interface for API
type MealAPI interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
}

// MealRepository is meal interface for repository
type MealRepository interface {
	Get(mealDate time.Time, id, clientID string) ([]models.GetMeal, int, error)
	GetByRange(startDate time.Time, endDate time.Time, id, clientID string) ([]models.GetMeal, int, error)
}

// MealService is meal interface for service
type MealService interface {
	Add(path url.PathClient, body models.AddMeal, user interface{}) ([]models.GetMeal, int, error)
}
