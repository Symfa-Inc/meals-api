package services

import (
	"errors"
	"net/http"

	"github.com/Aiscom-LLC/meals-api/domain"

	"github.com/Aiscom-LLC/meals-api/utils"

	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/gin-gonic/gin"
)

// AuthService struct
type AuthService struct{}

// NewAuthService returns pointer to Auth struct
// with all methods
func NewAuthService() *AuthService {
	return &AuthService{}
}

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

func (as AuthService) ForgotPassword(body swagger.ForgotPassword) (domain.User, string, int, error) {
	user, err := userRepo.GetByKey("email", body.Email)

	if err != nil {
		return user, "", http.StatusBadRequest, err
	}

	password := utils.GenerateString(10)
	hashPassword := utils.HashString(password)

	code, err := userRepo.UpdatePassword(user.ID, hashPassword)

	return user, password, code, err
}
