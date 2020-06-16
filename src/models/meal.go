package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Meal struct {
	Base
	Date       time.Time `json:"date,omitempty" binding:"required"`
	CateringID uuid.UUID `json:"-"`
} // @name MealsResponse
