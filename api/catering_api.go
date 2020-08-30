package api

import (
	"github.com/Aiscom-LLC/meals-api/schemes/response"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Catering struct
type Catering struct{}

// NewCatering returns pointer to catering struct
// with all methods
func NewCatering() *Catering {
	return &Catering{}
}

var cateringService = services.NewCateringService()

// Get return list of caterings
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Success 200 {object} response.GetCaterings "List of caterings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings [get]
func (ca Catering) Get(c *gin.Context) {
	var query types.PaginationQuery
	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	caterings, total, code, err := cateringService.Get(c, &query)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, response.GetCaterings{
		Items: caterings,
		Page:  query.Page,
		Total: total,
	})
}
