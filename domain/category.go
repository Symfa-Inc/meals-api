package domain

import (
	uuid "github.com/satori/go.uuid"
)

// Category struct
type Category struct {
	Base
	Name       string    `gorm:"type:varchar(150);not null" json:"name" binding:"required"`
	CateringID uuid.UUID `json:"-"`
} //@name Category
