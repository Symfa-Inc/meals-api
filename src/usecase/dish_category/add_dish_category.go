package dish_category

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/models"
	"go_api/src/repository"
	"go_api/src/schemes/request"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
	"strings"
	"time"
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
// @Router /caterings/{id}/dish-categories [post]
func AddDishCategory(c *gin.Context) {
	var body request.AddDishCategory
	var path types.PathId

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	fmt.Println(body.DeletedAt)

	if body.DeletedAt != nil {
		if body.DeletedAt.Sub(time.Now()).Hours() < 0 {
			utils.CreateError(http.StatusBadRequest, "can't create dish category with already passed deletedAt date", c)
			return
		}
	}

	cateringId, _ := uuid.FromString(path.ID)
	categoryModel := models.DishCategory{
		Base:       models.Base{DeletedAt: body.DeletedAt},
		Name:       strings.ToLower(body.Name),
		CateringID: cateringId,
	}

	dishResult, err := repository.CreateDishCategory(categoryModel)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, dishResult)
}
