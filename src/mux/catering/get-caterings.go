package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/types"
	"net/http"
)

// GetCaterings godoc
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {array} models.Catering "List of caterings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [get]
func GetCaterings(c *gin.Context) {
	var query types.PaginationQuery
	if err := c.ShouldBind(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	cateringsList, err := models.GetCateringsDB(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, cateringsList)
}