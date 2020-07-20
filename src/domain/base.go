package domain

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuidv4, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuidv4)
}
