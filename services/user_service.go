package services

import (
	"errors"
	"net/http"

	"github.com/Aiscom-LLC/meals-api/repository"

	"github.com/Aiscom-LLC/meals-api/api/swagger"

	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/utils"
)

// UserService struct
type UserService struct{}

// NewOrderService return pointer to order struct
// with all methods
func NewUserService() *UserService {
	return &UserService{}
}

var userRepository = repository.NewUserRepo()

func (u *UserService) ChangePassword(body swagger.UserPasswordUpdate, user interface{}) (int, error) {

	if len(body.NewPassword) < 10 {
		return http.StatusBadRequest, errors.New("password must contain at least 10 characters")
	}

	newPassword := utils.HashString(body.NewPassword)
	parsedUserID := user.(domain.User).ID

	if err := utils.CheckPasswordHash(body.NewPassword, user.(domain.User).Password); err {
		return http.StatusBadRequest, errors.New("passwords are the same")
	}

	if ok := utils.CheckPasswordHash(body.OldPassword, user.(domain.User).Password); !ok {
		return http.StatusBadRequest, errors.New("wrong password")
	}

	code, err := userRepository.UpdatePassword(parsedUserID, newPassword)

	return code, err
}
