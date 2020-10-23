package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// ClientCategory struct
type ClientCategory struct {
	Base
	Date       *time.Time `json:"date"`
	Name       string     `gorm:"type:varchar(150);not null" json:"name" binding:"required"`
	CateringID uuid.UUID  `json:"-"`
	ClientID   uuid.UUID  `json:"clientId"`
} //@name ClientCategory
