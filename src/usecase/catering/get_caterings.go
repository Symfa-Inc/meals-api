package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository"
	response2 "go_api/src/schemes/response"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// GetCaterings godoc
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} response.GetCaterings "List of caterings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [get]
func GetCaterings(c *gin.Context) {
	var query types.PaginationQuery

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	result, total, err := repository.GetCateringsDB(query)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}

	c.JSON(http.StatusOK, response2.GetCaterings{
		Items: result,
		Page:  query.Page,
		Total: total,
	})
}
