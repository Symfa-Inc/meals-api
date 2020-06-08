package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/types"
	"net/http"
)

// UpdateCatering godoc
// @Summary Returns updated catering object
// @Produce json
// @Accept json
// @Tags catering
// @Param id path string true "Catering ID"
// @Param body body catering.AddCateringScheme false "Catering Name"
// @Success 200 {object} models.Catering "Catering"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id} [put]
func UpdateCatering(c *gin.Context) {
	var path types.PathId
	var catering models.Catering
	if err := c.ShouldBindJSON(&catering); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	catering, err := models.UpdateCateringDB(path.ID, catering)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": path.ID,
		"name": catering.Name,
	})
}
