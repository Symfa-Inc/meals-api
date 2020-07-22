package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RequestBinderBody validates request body with passed model
func RequestBinderBody(model interface{}, c *gin.Context) error {
	if err := c.ShouldBindJSON(model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return err
	}
	return nil
}

// RequestBinderURI validates request uri with passed model
func RequestBinderURI(model interface{}, c *gin.Context) error {
	if err := c.ShouldBindUri(model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return err
	}
	return nil
}

// RequestBinderQuery validates request query with passed model
func RequestBinderQuery(model interface{}, c *gin.Context) error {
	if err := c.ShouldBindQuery(model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return err
	}
	return nil
}
