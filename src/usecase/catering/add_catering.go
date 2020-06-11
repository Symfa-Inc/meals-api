package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	catering "go_api/src/repository/catering"
	"net/http"
)

// AddCatering godoc
// @Summary Returns a ID of created catering
// @Produce json
// @Accept json
// @Tags catering
// @Param body body catering.AddCateringScheme false "Catering Name"
// @Success 201 {object} models.Catering
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [post]
func AddCatering(c *gin.Context) {
	var cateringModel models.Catering
	if err := c.ShouldBindJSON(&cateringModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	result, err := catering.CreateCatering(cateringModel)
	if err == nil {
		c.JSON(http.StatusCreated, gin.H{
			"id": result.ID,
			"name": result.Name,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "catering with that name already exist",
		})
		return
	}
}
