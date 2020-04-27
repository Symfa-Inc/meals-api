package users

import (
	"github.com/jinzhu/gorm"
	"go_api/src/types"
)

// User model
type User struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(20); unique_index"`
	LastName string `gorm:"type:varchar(20)"`
	Email string `gorm:"type:varchar(30);unique;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Role types.UserRole `sql:"type:user_roles"`
}