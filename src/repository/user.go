package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go_api/src/config"
	"go_api/src/domain"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
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

// Add adds new user for certain company passed in user struct
// returns user and error
func (ur UserRepo) Add(user domain.User) (domain.User, error) {
	if err := config.DB.Create(&user).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (ur UserRepo) Delete(companyID string, user domain.User) (int, error) {

	companyType := utils.DerefString(user.CompanyType)
	if companyType == types.CompanyTypesEnum.Catering {
		if userExist := config.DB.
			Model(&domain.User{}).
			Where("catering_id = ?", companyID).
			Update(&user).
			RowsAffected; userExist == 0 {
			return http.StatusNotFound, errors.New("user not found")
		}
		return 0, nil
	}

	if userExist := config.DB.
		Model(&domain.User{}).
		Where("client_id = ?", companyID).
		Update(&user).
		RowsAffected; userExist == 0 {
		return http.StatusBadRequest, errors.New("user not found")
	}
	return 0, nil
}

func (ur UserRepo) Update(companyID string, user domain.User) (domain.User, int, error) {
	companyType := utils.DerefString(user.CompanyType)
	var updatedUser domain.User

	if companyType == types.CompanyTypesEnum.Catering {
		if err := config.DB.
			Model(&domain.User{}).
			Where("catering_id = ?", companyID).
			Update(&user).
			Scan(&updatedUser).
			Error; err != nil {

			if gorm.IsRecordNotFoundError(err) {
				return domain.User{}, http.StatusNotFound, errors.New("user not found")
			}
			return domain.User{}, http.StatusBadRequest, err
		}
		return updatedUser, 0, nil
	}

	if err := config.DB.
		Model(&domain.User{}).
		Where("client_id = ?", companyID).
		Update(&user).
		Scan(&updatedUser).
		Error; err != nil {

		if gorm.IsRecordNotFoundError(err) {
			return domain.User{}, http.StatusNotFound, errors.New("user not found")
		}
		return domain.User{}, http.StatusBadRequest, err
	}
	return updatedUser, 0, nil
}
