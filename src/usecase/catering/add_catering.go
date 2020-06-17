package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/repository"
	"go_api/src/utils"
	"net/http"
)

// AddCatering godoc
// @Summary Returns a ID of created catering
// @Produce json
// @Accept json
// @Tags catering
// @Param body body request.AddCateringRequest false "Catering Name"
// @Success 201 {object} models.Catering
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [post]
func AddCatering(c *gin.Context) {
	var cateringModel models.Catering
	if err := utils.RequestBinderBody(&cateringModel, c); err != nil {
		return
	}

	result, err := repository.CreateCatering(cateringModel)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, "catering with that name already exist", c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   result.ID,
		"name": result.Name,
	})
}
