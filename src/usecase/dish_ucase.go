package usecase

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go_api/src/domain"
	"go_api/src/repository"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

type dish struct{}

func NewDish() *dish {
	return &dish{}
}

var dishRepo = repository.NewDishRepo()

// AddDish adds dish for catering with provided ID
// @Summary Add dish for certain category
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param payload body request.AddDish false "dish object"
// @Success 200 {object} domain.Dish false "dish object"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dishes [post]
func (d dish) Add(c *gin.Context) {
	var path types.PathId
	var body domain.Dish

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	body.CateringID, _ = uuid.FromString(path.ID)
	dish, err := dishRepo.Add(path.ID, body)
	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, dish)
}

// DeleteDish soft delete of dish
// @Summary Soft delete
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param dishId path string true "Dish ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/dishes/{dishId} [delete]
func (d dish) Delete(c *gin.Context) {
	var path types.PathDish

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := dishRepo.Delete(path); err != nil {
		utils.CreateError(http.StatusNotFound, err.Error(), c)
	}

	c.Status(http.StatusNoContent)
}

// GetDishes return list of dishes
// @Summary Returns list of dishes
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param categoryId query string true "Category ID"
// @Success 200 {array} domain.Dish "List of dishes"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dishes [get]
func (d dish) Get(c *gin.Context) {
	//var path types.PathDishGet
	var path types.PathId
	var query types.CategoryIdQuery
	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}
	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	dishes, err, code := dishRepo.Get(path.ID, query.CategoryId)

	if err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, dishes)
}

// UpdateDish updates dish
// @Summary Returns 204 if success and 4xx error if failed
// @Produce json
// @Accept json
// @Tags catering dishes
// @Param id path string true "Catering ID"
// @Param dishId path string true "Dish ID"
// @Param body body request.AddDish false "Dish object"
// @Success 204 "Successfully updated"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/dishes/{dishId} [put]
func (d dish) Update(c *gin.Context) {
	var path types.PathDish
	var body domain.Dish

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	if err, code := dishRepo.Update(path, body); err != nil {
		utils.CreateError(code, err.Error(), c)
		return
	}

	c.Status(http.StatusNoContent)
}
