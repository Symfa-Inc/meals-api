package swagger

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// AddMeal request scheme
type AddMeal struct {
	CateringID uuid.UUID `json:"-"`
	Date       time.Time `json:"date" binding:"required" example:"2020-06-20T00:00:00Z"`
	Dishes     []string  `json:"dishes" binding:"required"`
} // @name AddMealRequest

// AddMeal request scheme
type AddMealToDate struct {
	CateringID uuid.UUID `json:"-"`
	Date       time.Time `json:"date" binding:"required" example:"2020-06-20T00:00:00Z"`
	NewDate    time.Time `json:"new-date" binding:"required" example:"2020-06-20T00:00:00Z"`
} // @name AddMealToDate
