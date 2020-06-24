package repository

import (
	"go_api/src/config"
	"go_api/src/domain"
)

type userRepo struct{}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

// userRepo.GetByKey returns user
// and error if exist
func (ur userRepo) GetByKey(key, value string) (domain.User, error) {
	var user domain.User
	err := config.DB.Where(key+" = ?", value).First(&user).Error
	return user, err
}
