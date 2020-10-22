package models

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

// CopyMeal request scheme
type CopyMealToDates struct {
	CateringID uuid.UUID `json:"-"`
	Date       time.Time `json:"date" binding:"required" example:"2020-06-20T00:00:00Z"`
	ToDate     time.Time `json:"toDate" binding:"required" example:"2020-06-20T00:00:00Z"`
} // @name CopyMealToDate
