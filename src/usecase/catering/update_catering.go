package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/repository/catering"
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
	var cateringModel models.Catering
	if err := c.ShouldBindJSON(&cateringModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	result, err := catering.UpdateCateringDB(path.ID, cateringModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":   path.ID,
		"name": result.Name,
	})
}
