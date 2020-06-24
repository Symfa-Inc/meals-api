package domain

import "github.com/gin-gonic/gin"

type Auth struct{}

type AuthUsecase interface {
	IsAuthenticated(c *gin.Context)
}
