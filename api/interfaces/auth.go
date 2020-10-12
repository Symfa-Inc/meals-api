package domain

import (
	"github.com/Aiscom-LLC/meals-api/api/swagger"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository/models"
	"github.com/gin-gonic/gin"
)

// AuthAPI interface for auth API
type AuthAPI interface {
	IsAuthenticated(c *gin.Context)
}

// AuthService interface for auth service
type AuthService interface {
	IsAuthenticated(c *gin.Context) (models.UserClientCatering, int, error)
	ChangePassword(body swagger.UserPasswordUpdate, user interface{}) (int, error)
	RecoveryPassword(body swagger.RecoveryPassword) (domain.User, string, int, error)
}
