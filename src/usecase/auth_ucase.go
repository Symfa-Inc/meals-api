package usecase

import (
	"net/http"

	"github.com/Aiscom-LLC/meals-api/src/delivery/middleware"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/types"
	"github.com/Aiscom-LLC/meals-api/src/utils"
	"github.com/gin-gonic/gin"
)

// Auth struct
type Auth struct{}

// NewAuth returns pointer to Auth strcut
// with all methods
func NewAuth() *Auth {
	return &Auth{}
}

// IsAuthenticated check if user is authorized and
// if user exists
// @Summary Returns user info if authorized
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} response.UserResponse
// @Failure 401 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /is-authenticated [get]
func (a Auth) IsAuthenticated(c *gin.Context) {
	userRepo := repository.NewUserRepo()
	claims, err := middleware.Passport().CheckIfTokenExpire(c)

	if err != nil {
		utils.CreateError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	if int64(claims["exp"].(float64)) < middleware.Passport().TimeFunc().Unix() {
		_, _, _ = middleware.Passport().RefreshToken(c)
	}

	id := claims[middleware.IdentityKeyID]
	result, err := userRepo.GetByID(id.(string))

	if err != nil {
		utils.CreateError(http.StatusUnauthorized, "token is expired", c)
		return
	}

	if result.Status == &types.StatusTypesEnum.Deleted {
		utils.CreateError(http.StatusForbidden, "user was deleted", c)
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Returns info about user
// @Produce json
// @Accept json
// @Tags auth
// @Param body body request.LoginUserRequest false "User Credentials"
// @Success 200 {object} response.UserResponse
// @Failure 401 {object} types.Error "Error"
// @Router /login [post]
// nolint:deadcode
// nolint:unused
func login() {}

// @Summary Removes cookie if set
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} types.Error "Success"
// @Failure 401 {object} types.Error "Error"
// @Router /logout [get]
// nolint:deadcode
func logout() {}
