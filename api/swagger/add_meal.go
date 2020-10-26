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

// CopyMealToDate request scheme
type CopyMealToDate struct {
	CateringID uuid.UUID `json:"-"`
	Date       time.Time `json:"date" binding:"required" example:"2020-06-20T00:00:00Z"`
	ToDate     time.Time `json:"toDate" binding:"required" example:"2020-06-20T00:00:00Z"`
} // @name CopyMealToDate

// CopyMealToWeek request scheme
type CopyMealToWeek struct {
	CateringID uuid.UUID   `json:"-"`
	Date       []time.Time `json:"date" binding:"required" example:"2020-10-23T00:00:00Z,2020-10-24T00:00:00Z,2020-10-25T00:00:00Z,2020-10-26T00:00:00Z,2020-10-27T00:00:00Z"`
	ToWeek     []time.Time `json:"toDate" binding:"required" example:"2020-10-23T00:00:00Z,2020-10-30T00:00:00Z,2020-10-30T00:00:00Z,2020-11-01T00:00:00Z,2020-11-02T00:00:00Z"`
} // @name CopyMealToWeek
