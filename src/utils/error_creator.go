package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// CreateError creates an error
func CreateError(code int, err string, c *gin.Context) {
	c.JSON(code, gin.H{
		"code":  code,
		"error": err,
	})
	_ = c.AbortWithError(code, errors.New(err))
}
