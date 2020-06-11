package models

import (
	"go_api/src/types"
)

// User model
type User struct {
	Base
	FirstName string         `gorm:"type:varchar(20); unique_index" json:"firstname,omitempty"`
	LastName  string         `gorm:"type:varchar(20)" json:"lastname,omitempty"`
	Email     string         `gorm:"type:varchar(30);unique;not null" json:"email,omitempty"`
	Password  string         `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	Role      types.UserRole `sql:"type:user_roles" json:"role,omitempty"`
} //@name UsersResponse

