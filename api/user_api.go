package api

import (
	"net/http"

	"github.com/Aiscom-LLC/meals-api/services"

	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
)

// User struct
type User struct{}

// NewUser returns pointer to user struct
// with all methods
func NewUser() *User {
	return &User{}
}

// @Summary Returns error or 200 status code if success
// @Produce json
// @Accept json
// @Tags Users
// @Param body body swagger.UserPasswordUpdate false "User"
// @Success 200 {object} swagger.UserResponse false "User"
// @Failure 404 {object} Error "Error"
// @Failure 400 {object} Error "Error"
// @Failure 404 {object} Error "Error"
// @Router /users/change-password [put]
func (u User) ChangePassword(c *gin.Context) { //nolint:dupl
	var body swagger.UserPasswordUpdate

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	user, _ := c.Get("user")

	code, err := services.NewUserService().ChangePassword(body, user)
	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, "Password updated")
}
