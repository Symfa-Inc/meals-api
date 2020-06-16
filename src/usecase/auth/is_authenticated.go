package auth

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"go_api/src/delivery/middleware"
	"go_api/src/repository/user"
	"net/http"
)

// IsAuthenticated godoc
// @Summary Returns user info if authorized
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} models.User
// @Failure 401 {object} types.Error
// @Security ApiKeyAuth
// @Router /is-authenticated [get]
func IsAuthenticated(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[middleware.IdentityKeyID]
	result, _ := user.GetUserByKey("id", id.(string))
	c.JSON(http.StatusOK, result)
}
