package models

import (
	uuid "github.com/satori/go.uuid"
)

type DishCategory struct {
	Base
	Name       string    `gorm:"type:varchar(30);not null" json:"name" binding:"required"`
	CateringID uuid.UUID `json:"-"`
} //@name DishCategoryResponse
