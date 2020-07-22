package domain

import "github.com/gin-gonic/gin"

// Auth struct
type Auth struct{}

// AuthUsecase interface for auth usecase
type AuthUsecase interface {
	IsAuthenticated(c *gin.Context)
}
