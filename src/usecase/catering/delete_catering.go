package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository/catering"
	"go_api/src/types"
	"net/http"
)

// DeleteCatering godoc
// @Summary Soft delete
// @Tags catering
// @Produce json
// @Param id path string true "Catering ID"
// @Success 204
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id} [delete]
func DeleteCatering(c *gin.Context) {
	var path types.PathId
	if err := c.ShouldBindUri(&path); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	result := catering.DeleteCateringDB(path.ID)
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":  http.StatusNotFound,
			"error": "catering not found",
		})
		return
	}
	c.Status(http.StatusNoContent)
}
