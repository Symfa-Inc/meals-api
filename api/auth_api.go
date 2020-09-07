package api

import (
	"github.com/Aiscom-LLC/meals-api/services"
	"net/http"

	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
)

// Auth struct
type Auth struct{}

// NewAuth returns pointer to Auth struct
// with all methods
func NewAuth() *Auth {
	return &Auth{}
}

var authService = services.NewAuthService()

// IsAuthenticated check if user is authorized and
// if user exists
// @Summary Returns user info if authorized
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} swagger.UserResponse
// @Failure 401 {object} Error
// @Failure 404 {object} Error
// @Router /is-authenticated [get]
func (a Auth) IsAuthenticated(c *gin.Context) {
	user, code, err := authService.IsAuthenticated(c)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Returns info about user
// @Produce json
// @Accept json
// @Tags auth
// @Param body body swagger.LoginUserRequest false "User Credentials"
// @Success 200 {object} swagger.UserResponse
// @Failure 401 {object} Error "Error"
// @Router /login [post]
// nolint:deadcode, unused
func login() {}

// @Summary Removes cookie if set
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} Error "Success"
// @Failure 401 {object} Error "Error"
// @Router /logout [get]
// nolint:deadcode, unused
func logout() {}
