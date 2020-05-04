package auth

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	users "go_api/src/models/user"
	"go_api/src/mux/middleware"
	"go_api/src/schemes/response/auth"
	"net/http"
)

// IsAuthenticated godoc
// @Summary returns user info if user is authorized
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} auth.IsAuthenticated
// @Failure 404 {object} types.Error
// @Security ApiKeyAuth
// @Router /is-authenticated [get]
func IsAuthenticated(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[middleware.IdentityKey]
	user, _ := users.GetUserByKey("id", id.(string))
	c.JSON(http.StatusOK, auth.IsAuthenticated{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	})
}
