package domain

import (
	"time"

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
	Status     string    `json:"status" gorm:"default:'draft'"`
} // @name MealsResponse
