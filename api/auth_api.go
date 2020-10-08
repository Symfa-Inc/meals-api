package api

import (
	"net/http"

	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/services"

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

// @Summary Returns error or 200 status code if success
// @Produce json
// @Accept json
// @Tags auth
// @Param body body swagger.UserPasswordUpdate false "User"
// @Success 200 {object} swagger.UserResponse false "User"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Error"
// @Router /auth/change-password [put]
func (a Auth) ChangePassword(c *gin.Context) { //nolint:dupl
	var body swagger.UserPasswordUpdate

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	user, _ := c.Get("user")

	code, err := authService.ChangePassword(body, user)
	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, "Password updated")
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
