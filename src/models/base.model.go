package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid; primary_key" faker:"-" json:"id"`
	CreatedAt time.Time  `faker:"-" json:"createdAt"`
	UpdatedAt time.Time  `faker:"-" json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" faker:"-" json:"deletedAt"`
}

// BeforeCreate will set a UUID rather than numberic ID
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuidv4, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuidv4)
}
