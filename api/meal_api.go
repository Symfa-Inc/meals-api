package api

import (
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Meal struct
type Meal struct{}

// NewMeal return pointer to meal struct
// with all methods
func NewMeal() *Meal {
	return &Meal{}
}

var mealService = services.NewMealService

// Add Creates meal for certain client
// @Summary Creates meal for certain client
// @Tags catering meals
// @Produce json
// @Param id path string false "Catering ID"
// @Param clientId path string false "Client ID"
// @Param payload body request.AddMeal false "meal reading"
// @Success 201 {object} request.AddMeal "created meal"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/clients/{clientId}/meals [post]
func (m Meal) Add(c *gin.Context) {
	var path types.PathClient
	var body request.AddMeal

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	user, _ := c.Get("user")

	result, code, err := mealService().Add(path, body, user)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// Get returns array of meals
// @Summary GetByRange list of categories with dishes for passed meal ID
// @Tags catering meals
// @Produce json
// @Param startDate query string false "Meal Start Date in 2020-01-01T00:00:00Z format"
// @Param endDate query string false "Meal End Date in 2020-01-01T00:00:00Z format"
// @Param id path string false "Catering ID"
// @Param clientId path string false "Client ID"
// @Success 200 {array} response.GetMeal "dishes for passed day"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/meals [get]
func (m Meal) Get(c *gin.Context) {
	var query types.DateRangeQuery
	var path types.PathClient

	if err := utils.RequestBinderURI(&path, c); err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	result, code, err := mealService().Get(query, path)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, result)
}
