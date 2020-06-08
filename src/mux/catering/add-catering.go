package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"net/http"
)

// AddCatering godoc
// @Summary Returns a ID of created catering
// @Produce json
// @Accept json
// @Tags catering
// @Param body body catering.AddCateringScheme false "Catering Name"
// @Success 200 {object} response.ResponseID
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [post]
func AddCatering(c *gin.Context) {
	var catering models.Catering
	if err := c.ShouldBindJSON(&catering); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
	}
	catering, err := models.CreateCatering(catering)
	if err == nil {
		c.JSON(http.StatusCreated, gin.H{
			"id": catering.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "catering with that name already exist",
		})
	}
}
