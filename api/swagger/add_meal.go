package swagger

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// AddMeal request scheme
type AddMeal struct {
	CateringID uuid.UUID `json:"-"`
	Date       time.Time `json:"date" binding:"required" example:"2020-06-20T00:00:00Z"`
	Dishes     []string  `json:"dishes" binding:"required"`
} // @name AddMealRequest
