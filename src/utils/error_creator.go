package utils

import "github.com/gin-gonic/gin"

// Creates an error
func CreateError(code int, err string, c *gin.Context) {
	c.JSON(code, gin.H{
		"code": code,
		"error": err,
	})
}