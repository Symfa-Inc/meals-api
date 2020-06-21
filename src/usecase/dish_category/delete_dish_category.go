package dish_category

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// DeleteDishCategory soft delete of category reading
// @Summary Soft delete
// @Tags catering dish-categories
// @Produce json
// @Param id path string true "Catering ID"
// @Param categoryId path string true "Category ID"
// @Success 204
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id}/dish-categories/{categoryId} [delete]
func DeleteDishCategory(c *gin.Context) {
	var path types.PathDishCategory
	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	result := repository.DeleteDishCategoryDB(path)

	if result.RowsAffected == 0 {
		utils.CreateError(http.StatusNotFound, "category not found", c)
		return
	}

	c.Status(http.StatusNoContent)
}
