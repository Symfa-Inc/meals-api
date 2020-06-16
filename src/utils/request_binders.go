package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
// Validates request body with passed model
func RequestBinderBody(model interface{}, c *gin.Context) error {
	if err := c.ShouldBindJSON(model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"error": err.Error(),
		})
		return err
	}
	return nil
}

// Validates request uri with passed model
func RequestBinderUri(model interface{}, c *gin.Context) error {
	if err := c.ShouldBindUri(model); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return err
	}
	return nil
}

// Validates request query with passed model
func RequestBinderQuery(model interface{}, c *gin.Context) error {
	if err := c.ShouldBindQuery(model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"error": err.Error(),
		})
		return err
	}
	return nil
}
