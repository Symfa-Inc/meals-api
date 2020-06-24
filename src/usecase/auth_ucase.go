package usecase

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"go_api/src/delivery/middleware"
	"go_api/src/repository"
	"net/http"
)

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

// IsAuthenticated
// @Summary Returns user info if authorized
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} domain.User
// @Failure 401 {object} types.Error
// @Security ApiKeyAuth
// @Router /is-authenticated [get]
func (a auth) IsAuthenticated(c *gin.Context) {
	userRepo := repository.NewUserRepo()
	claims := jwt.ExtractClaims(c)
	id := claims[middleware.IdentityKeyID]
	result, _ := userRepo.GetByKey("id", id.(string))
	c.JSON(http.StatusOK, result)
}

// @Summary Returns info about user
// @Produce json
// @Accept json
// @Tags auth
// @Param body body request.LoginUserRequest false "User Credentials"
// @Success 200 {object} domain.User
// @Failure 401 {object} types.Error "Error"
// @Router /login [post]
func login() {}

// @Summary Removes cookie if set
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} types.Error "Success"
// @Failure 401 {object} types.Error "Error"
// @Router /logout [get]
func logout() {}

// @Summary Return JSON with code, expire date and new JWT
// @Produce json
// @Accept json
// @Tags auth
// @Success 200 {object} response.RefreshToken
// @Failure 401 {object} types.Error "Error"
// @Router /refresh-token [get]
func refreshToken() {}
