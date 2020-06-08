package models

import (
	"go_api/src/config"
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
}

//GetUserByKey returns user
func GetUserByKey(key, value string) (User, error) {
	var user User
	err := config.DB.Where(key+" = ?", value).First(&user).Error
	return user, err
}
