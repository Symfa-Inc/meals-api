package auth

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/mux/middleware"
	"go_api/src/schemes/response/auth"
	"net/http"
)

// IsAuthenticated godoc
// @Summary Returns user info if authorized
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} auth.IsAuthenticated
// @Failure 401 {object} types.Error
// @Security ApiKeyAuth
// @Router /is-authenticated [get]
func IsAuthenticated(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[middleware.IdentityKeyID]
	user, _ := models.GetUserByKey("id", id.(string))
	c.JSON(http.StatusOK, auth.IsAuthenticated{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	})
}
