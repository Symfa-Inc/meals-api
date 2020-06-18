package dish_category

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// GetDishCategories returns list of categories or error
// @Summary get list of categories
// @Tags catering dish-categories
// @Produce json
// @Param id path string false "Catering ID"
// @Success 200 {array} models.DishCategory "array of category readings"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/dish-category [get]
func GetDishCategories(c *gin.Context) {
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	categoriesResult, err := repository.GetDishCategoriesDB(path.ID)
	if err != nil {
		if err.Error() == "catering with that ID is not found" {
			utils.CreateError(http.StatusNotFound, err.Error(), c)
			return
		}
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, categoriesResult)
}
