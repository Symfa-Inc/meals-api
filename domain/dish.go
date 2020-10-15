package domain

import (
	uuid "github.com/satori/go.uuid"
)

// Dish struct used in DB
type Dish struct {
	Base
	Name       string       `json:"name" gorm:"not null" binding:"required"`
	Weight     float32      `json:"weight" gorm:"not null" binding:"required"`
	Price      float32      `json:"price" gorm:"not null" binding:"required"`
	Desc       string       `json:"desc"`
	Images     []ImageArray `json:"images"`
	CateringID uuid.UUID    `json:"-"`
	CategoryID uuid.UUID    `json:"categoryId,omitempty"`
} //@name DishRequest
