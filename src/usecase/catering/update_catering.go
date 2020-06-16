package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/repository/catering"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// UpdateCatering godoc
// @Summary Returns updated catering object
// @Produce json
// @Accept json
// @Tags catering
// @Param id path string true "Catering ID"
// @Param body body catering.AddCateringRequest false "Catering Name"
// @Success 200 {object} models.Catering "Catering"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id} [put]
func UpdateCatering(c *gin.Context) {
	var path types.PathId
	var cateringModel models.Catering

	if err := utils.RequestBinderBody(&cateringModel, c); err != nil {
		return
	}

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	result := catering.UpdateCateringDB(path.ID, cateringModel)

	cateringName := result.Value.(*models.Catering).Name

	if result.RowsAffected == 0 {
		if result.Error != nil {
			utils.CreateError(http.StatusBadRequest, result.Error.Error(), c)
			return
		}

		utils.CreateError(http.StatusNotFound, "catering not found", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   path.ID,
		"name": cateringName,
	})
}
