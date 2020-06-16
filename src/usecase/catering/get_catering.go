package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository/catering"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// GetCatering godoc
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param id path string true "Catering ID"
// @Success 200 {object} models.Catering "List of caterings"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id} [get]
func GetCatering(c *gin.Context) {
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	result, err := catering.GetCateringByKey("id", path.ID)

	if err != nil {
		if err.Error() == "record not found" {
			utils.CreateError(http.StatusNotFound, err.Error(), c)
			return
		} else {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}

	c.JSON(http.StatusOK, result)
}
