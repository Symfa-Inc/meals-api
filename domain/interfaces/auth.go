package domain

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/gin-gonic/gin"
)

// AuthAPI interface for auth API
type AuthAPI interface {
	IsAuthenticated(c *gin.Context)
}

// AuthService interface for auth service
type AuthService interface {
	IsAuthenticated(c *gin.Context) (domain.UserClientCatering, int, error)
}
