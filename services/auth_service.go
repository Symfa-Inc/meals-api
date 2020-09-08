package services

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/repository/enums"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthService struct
type AuthService struct{}

// NewAuthService returns pointer to Auth struct
// with all methods
func NewAuthService() *AuthService {
	return &AuthService{}
}

func (as *AuthService) IsAuthenticated(c *gin.Context) (domain.UserClientCatering, int, error) {
	userRepo := repository.NewUserRepo()
	claims, err := middleware.Passport().CheckIfTokenExpire(c)

	if err != nil {
		return domain.UserClientCatering{}, http.StatusUnauthorized, err
	}

	if int64(claims["exp"].(float64)) < middleware.Passport().TimeFunc().Unix() {
		_, _, _ = middleware.Passport().RefreshToken(c)
	}

	id := claims[middleware.IdentityKeyID]
	result, err := userRepo.GetByID(id.(string))

	if err != nil {
		return domain.UserClientCatering{}, http.StatusUnauthorized, errors.New("token is expired")
	}

	if result.Status == &enums.StatusTypesEnum.Deleted {
		return domain.UserClientCatering{}, http.StatusForbidden, errors.New("user was deleted")
	}
	return result, 0, nil
}
