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
	result := catering.UpdateCateringDB(path.ID, cateringModel)
	cateringName := result.Value.(*models.Catering).Name
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  http.StatusNotFound,
			"error": "catering not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":   path.ID,
		"name": cateringName,
	})
}
