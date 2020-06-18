package dish_category

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/models"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"strings"
)

// AddDishCategory add dish category in DB
// returns 201 if success and 4xx if request failed
// @Summary Returns created category
// @Produce json
// @Accept json
// @Tags catering dish-categories
// @Param id path string true "Catering ID"
// @Param body body request.AddDishCategory false "Category Name"
// @Success 201 {object} models.DishCategory "Dish Category"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dish-category [post]
func AddDishCategory(c *gin.Context) {
	var categoryModel models.DishCategory
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&categoryModel, c); err != nil {
		return
	}

	categoryModel.Name = strings.ToLower(categoryModel.Name)
	categoryModel.CateringID, _ = uuid.FromString(path.ID)

	dishResult, err := repository.CreateDishCategory(categoryModel)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, dishResult)
}
