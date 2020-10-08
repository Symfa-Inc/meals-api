package utils

import (
	"github.com/gin-gonic/gin"
)

// CreateError creates an error
func CreateError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"code":  code,
		"error": err.Error(),
	})
	_ = c.AbortWithError(code, err)
}
