package interfaces

import "github.com/gin-gonic/gin"

// Auth struct
type Auth struct{}

// AuthAPI interface for auth API
type AuthAPI interface {
	IsAuthenticated(c *gin.Context)
}

// AuthService interface for auth service
type AuthService interface {
	IsAuthenticated(c *gin.Context) (UserClientCatering, int, error)
}
