package api

import (
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Dish struct
type Dish struct{}

// NewDish return pointer to dish struct
// with all methods

func NewDish() *Dish {
	return &Dish{}
}

var dishRepo = repository.NewDishRepo()

// Add adds dish for catering with provided ID
// @Summary Add dish for certain category
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param payload body request.AddDish false "dish object"
// @Success 200 {object} domain.Dish false "dish object"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dishes [post]
func (d Dish) Add(c *gin.Context) {
	var path types.PathID
	var dish domain.Dish

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&dish, c); err != nil {
		return
	}

	dish.CateringID, _ = uuid.FromString(path.ID)

	err := dishRepo.Add(path.ID, &dish)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	dish.Images = make([]domain.ImageArray, 0)

	c.JSON(http.StatusOK, dish)
}

// Delete soft delete of dish
// @Summary Soft delete
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param dishId path string true "Dish ID"
// @Success 204 "Successfully deleted"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/dishes/{dishId} [delete]
func (d Dish) Delete(c *gin.Context) {
	var path types.PathDish

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := dishRepo.Delete(path); err != nil {
		utils.CreateError(http.StatusNotFound, err, c)
	}

	c.Status(http.StatusNoContent)
}

// Get return list of dishes
// @Summary Returns list of dishes
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param categoryID query string true "Category ID"
// @Success 200 {array} domain.Dish "List of dishes"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dishes [get]
func (d Dish) Get(c *gin.Context) {
	var path types.PathID
	var query types.CategoryIDQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}
	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	dishes, code, err := dishRepo.Get(path.ID, query.CategoryID)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	if len(dishes) == 0 {
		c.JSON(http.StatusOK, make([]string, 0))
		return
	}

	c.JSON(http.StatusOK, dishes)
}

// GetByID return dishes
// @Summary Returns dishes
// @Tags catering dishes
// @Produce json
// @Param id path string true "Catering ID"
// @Param dishId path string true "Dish ID"
// @Success 200 {array} domain.Dish "List of dishes"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/dishes/{dishId} [get]
func (d Dish) GetByID(c *gin.Context) {
	var path types.PathDishID

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	dish, code, err := dishRepo.FindByID(path.ID, path.DishID)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, dish)
}

// Update updates dish
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
func (d Dish) Update(c *gin.Context) {
	var path types.PathDish
	var dish domain.Dish

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&dish, c); err != nil {
		return
	}

	if code, err := dishRepo.Update(path, dish); err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}
