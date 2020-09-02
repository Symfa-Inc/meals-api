package api

import (
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
)

// Meal struct
type Meal struct {}

// NewMeal return pointer to meal struct
// with all methods
func NewMeal() *Meal {
	return &Meal{}
}

var mealRepo = repository.NewMealRepo()
var mealDishRepo = repository.NewMealDishesRepo()
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

	_, code, err := mealService().Add(path, body, user)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}


}