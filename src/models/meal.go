package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Meal struct {
	Base
	Date time.Time `json:"date,omitempty"`
	CateringID uuid.UUID `json:"cateringId,omitempty"`
}// @name MealsResponse