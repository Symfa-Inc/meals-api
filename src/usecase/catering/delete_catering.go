package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository/catering"
	"go_api/src/types"
	"go_api/src/utils"
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
	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	result := catering.DeleteCateringDB(path.ID)

	if result.RowsAffected == 0 {
		utils.CreateError(http.StatusNotFound, "catering not found", c)
		return
	}

	c.Status(http.StatusNoContent)
}
