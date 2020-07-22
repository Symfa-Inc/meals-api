package repository

import (
	"go_api/src/config"
	"go_api/src/domain"
)

// UserRepo struct
type UserRepo struct{}

// NewUserRepo returns pointer to user repository
// with all methods
func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// GetByKey returns user by key
// and error if exist
func (ur UserRepo) GetByKey(key, value string) (domain.User, error) {
	var user domain.User
	err := config.DB.Where(key+" = ?", value).First(&user).Error
	return user, err
}
