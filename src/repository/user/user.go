package user

import (
	"go_api/src/config"
	"go_api/src/models"
)

// GetUserByKey returns user
// and error if exists
func GetUserByKey(key, value string) (models.User, error) {
	var user models.User
	err := config.DB.Where(key+" = ?", value).First(&user).Error
	return user, err
}
