package services

import (
	"errors"
	"net/http"

	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
)

// AuthService struct
type AuthService struct{}

// NewAuthService returns pointer to Auth struct
// with all methods
func NewAuthService() *AuthService {
	return &AuthService{}
}

var userRepository = repository.NewUserRepo()

func (as *AuthService) IsAuthenticated(c *gin.Context) (models.UserClientCatering, int, error) {
	userRepo := repository.NewUserRepo()
	claims, err := middleware.Passport().CheckIfTokenExpire(c)

	if err != nil {
		return models.UserClientCatering{}, http.StatusUnauthorized, err
	}

	if int64(claims["exp"].(float64)) < middleware.Passport().TimeFunc().Unix() {
		_, _, _ = middleware.Passport().RefreshToken(c)
	}

	id := claims[middleware.IdentityKeyID]
	result, err := userRepo.GetByID(id.(string))

	if err != nil {
		return models.UserClientCatering{}, http.StatusUnauthorized, errors.New("token is expired")
	}

	if result.Status == &enums.StatusTypesEnum.Deleted {
		return models.UserClientCatering{}, http.StatusForbidden, errors.New("user was deleted")
	}
	return result, 0, nil
}

func (as *AuthService) ChangePassword(body swagger.UserPasswordUpdate, user interface{}) (int, error) {
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
