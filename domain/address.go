package domain

import (
	uuid "github.com/satori/go.uuid"
)

// Address struct
type Address struct {
	Base
	City     string    `json:"city" gorm:"not null" binding:"required"`
	Street   string    `json:"street" gorm:"not null" binding:"required"`
	House    string    `json:"house" gorm:"not null" binding:"required"`
	Floor    int       `json:"floor" gorm:"not null" binding:"required"`
	ClientID uuid.UUID `json:"-"`
} //@name AddressResponse
