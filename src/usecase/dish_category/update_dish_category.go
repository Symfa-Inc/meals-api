package dish_category

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// UpdateDishCategory updates dish category with new value provided in body
// @Summary Returns updated dish category object
// @Produce json
// @Accept json
// @Tags catering dish-categories
// @Param id path string true "Catering ID"
// @Param categoryId path string true "Category ID"
// @Param body body request.AddDishCategory false "new category name"
// @Success 200 {object} models.DishCategory "Dish Category"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/dish-category/{categoryId} [put]
func UpdateDishCategory(c *gin.Context) {
	var path types.PathDishCategory
	var body models.DishCategory

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	result, err := repository.UpdateDishCategoryDB(path, body)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if result.RowsAffected == 0 {
		if result.Error != nil {
			utils.CreateError(http.StatusBadRequest, result.Error.Error(), c)
			return
		}

		utils.CreateError(http.StatusNotFound, "category not found", c)
		return
	}

	categoryName := result.Value.(*models.DishCategory).Name

	c.JSON(http.StatusOK, gin.H{
		"id":   path.ID,
		"name": categoryName,
	})
}
