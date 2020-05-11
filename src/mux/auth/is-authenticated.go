package auth

import (
	"github.com/gin-gonic/gin"
	users "go_api/src/models/user"
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
	claims := middleware.ExtractClaims(c)
	id := claims[middleware.IdentityKeyID]
	user, _ := users.GetUserByKey("id", id.(string))
	c.JSON(http.StatusOK, auth.IsAuthenticated{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	})
}
